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

// Run starts the specified command and waits for it to complete
func Run(name string, arg ...string) (cmd *exec.Cmd, exitCode int, strout string, strerr string, err error) {
	var bufout, buferr bytes.Buffer

	cmd = execInit(name, arg...)

	cmd.Stdout = &bufout
	cmd.Stderr = &buferr

	err = cmd.Run()

	exitCode = cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()

	strout = strings.TrimSpace(bufout.String())
	strerr = strings.TrimSpace(buferr.String())

	return
}

// Start starts the specified command but does not wait for it to complete
//
// The Wait method will return the exit code and release associated resources
// once the command exits
func Start(name string, arg ...string) (cmd *exec.Cmd, err error) {
	cmd = execInit(name, arg...)

	err = cmd.Start()

	return cmd, err
}

// RunToLog run with log
func RunToLog(logger *logrus.Logger, level logrus.Level, msg string, name string, arg ...string) (cmd *exec.Cmd, exitCode int, strout string, strerr string, err error) {
	if logger == nil {
		err = errors.New("Invalid logger")
		return
	}

	cmd, exitCode, strout, strerr, err = Run(name, arg...)

	workDir, _ := filepath.Abs(cmd.Dir)

	execLogger := logger.WithFields(logrus.Fields{
		"path":              cmd.Path,
		"args":              strings.Join(arg, " |: "),
		"working directory": workDir,
		"exitcode":          cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(),
		"stdout":            strout,
		"stderr":            strerr,
		"internal error":    err,
	})

	levelLog(execLogger, level, msg)

	return
}

// StartToLog start with log
func StartToLog(logger *logrus.Logger, level logrus.Level, msg string, name string, arg ...string) (cmd *exec.Cmd, err error) {
	if logger == nil {
		err = errors.New("Invalid logger")
		return
	}

	cmd, err = Start(name, arg...)

	workDir, _ := filepath.Abs(cmd.Dir)

	execLogger := logger.WithFields(logrus.Fields{
		"path":              cmd.Path,
		"args":              strings.Join(arg, " |: "),
		"working directory": workDir,
		"internal error":    err,
	})

	levelLog(execLogger, level, msg)

	return
}

// RunToFile run and store output to file
func RunToFile(path string, name string, arg ...string) (cmd *exec.Cmd, exitCode int, strout string, strerr string, err error) {
	cmd, exitCode, strout, strerr, err = Run(name, arg...)

	file.Writeln(path, strout)

	if strerr != "" {
		file.Writeln(path+".err", strerr)
	}

	return
}

// RunToFileLog run with log and store output to file
func RunToFileLog(path string, logger *logrus.Logger, level logrus.Level, msg string, name string, arg ...string) (cmd *exec.Cmd, exitCode int, strout string, strerr string, err error) {
	cmd, exitCode, strout, strerr, err = RunToLog(logger, level, msg, name, arg...)

	file.Writeln(path, strout)

	if strerr != "" {
		file.Writeln(path+".err", strerr)
	}

	return
}

func execInit(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)

	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	return cmd
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
