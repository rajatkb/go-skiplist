package goskiplist

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func getElapsed(call func()) int64 {
	now := time.Now().UnixMilli()
	ch := make(chan bool)
	go func() {
		call()
		close(ch)
	}()
	<-ch
	return time.Now().UnixMilli() - now
}

func TestInsert(t *testing.T) {

	total := 30000
	list := CreateSkipList(int8(math.Log2(float64(total))))
	for num := range generateRandomNumber(0, 100, total) {
		data := []byte(fmt.Sprintf("value - %d", num))
		list.Insert(int64(num), &data)
	}

	value := -1
	for v := range list.Iterate() {
		if value <= int(v.key) {
			value = int(v.key)
		} else {
			t.Errorf(`Order for insertion is wrong`)
			break
		}
	}
}

func TestMultiInsert(t *testing.T) {

	total := 30000
	list := CreateSkipList(int8(math.Log2(float64(total))))
	batchSize := 100
	pairs := make([]Pair, batchSize)
	i := 0
	elapsed := getElapsed(func() {
		for num := range generateIncreasingNumbers(1, total) {
			pairs[i].key = int64(num)

			data := []byte(fmt.Sprintf("value - %d", num))
			pairs[i].value = &data
			i++
			if i == batchSize {
				i = 0
				list.BatchOrderedInsert(pairs)
			}
		}
	})

	value := -1
	for v := range list.Iterate() {
		if value <= int(v.key) {
			value = int(v.key)
		} else {
			t.Errorf(`Order for insertion is wrong`)
			break
		}
	}

	elapsed = int64(math.Max(float64(elapsed), 1.0))

	fmt.Printf("-------------Time benchamrk for Multi Insertion against map----------\n")
	fmt.Printf("Time taken for SkipList: %d ms, height : %d \n", elapsed, list.currentHeight)
	fmt.Printf("Operation per mili second SkipList : %d o/ms \n", int64(total)/elapsed)

}

func TestDelete(t *testing.T) {

	total := 30000
	list := CreateSkipList(int8(math.Log2(float64(total))))
	for num := range generateRandomNumber(0, 100, total) {
		data := []byte(fmt.Sprintf("value - %d", num))
		list.Insert(int64(num), &data)
		list.Delete(int64(num))
	}

	if !list.IsEmpty() {
		t.Errorf("List Deletion failed")
	}
}

func TestSearch(t *testing.T) {
	total := 30000
	list := CreateSkipList(int8(math.Log2(float64(total))))
	for num := range generateRandomNumber(200, 300, total) {
		data := []byte(fmt.Sprintf("value - %d", num))
		list.Insert(int64(num), &data)

	}
	testText := "I am 21"
	data := []byte(testText)
	list.Insert(21, &data)
	found, value := list.Search(21)
	if !found || string(*value) != testText {
		t.Errorf("List Deletion failed")
	}
}

func TestSize(t *testing.T) {
	list := CreateSkipList(4)
	total := 16
	numbers := make(map[int]bool, total)
	count := 0
	for num := range getStaticArray() {
		data := []byte(fmt.Sprintf("value - %d - %d", num, count))
		list.Insert(int64(num), &data)
		numbers[num] = true
		count++
		// fmt.Println(list.Stringify(true))
	}

	if int(list.Size()) != len(numbers) {
		t.Errorf("List size is incorrect after insertion %d", list.Size())
	}

	for k, _ := range numbers {
		// fmt.Println(k)
		list.Delete(int64(k))
		// fmt.Println(list.Stringify(true))
	}

	if list.Size() != 0 {
		// fmt.Println(list.Stringify(true))
		t.Errorf("List size is incorrect after deletion %d", list.Size())
	}
}

func TestCompareInsert(t *testing.T) {

	hashMap := make(map[int64]*[]byte)
	total := 90000
	list := CreateSkipList(int8(math.Log2(float64(total))))
	timeForList := getElapsed(func() {
		for num := range generateRandomNumber(200, 300, total) {
			data := []byte(fmt.Sprintf("value - %d", num))
			list.Insert(int64(num), &data)
		}
	})

	timeForMap := getElapsed(func() {
		for num := range generateRandomNumber(200, 300, total) {
			data := []byte(fmt.Sprintf("value - %d", num))
			hashMap[int64(num)] = &data
		}
	})
	fmt.Printf("-------------Time benchamrk for Insertion against map----------\n")
	fmt.Printf("Time taken for SkipList: %d , height : %d \n", timeForList, list.currentHeight)
	fmt.Printf("Time taken for Map : %d \n", timeForMap)
	fmt.Printf("Operation per mili second SkipList : %d o/ms \n", int64(total)/timeForList)
	fmt.Printf("Operation per mili second HashMap : %d o/ms \n", int64(total)/timeForMap)

}

func TestCompareSearch(t *testing.T) {

	hashMap := make(map[int64]*[]byte)
	total := 90000
	list := CreateSkipList(int8(math.Log2(float64(total))))
	numbers := make([]int, total)
	i := 0
	for num := range generateRandomNumber(200, 300, total) {
		data := []byte(fmt.Sprintf("value - %d", num))
		list.Insert(int64(num), &data)
		hashMap[int64(num)] = &data
		numbers[i] = num
		i++
	}

	timeForList := getElapsed(func() {
		for i := 0; i < total; i++ {
			_, _ = list.Search(int64(numbers[i]))
		}
	})

	timeForMap := getElapsed(func() {
		for i := 0; i < total; i++ {
			_, _ = hashMap[int64(numbers[i])]
		}
	})

	if timeForMap == 0 {
		timeForMap = 1
	}

	fmt.Printf("-------------Time benchamrk for Search against map----------\n")
	fmt.Printf("Time taken for SkipList: %d  height %d , \n", timeForList, list.currentHeight)
	fmt.Printf("Time taken for Map : %d \n", timeForMap)
	fmt.Printf("Operation per mili second SkipList : %d o/ms \n", int64(total)/(timeForList))
	fmt.Printf("Operation per mili second HashMap : %d o/ms \n", int64(total)/timeForMap)

}
