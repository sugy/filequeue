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
	"slices"
	"strings"

	maildir "github.com/emersion/go-maildir"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Queue struct is...
type Queue struct {
	Dir     maildir.Dir `yaml:"dir"`
	Payload Payload     `yaml:"payload"`
}

type Payload struct {
	Massage string `yaml:"message"`
	Type    string `yaml:"type"`
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
	q.Payload.Type, q.Payload.Massage = t, m
	log.Debug(fmt.Sprintf("Queue: %v", q))

	yamlBytes, err := yaml.Marshal(q)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error marshaling to YAML: %v", err))
	}

	d, err := maildir.NewDelivery(string(q.Dir))
	if err != nil {
		log.Fatal(err)
	}
	b, err := d.Write(yamlBytes)
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
	news, err := q.Dir.Unseen()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(fmt.Sprintf("new keys: %v", news))

	err = q.Dir.Walk(func(key string, flags []maildir.Flag) error {
		log.Debug(fmt.Sprintf("%v, %v", key, flags))

		if !slices.Contains(news, key) {
			return nil
		}
		log.Info(fmt.Sprintf("new key: %v, %v", key, flags))

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

		var tq Queue
		if err := yaml.Unmarshal([]byte(fileContent), &tq); err != nil {
			log.Fatal(fmt.Sprintf("Error Unmarshaling from YAML: %v\n", err))
		}
		log.Info(fmt.Sprintf("%v", tq))

		cmdStr := strings.Fields(os.ExpandEnv(tq.Payload.Massage))
		c := newCommand(cmdStr[0], cmdStr[1:])
		err = c.run()
		if err != nil {
			log.Fatal(err)
			return err
		}
		log.Info(fmt.Sprintf("execute command. exitCode: %d, stdout: '%s', stderr: '%s'\n", c.exitCode, c.stdout, c.stderr))

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
