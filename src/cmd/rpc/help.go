package rpc

import (
	"fmt"
	"rskcli/src/utils/color"
	"strings"
)

var brackets = map[bool]map[string]string{true: {"open": "<", "close": ">"}, false: {"open": "[", "close": "]"}}
var mandatoryTxt = map[bool]string{true: color.Blue("required:"), false: "optional:"}

func PrintMethodHelp(name string, description string, params []*DefParam, error error) {

	if error != nil {
		fmt.Println(color.Red(error.Error()))
	}

	var paramsLine strings.Builder
	var paramsDesc strings.Builder
	var noParams string

	for _, param := range params {
		// params usage line
		if paramsLine.Len() > 0 {
			paramsLine.WriteString(" ")
		}
		paramsLine.WriteString(brackets[param.Mandatory]["open"])
		paramsLine.WriteString(param.Name)
		paramsLine.WriteString(brackets[param.Mandatory]["close"])
		// params desc
		paramsDesc.WriteString("\t\t")
		paramsDesc.WriteString(brackets[param.Mandatory]["open"])
		paramsDesc.WriteString(param.Name)
		paramsDesc.WriteString(brackets[param.Mandatory]["close"])
		paramsDesc.WriteString(" ")
		paramsDesc.WriteString(mandatoryTxt[param.Mandatory])
		paramsDesc.WriteString(" ")
		paramsDesc.WriteString(param.Description)
		paramsDesc.WriteString("\n")
	}
	if len(params) == 0 {
		noParams = " (has no params)"
	}
	fmt.Println("\nusage: rsk " + name + " " + paramsLine.String() + "\n")
	fmt.Println("\t" + description + "\n")
	fmt.Println("\tparameter list:" + noParams)
	fmt.Print(paramsDesc.String())

}

func PrintHelp() {
	fmt.Println("help")
}
