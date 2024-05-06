// Package cmd implements CLI applications.
/*
Copyright © 2024 sugy <sugy.kz@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func getDefaultQueueDirPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "user homedir:", err)
		return filepath.Join(os.TempDir(), "filequeue")
	}
	return filepath.Join(home, "filequeue")
}
