package execute

import (
	"github.com/silver886/file"
)

// Start starts the specified command but does not wait for it to complete
func Start(hide bool, name string, arg ...string) (*Cmd, error) {
	cmd := New(name, arg...)
	if hide {
		cmd.Hide()
	}

	return cmd, cmd.Start()
}

// Run starts the specified command and waits for it to complete
func Run(hide bool, name string, arg ...string) (*Cmd, error) {
	cmd, err := Start(hide, name, arg...)
	if err != nil {
		return cmd, err
	}

	return cmd, cmd.Wait()
}

// RunToFile run and store output to file
func RunToFile(hide bool, path string, name string, arg ...string) (*Cmd, error) {
	cmd, err := Run(hide, name, arg...)

	if _, err := file.Writeln(path, cmd.OutString()); err != nil {
		return cmd, err
	}

	if errString := cmd.ErrString(); errString != "" {
		if _, err := file.Writeln(path+".err", errString); err != nil {
			return cmd, err
		}
	}

	return cmd, err
}
