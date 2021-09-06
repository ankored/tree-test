package persist

import (
	"errors"
	"fmt"
	"os"

	"github.com/ankored/tree-test/tree"
)

// File provides a persistence implementation that writes to disk
// using a write-ahead-log
type File struct {
	compacted *os.File // The compacted tree
	log       *os.File // The WAL file

	sizeLimit int64 // The size limit of the WAL, default is 10MB
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

	f.compacted, err = os.OpenFile(comFn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return nil, fmt.Errorf("error opening compacted file: %s", err)
	}

	f.log, err = os.OpenFile(logFn, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %s", err)
	}

	return f, nil
}

func (f *File) Record(op tree.Op) error {
	return errors.New("unimplemented")
}

func (f *File) Restore() (*tree.Tree, error) {
	return nil, errors.New("unimplemented")
}
