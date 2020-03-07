// +build windows

package execute

import (
	"fmt"
	"syscall"
)

// RawArgf set system process attribute with raw argument content formatter
func (cmd *Cmd) RawArgf(format string, content ...interface{}) *Cmd {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CmdLine: fmt.Sprintf(" "+format, content...),
	}
	return cmd
}
