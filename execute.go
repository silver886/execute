package execute

import (
	"strings"

	"github.com/sirupsen/logrus"
	"leoliu.io/file"
	"leoliu.io/logger"
)

var (
	intLog    bool
	intLogger *logger.Entry
)

// SetLogger set internal logger for logging
func SetLogger(extLogger *logger.Logger) {
	intLogger = extLogger.WithField("prefix", "execute")
	intLog = true
}

// ResetLogger reset internal logger
func ResetLogger() {
	intLogger = nil
	intLog = false
}

// Start starts the specified command but does not wait for it to complete
func Start(hide bool, name string, arg ...string) (cmd *Cmd, err error) {
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"hide":      hide,
				"command":   name,
				"arguments": arg,
			}),
		).Debugln("Start command . . .")
	}

	cmd = New(name, arg...)
	if hide {
		cmd.Hide()
	}
	err = cmd.Start()

	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"cmd":               cmd,
				"executable_path":   cmd.Path,
				"arguments":         strings.Join(arg, " |: "),
				"working_directory": cmd.WorkDir(),
				"internal_error":    err,
			}),
		).Debugln("Start command")
	}

	return
}

// Run starts the specified command and waits for it to complete
func Run(hide bool, name string, arg ...string) (cmd *Cmd, err error) {
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"hide":      hide,
				"command":   name,
				"arguments": arg,
			}),
		).Debugln("Run command . . .")
	}

	cmd, err = Start(hide, name, arg...)
	if err != nil {
		if intLog {
			intLogger.WithFields(
				logger.DebugInfo(1, logrus.Fields{
					"cmd":            cmd,
					"internal_error": err,
				}),
			).Errorln("Cannot start command")
		}
		return
	}

	err = cmd.Wait()

	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"cmd":            cmd,
				"exit_code":      cmd.ExitCode(),
				"stdout":         cmd.Strout(),
				"stderr":         cmd.Strerr(),
				"internal_error": err,
			}),
		).Debugln("Run command")
	}

	return
}

// RunToFile run and store output to file
func RunToFile(hide bool, path string, name string, arg ...string) (cmd *Cmd, err error) {
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"hide":      hide,
				"file_path": path,
				"command":   name,
				"arguments": arg,
			}),
		).Debugln("Run command to file . . .")
	}

	cmd, err = Run(hide, name, arg...)

	if _, err := file.Writeln(path, cmd.Strout()); err != nil {
		if intLog {
			intLogger.WithFields(
				logger.DebugInfo(1, logrus.Fields{
					"internal_error": err,
				}),
			).Errorln("Cannot write to file")
		}
	}

	if strerr := cmd.Strerr(); strerr != "" {
		if intLog {
			intLogger.WithFields(
				logger.DebugInfo(1, logrus.Fields{
					"file_path": path + ".err",
				}),
			).Debugln("Generate stderr file . . .")
		}
		if _, err := file.Writeln(path+".err", strerr); err != nil {
			if intLog {
				intLogger.WithFields(
					logger.DebugInfo(1, logrus.Fields{
						"internal_error": err,
					}),
				).Errorln("Cannot write to file")
			}
		}
	}

	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"cmd":            cmd,
				"internal_error": err,
			}),
		).Debugln("Run command to file")
	}

	return
}
