package plan

type (
	// OpQueueTree represents a tree of operation queues
	OpQueueTree struct {
		Root *TreeNode
	}

	// TreeNode represents a node in the OpQueueTree
	TreeNode struct {
		children []*TreeNode
		parent   *TreeNode
		val      *OpQueue
	}
)

// NewOpQueueTree constructs an empty OpQueueTree
func NewOpQueueTree() *OpQueueTree {
	return &OpQueueTree{nil}
}

// AddChild adds a node containing an OpQueue as a child of an existing
// tree node.
func (t *OpQueueTree) AddChild(s *OpQueue, parent *TreeNode) (child *TreeNode) {
	child = &TreeNode{[]*TreeNode{}, parent, s}

	parent.children = append(parent.children, child)

	return child
}

// Prune removes all empty queues that don't have children from the OpQueueTree.
func (t *OpQueueTree) Prune() {
	if t.Root == nil {
		return
	}

	t.Root.pruneRecursive()
}

// GetNLeaves returns at most n leaves of the tree, representing the number of
// queues that can be computed in parallel.
func (t *OpQueueTree) GetNLeaves(n int) []*TreeNode {
	if t.Root == nil {
		return []*TreeNode{}
	}

	return getNLeavesRecursive(n, t.Root, []*TreeNode{})
}

// ----------------------------------------------------------------------------
// TreeNode helper methods
func (node *TreeNode) pruneRecursive() bool {
	if len(node.children) == 0 && node.val.Len() == 0 {
		return true
	}

	remove := []int{}

	for i, _node := range node.children {
		canRemove := _node.pruneRecursive()

		if canRemove {
			remove = append(remove, i)
		}
	}

	node.removeChildren(remove)

	// check again to see if removing children had an effect
	if len(node.children) == 0 && node.val.Len() == 0 {
		return true
	}

	return false
}

// removeChildren removes a set of children from an array, by swapping them all to
// the end of an array, and slicing out the undesired elements
func (node *TreeNode) removeChildren(remove []int) {
	a := node.children

	for _, i := range remove {
		a[i] = a[len(a)-1]
	}

	a = a[:len(a)-len(remove)]

	node.children = a
}

func getNLeavesRecursive(n int, node *TreeNode, currNodes []*TreeNode) []*TreeNode {
	if len(node.children) == 0 {
		currNodes = append(currNodes, node)
	} else {
		for _, _node := range node.children {
			currNodes = getNLeavesRecursive(n, _node, currNodes)

			if len(currNodes) == n {
				return currNodes
			}
		}
	}

	return currNodes
}
