package pool

import (
	"log"
	"sync"
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
	var wg sync.WaitGroup
	input := make(chan Work, 10)
	output := make(chan Result, 10)
	pool := NewPool(5, input, output)

	pool.Start()

	for i := 0; i < 10; i++ {
		work := NewWork(i)
		input <- work
		wg.Add(1)
	}

	go func() {
		for {
			select {
			case res := <-output:
				log.Printf("output: %v", res)
				wg.Done()
			}
		}
	}()
	wg.Wait()

	log.Print("Job finished")
}
