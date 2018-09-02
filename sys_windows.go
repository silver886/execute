package execute

import (
	"os/exec"
	"syscall"
)

func sysAttribute(execCmd *exec.Cmd) {
	execCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
