package gow

import (
	"log"
	"sync"
)

type Dispatcher struct {
	WorkerNum int
	PoolChan  chan chan Work
	Input     chan Work
	Output    chan Result
	Quit      chan bool

	wg *sync.WaitGroup
}

func NewDispatcher(workerNum int, input chan Work, output chan Result, wg *sync.WaitGroup) *Dispatcher {
	return &Dispatcher{
		WorkerNum: workerNum,
		PoolChan:  make(chan chan Work),
		Input:     input,
		Output:    output,
		Quit:      make(chan bool),
		wg:        wg,
	}
}

func (d *Dispatcher) Dispatch() {
	var workers []Worker

	for i := 0; i < d.WorkerNum; i++ {
		worker := NewWorker(i, d.PoolChan, d.Output)
		workers = append(workers, worker)
		worker.Start(d.wg)
	}

	go func() {
		for {
			select {
			case work := <-d.Input:
				worker := <-d.PoolChan
				worker <- work
				d.wg.Add(1)
			case <-d.Quit:
				for _, worker := range workers {
					worker.Stop()
				}
			}
		}
	}()
	log.Println("dispatcher goroutine stops")
}
