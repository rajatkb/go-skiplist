package goskiplist

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
)

// Data node use for holding key value pair & the skip nodes if any
type DataNode[Dt any] struct {
	key           int64
	data          Dt
	maxLevel      int8
	skipNodesNext []*DataNode[Dt]
}

func (node DataNode[Dt]) supportLevel(level int) bool {
	return int(node.maxLevel) > level
}

func (node DataNode[Dt]) getNext(level int) *DataNode[Dt] {
	return (node.skipNodesNext)[level]
}

func (node DataNode[Dt]) setNext(level int, value *DataNode[Dt]) {
	(node.skipNodesNext)[level] = value
}

type SkipListImpl[Dt any] struct {
	headList      []*DataNode[Dt] // current head of the list
	currentHeight int8            // current max heigth the Skip List has reached
	maxHeight     int8            // max height allowed for this list
	size          int64           // current size
}

/*
random level , simulates coin flip and generates level count till max lavel
*/
func (list SkipListImpl[Dt]) getRandomLevel() int8 {
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
func (list *SkipListImpl[Dt]) traverseList(returnPath bool, shortCircuit bool, key int64, startHeight int, seedPendingOperation []([]*DataNode[Dt])) ([]([]*DataNode[Dt]), *DataNode[Dt]) {
	currentExploringHeight := startHeight
	var currentNode *DataNode[Dt]

	var pendingOperations []([]*DataNode[Dt]) = seedPendingOperation

	var skipPointerList []*DataNode[Dt] = list.headList

	if pendingOperations != nil {
		skipPointerList = pendingOperations[startHeight]
	}

	if returnPath && seedPendingOperation == nil {
		pendingOperations = make([]([]*DataNode[Dt]), list.currentHeight)
	}

	tempDataNode := &DataNode[Dt]{
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

		if returnPath {
			pendingOperations[currentExploringHeight] = currentNode.skipNodesNext
		}

		if shortCircuit && currentNode.getNext(currentExploringHeight) != nil && currentNode.getNext(currentExploringHeight).key == key {
			return pendingOperations, currentNode.getNext(currentExploringHeight)
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
func (list *SkipListImpl[Dt]) Delete(key int64) bool {

	pendingOperations, foundNode := list.traverseList(true, false, key, int(list.currentHeight-1), nil)

	if foundNode.getNext(0) != nil && foundNode.getNext(0).key == key && pendingOperations != nil {
		tempNode := &DataNode[Dt]{}
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
func (list SkipListImpl[Dt]) Search(key int64) (bool, Dt) {
	_, currentNode := list.traverseList(false, true, key, int(list.currentHeight-1), nil)
	if currentNode != nil && currentNode.key == key {
		return true, currentNode.data
	}

	var value Dt
	return false, value
}

/*
Requires an Ordered set of data. Otherwise SkipList will break
Uses the ordered pairs to skip large section of this list after
the first insert. Improves insertion performance by 30%
*/
func (list *SkipListImpl[Dt]) BatchOrderedInsert(pairs []Pair[Dt]) {

	level := list.getRandomLevel()
	start := 0
	shouldReturn, _ := list.checkCurrentHead(pairs[0].key, pairs[0].value, level)
	if shouldReturn {
		start = 1
	}

	var startSkipList [][]*DataNode[Dt]

	for i := start; i < len(pairs); i++ {
		pair := pairs[i]
		level := list.getRandomLevel()

		pathToProcess, foundNode := list.traverseList(true, false, pair.key, int(list.currentHeight-1), startSkipList)

		if foundNode != nil && foundNode.key == pair.key {
			foundNode.data = pair.value // just change the data for whatever already existed
			continue                    // if shortcircuit worked we will get an existing node
		}
		tempNode := &DataNode[Dt]{}

		startSkipList = pathToProcess

		newNode := &DataNode[Dt]{
			key:           pair.key,
			data:          pair.value,
			skipNodesNext: make([]*DataNode[Dt], level),
			maxLevel:      level,
		}

		for i := int(list.currentHeight) - 1; i >= 0; i-- {
			tempNode.skipNodesNext = (pathToProcess)[i]
			if newNode.supportLevel(i) {
				newNode.setNext(i, tempNode.getNext(i))
				tempNode.setNext(i, newNode)
				startSkipList[i] = newNode.skipNodesNext
			}
		}

		prevHeight := list.currentHeight
		list.currentHeight = int8(math.Max(float64(level), float64(list.currentHeight)))
		if list.currentHeight > prevHeight {
			startSkipList = append(startSkipList, make([][]*DataNode[Dt], list.currentHeight-prevHeight)...)
			for i := list.currentHeight - 1; i > prevHeight-1; i-- {
				list.headList[i] = newNode
				startSkipList[i] = newNode.skipNodesNext
			}
		}
		list.size++
	}

}

/*
Insert a key
*/
func (list *SkipListImpl[Dt]) Insert(key int64, value Dt) *DataNode[Dt] {

	level := list.getRandomLevel()

	shouldReturn, returnValue := list.checkCurrentHead(key, value, level)
	if shouldReturn {
		return returnValue
	}

	newNode := &DataNode[Dt]{
		key:           key,
		data:          value,
		skipNodesNext: make([]*DataNode[Dt], level),
		maxLevel:      level,
	}

	pathToProcess, foundNode := list.traverseList(true, true, key, int(list.currentHeight-1), nil)

	if foundNode != nil && foundNode.key == key {
		foundNode.data = value // just change the data for whatever already existed
		return foundNode       // if shortcircuit worked we will get an existing node
	}

	prevHeight := list.currentHeight

	tempNode := &DataNode[Dt]{}

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

func (list *SkipListImpl[Dt]) checkCurrentHead(key int64, value Dt, level int8) (bool, *DataNode[Dt]) {
	if list.headList[0] == nil {
		newNode := &DataNode[Dt]{
			key:           key,
			data:          value,
			maxLevel:      level,
			skipNodesNext: make([]*DataNode[Dt], level),
		}

		list.currentHeight = level
		for i := 0; i < int(level); i++ {
			list.headList[i] = newNode
		}
		list.size++
		return true, list.headList[0]
	}
	return false, nil
}

type Pair[Dt any] struct {
	key   int64
	value Dt
}

func (list SkipListImpl[Dt]) Iterate() chan Pair[Dt] {

	ch := make(chan Pair[Dt])

	currentNode := list.headList[0]

	go func() {

		if currentNode != nil {
			for currentNode != nil {
				ch <- Pair[Dt]{key: currentNode.key, value: currentNode.data}
				currentNode = currentNode.getNext(0)
			}
		}
		close(ch)
	}()

	return ch
}

func (list SkipListImpl[Dt]) IsEmpty() bool {
	return list.headList[0] == nil
}

func (list SkipListImpl[Dt]) Size() int64 {
	return list.size
}

func (list SkipListImpl[Dt]) CurrentMaxHeight() int8 {
	return list.currentHeight
}

func (list SkipListImpl[Dt]) iterateDataNode() chan *DataNode[Dt] {

	ch := make(chan *DataNode[Dt])

	currentNode := list.headList[0]

	go func() {

		if currentNode != nil {
			for currentNode != nil {
				ch <- currentNode
				currentNode = currentNode.getNext(0)
			}
		}
		close(ch)
	}()

	return ch
}

/*
used only in debug purposes
*/
func (list SkipListImpl[Dt]) stringify(withSkips bool) string {
	var str bytes.Buffer
	str.WriteString("[ \n")
	for v := range list.iterateDataNode() {
		str.WriteString(fmt.Sprintf("( %d -> %s  ) ", v.key, fmt.Sprint(v.data)))
		str.WriteString(" SKIP : ")
		for i := 0; i < int(v.maxLevel); i++ {
			if (v.skipNodesNext)[i] != nil {
				str.WriteString(fmt.Sprintf(" %d ,", v.skipNodesNext[i].key))
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
Creates a skip list for the specified height
*/
func CreateSkipList[Dt any](maxHeight int8) SkipList[Dt] {
	return &SkipListImpl[Dt]{
		headList:      make([]*DataNode[Dt], maxHeight),
		maxHeight:     maxHeight,
		currentHeight: 0,
	}
}
