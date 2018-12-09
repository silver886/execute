package execute

import (
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"leoliu.io/file"
)

// Cmd add some attributes on exec.Cmd
type Cmd struct {
	*exec.Cmd

	buildOut strings.Builder
	buildErr strings.Builder

	stroutIndex []int
	strerrIndex []int
}

// WorkDir get the working directory
func (cmd *Cmd) WorkDir() (workDir string) {
	workDir, _ = filepath.Abs(cmd.Dir)

	return
}

// ExitCode get the exit code
func (cmd *Cmd) ExitCode() (exitCode int) {
	if cmd.ProcessState == nil {
		return
	}
	exitCode = cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()

	return
}

// Strout get the stdout as string
func (cmd *Cmd) Strout() (strout string) {
	strout = strings.TrimSpace(cmd.buildOut.String())

	return
}

// StroutNext get unused stdout
func (cmd *Cmd) StroutNext() (strout string) {
	strout = cmd.Strout()
	history := len(cmd.stroutIndex)
	var pointer int
	if history == 0 {
		pointer = 0
	} else {
		pointer = cmd.stroutIndex[history-1]
	}
	cmd.stroutIndex = append(cmd.stroutIndex, len(strout))
	strout = strout[pointer:]

	return
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
func (cmd *Cmd) Strerr() (strerr string) {
	strerr = strings.TrimSpace(cmd.buildErr.String())

	return
}

// StrerrNext get unused stderr
func (cmd *Cmd) StrerrNext() (strerr string) {
	strerr = cmd.Strerr()
	history := len(cmd.strerrIndex)
	var pointer int
	if history == 0 {
		pointer = 0
	} else {
		pointer = cmd.strerrIndex[history-1]
	}
	cmd.strerrIndex = append(cmd.strerrIndex, len(strerr))
	strerr = strerr[pointer:]

	return
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

// Start starts the specified command but does not wait for it to complete
func (cmd *Cmd) Start() (err error) {

	err = cmd.Cmd.Start()

	return
}

// Run starts the specified command and waits for it to complete
func (cmd *Cmd) Run() (err error) {

	err = cmd.Start()
	if err != nil {
		return
	}

	err = cmd.Wait()

	return
}

// RunToFile run and store output to file
func (cmd *Cmd) RunToFile(path string) (err error) {

	err = cmd.Run()

	if _, err := file.Writeln(path, cmd.Strout()); err != nil {
	}

	if strerr := cmd.Strerr(); strerr != "" {
		if _, err := file.Writeln(path+".err", strerr); err != nil {
		}
	}

	return
}

// New create a execute.Cmd
func New(name string, arg ...string) (cmd *Cmd) {

	execCmd := exec.Command(name, arg...)

	cmd = &Cmd{
		Cmd: execCmd,
	}

	cmd.Stdout = &cmd.buildOut
	cmd.Stderr = &cmd.buildErr

	return
}
