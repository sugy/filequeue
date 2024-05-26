// Package filequeue ...
package filequeue

import "os/exec"

// CommandExecutor is an interface that wraps the exec.Command function.
type CommandExecutor interface {
	Command(name string, arg ...string) *exec.Cmd
}

type realCommandExecutor struct{}

func (r *realCommandExecutor) Command(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}
