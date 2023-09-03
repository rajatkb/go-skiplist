# go-skiplist
Skip List implementation in Golang. 


## API

In case you plan on using the API. 😵

### Highlights
- Completely for in memory use case.
- Generic supported implementation

- Insert

```go
   node := list.Insert(int64(num), []byte(fmt.Sprintf("value - %d", num)))
```

- Delete
```go
    deleted := list.Delete(int64(num))
```

- Search
```go
    found , value := list.Search(42)
```

- Iterate
```go
    for pair := range list.Iterate() {
	pair.key 
        pair.value
    }
```
- IsEmpty

```go
    isEmpty := list.IsEmpty()
```

- BatchOrderedInsert

```go
// experimental
// requirs an ordered slice of pairs
// imporeves insert performance on avg by 33%
// 
    var pairs []skiplist.Pair
    list.BatchOrderedInsert(pairs)

```

## Test

It's slow. Painfully slow. But it's ordered 🤣. Benchmarks based on logarithmic height
of total entries. posted results are consistently reproducible on a Ryzen 5600x machine with virtualized 2 Core + 4Gb in WSL 

```
=== RUN   TestCompareInsert
-------------Time benchamrk for Insertion against map----------
Time taken for SkipList: 76 , height : 16
Time taken for Map : 25
Operation per mili second SkipList : 1184 o/ms
Operation per mili second HashMap : 3600 o/ms
--- PASS: TestCompareInsert (0.10s)
=== RUN   TestCompareSearch
-------------Time benchamrk for Search against map----------
Time taken for SkipList: 29  height 16 ,
Time taken for Map : 3
Operation per mili second SkipList : 3103 o/ms
Operation per mili second HashMap : 30000 o/ms
--- PASS: TestCompareSearch (0.12s)
=== RUN   TestMultiInsert
-------------Time benchamrk for Multi Insertion against map----------
Time taken for SkipList: 51 ms, height : 16
Operation per mili second SkipList : 1764 o/ms , @ batch size : 500
--- PASS: TestMultiInsert (0.06s)
```

## TODO
- Add Batched Search and Delete also
- Add Range Scan operation 
- Add Batched Multi Operation support
- Add Doubly Linked List supports and Lookup from Tail to optimize for Long Range Scans
- Look into Compare and Swap implementation for SkipList , incorporate the same for thread safe version

