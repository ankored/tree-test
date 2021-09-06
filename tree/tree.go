// Package tree provides a tree implementation that persists to a store
package tree

import "sync"

type OpType string

const (
	OpAdd     OpType = "add"
	OpRemove  OpType = "remove"
	OpReplace OpType = "replace"
)

// Op is an action performed on the tree
type Op struct {
	Type  OpType
	Path  string
	Value string
}

// Persistence is how the tree records itself during actions to be restored from later
type Persistence interface {
	Record(op Op) error
	Restore() (Tree, error)
}

// Node is an entry in a tree. It has a path name and can contain child nodes.
type Node struct {
	Name   string
	Value  string
	Childs map[string]Node
}

// Tree is a root node with a lock aroud it
type Tree struct {
	root Node
	mu   sync.RWMutex
}

func New(root Node) *Tree {
	return &Tree{
		root: root,
	}
}
