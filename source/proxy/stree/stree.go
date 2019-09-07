package stree

//move [] leaves rather than ()

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/golang-collections/go-datastructures/queue"
)

var (
	err error

	defaultOrder = 4
	minOrder     = 3
	maxOrder     = 20
	defaultBack  = 4
	minBack      = 0
	maxBack      = 20

	order          = defaultOrder
	back           = defaultBack
	myqueue        *Node
	verbose_output = false
	version        = 0.1
)

type Tree struct {
	Order   int
	NumKeys int
	Root    *Node
}

type Record struct {
	Value []byte
}

type Node struct {
	Pointers     []interface{}
	Keys         []int
	Parent       *Node
	IsLeaf       bool
	Backpointers []interface{}
	Backkeys     []int
	Next         *Node
}

func NewTree() *Tree {
	t := &Tree{}
	t.Order = order
	t.NumKeys = 0
	return t
}

//ok
func (t *Tree) Insert(rowID int, pos int) error {
	t.NumKeys++
	if t.Root == nil {
		return t.startNewTree(rowID)
	}
	NodeNum := t.NodeSumKeys(t.Root)
	if pos > NodeNum+1 {
		pos = NodeNum + 1
	}
	var leaf *Node

	leaf, PreIndexInLeaf, err := t.findLeaf(pos-1, false)
	if err != nil {
		return errors.New("error happens, insertion stops")
	}

	if t.NodeNumKeys(leaf) < order {
		t.insertIntoLeaf(leaf, rowID, PreIndexInLeaf+1)
		return nil
	}
	return t.insertIntoLeafAfterSplitting(leaf, rowID, PreIndexInLeaf+1)
}

//return Leaf, indexInLeaf, rowID, err
//ok
func (t *Tree) Query(pos int, verbose bool) (*Node, int, int, error) {
	NodeNum := t.NodeSumKeys(t.Root)
	if pos > NodeNum {
		return nil, 0, 0, errors.New("position exceeds total tuple number")
	}
	if pos == 0 {
		return nil, 0, 0, errors.New("cannot query the 0th tuple")
	}
	Leaf, indexInLeaf, err := t.findLeaf(pos, verbose)
	if err != nil {
		return Leaf, indexInLeaf, 0, err
	}
	return Leaf, indexInLeaf, Leaf.Keys[indexInLeaf], err
}

//leftmovedone
func (t *Tree) Rank(Leaf *Node, rowID int) int {
	rank := 0
	for i := 0; i < order; i++ {
		if Leaf.Keys[i] == rowID {
			rank++
			break
		} else if Leaf.Keys[i] != 0 {
			rank++
		} else {
			break
		}
	}
	c := Leaf
	for c != t.Root {
		parent := c.Parent
		for i := 0; i < order; i++ {
			if parent.Pointers[i] != c && parent.Pointers[i] != nil {
				rank += parent.Keys[i]
			} else if parent.Pointers[i] == c {
				c = parent
				break
			}
		}
	}
	return rank
}

/*
func (t *Tree) Find(key int, verbose bool) (*Record, error) {
	i := 0
	c := t.findLeaf(key, verbose)
	if c == nil {
		return nil, errors.New("key not found")
	}
	for i = 0; i < c.NumKeys; i++ {
		if c.Keys[i] == key {
			break
		}
	}
	if i == c.NumKeys {
		return nil, errors.New("key not found")
	}

	r, _ := c.Pointers[i].(*Record)

	return r, nil
}
*/

/*
func (t *Tree) FindAndPrint(key int, verbose bool) {
	r, err := t.Find(key, verbose)

	if err != nil || r == nil {
		fmt.Printf("Record not found under key %d.\n", key)
	} else {
		fmt.Printf("Record at %d -- key %d, value %s.\n", r, key, r.Value)
	}
}

*/

/*
func (t *Tree) FindAndPrintRange(key_start, key_end int, verbose bool) {
	var i int
	array_size := key_end - key_start + 1
	returned_keys := make([]int, array_size)
	returned_pointers := make([]interface{}, array_size)
	num_found := t.findRange(key_start, key_end, verbose, returned_keys, returned_pointers)
	if num_found == 0 {
		fmt.Println("None found,\n")
	} else {
		for i = 0; i < num_found; i++ {
			c, _ := returned_pointers[i].(*Record)
			fmt.Printf("Key: %d  Location: %d  Value: %s\n",
				returned_keys[i],
				returned_pointers[i],
				c.Value)
		}
	}
}

*/
func (t *Tree) PrintSubTree(s *Node) {
	fmt.Printf("\n")
	var n *Node
	i := 0
	rank := 0
	new_rank := 0

	if s == nil {
		fmt.Printf("Empty subtree.\n")
		return
	}
	myqueue = nil
	enqueue(s)
	for myqueue != nil {
		n = dequeue()
		if n != nil {
			//if n.Parent != nil && n == n.Parent.Pointers[0] {
			if n.Parent != nil {
				rank = new_rank
				new_rank = t.pathToNode(n, s) // height
				if new_rank != rank {
					fmt.Printf("\n")
				}
			}
			for i = 0; i < order; i++ {
				if verbose_output {
					fmt.Printf("%d ", n.Pointers[i])
				}
				fmt.Printf("%d ", n.Keys[i])
			}
			if !n.IsLeaf {
				for i = 0; i < order; i++ {
					if n.Keys[i] != 0 {
						c, _ := n.Pointers[i].(*Node)
						enqueue(c)
					}
				}
			}
			fmt.Printf(" | ")
		}
	}
	fmt.Printf("\n\n")
}

//ok
func (t *Tree) PrintTree() {
	fmt.Printf("\n")
	var n *Node
	i := 0
	rank := 0
	new_rank := 0

	if t.Root == nil {
		fmt.Printf("Empty tree.\n")
		return
	}
	myqueue = nil
	enqueue(t.Root)
	for myqueue != nil {
		n = dequeue()
		if n != nil {

			//if n.Parent != nil && n == n.Parent.Pointers[0] {
			if n.Parent != nil {
				rank = new_rank
				new_rank = t.pathToRoot(n) // height
				if new_rank != rank {
					fmt.Printf("\n")
				}
			}
			for i = 0; i < order; i++ {
				if verbose_output {
					fmt.Printf("%d ", n.Pointers[i])
				}
				fmt.Printf("%d ", n.Keys[i])
			}
			if !n.IsLeaf {
				for i = 0; i < order; i++ {
					if n.Keys[i] != 0 {
						c, _ := n.Pointers[i].(*Node)
						enqueue(c)
					}
				}
			}
			fmt.Printf(" | ")
		}
	}
	fmt.Printf("\n\n")
}

/* // s-tree doesn't have pointers between leaves
func (t *Tree) PrintLeaves() {
	if t.Root == nil {
		fmt.Printf("Empty tree.\n")
		return
	}

	var i int
	c := t.Root
	for !c.IsLeaf {
		c, _ = c.Pointers[0].(*Node)
	}

	for {
		for i = 0; i < c.NumKeys; i++ {
			if verbose_output {
				fmt.Printf("%d ", c.Pointers[i])
			}
			fmt.Printf("%d ", c.Keys[i])
		}
		if verbose_output {
			fmt.Printf("%d ", c.Pointers[order-1])
		}
		if c.Pointers[order-1] != nil {
			fmt.Printf(" | ")
			c, _ = c.Pointers[order-1].(*Node)
		} else {
			break
		}
	}
	fmt.Printf("\n")
}
*/

func (t *Tree) NodeSumKeys(n *Node) int {
	if n.IsLeaf {
		return t.NodeNumKeys(n)
	}
	ans := 0
	for i := 0; i < order; i++ {
		if n.Keys[i] != 0 {
			ans += n.Keys[i]
		} else {
			break
		}
	}
	return ans
}

// ok (rank = pos)
func (t *Tree) DeleteByRank(rank int) error {
	t.NumKeys--
	NodeNum := t.NodeSumKeys(t.Root)
	if rank > NodeNum {
		rank = NodeNum
	}
	// return Leaf, indexInLeaf
	Leaf, indexInLeaf, err := t.findLeaf(rank, false)
	if err != nil {
		return err
	}
	t.deleteEntry(Leaf, indexInLeaf)
	return nil
}

//ok
func (t *Tree) DeleteByValue(Leaf *Node, rowID int) error {
	t.NumKeys--
	rank := t.Rank(Leaf, rowID)
	return t.DeleteByRank(rank)
}

func enqueue(new_node *Node) {
	/*for i := 0; i < order; i ++ {
		fmt.Printf("%d ", new_node.Keys[i])
	}
	fmt.Printf("\n")

	*/
	var c *Node
	if myqueue == nil {
		myqueue = new_node
		myqueue.Next = nil
	} else {
		c = myqueue
		for c.Next != nil {
			c = c.Next
		}
		c.Next = new_node
		new_node.Next = nil
	}
}

func dequeue() *Node {
	n := myqueue
	myqueue = myqueue.Next
	n.Next = nil
	return n
}

func (t *Tree) height() int {
	h := 0
	c := t.Root
	for !c.IsLeaf {
		c, _ = c.Pointers[0].(*Node)
		h++
	}
	return h
}

func (t *Tree) pathToRoot(child *Node) int {
	length := 0
	c := child
	for c != t.Root {
		c = c.Parent
		length += 1
	}
	return length
}

func (t *Tree) pathToNode(child *Node, to *Node) int {
	length := 0
	c := child
	for c != to {
		c = c.Parent
		length += 1
	}
	return length
}

func (t *Tree) PrintVerbose(c *Node) {
	if c.IsLeaf {
		fmt.Printf("Leaf [")
	} else {
		fmt.Printf("[")
	}
	for i := 0; i < order; i++ {
		if c.Keys[i] != 0 {
			fmt.Printf("%d ", c.Keys[i])
		}
	}
	fmt.Printf("|")
	for i := 0; i < back; i++ {
		if c.Backkeys[i] != 0 {
			fmt.Printf("%d ", c.Backkeys[i])
		}
	}
	fmt.Printf("]\n")
}

//leftmovedone
func (t *Tree) NodeNumKeys(c *Node) int {
	num := 0
	for i := 0; i < order; i++ {
		if c.Keys[i] != 0 {
			num++
		} else {
			break
		}
	}
	return num
}

// verbose: whether need output
// find the leaf according to position (rank)
// ok
// return Leaf, indexInLeaf
func (t *Tree) findLeaf(pos int, verbose bool) (*Node, int, error) {
	flag := false
	c := t.Root
	p := pos
	if pos == 0 {
		p = 1
		flag = true
	}
	if c == nil {
		if verbose {
			fmt.Printf("Empty tree.\n")
		}
		return c, 0, errors.New("error: Cannot find a tuple in an empty tree")
	}
	for !c.IsLeaf {
		if verbose {
			t.PrintVerbose(c)
		}
		i := 0
		for ; i < order; i++ {
			if p <= c.Keys[i] {
				c, _ = c.Pointers[i].(*Node)
				break
			}
			p -= c.Keys[i]
		}
		if verbose {
			fmt.Printf("%d ->\n", i)
		}
		//c, _ = c.Pointers[i].(*Node)
	}
	if verbose {
		t.PrintVerbose(c)
	}
	if flag == true {
		return c, -1, nil
	}
	j := 0
	i := 0
	for i = 0; i < order; i++ {
		if c.Keys[i] != 0 {
			j++
			if j == p {
				break
			}
		}
	}
	return c, i, nil
}

func cut(length int) int {
	if length%2 == 0 {
		return length / 2
	}

	return length/2 + 1
}

//
//	INSERTION
//
func makeRecord(value []byte) (*Record, error) {
	new_record := new(Record)
	if new_record == nil {
		return nil, errors.New("Error: Record creation.")
	} else {
		new_record.Value = value
	}
	return new_record, nil
}

func MakeNode() (*Node, error) {
	new_node := new(Node)
	if new_node == nil {
		return nil, errors.New("Error: Node creation.")
	}
	new_node.Keys = make([]int, order)
	if new_node.Keys == nil {
		return nil, errors.New("Error: New node keys array.")
	}
	new_node.Pointers = make([]interface{}, order)
	if new_node.Keys == nil {
		return nil, errors.New("Error: New node pointers array.")
	}
	new_node.Backpointers = make([]interface{}, back)
	if new_node.Backpointers == nil {
		return nil, errors.New("Error: New node backpointers array.")
	}
	new_node.Backkeys = make([]int, back)
	if new_node.Backkeys == nil {
		return nil, errors.New("Error: New node backkeys array.")
	}
	new_node.IsLeaf = false
	new_node.Parent = nil
	new_node.Next = nil
	return new_node, nil
}

//ok
func makeLeaf() (*Node, error) {
	leaf, err := MakeNode()
	if err != nil {
		return nil, err
	}
	leaf.IsLeaf = true
	return leaf, nil
}

//ok
func getLeftIndex(parent, left *Node) int {
	left_index := 0
	for left_index < order && parent.Pointers[left_index] != left {
		left_index += 1
	}
	return left_index
}

//ok
func (t *Tree) insertIntoLeaf(leaf *Node, key int, indexInLeaf int) {
	keyToStore := key
	var tmpKey int
	for i := indexInLeaf; i < order; i++ {
		if leaf.Keys[i] != 0 {
			tmpKey = leaf.Keys[i]
			leaf.Keys[i] = keyToStore
			keyToStore = tmpKey
		} else {
			leaf.Keys[i] = keyToStore
			break
		}
	}

	if leaf == t.Root {
		return
	}
	c := leaf
	parent := c.Parent
	for parent != nil {
		for i := 0; i < order; i++ {
			if reflect.DeepEqual(parent.Pointers[i], c) {
				parent.Keys[i]++
				c = parent
				parent = c.Parent
				break
			}
		}
	}
	return
}

//leftmovedone
//oksum
//t.insertIntoLeafAfterSplitting(leaf, rowID, PreIndexInLeaf+1)
func (t *Tree) insertIntoLeafAfterSplitting(leaf *Node, key int, indexInLeaf int) error {
	//fmt.Printf("Here is insertIntoLeafAfterSplitting, key = %d, indexInLeaf = %d\n", key, indexInLeaf)

	var new_leaf *Node
	var insertion_index, split, i, j int
	var err error
	var leftNumKeys, rightNumKeys int

	new_leaf, err = makeLeaf()
	if err != nil {
		return nil
	}

	temp_keys := make([]int, order+1)
	if temp_keys == nil {
		return errors.New("Error: Temporary keys array.")
	}
	temp_pointers := make([]interface{}, order+1)
	if temp_pointers == nil {
		return errors.New("Error: Temporary pointers array.")
	}

	insertion_index = indexInLeaf
	j = 0
	for i = 0; i < order; i++ {
		if j == insertion_index {
			j += 1
		}
		if leaf.Keys[i] != 0 {
			temp_keys[j] = leaf.Keys[i]
			temp_pointers[j] = leaf.Pointers[i]
			leaf.Keys[i] = 0
			leaf.Pointers[i] = nil
			j += 1
		} else {
			break
		}
	}

	temp_keys[insertion_index] = key
	temp_pointers[insertion_index] = nil
	split = cut(order)
	leftNumKeys = split

	for i = 0; i < split; i++ {
		leaf.Keys[i] = temp_keys[i]
		leaf.Pointers[i] = temp_pointers[i]
	}

	j = 0
	for i = split; i < order+1; i++ {
		new_leaf.Keys[j] = temp_keys[i]
		new_leaf.Pointers[j] = temp_pointers[j]
		j += 1
	}
	rightNumKeys = j
	new_leaf.Parent = leaf.Parent
	new_leaf.IsLeaf = true

	return t.insertIntoParent(leaf, leftNumKeys, new_leaf, rightNumKeys)
}

// insert right after left_index, modify right key and go upwards
func (t *Tree) insertIntoNode(n *Node, left_index int, sumKeysRight int, right *Node) {
	var i int
	key_to_store := sumKeysRight
	pointer_to_store := right
	for i = left_index + 1; i < order; i++ {
		if n.Keys[i] != 0 {
			tmp_key := n.Keys[i]
			tmp_pointer := n.Pointers[i]
			n.Keys[i] = key_to_store
			n.Pointers[i] = pointer_to_store
			key_to_store = tmp_key
			pointer_to_store = tmp_pointer.(*Node)
		} else {
			n.Keys[i] = key_to_store
			n.Pointers[i] = pointer_to_store
			break
		}
	}
	c := n
	parent := c.Parent
	for parent != nil {
		for i = 0; i < order; i++ {
			if reflect.DeepEqual(parent.Pointers[i], c) {
				parent.Keys[i] = t.NodeSumKeys(c)
				if c.IsLeaf {
					parent.Keys[i] = t.NodeNumKeys(c)
				}
				c = parent
				parent = c.Parent
				break
			}
		}
	}
}

//leftmovedone
//oksum
//return t.insertIntoNodeAfterSplitting(parent, left_index, numKeysRight, right)
func (t *Tree) insertIntoNodeAfterSplitting(parent *Node, left_index, sumKeysRight int, right *Node) error {
	var i, j, split, LeftSumKeys, RightSumKeys int
	var new_node, child *Node
	var temp_keys []int
	var temp_pointers []interface{}
	var err error

	temp_pointers = make([]interface{}, order+1)
	if temp_pointers == nil {
		return errors.New("Error: Temporary pointers array for splitting nodes.")
	}

	temp_keys = make([]int, order+1)
	if temp_keys == nil {
		return errors.New("Error: Temporary keys array for splitting nodes.")
	}

	j = 0
	for i = 0; i < order; i++ {
		if j == left_index+1 {
			j += 1
		}
		if parent.Keys[i] != 0 {
			temp_pointers[j] = parent.Pointers[i]
			temp_keys[j] = parent.Keys[i]
			parent.Keys[i] = 0
			parent.Pointers[i] = nil
			j++
		} else {
			break
		}
	}
	temp_pointers[left_index+1] = right
	temp_keys[left_index+1] = sumKeysRight

	split = cut(order)
	new_node, err = MakeNode()
	if err != nil {
		return err
	}
	for i = 0; i < split; i++ {
		LeftSumKeys += temp_keys[i]
		parent.Pointers[i] = temp_pointers[i]
		parent.Keys[i] = temp_keys[i]
	}
	j = 0
	for ; i < order+1; i++ {
		RightSumKeys += temp_keys[i]
		new_node.Pointers[j] = temp_pointers[i]
		new_node.Keys[j] = temp_keys[i]
		child, _ = new_node.Pointers[j].(*Node)
		child.Parent = new_node
		j += 1
	}
	new_node.Parent = parent.Parent

	return t.insertIntoParent(parent, LeftSumKeys, new_node, RightSumKeys)
}

//insert right after left
func (t *Tree) insertIntoParent(left *Node, sumKeysLeft int, right *Node, sumKeysRight int) error {
	parent := left.Parent

	if parent == nil {
		// insert left, right into a new root
		return t.insertIntoNewRoot(left, sumKeysLeft, right, sumKeysRight)
	}
	left_index := getLeftIndex(parent, left) //index of leaf, starting from 0
	parent.Keys[left_index] = sumKeysLeft

	if t.NodeNumKeys(parent) < order {
		// insert node "right" into parent, index is left_index + 1
		t.insertIntoNode(parent, left_index, sumKeysRight, right) // right key and keys upwards modified
		return nil
	}

	// split parent, and insert node "right" into parent, index is left_index + 1
	return t.insertIntoNodeAfterSplitting(parent, left_index, sumKeysRight, right)
}

//ok sum
func (t *Tree) insertIntoNewRoot(left *Node, sumKeysLeft int, right *Node, sumKeysRight int) error {
	t.Root, err = MakeNode()
	if err != nil {
		return err
	}
	t.Root.Keys[0] = sumKeysLeft
	t.Root.Pointers[0] = left
	t.Root.Keys[1] = sumKeysRight
	t.Root.Pointers[1] = right
	t.Root.Parent = nil
	left.Parent = t.Root
	right.Parent = t.Root
	return nil
}

//ok
func (t *Tree) startNewTree(rowID int) error {
	t.Root, err = makeLeaf()
	if err != nil {
		return err
	}
	t.Root.Keys[0] = rowID
	t.Root.Pointers[0] = nil
	t.Root.Parent = nil
	return nil
}

func getLeftNeighborIndex(n *Node) int {
	var i, nei int
	nei = -1
	parent := n.Parent
	for i = 0; i < order; i++ {
		if reflect.DeepEqual(parent.Pointers[i], n) {
			return nei
		}
		if parent.Pointers[i] != nil {
			nei = i
		}
	}

	return i
}

func (t *Tree) getRightNeighborIndex(n *Node) int {
	var i int
	parent := n.Parent
	for i = 0; i < order; i++ {
		if reflect.DeepEqual(parent.Pointers[i], n) {
			if parent.Pointers[i+1] != nil {
				return i + 1
			} else {
				return order
			}
		}
	}
	return -1
}

// ??
func removeEntryFromNode(n *Node, indexInNode int) *Node {
	n.Keys[indexInNode] = 0
	n.Pointers[indexInNode] = nil
	j := indexInNode
	for i := j + 1; i < order; i++ {
		if n.Keys[i] != 0 {
			n.Keys[j] = n.Keys[i]
			n.Pointers[j] = n.Pointers[i]
			n.Keys[i] = 0
			n.Pointers[i] = nil
			j++
		} else {
			break
		}
	}
	return n
}

// can optimize
func (t *Tree) adjustRoot() {
	var new_root *Node

	if t.NodeNumKeys(t.Root) > 1 {
		return
	}

	if !t.Root.IsLeaf {
		for i := 0; i < order; i++ { // don't need this loop
			if t.Root.Pointers[i] != nil {
				new_root, _ = t.Root.Pointers[i].(*Node)
				break
			}
		}
		new_root.Parent = nil
		t.Root = new_root
	}
	return
}

//ok
// t.coalesceNodes(n, neighbour, neighbour_index)
func (t *Tree) coalesceNodes(n, neighbour *Node, neighbour_index int, calledfromswap ...bool) *Node {
	var i, j, neighbour_insertion_index int
	var tmp *Node
	n_index := neighbour_index + 1
	//set neighbor to be the node on the left
	if neighbour_index == -1 {
		tmp = n
		n = neighbour
		neighbour = tmp
		n_index = 1
		neighbour_index = 0
	}
	neighbour_insertion_index = t.NodeNumKeys(neighbour)
	i = neighbour_insertion_index
	newParentKey := neighbour.Parent.Keys[neighbour_index] + neighbour.Parent.Keys[n_index] - 1
	if calledfromswap != nil && calledfromswap[0] == true {
		//newParentKey ++
		if neighbour.IsLeaf {
			newParentKey = t.NodeNumKeys(neighbour) + t.NodeNumKeys(n)
		} else {
			newParentKey = t.NodeSumKeys(neighbour) + t.NodeSumKeys(n)
		}

	}
	//sj := t.NodeNumKeys(neighbour)
	for j = 0; j < order; j++ {
		if n.Keys[j] != 0 {
			neighbour.Keys[i] = n.Keys[j]
			neighbour.Pointers[i] = n.Pointers[j]
			tmp, _ = neighbour.Pointers[i].(*Node)
			if tmp != nil {
				tmp.Parent = neighbour
			}
			i += 1
		} else {
			break
		}
	}
	neighbour.Parent.Keys[neighbour_index] = newParentKey

	t.deleteEntry(n.Parent, n_index)
	return neighbour
}

// ok
// redistributeNodes(n, neighbour, neighbour_index)
func (t *Tree) redistributeNodes(n, neighbour *Node, neighbour_index int) error {
	var i, leftNumKeys, rightNumKeys, leftIndex, leftSumKeys, rightSumKeys int

	min_keys := cut(order)
	if neighbour_index != -1 { // neighbor (more keys) is on the left
		leftIndex = neighbour_index
		for i = 0; i < min_keys; i++ {
			leftSumKeys += neighbour.Keys[i]
		}
		indexToGoRight := min_keys
		leftNumKeys = min_keys

		temp_keys := make([]int, order)
		if temp_keys == nil {
			return errors.New("Error: Temporary keys array.")
		}

		temp_pointers := make([]interface{}, order)
		if temp_pointers == nil {
			return errors.New("Error: Temporary pointers array.")
		}

		tmp_j := 0
		for i = indexToGoRight; i < order; i++ {
			if neighbour.Keys[i] != 0 {
				temp_pointers[tmp_j] = neighbour.Pointers[i]
				temp_keys[tmp_j] = neighbour.Keys[i]
				rightSumKeys += neighbour.Keys[i]
				tmp_j++
				neighbour.Keys[i] = 0
				neighbour.Pointers[i] = nil
			} else {
				break
			}
		}
		rightSumKeys += t.NodeSumKeys(n)
		for i = 0; i < order; i++ {
			if n.Keys[i] != 0 {
				temp_keys[tmp_j] = n.Keys[i]
				temp_pointers[tmp_j] = n.Pointers[i]
				tmp_j++
				n.Pointers[i] = nil
				n.Keys[i] = 0
			} else {
				break
			}
		}
		rightNumKeys = tmp_j
		for i = 0; i < order; i++ {
			if i == tmp_j {
				break
			}
			n.Pointers[i] = temp_pointers[i]
			n.Keys[i] = temp_keys[i]
			if n.Pointers[i] != nil {
				s := n.Pointers[i].(*Node)
				s.Parent = n
			}
		}
	} else { // neighbor is on the right (more keys) (neighbor_index == -1)
		leftIndex = 0
		j := 0
		for i = order - 1; i >= 0; i-- {
			if neighbour.Keys[i] != 0 {
				rightSumKeys += neighbour.Keys[i]
				j++
				if j == min_keys {
					break
				}
			}
		}
		indexToGoLeft := i - 1
		rightNumKeys = min_keys

		temp_keys := make([]int, order)
		if temp_keys == nil {
			return errors.New("Error: Temporary keys array.")
		}

		temp_pointers := make([]interface{}, order)
		if temp_pointers == nil {
			return errors.New("Error: Temporary pointers array.")
		}

		tmp_j := t.NodeNumKeys(n)
		leftSumKeys = t.NodeSumKeys(n)
		for i = 0; i <= indexToGoLeft; i++ {
			n.Pointers[tmp_j] = neighbour.Pointers[i]
			n.Keys[tmp_j] = neighbour.Keys[i]
			leftSumKeys += n.Keys[tmp_j]
			if n.Pointers[tmp_j] != nil {
				s := n.Pointers[tmp_j].(*Node)
				s.Parent = n
			}
			neighbour.Keys[i] = 0
			neighbour.Pointers[i] = nil
			tmp_j++
		}
		leftNumKeys = tmp_j
		k := indexToGoLeft + 1
		for i = 0; i < min_keys; i++ {
			neighbour.Keys[i] = neighbour.Keys[k]
			neighbour.Pointers[i] = neighbour.Pointers[k]
			neighbour.Keys[k] = 0
			neighbour.Pointers[k] = nil
			k++
		}
	}
	c := n.Parent
	if n.IsLeaf {
		c.Keys[leftIndex] = leftNumKeys
		c.Keys[leftIndex+1] = rightNumKeys
	} else {
		c.Keys[leftIndex] = leftSumKeys //!!!
		c.Keys[leftIndex+1] = rightSumKeys
	}

	parent := c.Parent
	for parent != nil {
		for i = 0; i < order; i++ {
			if reflect.DeepEqual(parent.Pointers[i], c) {
				parent.Keys[i]--
				c = c.Parent
				parent = c.Parent
				break
			}
		}
	}
	return nil
}

//ok
func (t *Tree) deleteEntry(n *Node, indexInNode int) {
	var min_keys, neighbour_index, capacity int
	var neighbour *Node

	n = removeEntryFromNode(n, indexInNode)

	if n == t.Root {
		t.adjustRoot()
		return
	}

	min_keys = cut(order)

	nNumKeys := t.NodeNumKeys(n)
	if nNumKeys >= min_keys {
		tmp := n
		for tmp != t.Root {
			for i := 0; i < order; i++ {
				if reflect.DeepEqual(tmp.Parent.Pointers[i], tmp) {
					tmp.Parent.Keys[i] = t.NodeSumKeys(tmp)
					if tmp.IsLeaf {
						tmp.Parent.Keys[i] = t.NodeNumKeys(tmp)
					}
					tmp = tmp.Parent
					break
				}
			}
		}
		return
	}

	neighbour_index = getLeftNeighborIndex(n)

	if neighbour_index == -1 {
		neighbour, _ = n.Parent.Pointers[1].(*Node)
	} else {
		neighbour, _ = n.Parent.Pointers[neighbour_index].(*Node)
	}

	capacity = order
	if t.NodeNumKeys(neighbour)+nNumKeys <= capacity {
		t.coalesceNodes(n, neighbour, neighbour_index)
		return
	} else {
		t.redistributeNodes(n, neighbour, neighbour_index)
		return
	}
}

//range leaves enqueue, encluding non-complete leaves
func (t *Tree) rangeLeafEnqueue(start int, end int, c *Node) (int, int, *queue.Queue) {
	LeafNum := t.NodeNumKeys(c)
	min_keys := cut(order)
	LeafNum = LeafNum / min_keys
	qu := queue.New(int64(LeafNum))
	var startIndexInLeaf, endIndexInLeaf int
	s := start
	e := end
	if c.IsLeaf {
		if s == 1 && e == order {
			qu.Put(c)
			startIndexInLeaf = 0
			endIndexInLeaf = order - 1
			return startIndexInLeaf, endIndexInLeaf, qu
		} else {
			qu.Put(c)
			return s - 1, e - 1, qu
		}
	}
	for i := 0; i < order; i++ {
		if c.Keys[i] < s {
			s -= c.Keys[i]
			e -= c.Keys[i]
			continue
		} else if s > 0 && c.Keys[i] >= s && e <= c.Keys[i] {
			return t.rangeLeafEnqueue(s, e, c.Pointers[i].(*Node))
		} else if s > 0 && c.Keys[i] >= s && e > c.Keys[i] {
			startIndexInLeaf = t.rangeStartLeaf(qu, c.Pointers[i].(*Node), s)
			s = 0
			e -= c.Keys[i]
		} else if s == 0 && c.Keys[i] < e {
			t.rangeWholeLeaf(qu, c.Pointers[i].(*Node))
			e -= c.Keys[i]
		} else if s == 0 && c.Keys[i] >= e {
			endIndexInLeaf = t.rangeEndLeaf(qu, c.Pointers[i].(*Node), e)
			break
		}
	}
	return startIndexInLeaf, endIndexInLeaf, qu
}

func (t *Tree) Reorder(start int, end int, new []int, c *Node) {
	//var startIndexInLeaf, endIndexInLeaf int
	myqueue = nil
	if c.IsLeaf {
		i := 0
		for j := start - 1; j < end; j++ {
			c.Keys[j] = new[i]
			i++
		}
		return
	}
	startIndexInLeaf, endIndexInLeaf, qu := t.rangeLeafEnqueue(start, end, c)
	startflag := false
	var lastnode interface{}
	newIndex := 0
	for qu.Len() != 0 {
		n, _ := qu.Get(1)
		ni := n[0]
		if n != nil {
			if lastnode != nil {
				if startflag == false {
					for i := startIndexInLeaf; i < order; i++ {
						if lastnode.(*Node).Keys[i] != 0 {
							lastnode.(*Node).Keys[i] = new[newIndex]
							newIndex++
						} else {
							break
						}
					}
					startIndexInLeaf = 0
					startflag = true
				} else {
					for i := 0; i < order; i++ {
						if lastnode.(*Node).Keys[i] != 0 {
							lastnode.(*Node).Keys[i] = new[newIndex]
							newIndex++
						} else {
							break
						}
					}
				}
			}
			lastnode = ni
		}
	}
	for i := startIndexInLeaf; i <= endIndexInLeaf; i++ {
		lastnode.(*Node).Keys[i] = new[newIndex]
		newIndex++
	}
}

func (t *Tree) rangeEndLeaf(qu *queue.Queue, n *Node, end int) int {
	if n.IsLeaf {
		qu.Put(n)
		return end - 1
	}
	var ans int
	e := end
	for i := 0; i < order; i++ {
		if n.Keys[i] < e {
			t.rangeWholeLeaf(qu, n.Pointers[i].(*Node))
			e -= n.Keys[i]
		} else {
			ans = t.rangeEndLeaf(qu, n.Pointers[i].(*Node), e)
			break
		}
	}
	return ans
}

func (t *Tree) rangeStartLeaf(qu *queue.Queue, n *Node, start int) int {
	if n.IsLeaf {
		qu.Put(n)
		return start - 1
	}
	var ans int
	s := start
	for i := 0; i < order; i++ {
		if n.Keys[i] == 0 {
			break
		}
		if n.Keys[i] < s {
			s -= n.Keys[i]
		} else if s != 0 && n.Keys[i] >= s {
			ans = t.rangeStartLeaf(qu, n.Pointers[i].(*Node), s)
			s = 0
		} else if s == 0 {
			t.rangeWholeLeaf(qu, n.Pointers[i].(*Node))
		}
	}
	return ans
}

func (t *Tree) rangeWholeLeaf(qu *queue.Queue, n *Node) {
	if n.IsLeaf {
		qu.Put(n)
	} else {
		for i := 0; i < order; i++ {
			if n.Keys[i] != 0 {
				t.rangeWholeLeaf(qu, n.Pointers[i].(*Node))
			} else {
				break
			}
		}
	}
}

func (t *Tree) Swap(start1 int, end1 int, start2 int, end2 int) {
	//!fmt.Printf("Here is Swap(%d, %d, %d, %d).\n", start1, end1, start2, end2)
	startIndexLeaf1, endIndexLeaf1, q1 := t.rangeLeafEnqueue(start1, end1, t.Root) // including non-complete leaves
	startIndexLeaf2, endIndexLeaf2, q2 := t.rangeLeafEnqueue(start2, end2, t.Root)
	numLeaf1 := q1.Len()
	numLeaf2 := q2.Len()
	inter1, _ := q1.Get(numLeaf1)
	inter2, _ := q2.Get(numLeaf2)
	longinter1 := make([]interface{}, numLeaf1+10)
	for i := 0; i < int(numLeaf1); i++ {
		longinter1[i] = inter1[i]
	}
	longinter2 := make([]interface{}, numLeaf2+10)
	for i := 0; i < int(numLeaf2); i++ {
		longinter2[i] = inter2[i]
	}
	t.swapLeaf(longinter1, longinter2, int(numLeaf1), int(numLeaf2), startIndexLeaf1, endIndexLeaf1, startIndexLeaf2, endIndexLeaf2)
}

//new swapLeaf
func (t *Tree) swapLeaf(inter1 []interface{}, inter2 []interface{}, numLeaf1 int, numLeaf2 int, startIndexLeaf1 int,
	endIndexLeaf1 int, startIndexLeaf2 int, endIndexLeaf2 int) {
	//!fmt.Printf("Here is swapLeaf, startindex:\n")
	//!fmt.Println(startIndexLeaf1, endIndexLeaf1, startIndexLeaf2, endIndexLeaf2)
	if numLeaf1 == 1 {
		t.movePointerAfter(inter1[0].(*Node), startIndexLeaf1, endIndexLeaf1, inter2[numLeaf2-1].(*Node), endIndexLeaf2)
		return
	} else if numLeaf2 == 1 {
		t.movePointerBefore(inter2[0].(*Node), startIndexLeaf2, endIndexLeaf2, inter1[0].(*Node), startIndexLeaf1)
		return
	}
	HeadAndTail := false // the tail of 1 is the same as the head of 2
	// in this case, this node follows inter2, and deal with it specially later
	var htOriEndIndexLeaf1, htOriStartIndexLeaf2 int
	var htNode *Node
	if reflect.DeepEqual(inter1[numLeaf1-1], inter2[0]) {
		htNode = inter2[0].(*Node)
		htOriEndIndexLeaf1 = endIndexLeaf1
		htOriStartIndexLeaf2 = startIndexLeaf2
		numLeaf1--
		//numLeaf2 --
		/*
			for i := 0; i < numLeaf2; i ++ {
				inter2[i] = inter2[i+1]
			}

		*/
		HeadAndTail = true
		endIndexLeaf1 = t.NodeNumKeys(inter1[numLeaf1-1].(*Node))
		//startIndexLeaf2 = 0

		/*
			//--------------------------------------new method for head tail----------------
			numKeyInEndLeftLeaf := endIndexLeaf1
			numKeysToBeMoved := htOriEndIndexLeaf1 + 1
			t.movePointerAfter(htNode, 0, htOriEndIndexLeaf1, inter1[numLeaf1-1].(*Node), endIndexLeaf1)
			if numKeyInEndLeftLeaf + numKeysToBeMoved > order {
				parent := inter1[numLeaf1-1].(*Node).Parent
				index := t.getIndexInParent(inter1[numLeaf1-1].(*Node))
				newNode := parent.Pointers[index+1]
				inter1[numLeaf1] = newNode
				numLeaf1++
				endIndexLeaf1 = t.NodeNumKeys(newNode.(*Node))-1
			} else {
				endIndexLeaf1 = t.NodeNumKeys(inter1[numLeaf1-1].(*Node))-1
			}

			fmt.Printf("Here is within headtail, tree:\n")
			t.PrintTree()
			//--------------------------------------new method for head tail----------------
		*/
	}
	/*
		fmt.Printf("swap Leaf 1:\n")
		for i := 0; i < numLeaf1; i ++ {
			t.printLeaf(inter1[i].(*Node))
		}
		fmt.Printf("\nswap Leaf 2:\n")
		for i := 0; i < numLeaf2; i ++ {
			t.printLeaf(inter2[i].(*Node))
		}
		fmt.Printf("\n")
	*/

	node1 := make([]*Node, numLeaf1)
	node2 := make([]*Node, numLeaf2)
	for i := 0; i < numLeaf1; i++ {
		node1[i] = inter1[i].(*Node)
	}
	for i := 0; i < numLeaf2; i++ {
		node2[i] = inter2[i].(*Node)
	}
	t.swapWholeLeafAndGoUpwards(node1, node2, numLeaf1, numLeaf2)
	/*
		fmt.Printf("Here is after swapWholeLeafAndGoUpwards, tree:\n")
		t.PrintTree()

	*/

	/*
		if HeadAndTail {
			s := t.NodeNumKeys(node1[numLeaf1-1])
			t.movePointerAfter(htNode, 0, htOriEndIndexLeaf1, node1[numLeaf1-1], s)
		}

	*/

	//move pointers
	if startIndexLeaf1 > 0 && startIndexLeaf2 > 0 {
		t.swapEndPointer(node1[0], startIndexLeaf1-1, node2[0], startIndexLeaf2-1)
	} else if startIndexLeaf1 > 0 && startIndexLeaf2 == 0 {
		new_n := t.movePointerBefore(node1[0], 0, startIndexLeaf1-1, node2[0], 0)
		if numLeaf1 == 2 && new_n != nil {
			startrowID := node1[1].Keys[endIndexLeaf1]
			flag1 := false
			for i := 0; i < order; i++ {
				if startrowID == new_n.Keys[i] {
					endIndexLeaf1 = i
					flag1 = true
					break
				}
			}
			if flag1 {
				node1[1] = new_n
				fmt.Printf("new_n:\n")
				t.printLeaf(new_n)
				fmt.Printf("starrowID = %d  endIndexLeaf1 = %d ", startrowID, endIndexLeaf1)
				fmt.Printf("\n")
			}

		}
	} else if startIndexLeaf1 == 0 && startIndexLeaf2 > 0 {
		new_n := t.movePointerBefore(node2[0], 0, startIndexLeaf2-1, node1[0], 0)
		if numLeaf2 == 2 && new_n != nil {
			startrowID := node2[1].Keys[endIndexLeaf2]
			for i := 0; i < order; i++ {
				if startrowID == new_n.Keys[i] {
					endIndexLeaf2 = i
					break
				}
			}
			node2[1] = new_n
			//	fmt.Printf("new_n:\n")
			//	t.printLeaf(new_n)
			//	fmt.Printf("   starrowID = %d  endIndexLeaf2 = %d ", startrowID, endIndexLeaf2)
			//	fmt.Printf("\n")
		}
	}
	/*
		fmt.Printf("Here is after move pointers in starting leaf, tree \n")
		t.PrintTree()

	*/

	numKeyEndLeaf1 := t.NodeNumKeys(node1[numLeaf1-1])
	numKeyEndLeaf2 := t.NodeNumKeys(node2[numLeaf2-1])
	if endIndexLeaf1 < numKeyEndLeaf1-1 && endIndexLeaf2 < numKeyEndLeaf2-1 {
		t.swapStartPointer(node1[numLeaf1-1], endIndexLeaf1+1, node2[numLeaf2-1], endIndexLeaf2+1)
	} else if endIndexLeaf1 == numKeyEndLeaf1-1 && endIndexLeaf2 < numKeyEndLeaf2-1 {
		t.movePointerAfter(node2[numLeaf2-1], endIndexLeaf2+1, numKeyEndLeaf2-1, node1[numLeaf1-1], endIndexLeaf1+1)
	} else if endIndexLeaf1 < numKeyEndLeaf1-1 && endIndexLeaf2 == numKeyEndLeaf2-1 {
		t.movePointerAfter(node1[numLeaf1-1], endIndexLeaf1+1, numKeyEndLeaf1-1, node2[numLeaf2-1], endIndexLeaf2+1)
	}
	/*
		fmt.Printf("Here is after move pointers in end leaf, tree \n")
		t.PrintTree()

	*/

	if HeadAndTail { // the tail of 1 is the same as the head of 2
		fmt.Printf("Here is ")
		numInHT := t.NodeNumKeys(htNode)
		num1 := t.NodeNumKeys(node1[numLeaf1-1])
		t.movePointerAfter(htNode, 0, htOriEndIndexLeaf1, node1[numLeaf1-1], num1)
		t.movePointerAfter(htNode, htOriStartIndexLeaf2, numInHT, node2[0], 0)

	}

}

func (t *Tree) OldswapLeaf(inter1 []interface{}, inter2 []interface{}, numLeaf1 int, numLeaf2 int, startIndexLeaf1 int,
	endIndexLeaf1 int, startIndexLeaf2 int, endIndexLeaf2 int) {
	fmt.Printf("Here is swapLeaf, startindex:\n")
	fmt.Println(startIndexLeaf1, endIndexLeaf1, startIndexLeaf2, endIndexLeaf2)
	startLeaf1Flag, startLeaf2Flag := true, true
	wholeLeaf1 := make([]*Node, numLeaf1)
	wholeLeaf2 := make([]*Node, numLeaf2)
	var w1, w2 int
	var startIndexEndLeaf1, startIndexEndLeaf2 int
	for i := 0; i < numLeaf1; i++ {
		if i == 0 && startIndexLeaf1 == 0 && numLeaf1 != 1 {
			wholeLeaf1[w1], _ = inter1[i].(*Node)
			startLeaf1Flag = false
			w1++
		} else if i > 0 && i < numLeaf1-1 {
			wholeLeaf1[w1], _ = inter1[i].(*Node)
			w1++
		} else if i == numLeaf1-1 {
			if numLeaf1 == 1 {
				startIndexEndLeaf1 = startIndexLeaf1
			} else {
				startIndexEndLeaf1 = 0
			}

			// endLeaf, _ := inter1[i].(*Node)
			/*
				numKeysInEndLeaf := t.NodeNumKeys(endLeaf)
				if (endIndexLeaf1 == numKeysInEndLeaf - 1  && numLeaf1 != 1) ||
					(endIndexLeaf1 == numKeysInEndLeaf - 1  && numLeaf1 == 1 && startIndexLeaf1 == 0){
					wholeLeaf1[w1] = endLeaf
					w1 ++
					endLeaf1Flag = false
				}
			*/
		}
	}
	for i := 0; i < numLeaf2; i++ {
		if i == 0 && startIndexLeaf2 == 0 && numLeaf2 != 1 {
			wholeLeaf2[w2], _ = inter2[i].(*Node)
			startLeaf2Flag = false
			w2++
		} else if i > 0 && i < numLeaf2-1 {
			wholeLeaf2[w2], _ = inter2[i].(*Node)
			w2++
		} else if i == numLeaf2-1 {
			if numLeaf2 == 1 {
				startIndexEndLeaf2 = startIndexLeaf2
			} else {
				startIndexEndLeaf2 = 0
			}

			/*
				endLeaf, _ := inter2[i].(*Node)
				numKeysInEndLeaf := t.NodeNumKeys(endLeaf)
				if (endIndexLeaf2 == numKeysInEndLeaf - 1  && numLeaf2 != 1) ||
					(endIndexLeaf2 == numKeysInEndLeaf - 1  && numLeaf2 == 1 && startIndexLeaf2 == 0) {
					wholeLeaf2[w2] = endLeaf
					w2 ++
					endLeaf2Flag = false
				}
			*/
		}
	}
	//	t.swapWholeLeafAndGoUpwards(wholeLeaf1, wholeLeaf2, w1, w2)

	fmt.Printf("Whole Leaf 1:\n")
	for i := 0; i < w1; i++ {
		t.printLeaf(wholeLeaf1[i])
	}
	fmt.Printf("\nWhole Leaf 2:\n")
	for i := 0; i < w2; i++ {
		t.printLeaf(wholeLeaf1[i])
	}

	if startLeaf1Flag == true && startLeaf2Flag == true {
		t.swapStartPointer(wholeLeaf1[0], startIndexLeaf1, wholeLeaf2[0], startIndexLeaf2)
	} else if startLeaf1Flag == false && startLeaf2Flag == true {
		tmp := t.NodeNumKeys(inter2[0].(*Node)) - 1
		t.movePointerBefore(inter2[0].(*Node), startIndexLeaf2, tmp, wholeLeaf2[0], 0)
	} else if startLeaf1Flag == true && startLeaf2Flag == false {
		tmp := t.NodeNumKeys(inter1[0].(*Node)) - 1
		t.movePointerBefore(inter1[0].(*Node), startIndexLeaf1, tmp, wholeLeaf1[0], 0)
	}

	i1 := endIndexLeaf1
	i2 := endIndexLeaf2
	endLeaf1 := inter1[numLeaf1-1].(*Node)
	endLeaf2 := inter2[numLeaf2-1].(*Node)
	for i1 >= startIndexEndLeaf1 && i2 >= startIndexEndLeaf2 {
		tmpk := endLeaf1.Keys[i1]
		endLeaf1.Keys[i1] = endLeaf2.Keys[i2]
		endLeaf2.Keys[i2] = tmpk
		i1--
		i2--
	}
	if i1 == startIndexEndLeaf1-1 && i2 == startIndexEndLeaf2-1 {
		return
	} else if i1 != startIndexEndLeaf1-1 {
		t.movePointerBefore(endLeaf1, 0, i1, endLeaf2, 0)
	} else if i2 != startIndexEndLeaf2-1 {
		t.movePointerBefore(endLeaf2, 0, i2, endLeaf1, 0)
	}
}

//done9.4
func (t *Tree) insertMulKeysIntoLeaf(keysToInsert []int, numKeysToInsert int, toLeaf *Node, toIndex int) {
	newKeys := make([]int, order)
	i := 0
	for ; i < numKeysToInsert; i++ {
		newKeys[i] = keysToInsert[i]
	}
	for j := 0; j < order; j++ {
		if toLeaf.Keys[j] != 0 {
			newKeys[i] = toLeaf.Keys[j]
			i++
		} else {
			break
		}
	}
	m := 0
	for k := toIndex; k < order; k++ {
		toLeaf.Keys[k] = newKeys[m]
		m++
		if m == i {
			break
		}
	}

	if toLeaf == t.Root {
		return
	}
	c := toLeaf
	parent := c.Parent
	for parent != nil {
		for i := 0; i < order; i++ {
			if reflect.DeepEqual(parent.Pointers[i], c) {
				parent.Keys[i] += numKeysToInsert
				c = parent
				parent = c.Parent
				break
			}
		}
	}

}

//done 9.4
func (t *Tree) insertMulKeysIntoLeafAfterSplitting(keysToInsert []int, numKeysToInsert int, toLeaf *Node, toIndex int) {
	newKeys := make([]int, order*2)
	i := 0
	for j := 0; j < toIndex; j++ {
		newKeys[i] = toLeaf.Keys[j]
		i++
	}
	for j := 0; j < numKeysToInsert; j++ {
		newKeys[i] = keysToInsert[j]
		i++
	}
	for j := toIndex; j < order; j++ {
		if toLeaf.Keys[j] != 0 {
			newKeys[i] = toLeaf.Keys[j]
			toLeaf.Keys[j] = 0
			i++
		} else {
			break
		}
	}

	leftNumKeys := cut(order)
	rightNumKeys := i - leftNumKeys
	for rightNumKeys > order {
		leftNumKeys++
		rightNumKeys--
	}

	new_leaf, _ := makeLeaf()
	y := 0
	for r := 0; r < leftNumKeys; r++ {
		toLeaf.Keys[r] = newKeys[y]
		y++
	}
	for r := 0; r < rightNumKeys; r++ {
		new_leaf.Keys[r] = newKeys[y]
		y++
		if y == i {
			break
		}
	}

	t.insertIntoParent(toLeaf, leftNumKeys, new_leaf, rightNumKeys)
}

//done 9.4
func (t *Tree) movePointerBefore(fromLeaf *Node, startIndex int, endIndex int, toLeaf *Node, toIndex int) *Node {
	keysToInsert := make([]int, order)
	for i := startIndex; i <= endIndex; i++ {
		keysToInsert[i-startIndex] = fromLeaf.Keys[i]
	}
	numKeysToInsert := endIndex - startIndex + 1
	//insert nobug!!!
	if numKeysToInsert+t.NodeNumKeys(toLeaf) <= order {
		t.insertMulKeysIntoLeaf(keysToInsert, numKeysToInsert, toLeaf, toIndex)
	} else {
		t.insertMulKeysIntoLeafAfterSplitting(keysToInsert, numKeysToInsert, toLeaf, toIndex)
	}

	//t.PrintTree()

	//delete
	min_key := cut(order)
	for i := startIndex; i <= endIndex; i++ {
		fromLeaf.Keys[i] = 0
	}
	// move keys in fromLeaf frontwards
	for i := endIndex + 1; i < order; i++ {
		if fromLeaf.Keys[i] != 0 {
			fromLeaf.Keys[i-endIndex-1] = fromLeaf.Keys[i]
			fromLeaf.Pointers[i-endIndex-1] = fromLeaf.Pointers[i]
			fromLeaf.Keys[i] = 0
			fromLeaf.Pointers[i] = nil
		}
	}
	curNumKeys := startIndex
	if fromLeaf == t.Root {
		t.adjustRoot()
		return nil
	}
	if t.NodeNumKeys(fromLeaf) >= min_key {
		if fromLeaf == t.Root {
			return nil
		}
		c := fromLeaf
		parent := c.Parent
		for parent != nil {
			for i := 0; i < order; i++ {
				if reflect.DeepEqual(parent.Pointers[i], c) {
					parent.Keys[i] -= numKeysToInsert
					c = parent
					parent = c.Parent
					break
				}
			}
		}
		return nil
	} else {
		neighbour_index := getLeftNeighborIndex(fromLeaf)
		var neighbour *Node
		if neighbour_index == -1 {
			neighbour, _ = fromLeaf.Parent.Pointers[1].(*Node)
		} else {
			neighbour, _ = fromLeaf.Parent.Pointers[neighbour_index].(*Node)
		}
		if t.NodeNumKeys(neighbour)+curNumKeys <= order {
			return t.coalesceNodes(fromLeaf, neighbour, neighbour_index, true)
		} else {
			t.redistributeNodes(fromLeaf, neighbour, neighbour_index)
			return nil
		}
	}

}

func (t *Tree) adjustAncestorKeys(n *Node) {
	c := n
	parent := c.Parent
	//index := t.getIndexInParent(c)
	//parent.Keys[index] = t.NodeNumKeys(c)
	//c = parent
	//parent = c.Parent
	for parent != nil {
		for i := 0; i < order; i++ {
			if reflect.DeepEqual(parent.Pointers[i], c) {
				parent.Keys[i] = t.NodeSumKeys(c)
				c = parent
				parent = c.Parent
				break
			}
		}
	}
}

func (t *Tree) adjustAncestorKeysAfterExchangingLeaves(nodeArray []*Node, numNodeInArray int) {
	newArray := make([]*Node, numNodeInArray)
	numNewArray := 0
	var lastParent, parent *Node
	for i := 0; i < numNodeInArray; i++ {
		parent = nodeArray[i].Parent
		if parent == nil { // root
			return
		}
		if !reflect.DeepEqual(lastParent, parent) {
			newArray[numNewArray] = parent
			lastParent = parent
			numNewArray++
		}
		index := t.getIndexInParent(nodeArray[i])
		parent.Keys[index] = t.NodeSumKeys(nodeArray[i])
	}
	t.adjustAncestorKeysAfterExchangingLeaves(newArray, numNewArray)
}

func (t *Tree) adjustOneEdgeNode(node *Node) {
	t.PrintTree()
	if node == t.Root {
		t.checkRootAfterMulDelete()
		return
	}
	if node.Parent == nil {
		return
	}
	min_key := cut(order)
	numKey := t.NodeNumKeys(node)
	parent := node.Parent
	if numKey >= min_key {
		t.adjustOneEdgeNode(parent)
		return
	}
	if numKey == 0 {
		t.adjustOneEdgeNode(parent)
		return
	}
	var neighbour_index int
	var neighbour *Node
	neighbour_index = getLeftNeighborIndex(node)
	if neighbour_index == -1 {
		neighbour_index = 1
		if parent.Keys[1] == 0 {
			t.adjustOneEdgeNode(parent)
			return
		}
	}
	neighbour, _ = node.Parent.Pointers[neighbour_index].(*Node)
	if t.NodeNumKeys(neighbour)+numKey <= order {
		t.coalesceNodes(node, neighbour, neighbour_index, true) //still too empty?????
		return
	} else {
		t.redistributeNodes(node, neighbour, neighbour_index)
		return
	}

}

//swap whole leaves and go upwards until the root
func (t *Tree) swapWholeLeafAndGoUpwards(wholeLeaf1 []*Node, wholeLeaf2 []*Node, w1 int, w2 int) /*(int, error, *Node, int, []int, []interface{})*/ {
	//!fmt.Printf("Here is swapWholeLeafAndGoUpwards\n")
	/*
		fmt.Printf("Left whole leaves:\n")
		for i := 0; i < w1; i ++ {
			t.printLeaf(wholeLeaf1[i])
		}
		fmt.Printf("\n")
		fmt.Printf("Right whole leaves:\n")
		for i := 0; i < w2; i ++ {
			t.printLeaf(wholeLeaf2[i])
		}
		fmt.Printf("\nBefore swapping, Tree:\n")
		t.PrintTree()

	*/
	var remain, exchangeNum int
	if w1 == w2 {
		exchangeNum = w1
		remain = 0
	} else if w1 < w2 {
		exchangeNum = w1
		remain = 1
	} else {
		exchangeNum = w2
		remain = -1 // left leaves more than right leaves
	}

	//exchange corresponding leaves
	var LeftEdgeLeaf, RightEdgeLeaf *Node
	LeftParentArray := make([]*Node, exchangeNum)
	RightParentArray := make([]*Node, exchangeNum)
	lp, rp := 0, 0
	var lastLeftParent *Node
	var lastRightParent *Node
	for i := 0; i < exchangeNum; i++ {
		LeftLeaf := wholeLeaf1[i]
		RightLeaf := wholeLeaf2[i]
		LeftParent := LeftLeaf.Parent
		RightParent := RightLeaf.Parent
		LeftEdgeLeaf, RightEdgeLeaf = t.exchangeLeaf(LeftLeaf, LeftParent, RightLeaf, RightParent)
		if !reflect.DeepEqual(LeftParent, lastLeftParent) {
			LeftParentArray[lp] = LeftParent
			lastLeftParent = LeftParent
			lp++
		}
		if !reflect.DeepEqual(RightParent, lastRightParent) {
			RightParentArray[rp] = RightParent
			lastRightParent = RightParent
			rp++
		}
	}

	// adjust keys in ancestors
	t.adjustAncestorKeysAfterExchangingLeaves(LeftParentArray, lp)
	t.adjustAncestorKeysAfterExchangingLeaves(RightParentArray, rp)
	/*
		fmt.Printf("\nTree after exchanging corresponding leaves:\n")
		t.PrintTree()
	*/

	//delete leaves
	if remain != 0 {
		//delete leaves without merging/redistribution/adjusting ancestor's keys
		var leafToDelete []interface{}
		var numLeafToDelete int
		var leftLeafParent, rightLeafParent *Node
		//var indexOfLeftLeafParent, indexOfRightLeafParent int
		if remain == -1 { // left leaves more
			numLeafToDelete = w1 - w2
			leafToDelete = make([]interface{}, numLeafToDelete)
			wl1 := exchangeNum
			for i := 0; i < numLeafToDelete; i++ {
				leafToDelete[i] = wholeLeaf1[wl1]
				wl1++
			}
			leftLeafParent = leafToDelete[0].(*Node).Parent
			//indexOfLeftLeafParent = t.getIndexInParent(leftLeafParent)
			rightLeafParent = leafToDelete[numLeafToDelete-1].(*Node).Parent
			//indexOfRightLeafParent = t.getIndexInParent(rightLeafParent)
		} else { // remain == 1, right leaves more
			numLeafToDelete = w2 - w1
			leafToDelete = make([]interface{}, numLeafToDelete)
			wl2 := exchangeNum
			for i := 0; i < numLeafToDelete; i++ {
				leafToDelete[i] = wholeLeaf2[wl2]
				wl2++
			}
			leftLeafParent = leafToDelete[0].(*Node).Parent
			//indexOfLeftLeafParent = t.getIndexInParent(leftLeafParent)
			rightLeafParent = leafToDelete[numLeafToDelete-1].(*Node).Parent
			//indexOfRightLeafParent = t.getIndexInParent(rightLeafParent)
		}
		leafParent := make([]interface{}, numLeafToDelete)
		lastp := leafToDelete[0].(*Node).Parent
		leafParent[0] = lastp
		lp := 1
		for i := 1; i < numLeafToDelete; i++ {
			curp := leafToDelete[i].(*Node).Parent
			if curp != lastp {
				leafParent[lp] = curp
				lp++
				lastp = curp
			}
		}
		t.deleteMulEntry(leafToDelete, numLeafToDelete)
		/*
			fmt.Printf("Here is after deleteSeparateLeaves, before adjust keys in ancestors\n")
			t.PrintTree()
		*/

		//adjust keys in leftLeafParent, rightLeafParent and ancestors
		adjust1 := leftLeafParent
		numKeys1 := t.NodeNumKeys(adjust1)
		if numKeys1 > 0 {
			t.adjustAncestorKeys(adjust1)
		} else {
			adjust1 = adjust1.Parent
			parent := adjust1.Parent
			for parent != nil {
				if t.NodeNumKeys(adjust1) > 0 {
					t.adjustAncestorKeys(adjust1)
					break
				} else {
					adjust1 = adjust1.Parent
					parent = adjust1.Parent
				}
			}
			//leftNeighborOfLeftParent := leftLeafParent.Parent.Pointers[lef]
		}
		adjust2 := rightLeafParent
		numKeys2 := t.NodeNumKeys(adjust2)
		if numKeys2 > 0 {
			t.adjustAncestorKeys(adjust2)
		} else {
			adjust2 = adjust2.Parent
			parent := adjust2.Parent
			for parent != nil {
				if t.NodeNumKeys(adjust2) > 0 {
					t.adjustAncestorKeys(adjust2)
					break
				} else {
					adjust2 = adjust2.Parent
					parent = adjust2.Parent
				}
			}
			//leftNeighborOfLeftParent := leftLeafParent.Parent.Pointers[lef]
		}
		/*
			fmt.Printf("Here is after adjust keys in ancestors\n")
			t.PrintTree()

		*/
		// merge/redistribute nodes
		t.adjustOneEdgeNode(leftLeafParent)  //!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		t.adjustOneEdgeNode(rightLeafParent) // bug
		/*
			fmt.Printf("Here is after adjust edge nodes\n")
			t.PrintTree()

			if LeftEdgeLeaf == nil && RightEdgeLeaf == nil {
				fmt.Printf("useless\n")
			}
		*/

		// remain != 0 , insert leaves

		// insert leaves
		if remain == -1 { // left leaves more
			numLeafToInsert := w1 - w2
			leafToInsert := make([]interface{}, numLeafToInsert)
			wl1 := exchangeNum
			for i := 0; i < numLeafToInsert; i++ {
				leafToInsert[i] = wholeLeaf1[wl1]
				wl1++
			}
			toParent := RightEdgeLeaf.Parent
			startIndexInParent := t.getIndexInParent(RightEdgeLeaf)
			if t.NodeNumKeys(toParent)+numLeafToInsert < order { // enough space
				t.insertMulLeavesIntoLeafParent(toParent, startIndexInParent, leafToInsert, numLeafToInsert)
			} else {
				t.insertMulLeavesIntoLeafParentAfterSplitting(toParent, startIndexInParent, leafToInsert, numLeafToInsert)
			}
		} else if remain == 1 { // right leaves more
			numLeafToInsert := w2 - w1
			leafToInsert := make([]interface{}, numLeafToInsert)
			wl2 := exchangeNum
			for i := 0; i < numLeafToInsert; i++ {
				leafToInsert[i] = wholeLeaf2[wl2]
				wl2++
			}
			toParent := LeftEdgeLeaf.Parent
			startIndexInParent := t.getIndexInParent(LeftEdgeLeaf) + 1
			if t.NodeNumKeys(toParent)+numLeafToInsert < order { // enough space
				t.insertMulLeavesIntoLeafParent(toParent, startIndexInParent, leafToInsert, numLeafToInsert)
			} else {
				t.insertMulLeavesIntoLeafParentAfterSplitting(toParent, startIndexInParent, leafToInsert, numLeafToInsert)
			}
		}

	} // if remain != 0
}

func (t *Tree) setNewRootAfterMulDelete(n *Node) {
	var new_root *Node
	new_root = n
	new_root.Parent = nil
	t.Root = new_root
}

func (t *Tree) checkRootAfterMulDelete() {
	numKeyInRoot := t.NodeNumKeys(t.Root)
	if numKeyInRoot > 1 || t.Root.IsLeaf {
		return
	}
	var new_root *Node
	new_root, _ = t.Root.Pointers[0].(*Node)
	new_root.Parent = nil
	t.Root = new_root
}

func (t *Tree) adjustEdgeNodesReachRoot(n *Node) {
	num := t.NodeNumKeys(n)
	if num > 1 {
		return
	}
	var new_root *Node
	new_root = n.Pointers[0].(*Node)
	new_root.Parent = nil
	t.Root = new_root
	t.adjustEdgeNodesReachRoot(new_root)
}

/*
func (t *Tree) adjustLeftEdgeNodes(leftLeafParent *Node) {
	thisNode := leftLeafParent
	numKeysInThisNode := t.NodeNumKeys(thisNode)
	parent := thisNode.Parent
	numKeysInParent := t.NodeNumKeys(parent)
	min_key := cut(order)
	var neighbour_index int
	var neighbour *Node

	if parent == nil { // root
		t.adjustEdgeNodesReachRoot(thisNode)
	}
	for numKeysInThisNode == 0 || numKeysInThisNode >= min_key {
		thisNode = thisNode.Parent
		numKeysInThisNode = t.NodeNumKeys(thisNode)
		parent = thisNode.Parent
		numKeysInParent = t.NodeNumKeys(parent)
		if parent == nil {
			t.adjustEdgeNodesReachRoot(thisNode)
			return
		}
	}

	//for numKeysInThisNode


}

*/

func (t *Tree) adjustEdgeNodes(leftLeafParent *Node, rightLeafParent *Node) {
	//adjust left side with its left neighbor, then adjust right side with its left neighbor
	//if left side has no left neighbor
	//adjust right side with its right neighbor, then adjust left side with its right neighbor
	//if left side has no left neighbor and right side has no right neighbor:
	//special case
	//root ??
	startNode := leftLeafParent
	parent := startNode.Parent
	min_key := cut(order)
	var neighbour_index int
	var neighbour *Node
	rightFirst := false
	found1, found2, found3, found4 := false, false, false, false
	needAdjustLeft := true  //bool (t.NodeNumKeys(leftLeafParent) > 0)
	needAdjustRight := true //bool (t.NodeNumKeys(rightLeafParent) > 0)
	//adjust left side
	if needAdjustLeft {
		for parent != nil {
			numKeyInStartNode := t.NodeNumKeys(startNode)
			if numKeyInStartNode < min_key && numKeyInStartNode > 0 {
				found1 = true
				neighbour_index = getLeftNeighborIndex(startNode)
				if neighbour_index == -1 {
					rightFirst = true
				} else {
					neighbour, _ = startNode.Parent.Pointers[neighbour_index].(*Node)
					if t.NodeNumKeys(neighbour)+numKeyInStartNode <= order {
						t.coalesceNodes(startNode, neighbour, neighbour_index)
						return
					} else {
						t.redistributeNodes(startNode, neighbour, neighbour_index)
						return
					}
				}
				break
			} else {
				startNode = startNode.Parent
				parent = startNode.Parent
			}
		}
		if found1 == false { // check root!
			t.checkRootAfterMulDelete()
		}
	}

	startNode = rightLeafParent
	parent = startNode.Parent
	if rightFirst == false && needAdjustRight { //adjust right side with left neighbor
		for parent != nil {
			numKeyInStartNode := t.NodeNumKeys(startNode)
			if numKeyInStartNode < min_key && numKeyInStartNode > 0 {
				found2 = true
				neighbour_index = getLeftNeighborIndex(startNode)
				if neighbour_index == -1 {
					fmt.Printf("error! This case is impossible!\n")
				} else {
					neighbour, _ = startNode.Parent.Pointers[neighbour_index].(*Node)
					if t.NodeNumKeys(neighbour)+numKeyInStartNode <= order {
						t.coalesceNodes(startNode, neighbour, neighbour_index)
						return
					} else {
						t.redistributeNodes(startNode, neighbour, neighbour_index)
						return
					}
				}
				break
			} else {
				startNode = startNode.Parent
				parent = startNode.Parent
			}
		}
		if found2 == false { // check root!
			t.checkRootAfterMulDelete()
			//return
		}
		return
	} else { //adjust right side with right neighbor first
		special := false
		startNode = rightLeafParent
		parent = startNode.Parent
		if needAdjustRight { //adjust right side with right neighbor
			for parent != nil {
				numKeyInStartNode := t.NodeNumKeys(startNode)
				if numKeyInStartNode < min_key && numKeyInStartNode > 0 {
					found3 = true
					neighbour_index = t.getRightNeighborIndex(startNode)
					if neighbour_index == order {
						special = true
					} else {
						neighbour, _ = startNode.Parent.Pointers[neighbour_index].(*Node)
						if t.NodeNumKeys(neighbour)+numKeyInStartNode <= order {
							t.coalesceNodes(startNode, neighbour, neighbour_index)
							//return
						} else {
							t.redistributeNodes(startNode, neighbour, neighbour_index)
							return
						}
					}
					break
				} else {
					startNode = startNode.Parent
					parent = startNode.Parent
				}
			}
			if found3 == false { // check root!
				t.checkRootAfterMulDelete()
			}
		}
		startNode = leftLeafParent
		parent = startNode.Parent
		if needAdjustLeft && !special {
			for parent != nil {
				numKeyInStartNode := t.NodeNumKeys(startNode)
				if numKeyInStartNode < min_key && numKeyInStartNode > 0 {
					found4 = true
					neighbour_index = t.getRightNeighborIndex(startNode)
					if neighbour_index == order {
						fmt.Printf("error! This case is impossible!\n")
					} else {
						neighbour, _ = startNode.Parent.Pointers[neighbour_index].(*Node)
						if t.NodeNumKeys(neighbour)+numKeyInStartNode <= order {
							t.coalesceNodes(startNode, neighbour, neighbour_index)
							return
						} else {
							t.redistributeNodes(startNode, neighbour, neighbour_index)
							return
						}
					}
					break
				} else {
					startNode = startNode.Parent
					parent = startNode.Parent
				}
			}
			if found4 == false { // check root!
				t.checkRootAfterMulDelete()
				//return
			}
		}

		if special {
			fmt.Printf("Here is special\n")
			if reflect.DeepEqual(leftLeafParent, rightLeafParent) { // the same node
				node := leftLeafParent
				parent := node.Parent
				for parent != nil {
					if t.NodeNumKeys(node) < min_key {
						t.setNewRootAfterMulDelete(node)
						break
					} else {
						node = node.Parent
						parent = node.Parent
					}
				}
			} else { // onle left two nodes, both too empty
				//	numLeft := t.NodeNumKeys(leftLeafParent)
				//	numRight := t.NodeNumKeys(rightLeafParent)

			}
		}

	}

}

//move pointers in node forward, if node is empty, return true, otherwise, return false
func (t *Tree) movePointerInNodeForward(node *Node) bool {
	tmpP := make([]interface{}, order)
	tmpK := make([]int, order)
	i_tmp := 0
	for i := 0; i < order; i++ {
		if node.Pointers[i] != nil {
			tmpP[i_tmp] = node.Pointers[i]
			tmpK[i_tmp] = node.Keys[i]
			i_tmp++
			node.Pointers[i] = nil
			node.Keys[i] = 0
		}
	}
	for i := 0; i < i_tmp; i++ {
		node.Pointers[i] = tmpP[i]
		node.Keys[i] = tmpK[i]
	}
	if i_tmp == 0 {
		return true
	} else {
		return false
	}
}

func (t *Tree) deleteMulEntry(entryToDelete []interface{}, numEntryToDelete int) {
	//fmt.Printf("Here is deleteMulEntry, tree:\n")
	//t.PrintTree()
	//leftParent := entryToDelete[0].(*Node).Parent
	//rightParent := entryToDelete[numEntryToDelete-1].(*Node).Parent
	//fmt.Printf("entryToDelete:\n")
	min_keys := cut(order)
	newEntryToDelete := make([]interface{}, numEntryToDelete/min_keys+2)
	i_newEntryToDelete := 0
	lastParent := entryToDelete[0].(*Node).Parent
	var curParent *Node
	for i := 0; i < numEntryToDelete; i++ {
		curEntry := entryToDelete[i].(*Node)
		curParent = curEntry.Parent
		index := t.getIndexInParent(curEntry)
		curParent.Pointers[index] = nil
		curParent.Keys[index] = 0
		if !reflect.DeepEqual(curParent, lastParent) {
			num := t.NodeNumKeys(lastParent)
			if num == 0 {
				newEntryToDelete[i_newEntryToDelete] = lastParent
				i_newEntryToDelete++
			}
		}
		lastParent = curParent
	}
	empty := t.movePointerInNodeForward(curParent)
	if empty == true {
		newEntryToDelete[i_newEntryToDelete] = lastParent
		i_newEntryToDelete++
	}

	if i_newEntryToDelete > 0 {
		t.deleteMulEntry(newEntryToDelete, i_newEntryToDelete)
	}
}

//toParent's ancestor key modified here
func (t *Tree) insertMulLeavesIntoLeafParentAfterSplitting(toParent *Node, startIndexInParent int,
	leafToInsert []interface{}, numLeafToInsert int) {
	/*
		fmt.Printf("Here is insertMulLeavesIntoLeafParentAfterSplitting, tree: \n")
		t.PrintTree()

	*/
	allLeaf := make([]interface{}, order+numLeafToInsert)
	j := 0
	for i := 0; i < startIndexInParent; i++ {
		allLeaf[j] = toParent.Pointers[i]
		toParent.Keys[i] = 0
		toParent.Pointers[i] = nil
		j++
	}
	for i := 0; i < numLeafToInsert; i++ {
		allLeaf[j] = leafToInsert[i].(*Node)
		j++
	}
	for i := startIndexInParent; i < order; i++ {
		if toParent.Pointers[i] == nil {
			break
		}
		allLeaf[j] = toParent.Pointers[i]
		toParent.Keys[i] = 0
		toParent.Pointers[i] = nil
		j++
	}
	totalNumLeaf := j // number of keys in toParent + leaftoinsert

	/*
		fmt.Printf("Here is insertMulLeavesIntoLeafParentAfterSplitting, allLeaf: \n")
		for i := 0; i < totalNumLeaf; i ++ {
			t.printLeaf(allLeaf[i].(*Node))
		}
		fmt.Printf("\n")
	*/

	numFullParent := totalNumLeaf / order
	remainNumLeaf := totalNumLeaf - numFullParent*order
	min_keys := cut(order)
	if remainNumLeaf == 0 || remainNumLeaf >= min_keys { // full nodes + remainnumleaf
		//set toParent , the left most one
		i := 0
		for ; i < order; i++ {
			toParent.Pointers[i] = allLeaf[i]
			toParent.Keys[i] = t.NodeNumKeys(allLeaf[i].(*Node))
			allLeaf[i].(*Node).Parent = toParent
		}
		/*
			fmt.Printf("toParent:\n")
			for l := 0; l < order; l ++ {
				fmt.Printf("%d " , toParent.Keys[l])
			}
			fmt.Printf("\n")

		*/
		parentNewSiblingNum := numFullParent - 1
		if remainNumLeaf != 0 {
			parentNewSiblingNum++
		}
		fmt.Printf("toParent's sibling:\n")
		parentNewSibling := make([]*Node, parentNewSiblingNum)
		for k := 0; k < parentNewSiblingNum; k++ {
			parentNewSibling[k], _ = MakeNode()
			for y := 0; y < order; y++ {
				parentNewSibling[k].Pointers[y] = allLeaf[i]
				parentNewSibling[k].Keys[y] = t.NodeNumKeys(allLeaf[i].(*Node))
				allLeaf[i].(*Node).Parent = parentNewSibling[k]
				fmt.Printf("%d ", parentNewSibling[k].Keys[y])
				i++
				if i == totalNumLeaf {
					break
				}
			}
			fmt.Printf("|")
		}
		fmt.Printf("parentNewSibling subtree : \n")
		for o := 0; o < parentNewSiblingNum; o++ {
			t.PrintSubTree(parentNewSibling[0])
		}

		//t.PrintTree()
		// insert parentNewSibling after toParent
		t.insertMulIntoParent(toParent, parentNewSibling, parentNewSiblingNum)
	} else {
		parentNewSiblingNum := numFullParent
		numFullParent--
		numkeysInSecondLastParent := min_keys
		numkeysInLastParent := totalNumLeaf - numFullParent*order - min_keys
		i_allleaf := 0
		// set keys/pinters in toParent
		if numFullParent == 0 { //toParent is the second last one
			for i := 0; i < min_keys; i++ {
				toParent.Pointers[i] = allLeaf[i_allleaf]
				toParent.Keys[i] = t.NodeNumKeys(allLeaf[i_allleaf].(*Node))
				allLeaf[i_allleaf].(*Node).Parent = toParent
				i_allleaf++
			}
		} else {
			for i := 0; i < order; i++ {
				toParent.Pointers[i] = allLeaf[i_allleaf]
				toParent.Keys[i] = t.NodeNumKeys(allLeaf[i_allleaf].(*Node))
				allLeaf[i_allleaf].(*Node).Parent = toParent
				i_allleaf++
			}
		}
		// set toParent's new siblings
		parentNewSibling := make([]*Node, parentNewSiblingNum)
		for k := 0; k < parentNewSiblingNum; k++ {
			parentNewSibling[k], _ = MakeNode()
			if k == parentNewSiblingNum-1 { // the last one
				for y := 0; y < numkeysInLastParent; y++ {
					parentNewSibling[k].Pointers[y] = allLeaf[i_allleaf]
					parentNewSibling[k].Keys[y] = t.NodeNumKeys(allLeaf[i_allleaf].(*Node))
					allLeaf[i_allleaf].(*Node).Parent = parentNewSibling[k]
					i_allleaf++
				}
			} else if k == parentNewSiblingNum-2 { // the second last one
				for y := 0; y < numkeysInSecondLastParent; y++ {
					parentNewSibling[k].Pointers[y] = allLeaf[i_allleaf]
					parentNewSibling[k].Keys[y] = t.NodeNumKeys(allLeaf[i_allleaf].(*Node))
					allLeaf[i_allleaf].(*Node).Parent = parentNewSibling[k]
					i_allleaf++
				}
			} else {
				for y := 0; y < order; y++ {
					parentNewSibling[k].Pointers[y] = allLeaf[i_allleaf]
					parentNewSibling[k].Keys[y] = t.NodeNumKeys(allLeaf[i_allleaf].(*Node))
					allLeaf[i_allleaf].(*Node).Parent = parentNewSibling[k]
					i_allleaf++
				}
			}
		}
		t.insertMulIntoParent(toParent, parentNewSibling, parentNewSiblingNum)
	}
}

//done 1
func (t *Tree) insertMulIntoNewRoot(toParent *Node, parentNewSibling []*Node, parentNewSiblingNum int) {
	/*
		fmt.Printf("Here is insertMulIntoNewRoot, tree: \n")
		t.PrintTree()

	*/
	totalNum := parentNewSiblingNum + 1
	numFullNodes := totalNum / order
	if numFullNodes == 0 || (totalNum == order) { // new root!
		t.Root, _ = MakeNode()
		t.Root.Keys[0] = t.NodeSumKeys(toParent)
		t.Root.Pointers[0] = toParent
		toParent.Parent = t.Root
		for i := 0; i < parentNewSiblingNum; i++ {
			t.Root.Pointers[i+1] = parentNewSibling[i]
			t.Root.Keys[i+1] = t.NodeSumKeys(parentNewSibling[i])
			parentNewSibling[i].Parent = t.Root
		}
		return
	}
	if totalNum%order == 0 { //all full nodes
		leftNode, _ := MakeNode()
		rightSiblingNum := numFullNodes - 1
		rightSibling := make([]*Node, rightSiblingNum)
		i_rightchild := 0
		// set keys/pointers in leftNode
		leftNode.Pointers[0] = toParent
		leftNode.Keys[0] = t.NodeSumKeys(toParent)
		toParent.Parent = leftNode
		for i := 1; i < order; i++ {
			leftNode.Pointers[i] = parentNewSibling[i_rightchild]
			leftNode.Keys[i] = t.NodeSumKeys(parentNewSibling[i_rightchild])
			parentNewSibling[i_rightchild].Parent = leftNode
			i_rightchild++
		}
		for i := 0; i < rightSiblingNum; i++ {
			rightSibling[i], _ = MakeNode()
			for j := 0; j < order; j++ {
				rightSibling[i].Pointers[j] = parentNewSibling[i_rightchild]
				rightSibling[i].Keys[j] = t.NodeSumKeys(parentNewSibling[i_rightchild])
				parentNewSibling[i_rightchild].Parent = rightSibling[i]
				i_rightchild++
			}
		}
		t.insertMulIntoNewRoot(leftNode, rightSibling, rightSiblingNum)
	} else { // not all full nodes
		leftNode, _ := MakeNode()
		leftNode.Pointers[0] = toParent
		leftNode.Keys[0] = t.NodeSumKeys(toParent)
		toParent.Parent = leftNode

		rightSiblingNum := numFullNodes
		numFullNodes--
		rightSibling := make([]*Node, rightSiblingNum)
		i_rightchild := 0
		remainer := totalNum / order
		min_keys := cut(order)
		var numkeysInSecondLast, numkeysInLast int
		if remainer >= min_keys {
			numkeysInSecondLast = order
			numkeysInLast = remainer
		} else {
			numkeysInSecondLast = min_keys
			numkeysInLast = totalNum - numFullNodes*order - min_keys
		}
		if numFullNodes == 0 { //leftNode is the second last one
			for i := 1; i < min_keys; i++ { // leftNode
				leftNode.Pointers[i] = parentNewSibling[i_rightchild]
				leftNode.Keys[i] = t.NodeNumKeys(parentNewSibling[i_rightchild])
				parentNewSibling[i_rightchild].Parent = leftNode
				i_rightchild++
			}
			rightSibling[0], _ = MakeNode()
			for i := 0; i < numkeysInLast; i++ {
				rightSibling[0].Pointers[i] = parentNewSibling[i_rightchild]
				rightSibling[0].Keys[i] = t.NodeSumKeys(parentNewSibling[i_rightchild])
				parentNewSibling[i_rightchild].Parent = rightSibling[i]
				i_rightchild++
			}
		} else { // leftnode is not the second last one
			for i := 1; i < order; i++ { // leftNode
				leftNode.Pointers[i] = parentNewSibling[i_rightchild]
				leftNode.Keys[i] = t.NodeNumKeys(parentNewSibling[i_rightchild])
				parentNewSibling[i_rightchild].Parent = leftNode
				i_rightchild++
			}
			for k := 0; k < rightSiblingNum; k++ {
				rightSibling[k], _ = MakeNode()
				if k == rightSiblingNum-1 { // the last one
					for y := 0; y < numkeysInLast; y++ {
						rightSibling[k].Pointers[y] = parentNewSibling[i_rightchild]
						rightSibling[k].Keys[y] = t.NodeNumKeys(parentNewSibling[i_rightchild])
						parentNewSibling[i_rightchild].Parent = rightSibling[k]
						i_rightchild++
					}
				} else if k == rightSiblingNum-2 { // the second last one
					for y := 0; y < numkeysInSecondLast; y++ {
						rightSibling[k].Pointers[y] = parentNewSibling[i_rightchild]
						rightSibling[k].Keys[y] = t.NodeNumKeys(parentNewSibling[i_rightchild])
						parentNewSibling[i_rightchild].Parent = rightSibling[k]
						i_rightchild++
					}
				} else {
					for y := 0; y < order; y++ {
						rightSibling[k].Pointers[y] = parentNewSibling[i_rightchild]
						rightSibling[k].Keys[y] = t.NodeNumKeys(parentNewSibling[i_rightchild])
						parentNewSibling[i_rightchild].Parent = rightSibling[k]
						i_rightchild++
					}
				}
			}
		}
		t.insertMulIntoNewRoot(leftNode, rightSibling, rightSiblingNum)
	}
}

//insert parentNewSibling after toParent
func (t *Tree) insertMulIntoParent(toParent *Node, parentNewSibling []*Node, parentNewSiblingNum int) {
	//fmt.Printf("Here is insertMulIntoParent \n")
	//t.PrintTree()//!!!!!!!!!!!!!!!!!!!!
	parent := toParent.Parent
	if parent == nil {
		t.insertMulIntoNewRoot(toParent, parentNewSibling, parentNewSiblingNum)
	}
	toParentIndex := t.getIndexInParent(toParent)
	parent.Keys[toParentIndex] = t.NodeSumKeys(toParent)
	//fmt.Printf("toParent: \n")
	//t.PrintSubTree(toParent)
	if t.NodeNumKeys(parent)+parentNewSiblingNum <= order {
		t.insertMulIntoNode(parent, toParentIndex+1, parentNewSibling, parentNewSiblingNum)
		//fmt.Printf("tree:\n")
	} else {
		t.insertMulIntoNodeAfterSplitting(parent, toParentIndex+1, parentNewSibling, parentNewSiblingNum)
	}
}

// insert parentNewSibling into parent, start index is toIndex
// need to split parent first
func (t *Tree) insertMulIntoNodeAfterSplitting(parent *Node, toIndex int, pointerToInsert []*Node, pointerToInsertNum int) {
	allPointers := make([]*Node, order+pointerToInsertNum)
	/*
		fmt.Printf("Here is insertMulIntoNodeAfterSplitting, tree: \n")
		t.PrintTree()

	*/
	j := 0
	for i := 0; i < toIndex; i++ {
		allPointers[j] = parent.Pointers[i].(*Node)
		j++
	}
	for i := 0; i < pointerToInsertNum; i++ {
		allPointers[j] = pointerToInsert[i]
		j++
	}
	for i := toIndex; i < order; i++ {
		if parent.Pointers[i] == nil {
			break
		}
		allPointers[j] = parent.Pointers[i].(*Node)
		j++
	}
	totalNumPointer := j
	numFullParent := totalNumPointer / order
	remainNumLeaf := totalNumPointer - numFullParent*order
	min_keys := cut(order)
	if remainNumLeaf == 0 || remainNumLeaf >= min_keys { // full nodes + remainnumleaf
		i := 0
		for ; i < order; i++ {
			parent.Pointers[i] = allPointers[i]
			parent.Keys[i] = t.NodeNumKeys(allPointers[i])
			allPointers[i].Parent = parent
		}
		parentNewSiblingNum := numFullParent - 1
		if remainNumLeaf != 0 {
			parentNewSiblingNum++
		}
		parentNewSibling := make([]*Node, parentNewSiblingNum)
		for k := 0; k < parentNewSiblingNum; k++ {
			parentNewSibling[k], _ = MakeNode()
			for y := 0; y < order; y++ {
				parentNewSibling[k].Pointers[y] = allPointers[i]
				parentNewSibling[k].Keys[y] = t.NodeNumKeys(allPointers[i])
				allPointers[i].Parent = parentNewSibling[k]
				i++
				if i == totalNumPointer {
					break
				}
			}
		}
		//t.PrintTree()
		t.insertMulIntoParent(parent.Parent, parentNewSibling, parentNewSiblingNum)
	} else {
		parentNewSiblingNum := numFullParent
		numFullParent--
		numkeysInSecondLastParent := min_keys
		numkeysInLastParent := totalNumPointer - numFullParent*order - min_keys
		i_allpointer := 0
		if numFullParent == 0 { //toParent is the second last one
			for i := 0; i < min_keys; i++ {
				parent.Pointers[i] = allPointers[i_allpointer]
				parent.Keys[i] = t.NodeNumKeys(allPointers[i_allpointer])
				allPointers[i_allpointer].Parent = parent
				i_allpointer++
			}
		} else {
			for i := 0; i < order; i++ {
				parent.Pointers[i] = allPointers[i_allpointer]
				parent.Keys[i] = t.NodeNumKeys(allPointers[i_allpointer])
				allPointers[i_allpointer].Parent = parent
				i_allpointer++
			}
		}

		//set parent's new siblings
		parentNewSibling := make([]*Node, parentNewSiblingNum)
		for k := 0; k < parentNewSiblingNum; k++ {
			parentNewSibling[k], _ = MakeNode()
			if k == parentNewSiblingNum-1 { // the last one
				for y := 0; y < numkeysInLastParent; y++ {
					parent.Pointers[y] = allPointers[i_allpointer]
					parent.Keys[y] = t.NodeNumKeys(allPointers[i_allpointer])
					allPointers[i_allpointer].Parent = parent
					i_allpointer++
				}
			} else if k == parentNewSiblingNum-2 { // the second last one
				for y := 0; y < numkeysInSecondLastParent; y++ {
					parent.Pointers[y] = allPointers[i_allpointer]
					parent.Keys[y] = t.NodeNumKeys(allPointers[i_allpointer])
					allPointers[i_allpointer].Parent = parent
					i_allpointer++
				}
			} else {
				for y := 0; y < order; y++ {
					parent.Pointers[y] = allPointers[i_allpointer]
					parent.Keys[y] = t.NodeNumKeys(allPointers[i_allpointer])
					allPointers[i_allpointer].Parent = parent
					i_allpointer++
				}
			}
		}
		t.insertMulIntoParent(parent.Parent, parentNewSibling, parentNewSiblingNum)
	}
}

//1done
// insert parentNewSibling into parent, start index is toIndex
func (t *Tree) insertMulIntoNode(parent *Node, toIndex int, parentNewSibling []*Node, parentNewSiblingNum int) {
	//fmt.Printf("Here is insertMulIntoNode, tree: \n")
	//	t.PrintTree()
	tmpNodes := make([]*Node, order)
	i := toIndex
	i_tmpNodes := 0
	for j := 0; j < parentNewSiblingNum; j++ {
		tmpNodes[i_tmpNodes] = parentNewSibling[j]
		i_tmpNodes++
	}
	for ; i < order; i++ {
		if parent.Pointers[i] == nil {
			break
		}
		tmpNodes[i_tmpNodes] = parent.Pointers[i].(*Node)
		parent.Pointers[i] = nil
		i_tmpNodes++
	}
	/*
		fmt.Printf("tmpNodes:\n")
		for h := 0; h < i_tmpNodes; h ++ {
			t.PrintSubTree(tmpNodes[h])
		}

	*/

	r := 0
	for k := toIndex; k < order; k++ {
		parent.Pointers[k] = tmpNodes[r]
		parent.Keys[k] = t.NodeSumKeys(tmpNodes[r])
		tmpNodes[r].Parent = parent
		r++
		if r == i_tmpNodes {
			break
		}
	}
	c := parent
	parent = c.Parent
	for parent != nil {
		for i = 0; i < order; i++ {
			if reflect.DeepEqual(parent.Pointers[i], c) {
				parent.Keys[i] += parentNewSiblingNum
				c = parent
				parent = c.Parent
				break
			}
		}
	}
}

//done
func (t *Tree) insertMulLeavesIntoLeafParent(toParent *Node, startIndexInParent int, leafToInsert []interface{}, numLeafToInsert int) {
	//move existing pointers afterwards
	for j := order - 1; j >= startIndexInParent; j-- {
		if toParent.Keys[j] == 0 {
			continue
		} else {
			toParent.Keys[j+numLeafToInsert] = toParent.Keys[j]
			toParent.Pointers[j+numLeafToInsert] = toParent.Pointers[j]
			toParent.Keys[j] = 0
			toParent.Pointers[j] = nil
		}
	}
	// insert new leaves
	j := 0
	for i := startIndexInParent; i < order; i++ {
		toParent.Keys[i] = t.NodeNumKeys(leafToInsert[j].(*Node))
		toParent.Pointers[i] = leafToInsert[j]
		leafToInsert[j].(*Node).Parent = toParent
		j++
	}

	//go upwards, modify keys in ancestors
	if toParent == t.Root {
		return
	}
	c := toParent
	parent := c.Parent
	for parent != nil {
		for i := 0; i < order; i++ {
			if reflect.DeepEqual(parent.Pointers[i], c) {
				parent.Keys[i] += numLeafToInsert
				c = parent
				parent = c.Parent
				break
			}
		}
	}
	return
}

func (t *Tree) movePointerAfter(fromLeaf *Node, startIndex int, endIndex int, toLeaf *Node, toIndex int) *Node {
	//attention! toIndex can = order
	keysToInsert := make([]int, order)
	for i := startIndex; i <= endIndex; i++ {
		keysToInsert[i-startIndex] = fromLeaf.Keys[i]
	}
	numKeysToInsert := endIndex - startIndex + 1
	//insert
	if numKeysToInsert+t.NodeNumKeys(toLeaf) <= order {
		t.insertMulKeysIntoLeaf(keysToInsert, numKeysToInsert, toLeaf, toIndex)
	} else {
		t.insertMulKeysIntoLeafAfterSplitting(keysToInsert, numKeysToInsert, toLeaf, toIndex)
	}

	/*
		fmt.Printf("Here is movePointerAfter, after insert, before deleting, tree:\n")
		t.PrintTree()

	*/
	//delete
	min_key := cut(order)
	for i := startIndex; i <= endIndex; i++ {
		fromLeaf.Keys[i] = 0
	}
	/*
		fmt.Printf("again tree:\n")
		t.PrintTree()

	*/

	for i := endIndex + 1; i < order; i++ {
		if fromLeaf.Keys[i] != 0 {
			fromLeaf.Keys[i-endIndex-1] = fromLeaf.Keys[i]
			fromLeaf.Pointers[i-endIndex-1] = fromLeaf.Pointers[i]
			fromLeaf.Keys[i] = 0
			fromLeaf.Pointers[i] = nil
		}
	}
	curNumKeys := startIndex
	if fromLeaf == t.Root {
		t.adjustRoot()
		return nil
	}
	if t.NodeNumKeys(fromLeaf) >= min_key {
		if fromLeaf == t.Root {
			return nil
		}
		c := fromLeaf
		parent := c.Parent
		for parent != nil {
			for i := 0; i < order; i++ {
				if reflect.DeepEqual(parent.Pointers[i], c) {
					parent.Keys[i] -= numKeysToInsert
					c = parent
					parent = c.Parent
					break
				}
			}
		}
		return nil
	} else {
		neighbour_index := getLeftNeighborIndex(fromLeaf)
		var neighbour *Node
		if neighbour_index == -1 {
			neighbour, _ = fromLeaf.Parent.Pointers[1].(*Node)
		} else {
			neighbour, _ = fromLeaf.Parent.Pointers[neighbour_index].(*Node)
		}
		if t.NodeNumKeys(neighbour)+curNumKeys <= order {
			return t.coalesceNodes(fromLeaf, neighbour, neighbour_index, true)

		} else {
			t.redistributeNodes(fromLeaf, neighbour, neighbour_index)
			return nil
		}
	}

}

// move pointers in fromLeaf from startIndex to endIndex to toLeaf, starting from toIndex
//if no need to split, return toLeaf, nil
//if split node, return toLeaf, newnode
func (t *Tree) OldmovePointerAfter(fromLeaf *Node, startIndex int, endIndex int, toLeaf *Node, toIndex int) (*Node, *Node) {
	//moveNum := endIndex - startIndex + 1
	pointers := make([]*Node, 2*order)
	j := startIndex
	k := toIndex
	r := 0
	for ; r < 2*order; r++ {
		if j <= endIndex {
			pointers[r] = fromLeaf.Pointers[j].(*Node)
			fromLeaf.Pointers[j] = nil
			j++
		} else if toLeaf.Pointers[k] != nil {
			pointers[r] = toLeaf.Pointers[k].(*Node)
			k++
		} else {
			break
		}
	}
	u := 0
	for s := endIndex + 1; s < order; s++ {
		if fromLeaf.Pointers[s] == nil {
			break
		}
		fromLeaf.Pointers[u] = fromLeaf.Pointers[s]
		fromLeaf.Pointers[s] = nil
		u++
	}
	if toIndex+1+r <= order { // don't need to split nodes
		k := toIndex
		for i := 0; i < r; i++ {
			toLeaf.Pointers[k] = pointers[i]
			k++
		}
		return toLeaf, nil
	} else { // need to split nodes
		var new_leaf *Node
		new_leaf, err = makeLeaf()
		capacity := cut(order)
		j := 0
		for i := toIndex; i <= capacity; i++ {
			toLeaf.Pointers[i] = pointers[j]
			j++
		}
		s := 0
		for ; j < r; j++ {
			new_leaf.Pointers[s] = pointers[j]
			s++
		}
		return toLeaf, new_leaf
	}
}

//done??
func (t *Tree) swapEndPointer(leftLeaf *Node, endIndexInLeftLeaf int, rightLeaf *Node, endIndexInRightLeaf int) {
	i := 0
	j := 0
	for ; i <= endIndexInLeftLeaf; i++ {
		if leftLeaf.Pointers[i] == nil || rightLeaf.Pointers[j] == nil {
			break
		} else {
			tmpPointer := leftLeaf.Pointers[i]
			leftLeaf.Pointers[i] = rightLeaf.Pointers[j]
			rightLeaf.Pointers[j] = tmpPointer
			j++
		}
	}
	if leftLeaf.Pointers[i] == nil {
		t.movePointerAfter(rightLeaf, j, endIndexInRightLeaf, leftLeaf, i)
	} else if rightLeaf.Pointers[j] == nil {
		t.movePointerAfter(leftLeaf, i, endIndexInLeftLeaf, rightLeaf, j)
	} else { // equal number of separate pointers in left/right sides
		return
	}
}

//done 9.4
func (t *Tree) swapStartPointer(leftLeaf *Node, indexInLeftLeaf int, rightLeaf *Node, indexInRightLeaf int) {
	j := indexInRightLeaf
	i := indexInLeftLeaf
	for ; i < order; i++ {
		if leftLeaf.Keys[i] == 0 || rightLeaf.Keys[j] == 0 {
			break
		}
		tmpPointer := leftLeaf.Pointers[i]
		leftLeaf.Pointers[i] = rightLeaf.Pointers[j]
		rightLeaf.Pointers[j] = tmpPointer
		tmpKey := leftLeaf.Keys[i]
		leftLeaf.Keys[i] = rightLeaf.Keys[j]
		rightLeaf.Keys[j] = tmpKey
		j++
	}
	/*
		fmt.Printf("Here is in swapStartPointer, before move separate ones\n")
		t.PrintTree()

	*/

	if leftLeaf.Keys[i] == 0 && rightLeaf.Keys[j] == 0 {
		return
	} else if leftLeaf.Keys[i] != 0 {
		numKeysLeft := t.NodeNumKeys(leftLeaf)
		// move pointers in leftLeaf index starting from i ending xx to rightLeaf index starting from j
		t.movePointerAfter(leftLeaf, i, numKeysLeft-1, rightLeaf, j)
	} else {
		numKeysRight := t.NodeNumKeys(rightLeaf)
		t.movePointerAfter(rightLeaf, j, numKeysRight-1, leftLeaf, i)
	}
}

func (t *Tree) printLeaf(n *Node) {
	for i := 0; i < order; i++ {
		fmt.Printf("%d ", n.Keys[i])
	}
	fmt.Printf("|")
}
func (t *Tree) getIndexInParent(n *Node) int {
	parent := n.Parent
	for i := 0; i < order; i++ {
		if reflect.DeepEqual(parent.Pointers[i], n) {
			return i
		}
	}
	return -1
}

/*
// swap whole leaves
// exchange corresponding leaves, change parent
// for remain leaves in one side, delete them
// if there are enough slots to directly insert remain leaves, insert them, return nil, nil, nil, nil
// otherwise, return nil, EdgeLeaf, keysToInsert [] int, pointersToInsert []interface{}
// we need to insert keys and pointers after EdgeLeaf
// int = 0: no remain keys
// int = -1: need to insert remain leaves into left side
// int = 1: need to insert remain leaves into right side
// *Node: insert after Leaf
// int: number of remain keys/pointers
func (t *Tree) swapWholeLeaf(wholeLeaf1 []*Node, wholeLeaf2 []*Node, w1 int, w2 int) (int, error, *Node, int, []int, []interface{}) {
	fmt.Printf("Here is swapWholeLeaf\n")
	for i := 0; i < w1; i ++ {
		t.printLeaf(wholeLeaf1[i])
	}
	fmt.Printf("\n")
	for i := 0; i < w2; i ++ {
		t.printLeaf(wholeLeaf2[i])
	}
	var remain, exchangeNum int
	if w1 == w2 {
		exchangeNum = w1
		remain = 0
	} else if w1 < w2 {
		exchangeNum = w1
		remain = 1
	} else{
		exchangeNum = w2
		remain = -1 // left leaves more than right leaves
	}
	var LeftEdgeLeaf, RightEdgeLeaf *Node
	for i := 0; i < exchangeNum; i ++ {
		LeftLeaf := wholeLeaf1[i]
		RightLeaf := wholeLeaf2[i]
		LeftParent := LeftLeaf.Parent
		RightParent := RightLeaf.Parent
		LeftEdgeLeaf, RightEdgeLeaf = t.exchangeLeaf(LeftLeaf, LeftParent, RightLeaf, RightParent)
	}

	if remain == 0 {
		return 0, nil, 0, nil, nil, nil
	} else if remain == -1 { //remain leaves in left side
		wl1 := exchangeNum
		for i := 0; i < w2-w1; i ++ {
			moveLeaf1[i] = wholeLeaf1[wl1]
			wl1++
		}
		return t.moveLeaf(-1, moveLeaf1, w2-w1, RightEdgeLeaf)
	} else if remain == 1 {
		moveLeaf2 := make([]*Node, w1-w2)
		wl2 := exchangeNum
		for i := 0; i < w1-w2; i ++ {
			moveLeaf2[i] = wholeLeaf2[wl2]
			wl2++
		}
		return t.moveLeaf(1, moveLeaf2, w1-w2, LeftEdgeLeaf)
	}
	return 0, errors.New("remain should be 0 or 1 or -1"), 0, nil, nil, nil
}

*/

/*
// should move leaves in moveLeaf[] after edgeLeaf
// lrflag = -1: move into left side, 1: move into right side
// return int = 0: no remain leaves, -1: remain to be inserted into left side, 1: remain to be inserted into right side
// return *Node: should insert after this leaf
//return int: number of remain keys/pointers to insert
func (t *Tree) moveLeaf(lrflag int, moveLeaf []*Node, num int, edgeLeaf *Node) (int, error, *Node, int, []int, []interface{}) {
	i := 0
	leaf := moveLeaf[i]
	parent := leaf.Parent
	//deletedLeafParent := parent
	indexInParent := t.getIndexInParent(leaf)
	if indexInParent == -1 {
		return 0, errors.New("does not match with parent, indexInLeaf is -1"),edgeLeaf, 0, nil, nil
	}
	// delete leaves
	for i <= num {
		for j := indexInParent; j < order; j ++ {
			parent.Keys[j] = 0
			parent.Pointers[j] = nil
			i ++
			if i == num {
				break
			}
		}
		leaf = moveLeaf[i]
		parent = leaf.Parent
		indexInParent = 0
	}
	//insert leaves
	edgeParent := edgeLeaf.Parent
	edgeIndexInParent := t.getIndexInParent(edgeLeaf)
	curParentKeysNum := t.NodeNumKeys(edgeParent)
	if num + curParentKeysNum <= order { // can insert leaves directly
		pointerToStore := moveLeaf[0]
		keyToStore := t.NodeNumKeys(moveLeaf[0])
		for i := indexInParent; i < order; i ++ {
			if edgeParent.Keys[i] != 0 {
				tmpPointer := edgeParent.Pointers[i]
				tmpKey := edgeParent.Keys[i]
				edgeParent.Pointers[i] = pointerToStore
				edgeParent.Keys[i] = keyToStore
				pointerToStore = tmpPointer.(*Node)
				keyToStore = tmpKey
			} else {
				edgeParent.Pointers[i] = pointerToStore
				edgeParent.Keys[i] = keyToStore
				break
			}
		}
		return 0, nil, nil, 0, nil, nil
	} else { // no enough space to insert leaves
		leafParentKeys := make( []int, curParentKeysNum + num + 5 )
		leafParentPointers := make( []interface{}, curParentKeysNum + num + 5)
		i := 0
		for ; i <= edgeIndexInParent; i ++ {
			leafParentKeys[i] = edgeParent.Keys[i]
			leafParentPointers[i] = edgeParent.Pointers[i]
		}
		for j := 0; j < num; j ++ {
			leafParentKeys[i] = t.NodeNumKeys(moveLeaf[j])
			leafParentPointers[i] = moveLeaf[j]
			i++
		}
		for j := edgeIndexInParent + 1; j < order; j ++ {
			if edgeParent.Keys[j] != 0 {
				leafParentKeys[i] = edgeParent.Keys[j]
				leafParentPointers[i] = edgeParent.Pointers[j]
				i++
			} else {
					break
			}
		}
		return lrflag, nil, edgeLeaf, i, leafParentKeys, leafParentPointers
	}
}


*/

//done
//exchange leaves, exchange parent info
func (t *Tree) exchangeLeaf(LeftLeaf *Node, LeftParent *Node, RightLeaf *Node, RightParent *Node) (*Node, *Node) {
	var index1, index2, key1, key2 int
	for i := 0; i < order; i++ {
		if LeftParent.Pointers[i] == LeftLeaf {
			index1 = i
			key1 = LeftParent.Keys[i]
			break
		}
	}
	for i := 0; i < order; i++ {
		if RightParent.Pointers[i] == RightLeaf {
			index2 = i
			key2 = RightParent.Keys[i]
			break
		}
	}
	LeftParent.Pointers[index1] = RightLeaf
	LeftParent.Keys[index1] = key2
	RightParent.Pointers[index2] = LeftLeaf
	RightParent.Keys[index2] = key1
	RightLeaf.Parent = LeftParent
	LeftLeaf.Parent = RightParent
	return RightLeaf, LeftLeaf
}
