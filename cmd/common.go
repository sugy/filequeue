// Package cmd implements CLI applications.
/*
Copyright © 2024 sugy <sugy.kz@gmail.com>
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func getDefaultQueueDirPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "user homedir:", err)
		path := filepath.Join(os.TempDir(), "filequeue")
		fmt.Fprintln(os.Stderr, "filequeue directory path:", path)
		return path
	}
	return filepath.Join(home, "filequeue")
}

func getDefaultPidFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		path := filepath.Join(os.TempDir(), "filequeue.pid")
		fmt.Fprintln(os.Stderr, "user homedir:", err)
		fmt.Fprintln(os.Stderr, "pidfile path:", path)
		return path
	}
	return filepath.Join(home, ".filequeue", "filequeue.pid")
}

func getStringFromIoReader(r io.Reader) (string, error) {
	scanner := bufio.NewScanner(r)
	var s []string
	for scanner.Scan() {
		s = append(s, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return strings.Join(s, "\n"), nil
}
