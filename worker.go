package gow

import (
	"log"
)

type Worker struct {
	ID         int
	WorkerChan chan chan Work
	WorkChan   chan Work
	ResultChan chan Result
	Quit       chan bool
}

func NewWorker(id int, workerChan chan chan Work, resultChan chan Result) Worker {
	return Worker{
		ID:         id,
		WorkerChan: workerChan,
		WorkChan:   make(chan Work),
		ResultChan: resultChan,
		Quit:       make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerChan <- w.WorkChan
			select {
			case work := <-w.WorkChan:
				log.Printf("Worker %d is running", w.ID)
				result := work.Execute()
				w.ResultChan <- result
			case <-w.Quit:
				close(w.WorkChan)
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	w.Quit <- true
}
