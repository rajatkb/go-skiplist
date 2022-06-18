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
of total entries

```
-------------Time benchamrk for Insertion against map----------
Time taken for SkipList: 321 , height : 16 
Time taken for Map : 73
Operation per mili second SkipList : 280 o/ms
Operation per mili second HashMap : 1232 o/ms
--- PASS: TestCompareInsert (0.39s)
=== RUN   TestCompareSearch
-------------Time benchamrk for Search against map----------
Time taken for SkipList: 23  height 16 , 
Time taken for Map : 11
Operation per mili second SkipList : 3913 o/ms
Operation per mili second HashMap : 8181 o/ms
--- PASS: TestCompareSearch (0.41s)
```

## TODO

- Investigate improvement for search and insert
- Look into Compare and Swap implementation for SkipList , incorporate the same for thread safe version
