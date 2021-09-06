// Package tree provides a tree implementation that persists to a store
package tree

import (
	"fmt"
	"strings"
	"sync"
)

type ErrNotFound struct {
	Path string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("could not find path: %s", e.Path)
}

type ErrNotDir struct {
	Path string
}

func (e ErrNotDir) Error() string {
	return fmt.Sprintf("path is not a directory: %s", e.Path)
}

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

type NodeType string

const (
	NodeTypeDir = "dir"
	NodeTypeVal = "val"
)

// Node is an entry in a tree. It has a path name and can contain child nodes.
type Node struct {
	Type   NodeType
	Value  string
	Childs map[string]*Node
}

// Tree is a root node with a lock aroud it
type Tree struct {
	Root *Node
	mu   sync.RWMutex
}

func New(root *Node) *Tree {
	if root.Type != NodeTypeDir {
		panic("root must be a directory")
	}

	return &Tree{
		Root: root,
	}
}

// AddChild adds a node to the tree at the given path
func (t *Tree) AddChild(path string, n *Node) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// TODO: Get the parent node being referenced by the path
}

func (t *Tree) Get(path string) (*Node, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	n, _, err := t.traverse(path)
	if err != nil {
		return nil, fmt.Errorf("error getting node: %w", err)
	}

	return n, nil
}

// Grabs the node referenced in the path (plus traversed nodes to get there),
// or returns a NotFound error where it failed to find a node
func (t *Tree) traverse(path string) (*Node, []*Node, error) {
	var fn func([]string, int, *Node, []*Node) (*Node, []*Node, error)
	fn = func(parts []string, i int, n *Node, traversed []*Node) (*Node, []*Node, error) {
		if i >= len(parts) {
			// This is base case, there's nothing more to recurse
			return n, traversed, nil
		}

		// We have more to traverse, make sure the current node is a directory
		if n.Type != NodeTypeDir {
			return nil, nil, ErrNotDir{
				Path: strings.Join(parts[0:i], "/"),
			}
		}

		// Find the child node from the current one
		c, ok := n.Childs[parts[i]]
		if !ok {
			return nil, nil, ErrNotFound{
				Path: strings.Join(parts[0:i+1], "/"),
			}
		}

		// Continue going by increasing i to look at the next path part
		// and use the child node as the next current node
		return fn(parts, i+1, c, append(traversed, n))
	}

	path = strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/")
	parts := strings.Split(path, "/")

	return fn(parts, 0, t.Root, make([]*Node, 0, len(parts)))
}

// Gets the parent node being referenced by the path
func parentInPath(path string) string {
	parts := strings.Split(path, "/")
	return "/" + strings.Join(parts[0:len(parts)-1], "/")
}
