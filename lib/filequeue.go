// Package filequeue ...
/*
Copyright © 2024 sugy <sugy.kz@gmail.com>
*/
package filequeue

import (
	"fmt"

	maildir "github.com/emersion/go-maildir"
	log "github.com/sirupsen/logrus"
)

// Filequeue struct is...
type Filequeue struct {
	Dir   maildir.Dir
	Queue string
	Type  string
}

// NewQueue is...
func NewQueue(d string) *Filequeue {
	f := &Filequeue{
		Dir: maildir.Dir(d),
	}
	return f
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
