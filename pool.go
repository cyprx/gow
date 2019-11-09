package gow

import (
	"log"

	"github.com/google/uuid"
)

type Pool struct {
	Name        string
	Size        int
	InputQueue  chan Work
	OutputQueue chan Result
	QuitChan    chan bool
}

type PoolConfig struct {
	Name string
	Size int
}

func NewPool(config *PoolConfig) *Pool {
	name := config.Name
	size := config.Size
	if name == "" {
		id := uuid.New()
		name = id.String()
	}
	if size == 0 {
		size = 10
	}
	return &Pool{
		Name:        name,
		Size:        size,
		InputQueue:  make(chan Work, 1),
		OutputQueue: make(chan Result, 1),
		QuitChan:    make(chan bool),
	}
}

func (p *Pool) Start() {
	log.Printf("Pool %s is starting", p.Name)

	dispatcher := NewDispatcher(p.Size, p.InputQueue, p.OutputQueue)
	dispatcher.Dispatch()
	<-p.QuitChan
	dispatcher.Close()
	log.Println("Pool closed")
}

func (p *Pool) Submit(w Work) {
	p.InputQueue <- w
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
