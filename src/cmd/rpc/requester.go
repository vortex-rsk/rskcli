package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"rskcli/src/utils/color"
	"time"
)

type Payload struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      string        `json:"id"`
}

func CallSimple(methodRpcName string, params []interface{}, serverUrl string) (*SimpleRPCResult, error, time.Duration) {

	var payload = Payload{"2.0", methodRpcName, params, "74"} // []interface{}{700000, "false"}
	jsonReq, errorMarshal := json.Marshal(payload)

	if errorMarshal != nil {
		fmt.Printf("Error parsing request json: %s.\n", color.Red(errorMarshal.Error()))
		fmt.Println(color.Red(string(jsonReq)))
	}

	//fmt.Println(string(jsonReq))

	start := time.Now()
	response, errorPost := http.Post(serverUrl, "application/json; charset=UTF-8", bytes.NewReader(jsonReq))
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

	//fmt.Println(string(body))

	result := &SimpleRPCResult{}

	errorUnmarshal := json.Unmarshal(body, result)

	if errorUnmarshal != nil {
		err := errors.New(fmt.Sprintf("Error parsing response json: %s.\n", color.Red(errorUnmarshal.Error())) + "\n" + color.Red(string(body)))
		fmt.Printf("Error parsing response json: %s.\n", color.Red(errorUnmarshal.Error()))
		return nil, err, elapsed
	}

	return result, nil, elapsed
}

func CallBlock(methodRpcName string, params []interface{}, serverUrl string) (interface{}, error, time.Duration) {

	var payload = Payload{"2.0", methodRpcName, params, "74"} // []interface{}{700000, "false"}
	jsonReq, errorMarshal := json.Marshal(payload)

	if errorMarshal != nil {
		fmt.Printf("Error parsing request json: %s.\n", color.Red(errorMarshal.Error()))
		fmt.Println(color.Red(string(jsonReq)))
	}

	fmt.Println(string(jsonReq))

	start := time.Now()
	response, errorPost := http.Post(serverUrl, "application/json; charset=UTF-8", bytes.NewReader(jsonReq))
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

	//fmt.Println(string(body))
	result := &BlockRpcResult{}
	errorUnmarshal := json.Unmarshal(body, result)

	if errorUnmarshal != nil {
		err := errors.New(fmt.Sprintf("Error parsing response json: %s.\n", color.Red(errorUnmarshal.Error())) + "\n" + color.Red(string(body)))
		fmt.Printf("Error parsing response json: %s.\n", color.Red(errorUnmarshal.Error()))
		return nil, err, elapsed
	}

	return result, nil, elapsed
}

func CallInternal(methodRpcName string, params []interface{}, server string) *SimpleRPCResult {

	var payload = Payload{"2.0", methodRpcName, params, "73"}
	jsonReq, errorMarshal := json.Marshal(payload)

	if errorMarshal != nil {
		fmt.Printf("Error parsing request json: %s.\n", color.Red(errorMarshal.Error()))
		fmt.Println(color.Red(string(jsonReq)))
	}

	response, errorPost := http.Post(server, "application/json; charset=UTF-8", bytes.NewReader(jsonReq))
	if errorPost != nil {
		return nil
	}

	if response.StatusCode == 403 { // Method Not Allowed
		return nil
	}

	body, errorRead := ioutil.ReadAll(response.Body)
	if errorRead != nil {
		return nil
	}

	result := &SimpleRPCResult{}
	errorUnmarshal := json.Unmarshal(body, result)
	if errorUnmarshal != nil {
		return nil
	}

	return result
}
