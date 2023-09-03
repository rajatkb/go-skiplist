package goskiplist

type SkipList[Dt any] interface {
	Delete(key int64) bool
	/*
		Search a key
	*/
	Search(key int64) (bool, Dt)

	/*
		Requires an Ordered set of data. Otherwise SkipList will break
		Uses the ordered pairs to skip large section of this list after
		the first insert. Improves insertion performance by 30%
	*/
	BatchOrderedInsert(pairs []Pair[Dt])
	/*
		Insert a key
	*/
	Insert(key int64, value Dt) *DataNode[Dt]
	Iterate() chan Pair[Dt]
	IsEmpty() bool
	Size() int64

	CurrentMaxHeight() int8
	// Scan(startKey int64, endKey int64) []Pair[Dt]
}
