package execute

import "strings"

// OutString get the stdout as string
func (cmd *Cmd) OutString() string {
	return strings.TrimSpace(cmd.OutBuffer.String())
}

// OutStringNext get unused stdout
func (cmd *Cmd) OutStringNext() string {
	outString := cmd.OutString()
	history := len(cmd.outStringIndex)
	var pointer int
	if history == 0 {
		pointer = 0
	} else {
		pointer = cmd.outStringIndex[history-1]
	}
	cmd.outStringIndex = append(cmd.outStringIndex, len(outString))
	return outString[pointer:]
}

// OutStringHistory get stdout history access by OutStringNext
func (cmd *Cmd) OutStringHistory() []string {
	outString, outStringFull := []string{}, cmd.OutString()
	for i, val := range cmd.outStringIndex {
		if i == 0 {
			outString = append(outString, outStringFull[:val])
		} else {
			outString = append(outString, outStringFull[cmd.outStringIndex[i-1]:val])
		}
	}
	return outString
}
