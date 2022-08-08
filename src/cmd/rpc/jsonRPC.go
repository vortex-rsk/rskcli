package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"rskcli/src/utils"
	"strconv"
)

//var NOPARAMS []DefParam

type SimpleRPCResult struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      string      `json:"id"`
	Result  interface{} `json:"result"`
	Error   struct {
		Message string `json:"message"`
	} `json:"error"`
}

type Transaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

type Block struct {
	Number               string        `json:"number"`
	Hash                 string        `json:"hash"`
	ParentHash           string        `json:"parentHash"`
	Nonce                string        `json:"nonce"`
	Sha3Uncles           string        `json:"sha3Uncles"`
	LogsBloom            string        `json:"logsBloom"`
	TransactionsRoot     string        `json:"transactionsRoot"`
	StateRoot            string        `json:"stateRoot"`
	ReceiptsRoot         string        `json:"receiptsRoot"`
	Miner                string        `json:"miner"`
	Difficulty           string        `json:"difficulty"`
	CumulativeDifficulty string        `json:"cumulativeDifficulty"`
	TotalDifficulty      string        `json:"totalDifficulty"`
	ExtraData            string        `json:"extraData"`
	Size                 string        `json:"size"`
	GasLimit             string        `json:"gasLimit"`
	GasUsed              string        `json:"gasUsed"`
	Timestamp            string        `json:"timestamp"`
	Transactions         []interface{} `json:"transactions"`
	Uncles               []string      `json:"uncles"`
}

type BlockRpcResult struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      string `json:"id"`
	Result  Block  `json:"result"`
	Error   struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (br *BlockRpcResult) UnmarshalJSON(buf []byte) error {
	var tmp map[string]interface{}
	errorUnm := json.Unmarshal(buf, &tmp)
	if errorUnm != nil {
		return errorUnm
	}

	if tmp["result"] == nil {
		return errors.New("null result from server")
	}

	//br.Id = strconv.FormatFloat(tmp["id"].(float64), 'E', -1, 32)
	br.Id = tmp["id"].(string)
	br.Jsonrpc = tmp["jsonrpc"].(string)

	if tmp["error"] != nil {
		errorMap := tmp["error"].(map[string]interface{})
		br.Error.Message = errorMap["message"].(string)
	}

	if tmp["result"] != nil {
		blockMap := tmp["result"].(map[string]interface{})

		br.Result.CumulativeDifficulty = utils.HexInt(blockMap["cumulativeDifficulty"])
		br.Result.Number = utils.HexInt(blockMap["number"])
		br.Result.Hash = blockMap["hash"].(string)
		br.Result.ParentHash = blockMap["parentHash"].(string)
		if blockMap["nonce"] != nil {
			br.Result.Nonce = blockMap["nonce"].(string)
		}
		br.Result.Sha3Uncles = blockMap["sha3Uncles"].(string)
		br.Result.LogsBloom = blockMap["logsBloom"].(string)
		br.Result.TransactionsRoot = blockMap["transactionsRoot"].(string)
		br.Result.StateRoot = blockMap["stateRoot"].(string)
		br.Result.ReceiptsRoot = blockMap["receiptsRoot"].(string)
		br.Result.Miner = blockMap["miner"].(string)
		br.Result.Difficulty = utils.HexInt(blockMap["difficulty"])
		br.Result.CumulativeDifficulty = utils.HexInt(blockMap["cumulativeDifficulty"])
		br.Result.TotalDifficulty = utils.HexInt(blockMap["totalDifficulty"])
		br.Result.ExtraData = blockMap["extraData"].(string)
		br.Result.Size = utils.HexInt(blockMap["size"])
		br.Result.GasLimit = utils.HexInt(blockMap["gasLimit"])
		br.Result.GasUsed = utils.HexInt(blockMap["gasUsed"])
		br.Result.Timestamp = utils.HexInt(blockMap["timestamp"])

		uncles := blockMap["uncles"].([]interface{})
		br.Result.Uncles = make([]string, len(uncles))
		for idx, val := range uncles {
			br.Result.Uncles[idx] = val.(string)
		}

		transactions := blockMap["transactions"].([]interface{})
		br.Result.Transactions = make([]interface{}, len(transactions))
		for idx, val := range transactions {
			if str, isString := val.(string); isString {
				br.Result.Transactions[idx] = &Transaction{Hash: str}
			} else if reflect.ValueOf(val).Kind() == reflect.Map {
				//fmt.Println(reflect.ValueOf(val).Kind())
				txMap := val.(map[string]interface{})
				br.Result.Transactions[idx] = &Transaction{
					BlockHash:        txMap["blockHash"].(string),
					BlockNumber:      utils.HexInt(txMap["blockNumber"]),
					From:             txMap["from"].(string),
					Gas:              utils.HexInt(txMap["gas"]),
					GasPrice:         utils.HexInt(txMap["gasPrice"]),
					Hash:             txMap["hash"].(string),
					Input:            txMap["input"].(string),
					Nonce:            utils.HexInt(txMap["nonce"]),
					To:               txMap["to"].(string),
					TransactionIndex: utils.HexInt(txMap["transactionIndex"]),
					Value:            utils.HexInt(txMap["value"]),
				}

				if txMap["v"] != nil {
					br.Result.Transactions[idx].(*Transaction).V = txMap["v"].(string)
				}
				if txMap["r"] != nil {
					br.Result.Transactions[idx].(*Transaction).R = txMap["r"].(string)
				}
				if txMap["s"] != nil {
					br.Result.Transactions[idx].(*Transaction).S = txMap["s"].(string)
				}
			} else {
				fmt.Println("===============================")
				fmt.Println(reflect.ValueOf(val).Kind())
				fmt.Println("===============================")
			}
		}

	}

	return nil
}

func UnmarshalResponse(bytes []byte) (*SimpleRPCResult, error) {

	// Construct an anonymous struct that has looser typing
	// than our output field. We use this as a temporary
	// placeholder to parse the contents and construct
	// a properly constructed final result
	raw := struct {
		Jsonrpc string      `json:"jsonrpc"`
		Id      string      `json:"id"`
		Result  interface{} `json:"result"`
	}{}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return nil, err
	}

	// Construct result instance with as much information as
	// possible from the raw data
	result := &SimpleRPCResult{
		Jsonrpc: raw.Jsonrpc,
		Id:      raw.Id,
	}

	// Populate result by converting the value into a string
	// depending on the type of the value received
	switch val := raw.Result.(type) {
	case bool:
		result.Result = strconv.FormatBool(val)
	case string:
		result.Result = val
	default:
		fmt.Println(val)
	}

	return result, nil
}

func GetRPCName(cmd *Command) string {
	return cmd.Names[len(cmd.Names)-1]
}
