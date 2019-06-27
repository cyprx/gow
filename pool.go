package gow

import (
	"sync"
)

type Pool struct {
	Size   int
	Input  chan Work
	Output chan Result
}

func NewPool(size int, input chan Work, output chan Result) *Pool {
	return &Pool{size, input, output}
}

func (p *Pool) Start() {
	var wg sync.WaitGroup
	dispatcher := NewDispatcher(5, p.Input, p.Output, &wg)

	dispatcher.Dispatch()

	wg.Wait()
}
