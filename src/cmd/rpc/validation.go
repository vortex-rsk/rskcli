package rpc

import (
	"errors"
	"fmt"
)

func validateMethodParams(params []*DefParam, args []string) error {
	if len(params) > 0 {
		var countMandatory int
		for _, param := range params {
			if param.Mandatory {
				countMandatory++
			}
		}
		if len(args)-1 < countMandatory {
			return errors.New(fmt.Sprintf("Method %s requires %d params", args[0], countMandatory))
		}
	}
	return nil
}
