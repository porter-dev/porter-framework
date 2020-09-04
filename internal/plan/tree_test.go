package plan

import (
	"testing"

	v "github.com/porterdev/ego/internal/value"
)

func TestEmpty(t *testing.T) {
	tree := NewOpQueueTree()

	// make sure empty prune doesn't crash
	tree.Prune()

	if leaves := tree.GetNLeaves(10); len(leaves) != 0 {
		t.Errorf("Expected 0 leaves, got %v leaves", len(leaves))
	}
}

// Constructs a tree that looks like:
//     root
//    /
//   c1
//
// Where c1 has a value of an empty queue
func initSimpleTree() (tree *OpQueueTree, children []*OpQueue) {
	tree = NewOpQueueTree()

	s1 := NewOpQueue()

	initQueue(s1)

	tree.Root = &TreeNode{
		val: s1,
	}

	s2 := NewOpQueue()

	tree.AddChild(s2, tree.Root)

	return tree, []*OpQueue{s2}
}

// Constructs a tree that looks like:
//     root
//    /    \
//   c1    c2
//  /  \  /  \
// c3 c4 c5  c6
//
// Where no children have empty queues
func initComplexTree() (tree *OpQueueTree, children []*OpQueue) {
	tree = NewOpQueueTree()

	s1 := NewOpQueue()

	initQueue(s1)

	tree.Root = &TreeNode{
		val: s1,
	}

	for i := 0; i < 6; i++ {
		s := NewOpQueue()

		initQueue(s)

		children = append(children, s)
	}

	tree.AddChild(children[0], tree.Root)
	tree.AddChild(children[1], tree.Root)
	tree.AddChild(children[2], tree.Root.children[0])
	tree.AddChild(children[3], tree.Root.children[0])
	tree.AddChild(children[4], tree.Root.children[1])
	tree.AddChild(children[5], tree.Root.children[1])

	return tree, children
}

func TestAddChild(t *testing.T) {
	tree, children := initSimpleTree()

	if len(tree.Root.children) != 1 {
		t.Errorf("Expected tree.Root to have 1 child, got %v", len(tree.Root.children))
	}

	s2 := children[0]

	child := tree.Root.children[0]

	if child.val != s2 {
		t.Errorf("Expected tree.Root.children[0].val to be %v, got %v", *s2, *child.val)
	}
}

func TestPruneSimple(t *testing.T) {
	tree, _ := initSimpleTree()

	if len(tree.Root.children) != 1 {
		t.Errorf("Expected tree.Root to have 1 child, got %v", len(tree.Root.children))
	}

	tree.Prune()

	if len(tree.Root.children) != 0 {
		t.Errorf("Expected tree.Root to have 0 children, got %v", len(tree.Root.children))
	}
}

func TestPruneDeepAllEmpty(t *testing.T) {
	tree, _ := initSimpleTree()

	s3 := NewOpQueue()
	s4 := NewOpQueue()

	tree.AddChild(s3, tree.Root.children[0])
	tree.AddChild(s4, tree.Root.children[0].children[0])

	tree.Prune()

	if len(tree.Root.children) != 0 {
		t.Errorf("Expected tree.Root to have 0 children, got %v", len(tree.Root.children))
	}
}

func TestPruneDeepNotAllEmpty(t *testing.T) {
	tree, _ := initSimpleTree()

	s3 := NewOpQueue()
	s4 := NewOpQueue()
	s5 := NewOpQueue()

	s5.Enqueue(&Operation{
		Op:   READ,
		Path: ".hello",
		Old:  v.String("hello"),
		New:  v.Value("there"),
	})

	tree.AddChild(s3, tree.Root.children[0])
	tree.AddChild(s4, tree.Root.children[0].children[0])
	tree.AddChild(s5, tree.Root.children[0].children[0])

	if len := len(tree.Root.children[0].children[0].children); len != 2 {
		t.Errorf("Expected tree.Root.children[0] to have 2 children, got %v", len)
	}

	tree.Prune()

	if len := len(tree.Root.children[0].children[0].children); len != 1 {
		t.Errorf("Expected tree.Root.children[0] to have 1 child, got %v", len)
	}

	if child := tree.Root.children[0].children[0].children[0]; *child.val != *s5 {
		t.Errorf("Expected leaf to be %v, was %v", s5, *child)
	}
}

func TestGetNLeaves(t *testing.T) {
	tree, children := initComplexTree()

	leaves1 := tree.GetNLeaves(3)

	if len := len(leaves1); len != 3 {
		t.Errorf("Expected GetNLeaves(3) to return 3 leaves, got %v", len)
	}

	if leaves1[0].val != children[2] {
		t.Errorf("Expected GetNLeaves(3)[0] to return %v, got %v", *children[2], *leaves1[0].val)
	}

	if leaves1[1].val != children[3] {
		t.Errorf("Expected GetNLeaves(3)[0] to return %v, got %v", *children[3], *leaves1[1].val)
	}

	if leaves1[2].val != children[4] {
		t.Errorf("Expected GetNLeaves(3)[0] to return %v, got %v", *children[4], *leaves1[2].val)
	}

	leaves2 := tree.GetNLeaves(4)

	if len := len(leaves2); len != 4 {
		t.Errorf("Expected GetNLeaves(4) to return 4 leaves, got %v", len)
	}

	leaves3 := tree.GetNLeaves(5)

	if len := len(leaves3); len != 4 {
		t.Errorf("Expected GetNLeaves(5) to return 4 leaves, got %v", len)
	}
}
