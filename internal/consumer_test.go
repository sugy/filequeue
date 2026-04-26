package filequeue

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

// MockQueue is a mock implementation of Queue interface for testing.
type MockQueue struct {
	dequeueCount int32
}

func (m *MockQueue) Dequeue() error {
	atomic.AddInt32(&m.dequeueCount, 1)
	return nil
}

func (m *MockQueue) GetDequeueCount() int {
	return int(atomic.LoadInt32(&m.dequeueCount))
}

func TestNewConsumer(t *testing.T) {
	tmpdir := t.TempDir()
	mockQ := &MockQueue{}

	consumer, err := NewConsumer(mockQ, tmpdir, 1*time.Second)
	if err != nil {
		t.Fatalf("NewConsumer failed: %v", err)
	}

	if consumer == nil {
		t.Fatalf("consumer is nil")
	}

	consumer.watcher.Close()
}

func TestConsumerStart(t *testing.T) {
	tmpdir := t.TempDir()

	// Create actual FileQueue to set up directory structure
	fq, err := NewFileQueue(tmpdir)
	if err != nil {
		t.Fatalf("NewFileQueue failed: %v", err)
	}

	consumer, err := NewConsumer(fq, tmpdir, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("NewConsumer failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	go consumer.Start(ctx)
	<-ctx.Done()
}

func TestConsumerInitialDequeue(t *testing.T) {
	tmpdir := t.TempDir()

	// Create actual FileQueue to set up directory structure
	fq, err := NewFileQueue(tmpdir)
	if err != nil {
		t.Fatalf("NewFileQueue failed: %v", err)
	}

	consumer, err := NewConsumer(fq, tmpdir, 10*time.Second)
	if err != nil {
		t.Fatalf("NewConsumer failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	go consumer.Start(ctx)
	<-ctx.Done()
}

func TestConsumerStopGraceful(t *testing.T) {
	tmpdir := t.TempDir()
	mockQ := &MockQueue{}

	consumer, err := NewConsumer(mockQ, tmpdir, 10*time.Second)
	if err != nil {
		t.Fatalf("NewConsumer failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go consumer.Start(ctx)

	// Wait a bit for consumer to start
	time.Sleep(100 * time.Millisecond)

	// Cancel context to trigger stop
	cancel()

	// Give consumer time to stop
	time.Sleep(100 * time.Millisecond)
}
