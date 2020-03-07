// +build !windows

package execute

// Hide does nothing on non Windows platforms
func (cmd *Cmd) Hide() *Cmd { return cmd }
