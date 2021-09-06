package persist

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ankored/tree-test/tree"
)

// File provides a persistence implementation that writes to disk
// using a write-ahead-log
type File struct {
	compacted *os.File // The compacted tree
	log       *os.File // The WAL file

	sizeLimit int64 // The size limit of the WAL in bytes, default is 10MB
	offset    int64 // The current offset for writing to the log file
}

type FileOpt func(*File)

// WithLogLimit makes an option that limits the log of the file
// to a given size, measured in KB.
func WithLogLimit(s int64) FileOpt {
	return func(f *File) {
		f.sizeLimit = s
	}
}

// NewFile creates a file persistence that writes to the root filename
// for storage a compacted version of the tree and logFn which stores the
// in-progress WAL.
func NewFile(comFn, logFn string, opts ...FileOpt) (*File, error) {
	var err error
	f := &File{
		sizeLimit: 1024 * 1024 * 10,
	}

	// Apply all options that were provided
	for _, opt := range opts {
		opt(f)
	}

	f.compacted, err = os.OpenFile(comFn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return nil, fmt.Errorf("error opening compacted file: %s", err)
	}

	f.log, err = os.OpenFile(logFn, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %s", err)
	}

	// Get size of the log to know where to start the offset
	stat, err := f.log.Stat()
	if err != nil {
		return nil, fmt.Errorf("error getting stats for log file: %s", err)
	}
	f.offset = stat.Size()

	return f, nil
}

// Close releases files
func (f *File) Close() error {
	var errs []error

	if err := f.compacted.Close(); err != nil {
		errs = append(errs, fmt.Errorf("error closing compaction file: %s", err))
	}

	if err := f.log.Close(); err != nil {
		errs = append(errs, fmt.Errorf("error closing log file: %s", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing files: %s", errs)
	}

	return nil
}

// Record appends to WAL to make sure the operation is persisted in case of a restart.
// When the file hits the limit, it also triggers a compaction and clears the log
// before returning.
func (f *File) Record(op tree.Op) error {
	byts, err := json.Marshal(op)
	if err != nil {
		return fmt.Errorf("error marshalling operation: %s", err)
	}

	// Append a delimiter so it can be restored line by line later
	byts = append(byts, []byte("\n")...)

	n, err := f.log.WriteAt(byts, f.offset)
	if err != nil {
		return fmt.Errorf("error writing to log: %s", err)
	}

	// Add the number of written bytes to update the log offset
	f.offset += int64(n)

	return nil
}

func (f *File) Restore() (*tree.Tree, error) {
	// Read the compaction file into a node to start the tree
	n := &tree.Node{}
	byts, err := ioutil.ReadAll(f.compacted)
	if err != nil {
		return nil, fmt.Errorf("error reading compaction file: %s", err)
	}
	if err := json.Unmarshal(byts, &n); err != nil {
		return nil, fmt.Errorf("error unmarshalling compaction contents into node: %s", err)
	}

	// Make the tree to apply operations to
	t := tree.New(n)

	// Read the log file line by line
	scanner := bufio.NewScanner(f.log)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		b := scanner.Bytes()

		// Try to get an operation out of it
		op := tree.Op{}
		if err := json.Unmarshal(b, &op); err != nil {
			return nil, fmt.Errorf("error unmarshalling operation: %s. Line: %s", err, string(b))
		}

		// Apply the operation
		if err := t.Put(op.Path, op.Node); err != nil {
			return nil, fmt.Errorf("error applying op %#v: %w", op, err)
		}
	}

	return t, nil
}
