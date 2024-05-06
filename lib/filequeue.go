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

	err := f.setupQueuedir()
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func (f *Filequeue) setupQueuedir() error {
	path := string(f.Dir)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.FileMode(0700))
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	if _, err := os.Stat(filepath.Join(path, "new")); errors.Is(err, os.ErrNotExist) {
		err := f.Dir.Init()
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	return nil
}

// Enqueue is...
func (f *Filequeue) Enqueue(t string, q string) error {
	log.Debug("enqueue!")
	f.Type, f.Queue = t, q
	log.Debug(fmt.Sprintf("Filequeue: %v", f))

	d, err := maildir.NewDelivery(string(f.Dir))
	if err != nil {
		log.Fatal(err)
	}
	b, err := d.Write([]byte(f.Queue))
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
func (f *Filequeue) Dequeue() error {
	log.Debug("dequeue!")
	return nil
}
