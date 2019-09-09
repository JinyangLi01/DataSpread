package encode

import (
	"log"
	"strconv"
	"strings"

	"../stree"

	// "io/ioutil"
	"os"
)

var myqueue *stree.Node

func enqueue(new_node *stree.Node) {
	var c *stree.Node
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

func dequeue() *stree.Node {
	n := myqueue
	myqueue = myqueue.Next
	n.Next = nil
	return n
}

func pathToRoot(child *stree.Node) int {
	length := 0
	c := child
	parent := c.Parent
	for parent != nil {
		c = c.Parent
		parent = c.Parent
		length += 1
	}
	return length
}

func WriteToFile(filePath string, t *stree.Tree) {
	log.Printf("Here is WriteToFile in encode.go\n")
	//writeData := EncodeStree(t)
	f, err := os.OpenFile(filePath, os.O_RDWR, 0755)
	if err == nil {
		f.Close()
		err = os.Remove(filePath)
		if err != nil {
			log.Fatal(err)
		}
	}

	encodedStree := EncodeStree(t)
	//fmt.Printf(encodedStree)

	f, err = os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	//s1 := "hahaha\nhehehe"
	f.WriteString(encodedStree)
	f.Close()

}

func ReadFromFile(filePath string) *stree.Tree {
	//writeData := EncodeStree(t)
	f, err := os.OpenFile(filePath, os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fi, _ := f.Stat()
	size := fi.Size()
	readdata := make([]byte, size)
	f.Read(readdata)
	//	log.Printf("Here is ReadFromFile, str:\n")
	str := string(readdata)
	//	log.Printf(str)
	//	log.Printf("\n")
	//fmt.Printf("Before DecodeStree\n")
	t := DecodeStree(str)
	//	log.Printf("Here is before return from ReadFromFile, tree:\n")
	//	t.PrintTree()
	return t
}

func EncodeStree(t *stree.Tree) string {
	var str string
	var n *stree.Node
	i := 0
	rank := 0
	new_rank := 0
	if t.Root == nil {
		return str
	}
	str = strconv.Itoa(t.Order)
	str += "\n"
	myqueue = nil
	enqueue(t.Root)
	for myqueue != nil {
		n = dequeue()
		if n != nil {
			//if n.Parent != nil && n == n.Parent.Pointers[0] {
			if n.Parent != nil {
				rank = new_rank
				new_rank = pathToRoot(n) // height
				if new_rank != rank {
					str = str + "\n"
				}
			}
			for i = 0; i < t.Order; i++ {
				str1 := strconv.Itoa(n.Keys[i])
				str = str + str1
				str = str + " "
				//fmt.Printf("%d ", n.Keys[i])
			}
			if !n.IsLeaf {
				for i = 0; i < t.Order; i++ {
					if n.Keys[i] != 0 {
						c, _ := n.Pointers[i].(*stree.Node)
						enqueue(c)
					} else {
						break
					}
				}
			}
			str = str + "| "
			//fmt.Printf(" | ")
		}
	}
	return str
}

func DecodeStree(str string) *stree.Tree {
	//	log.Printf("Here is DecodeStree\n")
	myqueue = nil
	t := &stree.Tree{}
	createNode, _ := stree.MakeNode()
	t.Root = createNode
	enqueue(createNode)
	//	log.Printf("str = %s\n", str)
	levels := strings.Split(str, "\n")
	//	log.Printf("levels = \n")
	//	log.Print(levels)
	t.Order, _ = strconv.Atoi(levels[0])
	numLevel := len(levels)
	//	log.Printf("numLevel = %d\n", numLevel)
	for i := 1; i < numLevel; i++ {
		//		log.Printf("level[%d] = %s\n", i, levels[i])
		nodes := strings.Split(levels[i], " | ")
		numNode := len(nodes)
		if nodes[numNode-1] == "" {
			numNode--
		}
		for j := 0; j < numNode; j++ {
			keys := strings.Split(nodes[j], " ")
			newNode := dequeue()
			for k := 0; k < t.Order; k++ {
				newNode.Keys[k], _ = strconv.Atoi(keys[k])
				if i == numLevel-1 { // leaf
					newNode.IsLeaf = true
				}
				if newNode.Keys[k] != 0 && i != numLevel-1 {
					createNode, _ := stree.MakeNode()
					newNode.Pointers[k] = createNode
					createNode.Parent = newNode
					enqueue(createNode)
				}
			}
		}
	}
	t.NumKeys = t.NodeSumKeys(t.Root)
	return t
}
