package utils

import (
	"fmt"
	"rskcli/src/utils/color"
	"strconv"
	"strings"
)

func IndexOf(arr []string, search string) int {
	for idx, val := range arr {
		if search == val {
			return idx
		}
	}
	return -1
}

func IndexOfContain(arr []string, search string) int {
	for idx, val := range arr {
		if strings.Contains(val, search) {
			return idx
		}
	}
	return -1
}

func HexInt(value interface{}) string {
	result, _ := strconv.ParseInt(strings.Replace(value.(string), "0x", "", -1), 16, 32)
	return strconv.FormatInt(result, 10)
}

func GetEnding(clean bool) string {
	if clean {
		return ""
	} else {
		return "\n"
	}
}

func PrintResult(str string, clean bool) {
	fmt.Print(color.Green(str + GetEnding(clean)))
}
