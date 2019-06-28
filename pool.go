package gow

import ()

type Pool struct {
	Size        int
	InputQueue  chan Work
	OutputQueue chan Result
	Quit        chan bool
}

type PoolConfig struct {
	Size           int
	InputQueueSize int
	OuputQueueSize int
}

func NewPool(config *PoolConfig) *Pool {
	size := config.Size
	inputSize := config.InputQueueSize
	outputSize := config.OuputQueueSize
	if config.Size == 0 {
		size = 10
	}
	if inputSize == 0 {
		inputSize = 10
	}
	if outputSize == 0 {
		outputSize = 10
	}
	return &Pool{
		Size:        size,
		InputQueue:  make(chan Work, inputSize),
		OutputQueue: make(chan Result, outputSize),
	}
}

func (p *Pool) Start() {
	dispatcher := NewDispatcher(p.Size, p.InputQueue, p.OutputQueue)
	dispatcher.Dispatch()
	<-p.Quit
}

func (p *Pool) Input() chan Work {
	return p.InputQueue
}

func (p *Pool) Output() chan Result {
	return p.OutputQueue
}
