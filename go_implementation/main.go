package main

import (
	"fmt"
	"math"
	"runtime"
)

const SliceLength = 10

func NewSlice() []int {
	a := make([]int, SliceLength)
	for i := 0; i < SliceLength; i++ {
		a[i] = i + 1
	}
	return a
}

func sum(a []int) int {
	res := 0
	for _, v := range a {
		res += v
	}
	return res
}

func add(start int, slice *[]int) int {
	length := len(*slice)
	isLengthEven := length%2 == 0

	if !isLengthEven {
		/*
		* якщо довжина слайсу - це непарне число,
		* а стартовий індекс знаходиться посередині,
		* то повертаємо число, яке знаходиться на стартовому індексі
		 */
		if float64(start) >= math.Floor(float64(length)/2) {
			return (*slice)[start]
		}
	}
	end := length - start - 1
	return (*slice)[start] + (*slice)[end]
}

func waveSum(slice []int, pool *ThreadPool, waveNum int) int {
	stopIndex := int(math.Ceil(float64(len(slice)) / 2))
	jobs := make([]*ThreadJob, stopIndex)

	for i := 0; i < stopIndex; i++ {
		jobs[i] = &ThreadJob{i, &slice, add}
	}
	res := pool.Process(jobs)

	fmt.Printf("Хвиля № %d, результат %+v\n", waveNum, res)

	if len(res) > 1 { // рекурсивний кейс
		waveNum++
		return waveSum(res, pool, waveNum)
	}
	return res[0] // термінальний кейс
}

func main() {
	numCPUs := runtime.NumCPU()
	slice := NewSlice()

	fmt.Printf("Максимальна к-ть потіків, для оптимальної роботи: %d\n", numCPUs)
	fmt.Printf("Довжина слайсу: %d\n", SliceLength)
	fmt.Printf("Слайс: %+v\n", slice)
	fmt.Printf("Синхронно порахований результат: %d\n", sum(slice))
	pool := NewThreadPool(numCPUs)
	pool.SpawnThreads()
	res := waveSum(slice, pool, 0)
	pool.StopThreads()

	fmt.Printf("Результат хвилевого алгоритму: %d", res)
}
