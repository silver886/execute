package execute

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
)

// Cmd add some attributes on exec.Cmd
type Cmd struct {
	*exec.Cmd

	wait struct {
		mutex sync.Mutex
		done  bool
		err   error
	}

	OutBuffer bytes.Buffer
	ErrBuffer bytes.Buffer

	outStringIndex []int
	errStringIndex []int
}

// New create a execute.Cmd
func New(name string, arg ...string) (cmd *Cmd) {
	cmd = &Cmd{
		Cmd: exec.Command(name, arg...),
	}
	cmd.Stdout = &cmd.OutBuffer
	cmd.Stderr = &cmd.ErrBuffer

	return
}

// WorkDir get the working directory
func (cmd *Cmd) WorkDir() string {
	workDir, _ := filepath.Abs(cmd.Dir)
	return workDir
}

// ExitCode get the exit code
func (cmd *Cmd) ExitCode() int {
	if cmd.ProcessState == nil {
		return -1
	}
	return cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
}
