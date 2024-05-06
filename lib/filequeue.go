// Package filequeue ...
/*
Copyright © 2024 sugy <sugy.kz@gmail.com>
*/
package filequeue

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Filequeue struct is...
type Filequeue struct {
	Queue string
	Type  string
}

// Enqueue is...
func (f *Filequeue) Enqueue(t string, q string) error {
	log.Debug(fmt.Sprintf("enqueue!: Type: %s, Queue: %s", t, q))
	return nil
}

// Dequeue is...
func (f *Filequeue) Dequeue() error {
	log.Debug("dequeue!")
	return nil
}

// NewQueue is...
func NewQueue() *Filequeue {
	return &Filequeue{}
}
