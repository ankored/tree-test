package persist_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ankored/tree-test/tree"
	"github.com/ankored/tree-test/tree/persist"
	"github.com/google/go-cmp/cmp"
)

func TestCreatesFilesIfDoNotExist(t *testing.T) {
	t.Cleanup(func() {
		os.Remove("./test_log")
		os.Remove("./test_cmp")
	})

	if _, err := persist.NewFile("./test_cmp", "./test_log"); err != nil {
		t.Fatalf(" error: %s", err)
	}

	// Should have made the files
	if _, err := os.Stat("./test_cmp"); err != nil {
		t.Fatalf("stats test_cmp: %s", err)
	}
	if _, err := os.Stat("./test_log"); err != nil {
		t.Fatalf("stats test_log: %s", err)
	}
}

func TestRecordAppendsToLog(t *testing.T) {
	t.Cleanup(func() {
		os.Remove("./test_log")
		os.Remove("./test_cmp")
	})

	f, err := persist.NewFile("./test_cmp", "./test_log")
	if err != nil {
		t.Fatalf(" error creating file persistor: %s", err)
	}

	want := tree.Op{
		Path: "/some/thing",
		Node: &tree.Node{
			Type:  tree.NodeTypeVal,
			Value: "honk",
		},
	}
	if err := f.Record(want); err != nil {
		t.Fatalf(" error recording operation: %s", err)
	}

	// Read the log for the new entry
	bytes, err := ioutil.ReadFile("./test_log")
	if err != nil {
		t.Fatalf("error reading log file: %s", err)
	}

	got := tree.Op{}
	if err := json.Unmarshal(bytes, &got); err != nil {
		t.Fatalf("error unmarshalling operation from bytes: %s", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Tree mismatch (-want +got):\n%s", diff)
	}
}

func TestRestoreFromFile(t *testing.T) {
	t.Cleanup(func() {
		os.Remove("./test_log")
		os.Remove("./test_cmp")
	})

	n := &tree.Node{
		Type: tree.NodeTypeDir,
		Childs: map[string]*tree.Node{
			"geese": &tree.Node{
				Type: tree.NodeTypeDir,
				Childs: map[string]*tree.Node{
					"goose": &tree.Node{
						Type:  tree.NodeTypeVal,
						Value: "honk",
					},
				},
			},
		},
	}
	byts, _ := json.Marshal(n)

	fi, _ := os.OpenFile("./test_cmp", os.O_RDWR|os.O_CREATE, 0755)
	if _, err := fi.Write(byts); err != nil {
		t.Fatalf("error writing to compaction file: %s", err)
	}

	f, err := persist.NewFile("./test_cmp", "./test_log")
	if err != nil {
		t.Fatalf(" error creating file persistor: %s", err)
	}

	op := tree.Op{
		Path: "/geese/goose",
		Node: &tree.Node{
			Type:  tree.NodeTypeVal,
			Value: "honk",
		},
	}
	if err := f.Record(op); err != nil {
		t.Fatalf("error recording operation: %s", err)
	}

	tr, err := f.Restore()
	if err != nil {
		t.Fatalf("error restoring tree from compaction file: %s", err)
	}

	if diff := cmp.Diff(&tree.Node{
		Type: tree.NodeTypeDir,
		Childs: map[string]*tree.Node{
			"geese": &tree.Node{
				Type: tree.NodeTypeDir,
				Childs: map[string]*tree.Node{
					"goose": &tree.Node{
						Type:  tree.NodeTypeVal,
						Value: "honk",
					},
				},
			},
		},
	}, tr.Root); diff != "" {
		t.Errorf("Tree mismatch (-want +got):\n%s", diff)
	}
}
