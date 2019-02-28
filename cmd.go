package execute

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/silver886/file"
)

// Cmd add some attributes on exec.Cmd
type Cmd struct {
	*exec.Cmd

	Out bytes.Buffer
	Err bytes.Buffer

	stroutIndex []int
	strerrIndex []int
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

// Strout get the stdout as string
func (cmd *Cmd) Strout() string {
	return strings.TrimSpace(cmd.Out.String())
}

// StroutNext get unused stdout
func (cmd *Cmd) StroutNext() string {
	strout := cmd.Strout()
	history := len(cmd.stroutIndex)
	var pointer int
	if history == 0 {
		pointer = 0
	} else {
		pointer = cmd.stroutIndex[history-1]
	}
	cmd.stroutIndex = append(cmd.stroutIndex, len(strout))
	return strout[pointer:]
}

// StroutHistory get stdout history access by StroutNext
func (cmd *Cmd) StroutHistory() (strout []string) {
	stroutFull := cmd.Strout()
	for i, val := range cmd.stroutIndex {
		if i == 0 {
			strout = append(strout, stroutFull[:val])
		} else {
			strout = append(strout, stroutFull[cmd.stroutIndex[i-1]:val])
		}
	}

	return
}

// Strerr get the stderr as string
func (cmd *Cmd) Strerr() string {
	return strings.TrimSpace(cmd.Err.String())
}

// StrerrNext get unused stderr
func (cmd *Cmd) StrerrNext() string {
	strerr := cmd.Strerr()
	history := len(cmd.strerrIndex)
	var pointer int
	if history == 0 {
		pointer = 0
	} else {
		pointer = cmd.strerrIndex[history-1]
	}
	cmd.strerrIndex = append(cmd.strerrIndex, len(strerr))
	return strerr[pointer:]
}

// StrerrHistory get stderr history access by StrerrNext
func (cmd *Cmd) StrerrHistory() (strerr []string) {
	strerrFull := cmd.Strerr()
	for i, val := range cmd.strerrIndex {
		if i == 0 {
			strerr = append(strerr, strerrFull[:val])
		} else {
			strerr = append(strerr, strerrFull[cmd.strerrIndex[i-1]:val])
		}
	}

	return
}

// RunToFile run and store output to file
func (cmd *Cmd) RunToFile(path string) error {
	cmd.Run()

	if _, err := file.Writeln(path, cmd.Strout()); err != nil {
		return err
	}

	if strerr := cmd.Strerr(); strerr != "" {
		if _, err := file.Writeln(path+".err", strerr); err != nil {
			return err
		}
	}

	return nil
}

// New create a execute.Cmd
func New(name string, arg ...string) (cmd *Cmd) {
	cmd = &Cmd{
		Cmd: exec.Command(name, arg...),
	}
	cmd.Stdout = &cmd.Out
	cmd.Stderr = &cmd.Err

	return
}
