package utils

import (
	"fmt"
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
	return fmt.Sprintf("%-19s%s\n", label, value)
}
