// Package filequeue ...
/*
Copyright © 2024 sugy <sugy.kz@gmail.com>
*/
package filequeue

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	maildir "github.com/emersion/go-maildir"
	log "github.com/sirupsen/logrus"
)

// Queue struct is...
type Queue struct {
	Dir     maildir.Dir
	Massage string
	Type    string
}

// NewQueue is...
func NewQueue(d string) *Queue {
	q := &Queue{
		Dir: maildir.Dir(d),
	}

	err := q.setupQueuedir()
	if err != nil {
		log.Fatal(err)
	}
	return q
}

func (q *Queue) setupQueuedir() error {
	path := string(q.Dir)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.FileMode(0700))
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	if _, err := os.Stat(filepath.Join(path, "new")); errors.Is(err, os.ErrNotExist) {
		err := q.Dir.Init()
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	return nil
}

// Enqueue is...
func (q *Queue) Enqueue(t string, m string) error {
	log.Debug("enqueue!")
	q.Type, q.Massage = t, m
	log.Debug(fmt.Sprintf("Queue: %v", q))

	d, err := maildir.NewDelivery(string(q.Dir))
	if err != nil {
		log.Fatal(err)
	}
	b, err := d.Write([]byte(q.Massage))
	if err != nil {
		log.Fatal(err)
	}
	log.Debug(fmt.Sprintf("deliverd bytes: %v", b))
	err = d.Close()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// Dequeue is...
func (q *Queue) Dequeue() error {
	log.Debug("dequeue!")
	return nil
}
