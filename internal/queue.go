// Package filequeue ...
/*
Copyright © 2024 sugy <sugy.kz@gmail.com>
*/
package filequeue

import (
	"encoding/base64"
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

// FileQueue struct is...
type FileQueue struct {
	Dir maildir.Dir
}

// queue struct is...
type queue struct {
	Payload payload `yaml:"payload"`
}

// payload struct is...
type payload struct {
	Massage string `yaml:"message"`
	Kind    string `yaml:"kind"`
}

// NewFileQueue is...
func NewFileQueue(d string) *FileQueue {
	f := &FileQueue{
		Dir: maildir.Dir(d),
	}

	err := f.setupFileQueuedir()
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func (f *FileQueue) setupFileQueuedir() error {
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
func (f *FileQueue) Enqueue(k string, m string) error {
	log.Debug("enqueue!")
	var q queue
	q.Payload.Kind = k
	q.Payload.Massage = base64.StdEncoding.EncodeToString([]byte(m))
	log.Debug(fmt.Sprintf("Queue: %v", q))

	yamlBytes, err := yaml.Marshal(q)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error marshaling to YAML: %v", err))
	}

	d, err := maildir.NewDelivery(string(f.Dir))
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
func (f *FileQueue) Dequeue() error {
	log.Debug("dequeue!")
	news, err := f.Dir.Unseen()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(fmt.Sprintf("new keys: %v", news))

	err = f.Dir.Walk(func(key string, flags []maildir.Flag) error {
		log.Debug(fmt.Sprintf("%v, %v", key, flags))

		if !slices.Contains(news, key) {
			return nil
		}
		log.Info(fmt.Sprintf("new key: %v, %v", key, flags))

		rc, err := f.Dir.Open(key)
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

		var q queue
		if err := yaml.Unmarshal([]byte(fileContent), &q); err != nil {
			log.Fatal(fmt.Sprintf("Error Unmarshaling from YAML: %v\n", err))
		}
		log.Debug(fmt.Sprintf("%v", q))
		msg, err := base64.StdEncoding.DecodeString(q.Payload.Massage)
		if err != nil {
			log.Fatal(fmt.Sprintf("Error: base64 decoding: %v", err))
			return err
		}

		cmdStr := strings.Fields(os.ExpandEnv(string(msg)))
		exec := newExecute(cmdStr[0], cmdStr[1:])
		err = exec.run()
		if err != nil {
			log.Fatal(err)
			return err
		}
		log.Info(fmt.Sprintf("execute command. exitCode: %d, stdout: '%s', stderr: '%s'\n",
			exec.exitCode, exec.stdout, exec.stderr))

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
