package main

import (
	"sync"
)

type ThreadPool struct {
	size        int       // розмір пулу, вона ж к-ть потоків, якими оперує пул
	pool        []*Thread // пул потоків
	jobsQueue   []*ThreadJob
	resultQueue []int
	mu          *sync.Mutex
	jobWg       *sync.WaitGroup // вейтгрупа для очікування завершення роботи потоків
	stopWg      *sync.WaitGroup // вейтгрупа для очікування закриття потоків
}

func NewThreadPool(size int) *ThreadPool {
	mu := new(sync.Mutex)
	return &ThreadPool{
		size:        size,
		pool:        make([]*Thread, size),
		jobsQueue:   nil,
		resultQueue: nil,
		mu:          mu,
		jobWg:       new(sync.WaitGroup),
		stopWg:      new(sync.WaitGroup),
	}
}

func (p *ThreadPool) SpawnThreads() {
	// ініціалізуємо та запускаємо потоки
	for i := 0; i < p.size; i++ {
		p.stopWg.Add(1)
		thread := NewThread(i, p.nextJob, p.receiveResult, p.mu, p.jobWg, p.stopWg)
		p.pool[i] = thread // додаємо потік в пул
		go thread.Run()    // запускаємо потік, який очікує отримання роботи
	}
}

func (p *ThreadPool) StopThreads() {
	for _, thread := range p.pool {
		thread.CanWork = false
	}
	p.stopWg.Wait()
}

func (p *ThreadPool) Process(jobs []*ThreadJob) []int {
	p.jobWg.Add(len(jobs))
	p.resetQueues(jobs)
	p.jobWg.Wait() // чекаємо, поки потоки завершать роботу
	return p.resultQueue
}

func (p *ThreadPool) nextJob() *ThreadJob {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.jobsQueue) > 0 {
		var job *ThreadJob
		job, p.jobsQueue = p.jobsQueue[0], p.jobsQueue[1:]
		return job
	}
	return nil
}

func (p *ThreadPool) receiveResult(res int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.resultQueue = append(p.resultQueue, res)
}

func (p *ThreadPool) resetQueues(jobs []*ThreadJob) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.resultQueue = make([]int, 0)
	p.jobsQueue = jobs
}
