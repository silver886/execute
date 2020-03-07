package execute

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/silver886/file"
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

// RunToFile run and store output to file
func (cmd *Cmd) RunToFile(path string) error {
	cmd.Run()

	if _, err := file.Writeln(path, cmd.OutString()); err != nil {
		return err
	}

	if errString := cmd.ErrString(); errString != "" {
		if _, err := file.Writeln(path+".err", errString); err != nil {
			return err
		}
	}

	return nil
}

// Wait waits for the command to exit
func (cmd *Cmd) Wait() error {
	cmd.wait.mutex.Lock()

	if !cmd.wait.done {
		cmd.wait.err = cmd.Cmd.Wait()
		cmd.wait.done = true
	}

	cmd.wait.mutex.Unlock()

	return cmd.wait.err
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
