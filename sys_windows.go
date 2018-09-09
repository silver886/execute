package execute

import (
	"syscall"
)

// Hide set system process attribute with hide window
func (cmd *Cmd) Hide() {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
