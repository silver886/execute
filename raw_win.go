// +build windows

package execute

import (
	"fmt"
	"syscall"
)

// RawArgf set system process attribute with raw argument content formatter
func (cmd *Cmd) RawArgf(format string, content ...interface{}) *Cmd {
	if cmdLine := fmt.Sprintf(" "+format, content...); cmd.SysProcAttr != nil {
		cmd.SysProcAttr.CmdLine = cmdLine
	} else {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CmdLine: cmdLine,
		}
	}
	return cmd
}
