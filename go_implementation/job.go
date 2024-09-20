package main

type ThreadJob struct {
	start int
	slice *[]int
	/*
	* jobFn - функція, яка буде виконана в окремому потоці
	* призначена для ізоляції "бізнес-логіки"
	 */
	jobFn func(int, *[]int) int
}
