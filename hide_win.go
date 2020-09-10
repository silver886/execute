// +build windows

package execute

import (
	"syscall"
)

// Hide set system process attribute with hide window
func (cmd *Cmd) Hide() *Cmd {
	if cmd.SysProcAttr != nil {
		cmd.SysProcAttr.HideWindow = true
	} else {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
	}
	return cmd
}
