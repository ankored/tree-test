package tree_test

import (
	"errors"
	"testing"

	"github.com/ankored/tree-test/tree"
	"github.com/google/go-cmp/cmp"
)

func TestAddCreatesNode(t *testing.T) {
	tr := tree.New(&tree.Node{
		Type: tree.NodeTypeDir,
	})

	err := tr.AddChild("/child", &tree.Node{
		Type:  tree.NodeTypeVal,
		Value: "something",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if diff := cmp.Diff(&tree.Node{
		Type: tree.NodeTypeDir,
		Childs: map[string]*tree.Node{
			"child": &tree.Node{
				Type:  tree.NodeTypeVal,
				Value: "something",
			},
		},
	}, tr.Root); diff != "" {
		t.Errorf("Tree mismatch (-want +got):\n%s", diff)
	}
}

func TestAddRequiresParentToBeDir(t *testing.T) {
	tr := tree.New(&tree.Node{
		Type: tree.NodeTypeDir,
		Childs: map[string]*tree.Node{
			"child": &tree.Node{
				Type: tree.NodeTypeVal,
			},
		},
	})

	err := tr.AddChild("/child/grandchild", &tree.Node{
		Type:  tree.NodeTypeVal,
		Value: "something",
	})
	var dirErr tree.ErrNotDir
	if !errors.As(err, &dirErr) {
		t.Fatalf("expected ErrNotDir, instead got: %s", err)
	}
}

func TestAddRequiresAllParentsToExist(t *testing.T) {
	tr := tree.New(&tree.Node{
		Type: tree.NodeTypeDir,
		Childs: map[string]*tree.Node{
			"child": &tree.Node{
				Type: tree.NodeTypeDir,
			},
		},
	})

	err := tr.AddChild("/child/not-here/child", &tree.Node{
		Type:  tree.NodeTypeVal,
		Value: "something",
	})

	var errNF tree.ErrNotFound
	if !errors.As(err, &errNF) {
		t.Fatalf("expected ErrNotFound, instead got: %s", err)
	}

	if errNF.Path != "child/not-here" {
		t.Fatalf("expected error path to be %s, instead got: %s", "child/not-here", errNF.Path)
	}
}

func TestAddNilRemovesNode(t *testing.T) {
	tr := tree.New(&tree.Node{
		Type: tree.NodeTypeDir,
		Childs: map[string]*tree.Node{
			"child": &tree.Node{
				Type: tree.NodeTypeVal,
			},
		},
	})

	err := tr.AddChild("/child", nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if diff := cmp.Diff(&tree.Node{
		Type:   tree.NodeTypeDir,
		Childs: map[string]*tree.Node{},
	}, tr.Root); diff != "" {
		t.Errorf("Tree mismatch (-want +got):\n%s", diff)
	}
}

func TestGetNode(t *testing.T) {
	tr := tree.New(&tree.Node{
		Type:  tree.NodeTypeDir,
		Value: "root",
		Childs: map[string]*tree.Node{
			"child": &tree.Node{
				Type:  tree.NodeTypeDir,
				Value: "child",
				Childs: map[string]*tree.Node{
					"grandchild": &tree.Node{
						Type:  tree.NodeTypeVal,
						Value: "honk",
					},
				},
			},
		},
	})

	n, err := tr.Get("/child/grandchild")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if n.Value != "honk" {
		t.Fatalf("expected value: honk, instead got: %s", n.Value)
	}
}
