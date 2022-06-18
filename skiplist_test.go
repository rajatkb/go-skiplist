package goskiplist

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func generateRandomNumber(low int, high int, count int) chan int {
	ret := make(chan int)
	go func() {
		for i := 0; i < count; i++ {
			ret <- int(rand.Float64()*(float64(high)-float64(low)) + float64(low))
		}
		close(ret)
	}()

	return ret
}

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
	list := CreateSkipList(5)
	total := 30000
	for num := range generateRandomNumber(0, 100, total) {
		list.Insert(int64(num), []byte(fmt.Sprintf("value - %d", num)))
	}
	count := 0
	for range list.Iterate() {
		count++
	}
	if count != total {
		t.Errorf(`Total count not same`)
	}
}

func TestDelete(t *testing.T) {
	list := CreateSkipList(5)
	total := 30000
	for num := range generateRandomNumber(0, 100, total) {
		list.Insert(int64(num), []byte(fmt.Sprintf("value - %d", num)))
		list.Delete(int64(num))
	}

	if !list.IsEmpty() {
		t.Errorf("List Deletion failed")
	}
}

func TestSearch(t *testing.T) {
	list := CreateSkipList(5)
	total := 30000
	for num := range generateRandomNumber(200, 300, total) {
		list.Insert(int64(num), []byte(fmt.Sprintf("value - %d", num)))

	}
	testText := "I am 21"
	list.Insert(21, []byte(testText))
	found, value := list.Search(21)
	if !found || string(value) != testText {
		t.Errorf("List Deletion failed")
	}
}

func TestCompareInsert(t *testing.T) {
	list := CreateSkipList(5)
	hashMap := make(map[int64][]byte)
	total := 90000
	timeForList := getElapsed(func() {
		for num := range generateRandomNumber(200, 300, total) {
			list.Insert(int64(num), []byte(fmt.Sprintf("value - %d", num)))
		}
	})

	timeForMap := getElapsed(func() {
		for num := range generateRandomNumber(200, 300, total) {
			hashMap[int64(num)] = []byte(fmt.Sprintf("value - %d", num))
		}
	})
	fmt.Printf("-------------Time benchamrk for Insertion against map----------\n")
	fmt.Printf("Time taken for SkipList: %d \n", timeForList)
	fmt.Printf("Time taken for Map : %d \n", timeForMap)
	fmt.Printf("Operation per mili second SkipList : %d o/ms \n", int64(total)/timeForList)
	fmt.Printf("Operation per mili second HashMap : %d o/ms \n", int64(total)/timeForMap)

}

func TestCompareSearch(t *testing.T) {
	list := CreateSkipList(5)
	hashMap := make(map[int64][]byte)
	total := 90000
	numbers := make([]int, total)
	i := 0
	for num := range generateRandomNumber(200, 300, total) {
		list.Insert(int64(num), []byte(fmt.Sprintf("value - %d", num)))
		hashMap[int64(num)] = []byte(fmt.Sprintf("value - %d", num))
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
	fmt.Printf("Time taken for SkipList: %d \n", timeForList)
	fmt.Printf("Time taken for Map : %d \n", timeForMap)
	fmt.Printf("Operation per mili second SkipList : %d o/ms \n", int64(total)/(timeForList))
	fmt.Printf("Operation per mili second HashMap : %d o/ms \n", int64(total)/timeForMap)

}
