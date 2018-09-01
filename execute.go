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
func (cmd *Cmd) RunToLog(logger *logrus.Logger, level logrus.Level, msg string) (err error) {
	if logger == nil {
		err = errors.New("Invalid logger")
		return
	}

	err = cmd.Run()

	execLogger := logger.WithFields(logrus.Fields{
		"path":              cmd.Path,
		"args":              strings.Join(cmd.Args, " |: "),
		"working directory": cmd.WorkDir,
		"exitcode":          cmd.ExitCode(),
		"stdout":            cmd.Strout(),
		"stderr":            cmd.Strerr(),
		"internal error":    err,
	})

	levelLog(execLogger, level, msg)

	return
}

// StartToLog start with log
func (cmd *Cmd) StartToLog(logger *logrus.Logger, level logrus.Level, msg string) (err error) {
	if logger == nil {
		err = errors.New("Invalid logger")
		return
	}

	err = cmd.Start()

	execLogger := logger.WithFields(logrus.Fields{
		"path":              cmd.Path,
		"args":              strings.Join(cmd.Args, " |: "),
		"working directory": cmd.WorkDir,
		"internal error":    err,
	})

	levelLog(execLogger, level, msg)

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
func (cmd *Cmd) RunToFileLog(path string, logger *logrus.Logger, level logrus.Level, msg string) (err error) {
	err = cmd.RunToLog(logger, level, msg)

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
func RunToLog(logger *logrus.Logger, level logrus.Level, msg string, name string, arg ...string) (cmd *Cmd, err error) {
	if logger == nil {
		err = errors.New("Invalid logger")
		return
	}

	cmd, err = Run(name, arg...)

	execLogger := logger.WithFields(logrus.Fields{
		"path":              cmd.Path,
		"args":              strings.Join(arg, " |: "),
		"working directory": cmd.WorkDir,
		"exitcode":          cmd.ExitCode(),
		"stdout":            cmd.Strout(),
		"stderr":            cmd.Strerr(),
		"internal error":    err,
	})

	levelLog(execLogger, level, msg)

	return
}

// StartToLog start with log
func StartToLog(logger *logrus.Logger, level logrus.Level, msg string, name string, arg ...string) (cmd *Cmd, err error) {
	if logger == nil {
		err = errors.New("Invalid logger")
		return
	}

	cmd, err = Start(name, arg...)

	execLogger := logger.WithFields(logrus.Fields{
		"path":              cmd.Path,
		"args":              strings.Join(arg, " |: "),
		"working directory": cmd.WorkDir,
		"internal error":    err,
	})

	levelLog(execLogger, level, msg)

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
func RunToFileLog(path string, logger *logrus.Logger, level logrus.Level, msg string, name string, arg ...string) (cmd *Cmd, err error) {
	cmd, err = RunToLog(logger, level, msg, name, arg...)

	file.Writeln(path, cmd.Strout())

	if strerr := cmd.Strerr(); strerr != "" {
		file.Writeln(path+".err", strerr)
	}

	return
}

func levelLog(entry *logrus.Entry, level logrus.Level, msg string) {
	switch level {
	case logrus.DebugLevel:
		entry.Debugln(msg)
	case logrus.InfoLevel:
		entry.Infoln(msg)
	case logrus.WarnLevel:
		entry.Warnln(msg)
	case logrus.ErrorLevel:
		entry.Errorln(msg)
	case logrus.FatalLevel:
		entry.Fatalln(msg)
	case logrus.PanicLevel:
		entry.Panicln(msg)
	default:
		entry.Debugln(msg)
	}
}
