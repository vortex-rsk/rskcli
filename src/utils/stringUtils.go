package utils

import (
	"fmt"
	"rskcli/src/utils/color"
	"strings"
)

func Pad(str string, size int) string {

	retStr := str
	if len(str) < size {
		retStr = str + strings.Repeat(" ", size-len(str))
	}
	return retStr
}

func Line(label string, value string) string {
	return fmt.Sprintf("%s%s\n", color.Blue(Pad(label, 23)), color.Green(value))
}

func TxLine(label string, value string) string {
	return fmt.Sprintf("%s%s\n", color.Blue(Pad(label, 12)), color.Green(value))
}
