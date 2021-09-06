package tree

import (
	"github.com/ankored/tree-test/tree"
)

// NoopPersist implements a persistence but does only noops when
// called.
type NoopPersist struct{}

func (NoopPersist) Record(tree.Op) error {
	return nil
}

func (NoopPersist) Restore() (*tree.Tree, error) {
	// Does nothing, just return an empty tree
	return tree.New(&tree.Node{
		Childs: map[string]*tree.Node{},
	}), nil
}
