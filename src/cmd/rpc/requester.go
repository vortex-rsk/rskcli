package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"rskcli/src/utils"
	"rskcli/src/utils/color"
	"time"
)

type Payload struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      string        `json:"id"`
}

const INTERNAL = "INTERNAL"
const STANDARD = "STANDARD"

func CallInternal(returnType string, methodRpcName string, params []interface{}, ctx *Context) (interface{}, error, time.Duration) {
	return CallHttp(returnType, methodRpcName, params, ctx, INTERNAL)
}

func Call(returnType string, methodRpcName string, params []interface{}, ctx *Context) (interface{}, error, time.Duration) {
	return CallHttp(returnType, methodRpcName, params, ctx, STANDARD)
}

func CallHttp(returnType string, methodRpcName string, params []interface{}, ctx *Context, mode string) (interface{}, error, time.Duration) {

	var payload = Payload{"2.0", methodRpcName, params, "74"} // []interface{}{700000, "false"}
	jsonReq, errorMarshal := json.Marshal(payload)

	if errorMarshal != nil {
		fmt.Printf("Error parsing request json: %s.\n", color.Red(errorMarshal.Error()))
		fmt.Println(color.Red(string(jsonReq)))
	}

	if ctx.Flags["jsonreq"] && mode != INTERNAL {
		fmt.Println(string(jsonReq))
	}

	start := time.Now()
	response, errorPost := http.Post(ctx.Get("serverUrl"), "application/json; charset=UTF-8", bytes.NewReader(jsonReq))
	elapsed := time.Since(start)
	if errorPost != nil {
		err := errors.New(fmt.Sprintf("Error: %s.\n", color.Red(errorPost.Error())))
		return nil, err, elapsed
	}

	if response.StatusCode == 403 { // Method Not Allowed
		err := errors.New(fmt.Sprintf(color.Red("Error: Method Not Allowed")))
		return nil, err, elapsed
	}

	body, errorRead := ioutil.ReadAll(response.Body)
	if errorRead != nil {
		err := errors.New(fmt.Sprintf("Error: %s.\n", color.Red(errorRead.Error())))
		return nil, err, elapsed
	}

	if (ctx.Flags["jsonresp"] || ctx.Flags["json"]) && mode != INTERNAL {
		fmt.Print(string(body) + utils.GetEnding(ctx.Flags["clean"]))
	}

	var result interface{}
	if returnType == "SimpleRpcResult" {
		result = &SimpleRpcResult{}
	} else if returnType == "BlockRpcResult" {
		result = &BlockRpcResult{}
	} else if returnType == "TransactionRpcResult" {
		result = &TransactionRpcResult{}
	} else {
		fmt.Println("returnType is wrong: " + returnType)
		os.Exit(1)
	}

	errorUnmarshal := json.Unmarshal(body, result)

	if errorUnmarshal != nil {
		err := errors.New(fmt.Sprintf("Error parsing response json: %s.\n", color.Red(errorUnmarshal.Error())) + "\n" + color.Red(string(body)))
		fmt.Printf("Error parsing response json: %s.\n", color.Red(errorUnmarshal.Error()))
		return nil, err, elapsed
	}

	return result, nil, elapsed
}
