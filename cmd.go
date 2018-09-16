package execute

import (
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

	buildOut strings.Builder
	buildErr strings.Builder

	stroutIndex []int
	strerrIndex []int
}

// WorkDir get the working directory
func (cmd *Cmd) WorkDir() (workDir string) {
	workDir, _ = filepath.Abs(cmd.Dir)

	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"cmd":               cmd,
				"working_directory": workDir,
			}),
		).Debugln("Get the working directory")
	}

	return
}

// ExitCode get the exit code
func (cmd *Cmd) ExitCode() (exitCode int) {
	if cmd.ProcessState == nil {
		if intLog {
			intLogger.WithFields(
				logger.DebugInfo(1, logrus.Fields{
					"cmd": cmd,
				}),
			).Errorln("Cannot get the exit code")
		}
		return
	}
	exitCode = cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()

	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"cmd":       cmd,
				"exit_code": exitCode,
			}),
		).Debugln("Get the exit code")
	}

	return
}

// Strout get the stdout as string
func (cmd *Cmd) Strout() (strout string) {
	strout = strings.TrimSpace(cmd.buildOut.String())

	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"cmd":    cmd,
				"stdout": strout,
			}),
		).Debugln("Get the stdout as string")
	}

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

	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"cmd":           cmd,
				"unused_stdout": strout,
			}),
		).Debugln("Get unused stdout")
	}

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

	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"cmd":            cmd,
				"stdout_history": strout,
			}),
		).Debugln("Get get stdout history")
	}

	return
}

// Strerr get the stderr as string
func (cmd *Cmd) Strerr() (strerr string) {
	strerr = strings.TrimSpace(cmd.buildErr.String())

	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"cmd":    cmd,
				"stderr": strerr,
			}),
		).Debugln("Get the stderr as string")
	}

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

	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"cmd":           cmd,
				"unused_stderr": strerr,
			}),
		).Debugln("Get unused stderr")
	}

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

	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"cmd":            cmd,
				"stderr_history": strerr,
			}),
		).Debugln("Get get stderr history")
	}

	return
}

// Start starts the specified command but does not wait for it to complete
func (cmd *Cmd) Start() (err error) {
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"cmd": cmd,
			}),
		).Debugln("Start command . . .")
	}

	err = cmd.Cmd.Start()

	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"cmd":               cmd,
				"executable_path":   cmd.Path,
				"arguments":         strings.Join(cmd.Args, " |: "),
				"working_directory": cmd.WorkDir(),
				"internal_error":    err,
			}),
		).Debugln("Start command")
	}

	return
}

// Run starts the specified command and waits for it to complete
func (cmd *Cmd) Run() (err error) {
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"cmd": cmd,
			}),
		).Debugln("Run command . . .")
	}

	err = cmd.Start()
	if err != nil {
		if intLog {
			intLogger.WithFields(
				logger.DebugInfo(1, logrus.Fields{
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
func (cmd *Cmd) RunToFile(path string) (err error) {
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"cmd":       cmd,
				"file_path": path,
			}),
		).Debugln("Run command to file . . .")
	}

	err = cmd.Run()

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
				"internal_error": err,
			}),
		).Debugln("Run command to file")
	}

	return
}

// New create a execute.Cmd
func New(name string, arg ...string) (cmd *Cmd) {
	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"command":   name,
				"arguments": arg,
			}),
		).Debugln("New command . . .")
	}

	execCmd := exec.Command(name, arg...)

	cmd = &Cmd{
		Cmd: execCmd,
	}

	cmd.Stdout = &cmd.buildOut
	cmd.Stderr = &cmd.buildErr

	if intLog {
		intLogger.WithFields(
			logger.DebugInfo(1, logrus.Fields{
				"cmd": cmd,
			}),
		).Debugln("New command")
	}

	return
}
