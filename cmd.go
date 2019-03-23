package execute

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
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

// OutString get the stdout as string
func (cmd *Cmd) OutString() string {
	return strings.TrimSpace(cmd.OutBuffer.String())
}

// OutStringNext get unused stdout
func (cmd *Cmd) OutStringNext() string {
	outString := cmd.OutString()
	history := len(cmd.outStringIndex)
	var pointer int
	if history == 0 {
		pointer = 0
	} else {
		pointer = cmd.outStringIndex[history-1]
	}
	cmd.outStringIndex = append(cmd.outStringIndex, len(outString))
	return outString[pointer:]
}

// OutStringHistory get stdout history access by OutStringNext
func (cmd *Cmd) OutStringHistory() (outString []string) {
	outStringFull := cmd.OutString()
	for i, val := range cmd.outStringIndex {
		if i == 0 {
			outString = append(outString, outStringFull[:val])
		} else {
			outString = append(outString, outStringFull[cmd.outStringIndex[i-1]:val])
		}
	}

	return
}

// ErrString get the stderr as string
func (cmd *Cmd) ErrString() string {
	return strings.TrimSpace(cmd.ErrBuffer.String())
}

// ErrStringNext get unused stderr
func (cmd *Cmd) ErrStringNext() string {
	errString := cmd.ErrString()
	history := len(cmd.errStringIndex)
	var pointer int
	if history == 0 {
		pointer = 0
	} else {
		pointer = cmd.errStringIndex[history-1]
	}
	cmd.errStringIndex = append(cmd.errStringIndex, len(errString))
	return errString[pointer:]
}

// ErrStringHistory get stderr history access by ErrStringNext
func (cmd *Cmd) ErrStringHistory() (errString []string) {
	errStringFull := cmd.ErrString()
	for i, val := range cmd.errStringIndex {
		if i == 0 {
			errString = append(errString, errStringFull[:val])
		} else {
			errString = append(errString, errStringFull[cmd.errStringIndex[i-1]:val])
		}
	}

	return
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
