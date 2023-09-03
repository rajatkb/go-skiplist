package goskiplist

import "math/rand"

func getStaticArray() chan int {
	ret := make(chan int)
	arr := []int{5, 9, 8, 6, 1, 2, 3, 4, 4, 2, 5, 10}
	go func() {
		for i := 0; i < len(arr); i++ {
			ret <- arr[i]
		}
		close(ret)
	}()

	return ret
}

func generateIncreasingNumbers(low int, high int) chan int {
	ret := make(chan int, high-low+1)

	go func() {
		for i := low; i <= high; i++ {
			ret <- i
		}
		close(ret)
	}()

	return ret
}

func generateRandomNumber(low int, high int, count int) chan int {
	ret := make(chan int, count)
	go func() {
		for i := 0; i < count; i++ {
			ret <- int(rand.Float64()*(float64(high)-float64(low)) + float64(low))
		}
		close(ret)
	}()

	return ret
}
