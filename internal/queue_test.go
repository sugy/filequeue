package filequeue

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNewFileQueue(t *testing.T) {
	tests := []struct {
		name string
		dir  string
	}{
		{name: "Create New Queue", dir: "filequeue_test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := os.MkdirTemp("", tt.dir)
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)

			if _, err := NewFileQueue(dir); err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if _, err := os.Stat(filepath.Join(dir, "new")); os.IsNotExist(err) {
				t.Fatalf("Expected 'new' directory to be created")
			}
		})
	}
}

func TestEnqueue(t *testing.T) {
	tests := []struct {
		name    string
		kind    string
		message string
		wantErr bool
	}{
		{name: "Valid Exec Message", kind: "exec", message: "echo hello", wantErr: false},
		{name: "Valid Clipboard Message", kind: "clipboard", message: "copy this", wantErr: false},
		{name: "Invalid Kind", kind: "invalid", message: "invalid kind", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := os.MkdirTemp("", "filequeue_test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)

			fq, err := NewFileQueue(dir)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			err = fq.Enqueue(tt.kind, tt.message)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Expected error: %v, got: %v", tt.wantErr, err)
			}

			if !tt.wantErr {
				keys, err := fq.Dir.Unseen()
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}

				if len(keys) == 0 {
					t.Fatalf("Expected at least one message in the queue")
				}

				rc, err := fq.Dir.Open(keys[0])
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}
				defer rc.Close()

				var q queue
				fileContent, err := io.ReadAll(rc)
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}

				if err := yaml.Unmarshal(fileContent, &q); err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}

				if q.Payload.Kind != tt.kind {
					t.Fatalf("Expected '%v' kind, got %v", tt.kind, q.Payload.Kind)
				}
			}
		})
	}
}

func TestDequeue(t *testing.T) {
	tests := []struct {
		name    string
		kind    string
		message string
		wantErr bool
	}{
		{name: "Dequeue Exec Message", kind: "exec", message: "echo hello", wantErr: false},
		{name: "Dequeue Clipboard Message", kind: "clipboard", message: "copy this", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := os.MkdirTemp("", "filequeue_test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)

			fq, err := NewFileQueue(dir)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			err = fq.Enqueue(tt.kind, tt.message)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			err = fq.Dequeue()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Expected error: %v, got: %v", tt.wantErr, err)
			}

			keys, err := fq.Dir.Unseen()
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if len(keys) != 0 {
				t.Fatalf("Expected no messages in the queue")
			}
		})
	}
}

func TestPurge(t *testing.T) {
	tests := []struct {
		name    string
		kind    string
		message string
		wantErr bool
	}{
		{name: "Purge Exec Message", kind: "exec", message: "echo hello", wantErr: false},
		{name: "Purge Clipboard Message", kind: "clipboard", message: "copy this", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := os.MkdirTemp("", "filequeue_test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)

			fq, err := NewFileQueue(dir)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			err = fq.Enqueue(tt.kind, tt.message)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			err = fq.Purge()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Expected error: %v, got: %v", tt.wantErr, err)
			}

			keys, err := fq.Dir.Unseen()
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if len(keys) != 0 {
				t.Error(keys)
				t.Fatalf("Expected no messages in the queue")
			}
		})
	}
}
