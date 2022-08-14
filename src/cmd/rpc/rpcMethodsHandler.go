package rpc

import (
	"errors"
	"time"
)

func init() {

	AddHandler(&Handler{
		Methods: methods,
		Run: func(method *Method, ctx *Context) {

			var result interface{}
			var duration time.Duration
			var errReq error

			errValParam := validateMethodParams(method.Params, ctx.CommandArgs)

			params, paramParseError := parseMethodParams(method.Params, method.GetRPCName(), method.Names[0], ctx.CommandArgs[1:])

			if errValParam != nil || paramParseError != nil {
				err := getErrorToShow(errValParam, paramParseError)
				PrintMethodHelp(method.Names[0], method.Description, method.Params, err)
				PrintFooter(duration, ctx.Get("serverName"), ctx)
				return
			}

			result, errReq, duration = Call(method.ReturnType, method.GetRPCName(), params, ctx)

			// check for errors message in the request or response
			if errResp := checkResponseError(result); errReq != nil || errResp != nil {
				err := getErrorToShow(errReq, errResp)
				PrintMethodHelp(method.Names[0], method.Description, method.Params, err)
				PrintFooter(duration, ctx.Get("serverName"), ctx)
				return
			}

			// all right. Show results
			if !ctx.Flags["json"] {
				method.ProcessResult(&result, ctx)
				if !ctx.Flags["nofooter"] && !ctx.Flags["clean"] {
					PrintFooter(duration, ctx.Get("serverName"), ctx)
				}
			}
		},
	})
}

func getErrorToShow(err1 error, err2 error) error {
	if err1 != nil {
		return err1
	} else {
		return err2
	}
}

func checkResponseError(result interface{}) error {
	if result != nil {
		switch result.(type) {
		case *SimpleRpcResult:
			if len(result.(*SimpleRpcResult).Error.Message) > 0 {
				return errors.New(result.(*SimpleRpcResult).Error.Message)
			}
		case *BlockRpcResult:
			if len(result.(*BlockRpcResult).Error.Message) > 0 {
				return errors.New(result.(*BlockRpcResult).Error.Message)
			}
		}
	}
	return nil
}
