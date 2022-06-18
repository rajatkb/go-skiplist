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

It's slow. Painfully slow. But it's ordered 🤣

```
-------------Time benchamrk for Insertion against map----------
Time taken for SkipList: 6467 
Time taken for Map : 71
Operation per mili second SkipList : 13 o/ms
Operation per mili second HashMap : 1267 o/ms
--- PASS: TestCompareInsert (6.54s)
=== RUN   TestCompareSearch
-------------Time benchamrk for Search against map----------
Time taken for SkipList: 11861 
Time taken for Map : 1
Operation per mili second SkipList : 7 o/ms
Operation per mili second HashMap : 90000 o/ms
--- PASS: TestCompareSearch (18.73s)
```
