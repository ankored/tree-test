package persist

import (
	"github.com/ankored/tree-test/tree"
)

// Noo implements a persistence but does only no-ops when
// called.
type Noop struct{}

func (Noop) Record(tree.Op) error {
	return nil
}

func (Noop) Restore() (*tree.Tree, error) {
	// Does nothing, just return an empty tree
	return tree.New(&tree.Node{
		Childs: map[string]*tree.Node{},
	}), nil
}
