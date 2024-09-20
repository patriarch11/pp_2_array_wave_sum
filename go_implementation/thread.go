package main

import (
	"sync"
)

type Thread struct {
	CanWork     bool
	id          int
	jobQueue    func() *ThreadJob
	resultQueue func(int)
	resultChan  chan<- int // канал лише для запису результату
	mu          *sync.Mutex
	jobWg       *sync.WaitGroup
	stopWg      *sync.WaitGroup
}

func NewThread(
	id int,
	jobQueue func() *ThreadJob,
	resultQueue func(int),
	mu *sync.Mutex,
	jobWg *sync.WaitGroup,
	stopWg *sync.WaitGroup,
) *Thread {
	return &Thread{
		id:          id,
		jobQueue:    jobQueue,
		resultQueue: resultQueue,
		mu:          mu,
		jobWg:       jobWg,
		stopWg:      stopWg,
		CanWork:     true,
	}
}

func (t *Thread) Run() {
	defer t.stopWg.Done()
	for {
		if !t.CanWork {
			break
		}
		job := t.jobQueue()
		if job != nil {
			//fmt.Printf("Потік № %d почав роботу з індексу %d\n", t.id, job.start)
			//time.Sleep(2 * time.Second)
			result := job.jobFn(job.start, job.slice)
			//fmt.Printf("Потік № %d відпрацював з результатом: %d\n", t.id, result)
			t.resultQueue(result)
			t.jobWg.Done()
		}
	}
	//fmt.Printf("Потік № %d зупинено\n", t.id)
}
