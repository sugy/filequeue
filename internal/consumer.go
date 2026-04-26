package filequeue

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

// Consumer watches a queue directory and dequeues items as they arrive.
type Consumer struct {
	queue    Queue
	queueDir string
	watcher  *fsnotify.Watcher
	interval time.Duration
	wg       sync.WaitGroup
	done     chan struct{}
}

// NewConsumer creates a new Consumer instance.
func NewConsumer(q Queue, queueDir string, interval time.Duration) (*Consumer, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create fsnotify watcher: %w", err)
	}

	return &Consumer{
		queue:    q,
		queueDir: queueDir,
		watcher:  watcher,
		interval: interval,
		done:     make(chan struct{}),
	}, nil
}

// Start begins watching the queue directory and processing items.
// It blocks until ctx is canceled or an error occurs.
func (c *Consumer) Start(ctx context.Context) error {
	newDir := filepath.Join(c.queueDir, "new")
	if err := c.watcher.Add(newDir); err != nil {
		return fmt.Errorf("failed to watch directory: %w", err)
	}
	log.Info(fmt.Sprintf("watching directory: %s", newDir))

	// Initial dequeue in case items are already queued
	c.dequeueWithRecovery()

	go c.watchLoop(ctx)
	<-ctx.Done()

	return c.Stop()
}

// Stop gracefully shuts down the consumer.
func (c *Consumer) Stop() error {
	close(c.done)
	c.watcher.Close()

	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Info("consumer stopped gracefully")
	case <-time.After(30 * time.Second):
		log.Warn("consumer shutdown timeout (30s), force exiting")
	}

	return nil
}

// watchLoop monitors the queue directory for new items.
func (c *Consumer) watchLoop(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case event, ok := <-c.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				log.Debug(fmt.Sprintf("file created: %s", event.Name))
				c.dequeueWithRecovery()
			}

		case err, ok := <-c.watcher.Errors:
			if !ok {
				return
			}
			log.Error(fmt.Sprintf("watcher error: %v", err))

		case <-ticker.C:
			c.dequeueWithRecovery()

		case <-c.done:
			return

		case <-ctx.Done():
			return
		}
	}
}

// dequeueWithRecovery calls Dequeue but recovers from panics (e.g., log.Fatal).
func (c *Consumer) dequeueWithRecovery() {
	c.wg.Add(1)
	defer c.wg.Done()

	defer func() {
		if r := recover(); r != nil {
			log.Error(fmt.Sprintf("dequeue panic: %v", r))
		}
	}()

	if err := c.queue.Dequeue(); err != nil {
		log.Error(fmt.Sprintf("dequeue error: %v", err))
	}
}
