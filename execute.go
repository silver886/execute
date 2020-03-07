package execute

import (
	"github.com/silver886/file"
)

// Start starts the command but does not wait for it to complete
func (cmd *Cmd) Start() (*Cmd, error) {
	return cmd, cmd.Cmd.Start()
}

// Wait waits for the command to exit
func (cmd *Cmd) Wait() (*Cmd, error) {
	cmd.wait.mutex.Lock()

	if !cmd.wait.done {
		cmd.wait.err = cmd.Cmd.Wait()
		cmd.wait.done = true
	}

	cmd.wait.mutex.Unlock()

	return cmd, cmd.wait.err
}

// Run starts the command and waits for it to complete
func (cmd *Cmd) Run() (*Cmd, error) {
	if _, err := cmd.Start(); err != nil {
		return cmd, err
	}
	return cmd.Wait()
}

// RunToFile run and store output to file
func (cmd *Cmd) RunToFile(path string) (*Cmd, error) {
	if _, err := cmd.Run(); err != nil {
		return cmd, err
	}

	if _, err := file.Writeln(path, cmd.OutString()); err != nil {
		return cmd, err
	}

	if errString := cmd.ErrString(); errString != "" {
		if _, err := file.Writeln(path+".err", errString); err != nil {
			return cmd, err
		}
	}

	return cmd, nil
}
