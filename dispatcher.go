package gow

import (
	"log"
	"time"
)

type Dispatcher struct {
	WorkerNum int
	PoolChan  chan chan Work
	Input     chan Work
	Output    chan Result
	Quit      chan bool
}

func NewDispatcher(workerNum int, input chan Work, output chan Result) *Dispatcher {
	return &Dispatcher{
		WorkerNum: workerNum,
		PoolChan:  make(chan chan Work),
		Input:     input,
		Output:    output,
		Quit:      make(chan bool),
	}
}

func (d *Dispatcher) Dispatch() {
	var workers []Worker

	for i := 0; i < d.WorkerNum; i++ {
		worker := NewWorker(i, d.PoolChan, d.Output)
		workers = append(workers, worker)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-d.Input:
				worker := <-d.PoolChan
				worker <- work
			case <-d.Quit:
				for _, worker := range workers {
					worker.Stop()
				}
			default:
				log.Println("Waiting for job assigned")
				time.Sleep(2 * time.Second)
			}

		}
	}()
	log.Println("dispatcher goroutine stops")
}
