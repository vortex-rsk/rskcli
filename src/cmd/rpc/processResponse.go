package rpc

import (
	"fmt"
	"rskcli/src/utils"
	"rskcli/src/utils/color"
	"strconv"
	"time"
)

func processSimple(resultArg *interface{}) {

	result := (*resultArg).(*SimpleRPCResult)

	if result == nil || result.Result == nil {
		return
	}

	fmt.Println(color.Green(result.Result.(string)))
}

func processHexInt(resultArg *interface{}) {

	result := (*resultArg).(*SimpleRPCResult)

	if result == nil || result.Result == nil {
		return
	}

	fmt.Println(color.Green(utils.HexInt(result.Result)))

}

func processBool(resultArg *interface{}) {

	result := (*resultArg).(*SimpleRPCResult)

	if result.Result == nil {
		return
	}

	fmt.Println(color.Green(strconv.FormatBool(result.Result.(bool))))
}

func processSyncing(resultArg *interface{}) {

	result := (*resultArg).(*SimpleRPCResult)

	if result.Result == nil {
		return
	}

	switch val := result.Result.(type) {
	case bool:
		fmt.Println(color.Green(strconv.FormatBool(val)))
	case string:
		fmt.Println(color.Green(val))
	case map[string]interface{}:
		var valuesMap map[string]interface{} = result.Result.(map[string]interface{})
		for key, ele := range valuesMap {
			fmt.Println(color.Green(key + ": " + utils.HexInt(ele)))
		}
	default:
		fmt.Println("default")
		fmt.Println(result.Result)
	}

}

func processList(resultArg *interface{}) {

	result := (*resultArg).(*SimpleRPCResult)

	if result.Result == nil {
		return
	}

	var values []interface{} = result.Result.([]interface{})
	for _, val := range values {
		fmt.Println(color.Green(val.(string)))
	}
}

func processFloat64String(resultArg *interface{}) {
	result := (*resultArg).(*SimpleRPCResult)
	fmt.Println(color.Green(fmt.Sprintf("%v", result.Result.(float64))))
}

func processBlock(resultArg *interface{}) {

	result := (*resultArg).(*BlockRpcResult)

	if len(result.Error.Message) > 0 || len(result.Result.Number) == 0 {
		return
	}

	block := result.Result

	var dateTime string
	tsInt64, err := strconv.ParseInt(block.Timestamp, 10, 32)

	if err == nil {
		dateTime = " | " + time.Unix(tsInt64, 0).Format("2006-01-02 15:04:05")
	} else {
		fmt.Println(err)
	}

	msg := utils.Line("number: ", block.Number) +
		utils.Line("hash: ", block.Hash) +
		utils.Line("parent hash: ", block.ParentHash) +
		utils.Line("nonce: ", block.Nonce) +
		utils.Line("logsBloom: ", block.LogsBloom) +
		utils.Line("transactions root: ", block.TransactionsRoot) +
		utils.Line("stateRoot: ", block.StateRoot) +
		utils.Line("receipts root: ", block.ReceiptsRoot) +
		utils.Line("miner: ", block.Miner) +
		utils.Line("difficulty: ", block.Difficulty) +
		utils.Line("cumulative difficulty: ", block.CumulativeDifficulty) +
		utils.Line("total difficulty: ", block.TotalDifficulty) +
		utils.Line("extra data: ", block.ExtraData) +
		utils.Line("size: ", block.Size) +
		utils.Line("gas limit: ", block.GasLimit) +
		utils.Line("gas used: ", block.GasUsed) +
		utils.Line("timestamp: ", block.Timestamp+dateTime) +
		utils.Pad("uncles("+strconv.Itoa(len(block.Uncles))+"):", 19) + "[\n"

	for idx, val := range block.Uncles {
		if idx > 0 {
			msg += "\n"
		}
		msg += "\t" + val
	}
	msg += "\n]\n"

	msg += utils.Pad("transactions("+strconv.Itoa(len(block.Transactions))+"):", 19) + "[\n"

	for idx, val := range block.Transactions {
		if idx > 0 {
			msg += "\n"
		}
		msg += formatTransaction(val.(*Transaction), "\t")
	}

	msg += "\n]"

	fmt.Println(color.Green(msg))
}

func formatTransaction(tx *Transaction, prefix string) string {

	var msg string

	if len(tx.Input) == 0 {
		msg = prefix + tx.Hash
	} else {
		msg = utils.Line(prefix+"from: ", tx.From) +
			utils.Line(prefix+"to: ", tx.To) +
			utils.Line(prefix+"value: ", utils.HexInt(tx.Value)) +
			utils.Line(prefix+"gas: ", utils.HexInt(tx.Gas)) +
			utils.Line(prefix+"gas price: ", utils.HexInt(tx.GasPrice)) +
			utils.Line(prefix+"hash: ", tx.Hash) +
			utils.Line(prefix+"nonce: ", tx.Nonce) +
			utils.Line(prefix+"index: ", tx.TransactionIndex)
		//utils.Line(prefix+"input: ", tx.Input) +
		//utils.Line(prefix+"(V,R,S): ", fmt.Sprintf("%s,%s,%s", tx.V, tx.R, tx.S))
	}

	return msg
}
