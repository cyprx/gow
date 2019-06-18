package main

import (
	"log"
	"sync"
	"time"
)

type Job interface {
	Run() interface{}
}

type Result struct {
	data interface{}
}

type Work struct {
	ID  int
	Job Job
}

func NewWork(job Job) Work {
}

func (w Work) DoWork() Result {
	data := w.Job.Run()
	return Result{data}
}

type Worker struct {
	ID            int
	WorkerChannel chan chan Work
	WorkChan      chan Work
	Output        chan Result
}

func NewWorker(ID int, workerChannel chan chan Work, output chan Result) *Worker {
	return &Worker{ID, workerChannel, make(chan Work), output}
}

func (w *Worker) run() {
	go func() {
		for {
			w.WorkerChannel <- w.WorkChan
			select {
			case work := <-w.WorkChan:
				log.Printf("Worker %d is working", w.ID)
				w.Output <- work.DoWork()
			}
		}
	}()
}

func main() {
	var wg sync.WaitGroup

	workerChannel := make(chan chan Work)
	workQueue := make(chan Work, 10)
	output := make(chan Result)

	for i := 0; i < 3; i++ {
		worker := NewWorker(i, workerChannel, output)
		worker.run()
	}

	go func() {
		for {
			select {
			case work := <-workQueue:
				worker := <-workerChannel
				worker <- work
			}
		}
	}()

	for i := 0; i < 5; i++ {
		work := Work{i, 1 * time.Second}
		workQueue <- work
	}

	for len(output) > 0 {
	}

	wg.Wait()
}
