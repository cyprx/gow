package gow

import (
	"log"
)

type Pool struct {
	Size        int
	InputQueue  chan Work
	OutputQueue chan Result
	QuitChan    chan bool
}

type PoolConfig struct {
	Size            int
	InputQueueSize  int
	OutputQueueSize int
}

func NewPool(config *PoolConfig) *Pool {
	size := config.Size
	inputSize := config.InputQueueSize
	outputSize := config.OutputQueueSize
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
		QuitChan:    make(chan bool),
	}
}

func (p *Pool) Start() {
	dispatcher := NewDispatcher(p.Size, p.InputQueue, p.OutputQueue)
	dispatcher.Dispatch()
	log.Println("Waiting for quit command")
	<-p.QuitChan
	dispatcher.Close()
	log.Println("Pool closed")

}

func (p *Pool) Input() chan Work {
	return p.InputQueue
}

func (p *Pool) Output() chan Result {
	return p.OutputQueue
}

func (p *Pool) Close() {
	p.QuitChan <- true
}
