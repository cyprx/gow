package gow

import (
	"log"
	"testing"
	"time"
)

type sleepJob struct {
	id int
}

func NewWork(id int) Work {
	return &sleepJob{id}
}

func (w *sleepJob) Execute() interface{} {
	time.Sleep(2 * time.Second)
	return true
}

func TestPool(T *testing.T) {
	pool := NewPool(&PoolConfig{Size: 2})

	go func() {
		for i := 0; i < 10; i++ {
			work := NewWork(i)
			pool.Input() <- work
		}
	}()

	go func() {
		for {
			select {
			case res := <-pool.Output():
				log.Printf("output: %v", res)
			}
		}
	}()

	go func() {
		time.Sleep(10 * time.Second)
		pool.Close()
	}()
	pool.Start()

	log.Print("Job finished")
}
