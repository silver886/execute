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
}

// WorkDir get the working directory
func (cmd *Cmd) WorkDir() (workDir string) {
	workDir, _ = filepath.Abs(cmd.Dir)
	return
}

// ExitCode get the exit code
func (cmd *Cmd) ExitCode() (exitCode int) {
	exitCode = cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	return
}

// Strout get the stdout as string
func (cmd *Cmd) Strout() (strout string) {
	strout = strings.TrimSpace(cmd.bufout.String())
	return
}

// Strerr get the stderr as string
func (cmd *Cmd) Strerr() (strerr string) {
	strerr = strings.TrimSpace(cmd.buferr.String())
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
		"working directory": cmd.WorkDir(),
		"exitcode":          cmd.ExitCode(),
		"stdout":            cmd.Strout(),
		"stderr":            cmd.Strerr(),
		"internal error":    err,
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
		"working directory": cmd.WorkDir(),
		"internal error":    err,
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

	execCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	cmd = &Cmd{
		Cmd: execCmd,
	}

	cmd.Stdout = &cmd.bufout
	cmd.Stderr = &cmd.buferr

	return cmd
}

// Run starts the specified command and waits for it to complete
func Run(name string, arg ...string) (cmd *Cmd, err error) {
	cmd = New(name, arg...)
	err = cmd.Run()
	return
}

// Start starts the specified command but does not wait for it to complete
//
// The Wait method will return the exit code and release associated resources
// once the command exits
func Start(name string, arg ...string) (cmd *Cmd, err error) {
	cmd = New(name, arg...)
	err = cmd.Start()
	return cmd, err
}

// RunToLog run with log
func RunToLog(extLogger *logrus.Logger, level logrus.Level, msg string, name string, arg ...string) (cmd *Cmd, err error) {
	if extLogger == nil {
		err = errors.New("Invalid logger")
		return
	}

	cmd, err = Run(name, arg...)

	execLogger := extLogger.WithFields(logrus.Fields{
		"path":              cmd.Path,
		"args":              strings.Join(arg, " |: "),
		"working directory": cmd.WorkDir(),
		"exitcode":          cmd.ExitCode(),
		"stdout":            cmd.Strout(),
		"stderr":            cmd.Strerr(),
		"internal error":    err,
	})

	logger.LevelLog(execLogger, level, msg)

	return
}

// StartToLog start with log
func StartToLog(extLogger *logrus.Logger, level logrus.Level, msg string, name string, arg ...string) (cmd *Cmd, err error) {
	if extLogger == nil {
		err = errors.New("Invalid logger")
		return
	}

	cmd, err = Start(name, arg...)

	execLogger := extLogger.WithFields(logrus.Fields{
		"path":              cmd.Path,
		"args":              strings.Join(arg, " |: "),
		"working directory": cmd.WorkDir(),
		"internal error":    err,
	})

	logger.LevelLog(execLogger, level, msg)

	return
}

// RunToFile run and store output to file
func RunToFile(path string, name string, arg ...string) (cmd *Cmd, err error) {
	cmd, err = Run(name, arg...)

	file.Writeln(path, cmd.Strout())

	if strerr := cmd.Strerr(); strerr != "" {
		file.Writeln(path+".err", strerr)
	}

	return
}

// RunToFileLog run with log and store output to file
func RunToFileLog(path string, extLogger *logrus.Logger, level logrus.Level, msg string, name string, arg ...string) (cmd *Cmd, err error) {
	cmd, err = RunToLog(extLogger, level, msg, name, arg...)

	file.Writeln(path, cmd.Strout())

	if strerr := cmd.Strerr(); strerr != "" {
		file.Writeln(path+".err", strerr)
	}

	return
}
