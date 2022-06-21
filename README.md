# go-skiplist
Skip List implementation in Golang. I came across the idea and wanted to try out and implementatin for same. It's work in progress the operations 
per second needs work

## API

In case you plan on using the API. 😵

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

```golang
    isEmpty := list.IsEmpty()
```


## Test

It's slow. Painfully slow. But it's ordered 🤣. Benchmarks based on logarithmic height
of total entries. Posted results are the best of 5 consequtive runs.

```
-------------Time benchamrk for Insertion against map----------
Time taken for SkipList: 155 , height : 10 
Time taken for Map : 68
Operation per mili second SkipList : 580 o/ms
Operation per mili second HashMap : 1323 o/ms
-------------Time benchamrk for Search against map----------
Time taken for SkipList: 19  height 7 , 
Time taken for Map : 2
Operation per mili second SkipList : 4736 o/ms
Operation per mili second HashMap : 45000 o/ms
```

## TODO

- Provide Batch Insert / Delete / Search - motivation being while traversing the list a sorted batch will have 
  beter total time of operation than doing these individually. Motivation is to cutdown the time of traversal
  ```
  e.g 
  list.BatchInsert([6 , 9 , 10])

  1 -> 2 -> 3 -> 4 -> 8 -> 11
  
  insert 6
  search from 1
  4 -> () -> 8

  insert 9
  search from 6 or 4 based on previous top node which will need editing

  For larger list such kind of batches should reduce average complexity
  of insertion

  ```
- Look into Compare and Swap implementation for SkipList , incorporate the same for thread safe version
