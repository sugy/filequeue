package filequeue

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWritePid(t *testing.T) {
	tmpdir := t.TempDir()
	pidFile := filepath.Join(tmpdir, "test.pid")

	err := WritePid(pidFile)
	if err != nil {
		t.Fatalf("WritePid failed: %v", err)
	}

	if _, err := os.Stat(pidFile); err != nil {
		t.Fatalf("PID file was not created: %v", err)
	}

	// Clean up for next test
	os.Remove(pidFile)
}

func TestReadPid(t *testing.T) {
	tmpdir := t.TempDir()
	pidFile := filepath.Join(tmpdir, "test.pid")

	err := WritePid(pidFile)
	if err != nil {
		t.Fatalf("WritePid failed: %v", err)
	}

	pid, err := ReadPid(pidFile)
	if err != nil {
		t.Fatalf("ReadPid failed: %v", err)
	}

	if pid != os.Getpid() {
		t.Errorf("Read PID %d, expected %d", pid, os.Getpid())
	}

	os.Remove(pidFile)
}

func TestDuplicateWritePid(t *testing.T) {
	tmpdir := t.TempDir()
	pidFile := filepath.Join(tmpdir, "test.pid")

	err := WritePid(pidFile)
	if err != nil {
		t.Fatalf("First WritePid failed: %v", err)
	}

	// Try to write again with same process (should fail)
	err = WritePid(pidFile)
	if err == nil {
		t.Fatalf("Second WritePid should have failed but didn't")
	}

	if err.Error() != "process already running with PID "+os.ExpandEnv("$") {
		// Just check that error contains "already running"
		if !containsString(err.Error(), "already running") {
			t.Errorf("Expected 'already running' error, got: %v", err)
		}
	}

	os.Remove(pidFile)
}

func TestRemovePid(t *testing.T) {
	tmpdir := t.TempDir()
	pidFile := filepath.Join(tmpdir, "test.pid")

	err := WritePid(pidFile)
	if err != nil {
		t.Fatalf("WritePid failed: %v", err)
	}

	err = RemovePid(pidFile)
	if err != nil {
		t.Fatalf("RemovePid failed: %v", err)
	}

	if _, err := os.Stat(pidFile); err == nil {
		t.Fatalf("PID file was not removed")
	}
}

func TestIsRunning(t *testing.T) {
	currentPid := os.Getpid()

	if !IsRunning(currentPid) {
		t.Errorf("IsRunning failed for current process")
	}

	// Test with invalid PID
	invalidPid := 999999
	if IsRunning(invalidPid) {
		t.Logf("IsRunning returned true for invalid PID (may be valid on this system)")
	}
}

func containsString(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
