// Package filequeue ...
/*
Copyright © 2024 sugy <sugy.kz@gmail.com>
*/
package filequeue

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// command struct ...
type command struct {
	path     string
	args     []string
	exitCode int
	stdout   string
	stderr   string
}

// newCommand ...
func newCommand(path string, args []string) *command {
	return &command{
		path: path,
		args: args,
	}
}

// run executes the command and returns err
func (c *command) run() error {
	log.Debug(fmt.Sprintf("command: %v %v", c.path, strings.Join(c.args, " ")))
	cmd := exec.Command(c.path, c.args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := cmd.ProcessState.ExitCode()

	if err != nil {
		return fmt.Errorf(fmt.Sprintf("failed to execute command. exitCode: %d, Stdout: '%s', Stderr: '%s'\n",
			exitCode, stdout.String(), stderr.String()))
	}

	c.exitCode = exitCode
	c.stdout = stdout.String()
	c.stderr = stderr.String()
	return err
}
