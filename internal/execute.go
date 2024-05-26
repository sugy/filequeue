// Package filequeue ...
/*
Copyright © 2024 sugy <sugy.kz@gmail.com>
*/
package filequeue

import (
	"bytes"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// execute struct ...
type execute struct {
	cmdExec  CommandExecutor
	path     string
	args     []string
	exitCode int
	stdout   string
	stderr   string
	stdin    string
}

// newExecute ...
func newExecute(path string, args []string) *execute {
	return &execute{
		cmdExec: &realCommandExecutor{},
		path:    path,
		args:    args,
	}
}

// run executes the command and returns err
func (c *execute) run() error {
	log.Debug(fmt.Sprintf("execute: %v %v", c.path, strings.Join(c.args, " ")))

	cmd := c.cmdExec.Command(c.path, c.args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	var err error

	if len(c.stdin) != 0 {
		var stdin bytes.Buffer
		stdin.WriteString(c.stdin)
		cmd.Stdin = &stdin

		if err = cmd.Start(); err != nil {
			return err
		}
		if err = cmd.Wait(); err != nil {
			return err
		}
	} else {
		err = cmd.Run()
	}

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
