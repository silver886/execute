// +build !windows

package execute

// RawArgf does nothing on non Windows platforms
func (cmd *Cmd) RawArgf(format string, content ...interface{}) *Cmd { return cmd }
