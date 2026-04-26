package filequeue

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// WritePid writes the current process ID to the specified file atomically.
// Returns an error if the file already exists (process already running).
func WritePid(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create pidfile directory: %w", err)
	}

	pid := os.Getpid()
	pidStr := strconv.Itoa(pid)

	// Try to create the file exclusively (fail if already exists)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			// File exists, check if the process is still running
			existingPid, readErr := ReadPid(path)
			if readErr != nil {
				return fmt.Errorf("pidfile exists but cannot be read: %w", readErr)
			}
			if IsRunning(existingPid) {
				return fmt.Errorf("process already running with PID %d", existingPid)
			}
			// Process is dead (stale pidfile), remove it and retry
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove stale pidfile: %w", err)
			}
			return WritePid(path)
		}
		return fmt.Errorf("failed to create pidfile: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(pidStr); err != nil {
		os.Remove(path)
		return fmt.Errorf("failed to write PID to file: %w", err)
	}

	return nil
}

// ReadPid reads the process ID from the specified file.
func ReadPid(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read pidfile: %w", err)
	}

	pidStr := strings.TrimSpace(string(data))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, fmt.Errorf("invalid PID in pidfile: %w", err)
	}

	return pid, nil
}

// RemovePid removes the PID file.
func RemovePid(path string) error {
	return os.Remove(path)
}

// IsRunning checks if a process with the given PID is still running.
func IsRunning(pid int) bool {
	// Send signal 0 to check if process exists
	err := syscall.Kill(pid, 0)
	return err == nil
}
