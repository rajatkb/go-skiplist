package goskiplist

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
)

// Data node use for holding key value pair & the skip nodes if any
type DataNode struct {
	key           int64
	data          *[]byte
	maxLevel      int8
	skipNodesNext []*DataNode
}

func (node DataNode) supportLevel(level int) bool {
	return int(node.maxLevel) > level
}

func (node DataNode) getNext(level int) *DataNode {
	return (node.skipNodesNext)[level]
}

func (node DataNode) setNext(level int, value *DataNode) {
	(node.skipNodesNext)[level] = value
}

type SkipList struct {
	headList      []*DataNode
	currentHeight int8
	maxHeight     int8
	size          int64
}

/*
	random level , simulates coin flip and generates level count till max lavel
*/
func (list SkipList) getRandomLevel() int8 {
	count := int8(1)
	for rand.Float32() > 0.5 && count < list.maxHeight {
		count++
	}
	return count
}

/*
	traverses list to find an entry right before the key
	- useful for search
	- useful for insertion
	- useful for deletion

	Concurrency solution


*/
func (list *SkipList) traverseList(returnPath bool, shortCircuit bool, key int64) ([]([]*DataNode), *DataNode) {
	currentExploringHeight := int(list.currentHeight - 1)
	var currentNode *DataNode

	skipPointerList := list.headList

	var pendingOperations []([]*DataNode) = nil

	if returnPath {
		pendingOperations = make([]([]*DataNode), list.currentHeight)
	}

	tempDataNode := &DataNode{
		skipNodesNext: skipPointerList,
	}

	for currentExploringHeight != -1 {
		// if currentExploringHeight >= len(skipPointerList) {
		// 	return fmt.Errorf("Broken list , exploring height greater then current nodes height"), false
		// }

		// since operation is done in the nextPointer list
		// we create a temp node and keep changing it's skip pointer
		tempDataNode.skipNodesNext = skipPointerList

		currentNode = tempDataNode // temp extra node for start

		for currentNode.getNext(currentExploringHeight) != nil && currentNode.getNext(currentExploringHeight).key < key { // is next node lesser than key
			currentNode = currentNode.getNext(currentExploringHeight) // goto next node
		}

		if shortCircuit && currentNode.getNext(currentExploringHeight) != nil && currentNode.getNext(currentExploringHeight).key == key {
			return nil, currentNode.getNext(currentExploringHeight)
		}

		if returnPath {
			pendingOperations[currentExploringHeight] = currentNode.skipNodesNext
		}

		currentExploringHeight--

		skipPointerList = currentNode.skipNodesNext

	}
	if pendingOperations == nil {
		return nil, currentNode
	}
	return pendingOperations, currentNode
}

/*
	Deletes a key
*/
func (list *SkipList) Delete(key int64) bool {

	pendingOperations, foundNode := list.traverseList(true, false, key)

	if foundNode.getNext(0) != nil && foundNode.getNext(0).key == key && pendingOperations != nil {
		tempNode := &DataNode{}
		for i := int(list.currentHeight) - 1; i >= 0; i-- {
			// TODO pull outside
			tempNode.skipNodesNext = pendingOperations[i]
			currentNode := tempNode
			if currentNode.getNext(i) != nil {
				if currentNode.getNext(i).supportLevel(i) && currentNode.getNext(i).key == key {
					currentNode.setNext(i, currentNode.getNext(i).getNext(i))
				}
			}

		}
		list.size--

		return true
	}

	return false
}

/*
	Search a key
*/
func (list SkipList) Search(key int64) (bool, *[]byte) {
	_, currentNode := list.traverseList(false, true, key)
	if currentNode != nil && currentNode.key == key {
		return true, currentNode.data
	}
	return false, nil
}

/*
	Insert a key
*/
func (list *SkipList) Insert(key int64, value *[]byte) *DataNode {

	level := list.getRandomLevel()

	if list.headList[0] == nil {
		newNode := &DataNode{
			key:           key,
			data:          value,
			maxLevel:      level,
			skipNodesNext: make([]*DataNode, level),
		}

		list.currentHeight = level
		for i := 0; i < int(level); i++ {
			list.headList[i] = newNode
		}
		list.size++
		return list.headList[0]
	}

	newNode := &DataNode{
		key:           key,
		data:          value,
		skipNodesNext: make([]*DataNode, level),
		maxLevel:      level,
	}

	pathToProcess, foundNode := list.traverseList(true, true, key)

	if foundNode != nil && foundNode.key == key {
		foundNode.data = value // just change the data for whatever already existed
		return foundNode       // if shortcircuit worked we will get an existing node
	}

	if pathToProcess == nil {
		return nil
	}

	prevHeight := list.currentHeight

	tempNode := &DataNode{}

	for i := int(prevHeight) - 1; i >= 0; i-- {
		tempNode.skipNodesNext = (pathToProcess)[i]
		if newNode.supportLevel(i) {
			newNode.setNext(i, tempNode.getNext(i))
			tempNode.setNext(i, newNode)
		}
	}

	list.currentHeight = int8(math.Max(float64(level), float64(list.currentHeight)))
	if list.currentHeight > prevHeight {
		for i := list.currentHeight - 1; i > prevHeight-1; i-- {
			list.headList[i] = newNode
		}
	}
	list.size++
	return newNode
}

type Pair struct {
	key   int64
	value *([]byte)
	node  *DataNode
}

func (list SkipList) Iterate() chan Pair {
	type arrayByte = []byte
	ch := make(chan Pair)

	currentNode := list.headList[0]

	go func() {

		if currentNode != nil {
			for currentNode != nil {
				ch <- Pair{key: currentNode.key, value: currentNode.data, node: currentNode}
				currentNode = currentNode.getNext(0)
			}
		}
		close(ch)
	}()

	return ch
}

func (list SkipList) IsEmpty() bool {
	return list.headList[0] == nil
}

func (list SkipList) Size() int64 {
	return list.size
}

func (list SkipList) Stringify(withSkips bool) string {
	var str bytes.Buffer
	str.WriteString("[ \n")
	for v := range list.Iterate() {
		str.WriteString(fmt.Sprintf("( %d -> %s  ) ", v.key, string(*v.value)))
		str.WriteString(" SKIP : ")
		for i := 0; i < int(v.node.maxLevel); i++ {
			if (v.node.skipNodesNext)[i] != nil {
				str.WriteString(fmt.Sprintf(" %d ,", (v.node.skipNodesNext)[i].key))
			} else {
				str.WriteString(fmt.Sprintf(" nil "))
			}
		}

		str.WriteString(" \n")
	}
	str.WriteString("] \n")
	return str.String()
}

/*
	Create a maxHeight
*/
func CreateSkipList(maxHeight int8) *SkipList {
	return &SkipList{
		headList:      make([]*DataNode, maxHeight),
		maxHeight:     maxHeight,
		currentHeight: 0,
	}
}
