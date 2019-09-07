package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"../stree"
)

func fromGOB64(binary []byte) stree.Tree {
	t := stree.Tree{}
	r := bytes.NewReader(binary)
	d := gob.NewDecoder(r)
	err := d.Decode(&t)
	if err != nil {
		fmt.Println(`failed gob Decode`, err)
	}
	return t
}

func fetchStree(path string) stree.Tree {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("StreeFile does not exist.")
	} else {
		log.Printf("StreeFile exists, open it.")
	}
	streeData, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	stree := fromGOB64(streeData)
	return stree
}
func toGOB64(t stree.Tree) []byte {
	gob.Register(stree.Node{})
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(t)
	if err != nil {
		fmt.Println(`failed gob Encode`, err)
	}
	return b.Bytes()
}
func saveStree(t stree.Tree) {
	fmt.Println("Here is saveStree.")
	serializedStree := toGOB64(t)
	//var rootDirectory = flag.String("root", "/home/ljy/Desktop/grid-home/home/source/proxy/userdata/workspace-TESTUUID/", "root directory for user files")
	var rootDirectory = "/home/ljy/Desktop/"
	fmt.Printf("rootDirectory = %s.", rootDirectory)
	err := ioutil.WriteFile(rootDirectory+"streeData.serialized", serializedStree, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	t := stree.NewTree()
	for i := 1; i <= 30; i++ {
		rowID := 30 - i + 1
		pos := 1
		t.Insert(rowID, pos)
	}
	t.PrintTree()
	saveStree(*t)
	ft := fetchStree("/home/ljy/Desktop/streeData.serialized")
	ft.PrintTree()
}
