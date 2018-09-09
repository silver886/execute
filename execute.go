// Package execute provides string and file output for os.exec
package execute

import (
	"bytes"
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
	"leoliu.io/file"
	"leoliu.io/logger"
)

// Cmd add some attributes on exec.Cmd
type Cmd struct {
	*exec.Cmd

	bufout bytes.Buffer
	buferr bytes.Buffer

	stroutIndex []int
	strerrIndex []int
}

// WorkDir get the working_directory
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
	strout = strings.TrimSpace(cmd.bufout.String())
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
	strerr = strings.TrimSpace(cmd.buferr.String())
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

// RunToLog run with log
func (cmd *Cmd) RunToLog(extLogger *logrus.Logger, level logrus.Level, msg string) (err error) {
	if extLogger == nil {
		err = errors.New("Invalid logger")
		return
	}

	err = cmd.Run()

	execLogger := extLogger.WithFields(logrus.Fields{
		"path":              cmd.Path,
		"args":              strings.Join(cmd.Args, " |: "),
		"working_directory": cmd.WorkDir(),
		"exitcode":          cmd.ExitCode(),
		"stdout":            cmd.Strout(),
		"stderr":            cmd.Strerr(),
		"internal_error":    err,
	})

	logger.LevelLog(execLogger, level, msg)

	return
}

// StartToLog start with log
func (cmd *Cmd) StartToLog(extLogger *logrus.Logger, level logrus.Level, msg string) (err error) {
	if extLogger == nil {
		err = errors.New("Invalid logger")
		return
	}

	err = cmd.Start()

	execLogger := extLogger.WithFields(logrus.Fields{
		"path":              cmd.Path,
		"args":              strings.Join(cmd.Args, " |: "),
		"working_directory": cmd.WorkDir(),
		"internal_error":    err,
	})

	logger.LevelLog(execLogger, level, msg)

	return
}

// RunToFile run and store output to file
func (cmd *Cmd) RunToFile(path string) (err error) {
	err = cmd.Run()

	file.Writeln(path, cmd.Strout())

	if strerr := cmd.Strerr(); strerr != "" {
		file.Writeln(path+".err", strerr)
	}

	return
}

// RunToFileLog run with log and store output to file
func (cmd *Cmd) RunToFileLog(path string, extLogger *logrus.Logger, level logrus.Level, msg string) (err error) {
	err = cmd.RunToLog(extLogger, level, msg)

	file.Writeln(path, cmd.Strout())

	if strerr := cmd.Strerr(); strerr != "" {
		file.Writeln(path+".err", strerr)
	}

	return
}

// New create a execute.Cmd
func New(name string, arg ...string) (cmd *Cmd) {
	execCmd := exec.Command(name, arg...)

	cmd = &Cmd{
		Cmd: execCmd,
	}

	cmd.Stdout = &cmd.bufout
	cmd.Stderr = &cmd.buferr

	return
}

// Run starts the specified command and waits for it to complete
func Run(hide bool, name string, arg ...string) (cmd *Cmd, err error) {
	cmd = New(name, arg...)
	if hide {
		cmd.Hide()
	}
	err = cmd.Run()
	return
}

// Start starts the specified command but does not wait for it to complete
//
// The Wait method will return the exit code and release associated resources
// once the command exits
func Start(hide bool, name string, arg ...string) (cmd *Cmd, err error) {
	cmd = New(name, arg...)
	if hide {
		cmd.Hide()
	}
	err = cmd.Start()
	return
}

// RunToLog run with log
func RunToLog(hide bool, extLogger *logrus.Logger, level logrus.Level, msg string, name string, arg ...string) (cmd *Cmd, err error) {
	if extLogger == nil {
		err = errors.New("Invalid logger")
		return
	}

	cmd = New(name, arg...)
	if hide {
		cmd.Hide()
	}
	err = cmd.Run()

	execLogger := extLogger.WithFields(logrus.Fields{
		"path":              cmd.Path,
		"args":              strings.Join(arg, " |: "),
		"working_directory": cmd.WorkDir(),
		"exitcode":          cmd.ExitCode(),
		"stdout":            cmd.Strout(),
		"stderr":            cmd.Strerr(),
		"internal_error":    err,
	})

	logger.LevelLog(execLogger, level, msg)

	return
}

// StartToLog start with log
func StartToLog(hide bool, extLogger *logrus.Logger, level logrus.Level, msg string, name string, arg ...string) (cmd *Cmd, err error) {
	if extLogger == nil {
		err = errors.New("Invalid logger")
		return
	}

	cmd = New(name, arg...)
	if hide {
		cmd.Hide()
	}
	err = cmd.Run()

	execLogger := extLogger.WithFields(logrus.Fields{
		"path":              cmd.Path,
		"args":              strings.Join(arg, " |: "),
		"working_directory": cmd.WorkDir(),
		"internal_error":    err,
	})

	logger.LevelLog(execLogger, level, msg)

	return
}

// RunToFile run and store output to file
func RunToFile(hide bool, path string, name string, arg ...string) (cmd *Cmd, err error) {
	cmd = New(name, arg...)
	if hide {
		cmd.Hide()
	}
	err = cmd.Run()

	file.Writeln(path, cmd.Strout())

	if strerr := cmd.Strerr(); strerr != "" {
		file.Writeln(path+".err", strerr)
	}

	return
}

// RunToFileLog run with log and store output to file
func RunToFileLog(hide bool, path string, extLogger *logrus.Logger, level logrus.Level, msg string, name string, arg ...string) (cmd *Cmd, err error) {
	cmd, err = RunToLog(hide, extLogger, level, msg, name, arg...)

	file.Writeln(path, cmd.Strout())

	if strerr := cmd.Strerr(); strerr != "" {
		file.Writeln(path+".err", strerr)
	}

	return
}
