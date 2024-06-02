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
	exec *execute
}

// newClipboard ...
func newClipboard(cmd string) *clipboard {
	if len(cmd) == 0 {
		cmd = clipboardCopyCmd(runtime.GOOS)
	}

	return &clipboard{
		exec: newExecute(cmd, []string{}),
	}
}

// clipboardCopyCmd ...
func clipboardCopyCmd(os string) string {
	switch os {
	case "darwin":
		return "pbcopy"
	case "linux":
		return "cat"
	}
	return ""
}

// copy ...
func (c *clipboard) copy(txt []byte) error {
	c.exec.stdin = string(txt)
	err := c.exec.run()
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Info(fmt.Sprintf("execute command. exitCode: %d, stdout: '%s', stderr: '%s'\n",
		c.exec.exitCode, c.exec.stdout, c.exec.stderr))

	return nil
}
