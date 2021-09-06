package tree_test

import (
	"testing"

	"github.com/ankored/tree-test/tree"
)

// func TestAddCreatesNode(t *testing.T) {
// 	tr := tree.New(&tree.Node{
// 		Type: tree.NodeTypeDir,
// 	})
//
// 	tr.AddChild("/child", &tree.Node{
// 		Type:  tree.NodeTypeVal,
// 		Value: "something",
// 	})
//
// 	if diff := cmp.Diff(&tree.Node{
// 		Type: tree.NodeTypeDir,
// 		Childs: map[string]*tree.Node{
// 			"child": &tree.Node{
// 				Type:  tree.NodeTypeVal,
// 				Value: "something",
// 			},
// 		},
// 	}, tr.Root); diff != "" {
// 		t.Errorf("Tree mismatch (-want +got):\n%s", diff)
// 	}
// }

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
