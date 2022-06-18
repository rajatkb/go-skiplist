package goskiplist

import (
	"math"
	"math/rand"
)

// Data node use for holding key value pair & the skip nodes if any
type DataNode struct {
	key           int64
	data          []byte
	skipNodesNext *[]*DataNode
}

func (node DataNode) supportLevel(level int) bool {
	return len(*node.skipNodesNext) > level
}

func (node DataNode) getNext(level int) *DataNode {
	return (*node.skipNodesNext)[level]
}

func (node DataNode) setNext(level int, value *DataNode) {
	(*node.skipNodesNext)[level] = value
}

type SkipList struct {
	headList      []*DataNode
	currentHeight int8
	maxHeight     int8
}

func empytyFunc(d *DataNode, height int) {}

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
*/
func (list *SkipList) traverseList(returnPath bool, shortCircuit bool, key int64, call func(*DataNode, int)) (*[]*DataNode, *DataNode) {
	currentExploringHeight := int(list.currentHeight - 1)
	var currentNode *DataNode

	skipPointerList := &list.headList

	var pendingOperations []*DataNode = nil

	if returnPath {
		pendingOperations = make([]*DataNode, list.currentHeight)
	}

	for currentExploringHeight != -1 {
		// if currentExploringHeight >= len(skipPointerList) {
		// 	return fmt.Errorf("Broken list , exploring height greater then current nodes height"), false
		// }

		tempDataNode := &DataNode{
			skipNodesNext: skipPointerList,
		}
		currentNode = tempDataNode // temp extra node for start

		for currentNode.getNext(currentExploringHeight) != nil && currentNode.getNext(currentExploringHeight).key < key { // is next node lesser than key
			currentNode = currentNode.getNext(currentExploringHeight) // goto next node
		}

		if shortCircuit && currentNode.getNext(currentExploringHeight) != nil && currentNode.getNext(currentExploringHeight).key == key {
			return nil, currentNode.getNext(currentExploringHeight)
		}
		if returnPath {
			pendingOperations[currentExploringHeight] = currentNode
		}

		call(currentNode, currentExploringHeight)

		// short circuit if the node is same

		currentExploringHeight--

		skipPointerList = currentNode.skipNodesNext
	}
	if pendingOperations == nil {
		return nil, currentNode
	}
	return &pendingOperations, currentNode
}

/*
	Deletes a key
*/
func (list *SkipList) Delete(key int64) bool {

	pendingOperations, currentNode := list.traverseList(true, false, key, empytyFunc)

	if currentNode.getNext(0) != nil && currentNode.getNext(0).key == key {

		for i := int(list.currentHeight) - 1; i >= 0; i-- {
			currentNode := (*pendingOperations)[i]
			if currentNode != nil && currentNode.getNext(i) != nil && currentNode.getNext(i).supportLevel(i) {
				if currentNode.getNext(i).supportLevel(i) {
					currentNode.setNext(i, currentNode.getNext(i).getNext(i))
				} else {
					currentNode.setNext(i, nil)
				}
			}
		}
		return true
	}

	return false
}

/*
	Search a key
*/
func (list SkipList) Search(key int64) (bool, []byte) {
	_, currentNode := list.traverseList(false, true, key, empytyFunc)
	if currentNode != nil && currentNode.key == key {
		return true, currentNode.data
	}
	return false, nil
}

/*
	Insert a key
*/
func (list *SkipList) Insert(key int64, value []byte) *DataNode {

	level := list.getRandomLevel()

	if list.headList[0] == nil {
		arr := make([]*DataNode, level)
		newNode := &DataNode{
			key:           key,
			data:          value,
			skipNodesNext: &arr,
		}

		list.currentHeight = level
		for i := 0; i < int(level); i++ {
			list.headList[i] = newNode
		}
		return list.headList[0]
	}

	arr := make([]*DataNode, level)
	newNode := &DataNode{
		key:           key,
		data:          value,
		skipNodesNext: &arr,
	}

	list.traverseList(false, false, key, func(currentNode *DataNode, currentHeight int) {
		if currentNode != nil && newNode.supportLevel(currentHeight) {
			newNode.setNext(currentHeight, currentNode.getNext(currentHeight))
			currentNode.setNext(currentHeight, newNode)
		}
	})

	prevHeight := list.currentHeight
	list.currentHeight = int8(math.Max(float64(level), float64(list.currentHeight)))
	if list.currentHeight > prevHeight {
		for i := list.currentHeight - 1; i > prevHeight-1; i-- {
			list.headList[i] = newNode
		}
	}

	return nil
}

type Pair struct {
	key   int64
	value *([]byte)
}

func (list SkipList) Iterate() chan Pair {
	type arrayByte = []byte
	ch := make(chan Pair)

	currentNode := list.headList[0]
	if currentNode != nil {
		go func() {
			for currentNode != nil {
				ch <- Pair{key: currentNode.key, value: &(currentNode.data)}
				currentNode = currentNode.getNext(0)
			}
			close(ch)
		}()
	}

	return ch
}

func (list SkipList) IsEmpty() bool {
	return list.headList[0] == nil
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
