package execute

import "strings"

// ErrString get the stderr as string
func (cmd *Cmd) ErrString() string {
	return strings.TrimSpace(cmd.ErrBuffer.String())
}

// ErrStringNext get unused stderr
func (cmd *Cmd) ErrStringNext() string {
	errString := cmd.ErrString()
	history := len(cmd.errStringIndex)
	var pointer int
	if history == 0 {
		pointer = 0
	} else {
		pointer = cmd.errStringIndex[history-1]
	}
	cmd.errStringIndex = append(cmd.errStringIndex, len(errString))
	return errString[pointer:]
}

// ErrStringHistory get stderr history access by ErrStringNext
func (cmd *Cmd) ErrStringHistory() []string {
	errString, errStringFull := []string{}, cmd.ErrString()
	for i, val := range cmd.errStringIndex {
		if i == 0 {
			errString = append(errString, errStringFull[:val])
		} else {
			errString = append(errString, errStringFull[cmd.errStringIndex[i-1]:val])
		}
	}
	return errString
}
