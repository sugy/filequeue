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
	Path   string
	Args   []string
	Output string
}

// run executes the command and returns err
func (c *command) run() error {
	log.Debug(fmt.Sprintf("command: %v %v", c.Path, strings.Join(c.Args, " ")))
	cmd := exec.Command(c.Path, c.Args...)
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

	c.Output = stdout.String()
	return err
}
