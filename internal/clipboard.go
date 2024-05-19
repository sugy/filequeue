// Package filequeue ...
/*
Copyright © 2024 sugy <sugy.kz@gmail.com>
*/
package filequeue

import (
	"fmt"
	"runtime"

	log "github.com/sirupsen/logrus"
)

// clipboard struct ...
type clipboard struct {
	cmd string
}

// newClipboard ...
func newClipboard(cmd string) *clipboard {

	if len(cmd) == 0 {
		switch runtime.GOOS {
		case "darwin":
			cmd = "pbcopy"
		case "linux":
			cmd = "cat"
		}
	}

	return &clipboard{
		cmd: cmd,
	}
}

// copy ...
func (c *clipboard) copy(txt []byte) error {
	exec := newExecute(c.cmd, []string{})
	exec.stdin = string(txt)
	err := exec.run()
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Info(fmt.Sprintf("execute command. exitCode: %d, stdout: '%s', stderr: '%s'\n",
		exec.exitCode, exec.stdout, exec.stderr))

	return nil
}
