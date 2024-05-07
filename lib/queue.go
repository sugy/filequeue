// Package filequeue ...
/*
Copyright © 2024 sugy <sugy.kz@gmail.com>
*/
package filequeue

import (
	"errors"
	"fmt"
	"io"
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
	keys, err := q.Dir.Unseen()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(fmt.Sprintf("new keys: %v", keys))

	err = q.Dir.Walk(func(key string, flags []maildir.Flag) error {
		log.Info(fmt.Sprintf("%v, %v", key, flags))

		rc, err := q.Dir.Open(key)
		if err != nil {
			log.Fatal(fmt.Sprintf("Error opening file: %v\n", err))
			return err
		}
		defer rc.Close()

		var fileContent string
		buf := make([]byte, 256)
		for {
			n, err := rc.Read(buf)
			if err != nil && err != io.EOF {
				log.Fatal(fmt.Sprintf("Error reading file: %v\n", err))
				break
			}
			if n == 0 {
				break
			}
			fileContent += string(buf[:n])
		}

		log.Info(fmt.Sprintf("File content:\n%s\n", fileContent))

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
