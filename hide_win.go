// +build windows

package execute

import (
	"syscall"
)

// Hide set system process attribute with hide window
func (cmd *Cmd) Hide() *Cmd {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	return cmd
}
