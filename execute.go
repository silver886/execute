package execute

import (
	"leoliu.io/file"
)

// Start starts the specified command but does not wait for it to complete
func Start(hide bool, name string, arg ...string) (cmd *Cmd, err error) {

	cmd = New(name, arg...)
	if hide {
		cmd.Hide()
	}
	err = cmd.Start()

	return
}

// Run starts the specified command and waits for it to complete
func Run(hide bool, name string, arg ...string) (cmd *Cmd, err error) {

	cmd, err = Start(hide, name, arg...)
	if err != nil {
		return
	}

	err = cmd.Wait()

	return
}

// RunToFile run and store output to file
func RunToFile(hide bool, path string, name string, arg ...string) (cmd *Cmd, err error) {

	cmd, err = Run(hide, name, arg...)

	if _, err := file.Writeln(path, cmd.Strout()); err != nil {
	}

	if strerr := cmd.Strerr(); strerr != "" {
		if _, err := file.Writeln(path+".err", strerr); err != nil {
		}
	}

	return
}
