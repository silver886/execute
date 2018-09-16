package execute

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
	"leoliu.io/file"
	"leoliu.io/logger"
)

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
func RunToLog(hide bool, extLogger *logger.Logger, level logrus.Level, msg string, name string, arg ...string) (cmd *Cmd, err error) {
	if extLogger == nil {
		err = errors.New("Invalid logger")
		return
	}

	cmd = New(name, arg...)
	if hide {
		cmd.Hide()
	}
	err = cmd.Run()

	execLogger := extLogger.WithFields(
		logger.DebugInfo(1, logrus.Fields{
			"path":              cmd.Path,
			"args":              strings.Join(arg, " |: "),
			"working_directory": cmd.WorkDir(),
			"exitcode":          cmd.ExitCode(),
			"stdout":            cmd.Strout(),
			"stderr":            cmd.Strerr(),
			"internal_error":    err,
		}))

	logger.LevelLog(execLogger, level, msg)

	return
}

// StartToLog start with log
func StartToLog(hide bool, extLogger *logger.Logger, level logrus.Level, msg string, name string, arg ...string) (cmd *Cmd, err error) {
	if extLogger == nil {
		err = errors.New("Invalid logger")
		return
	}

	cmd = New(name, arg...)
	if hide {
		cmd.Hide()
	}
	err = cmd.Start()

	execLogger := extLogger.WithFields(
		logger.DebugInfo(1, logrus.Fields{
			"path":              cmd.Path,
			"args":              strings.Join(arg, " |: "),
			"working_directory": cmd.WorkDir(),
			"internal_error":    err,
		}))

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
func RunToFileLog(hide bool, path string, extLogger *logger.Logger, level logrus.Level, msg string, name string, arg ...string) (cmd *Cmd, err error) {
	cmd, err = RunToLog(hide, extLogger, level, msg, name, arg...)

	file.Writeln(path, cmd.Strout())

	if strerr := cmd.Strerr(); strerr != "" {
		file.Writeln(path+".err", strerr)
	}

	return
}
