package six

// Theres an assumption that no two nodes will have the same name throughout this implementation

type Tree struct {
	Root         *TreeNode
	LeavesByName map[string]*TreeNode
}

func NewTree() *Tree {
	return &Tree{
		LeavesByName: make(map[string]*TreeNode),
	}
}

func (t *Tree) AddRelationship(parentnode string, childnode string) {
	// lookup nodes to see if they are in tree
	var parent *TreeNode
	var child *TreeNode
	parent, ok := t.LeavesByName[parentnode]
	if !ok {
		parent = NewNode(parentnode)
		t.LeavesByName[parentnode] = parent
	}
	child, ok = t.LeavesByName[childnode]
	if !ok {
		child = NewNode(childnode)
		t.LeavesByName[childnode] = child
	}
	// add child to parent
	parent.AddChild(child)
	// update tree root if needed
	if t.Root == nil || child.Name == t.Root.Name {
		t.Root = parent
	}
}

func (t Tree) FindPath(startnode string, endnode string) []string {
	queue := newQueue()
	visited := make(map[string]bool)
	path := make([]string, 0)
	path = append(path, startnode)

	// add start to queue, then do BFS to find shortest path to endnode
	start := t.LeavesByName[startnode]
	queue.enqueue(start, path)

	for queue.size() > 0 {
		curnode := queue.dequeue()
		// mark current node as visited
		visited[curnode.node.Name] = true
		// check if we've reached our destination
		if curnode.node.Name == endnode {
			return curnode.path
		}

		neighbors := curnode.node.getNeighbors()
		for _, n := range neighbors {
			// enqueue the neighbors if they haven't been visited
			// AND they're not currently in queue
			if !visited[n.Name] && !queue.inQueue(n) {
				queue.enqueue(n, append(curnode.path, n.Name))
			}
		}
	}

	return nil
}

type TreeNode struct {
	Name     string
	Parent   *TreeNode
	Children map[string]*TreeNode
}

func NewNode(name string) *TreeNode {
	return &TreeNode{
		Name:     name,
		Children: make(map[string]*TreeNode),
	}
}

func (tn *TreeNode) AddChild(c *TreeNode) {
	_, ok := tn.Children[c.Name]
	if !ok {
		tn.Children[c.Name] = c
	}
	c.Parent = tn
}

func (tn TreeNode) getNeighbors() []*TreeNode {
	res := make([]*TreeNode, 0)

	if tn.Parent != nil {
		res = append(res, tn.Parent)
	}

	for _, n := range tn.Children {
		res = append(res, n)
	}

	return res
}

func GetOrbitCount(t *Tree) int {
	count := 0

	visitedMap := make(map[string]bool)

	for nodename, node := range t.LeavesByName {
		// check if node was visited
		visited := visitedMap[nodename]
		if !visited {
			// iterate through all parents, upping orbital count on the way
			// until reaching tree root
			// Note the number of parents along the way that have not been visited.
			orbitalCount := 0
			n := node

			parentVisitedDepth := 0
			for n.Parent != nil {
				orbitalCount++
				v := visitedMap[n.Name]
				if !v {
					parentVisitedDepth++
					visitedMap[n.Name] = true
				}
				n = n.Parent
			}
			// up count by orbitalCount + (orbitalCount -1 ) + ... (orbitalCount - parentVisitedDepth)
			// to count the orbits of the parents as well as node
			res := 0
			for i := 0; i < parentVisitedDepth; i++ {
				res += orbitalCount - i
			}

			count += res

		}
		visitedMap[nodename] = true
	}
	return count
}

func GetOrbitalTransfers(t *Tree, startSatellite string, endSiblingSatellite string) (int, []string) {
	startNode := t.LeavesByName[startSatellite]
	endNode := t.LeavesByName[endSiblingSatellite]

	startPath := startNode.Parent
	endPath := endNode.Parent

	path := t.FindPath(startPath.Name, endPath.Name)

	return len(path) - 1, path
}

type nodeQueue struct {
	queue      []*nodeQueueItem
	queueItems map[string]bool
}

type nodeQueueItem struct {
	node *TreeNode
	path []string
}

func newQueue() *nodeQueue {
	return &nodeQueue{
		queue:      make([]*nodeQueueItem, 0),
		queueItems: make(map[string]bool),
	}
}

func (n *nodeQueue) enqueue(node *TreeNode, path []string) {
	i := &nodeQueueItem{node: node, path: path}
	n.queue = append(n.queue, i)
	n.queueItems[node.Name] = true
}

func (n *nodeQueue) dequeue() *nodeQueueItem {
	// remove from front of queue, remove from nodeItems
	res := n.queue[0]
	delete(n.queueItems, res.node.Name)
	n.queue = n.queue[1:]
	return res
}

func (n nodeQueue) size() int {
	return len(n.queue)
}

func (n *nodeQueue) inQueue(node *TreeNode) bool {
	return n.queueItems[node.Name]
}
