package rpc

import (
	"fmt"
	"rskcli/src/utils"
	"rskcli/src/utils/color"
	"strconv"
	"time"
)

func processSimple(resultArg *interface{}, ctx *Context) {

	result := (*resultArg).(*SimpleRpcResult)

	if result == nil || result.Result == nil {
		return
	}

	utils.PrintResult(result.Result.(string), ctx.Flags["clean"])
}

func processHexInt(resultArg *interface{}, ctx *Context) {

	result := (*resultArg).(*SimpleRpcResult)

	if result == nil || result.Result == nil {
		return
	}

	utils.PrintResult(utils.HexInt(result.Result), ctx.Flags["clean"])
}

func processBool(resultArg *interface{}, ctx *Context) {

	result := (*resultArg).(*SimpleRpcResult)

	if result.Result == nil {
		return
	}

	utils.PrintResult(strconv.FormatBool(result.Result.(bool)), ctx.Flags["clean"])
}

func processSyncing(resultArg *interface{}, ctx *Context) {

	result := (*resultArg).(*SimpleRpcResult)

	if result.Result == nil {
		return
	}

	switch val := result.Result.(type) {
	case bool:
		utils.PrintResult(strconv.FormatBool(val), ctx.Flags["clean"])
	case string:
		utils.PrintResult(val, ctx.Flags["clean"])
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

func processList(resultArg *interface{}, ctx *Context) {

	result := (*resultArg).(*SimpleRpcResult)

	if result.Result == nil {
		return
	}

	var values []interface{} = result.Result.([]interface{})
	for _, val := range values {
		fmt.Println(color.Green(val.(string)))
	}
}

func processFloat64String(resultArg *interface{}, ctx *Context) {
	result := (*resultArg).(*SimpleRpcResult)
	utils.PrintResult(fmt.Sprintf("%v", result.Result.(float64)), ctx.Flags["clean"])
}

func processBlock(resultArg *interface{}, ctx *Context) {

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
		utils.Line("timestamp: ", block.Timestamp+dateTime) +
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
		utils.Line("minimum gas price: ", block.MinimumGasPrice) +
		utils.Line("bitcoin merged mining header: ", block.BitcoinMergedMiningHeader) +
		utils.Line("bitcoin merged mining coinbase transaction: ", block.BitcoinMergedMiningCoinbaseTransaction) +
		utils.Line("bitcoin merged mining merkle proof: ", block.BitcoinMergedMiningMerkleProof) +
		utils.Line("hash for merged mining: ", block.HashForMergedMining) +
		utils.Line("paid fees: ", block.PaidFees) +
		color.Blue(utils.Pad("uncles("+strconv.Itoa(len(block.Uncles))+"):", 23)+"[\n")

	for idx, val := range block.Uncles {
		if idx > 0 {
			msg += "\n"
		}
		msg += "\t" + color.Green(val)
	}
	msg += color.Blue("\n]\n")

	msg += color.Blue(utils.Pad("transactions("+strconv.Itoa(len(block.Transactions))+"):", 19) + "[\n")

	for idx, val := range block.Transactions {
		if idx > 0 {
			msg += "\n"
		}
		msg += formatTransaction(val.(*Transaction), "\t", "SHORT")
	}

	msg += color.Blue("\n]")

	fmt.Println(msg)
}

func formatTransaction(tx *Transaction, prefix string, mode string) string {

	var msg string

	if len(tx.Input) == 0 {
		msg = color.Green(prefix + tx.Hash)
	} else {

		if mode == "FULL" {
			msg += utils.TxLine(prefix+"block number: ", utils.HexInt(tx.BlockNumber)) +
				utils.TxLine(prefix+"block hash: ", tx.BlockHash)
		}

		msg += utils.TxLine(prefix+"hash: ", tx.Hash) +
			utils.TxLine(prefix+"index: ", utils.HexInt(tx.TransactionIndex)) +
			utils.TxLine(prefix+"from: ", tx.From) +
			utils.TxLine(prefix+"to: ", tx.To) +
			utils.TxLine(prefix+"value: ", utils.HexInt(tx.Value)) +
			utils.TxLine(prefix+"gas: ", utils.HexInt(tx.Gas)) +
			utils.TxLine(prefix+"gas price: ", utils.HexInt(tx.GasPrice)) +
			utils.TxLine(prefix+"nonce: ", utils.HexInt(tx.Nonce)) +
			utils.TxLine(prefix+"input: ", tx.Input) +
			utils.TxLine(prefix+"V: ", tx.V) +
			utils.TxLine(prefix+"R: ", tx.R) +
			utils.TxLine(prefix+"S: ", tx.S)
	}

	return msg
}

func processTransaction(resultArg *interface{}, ctx *Context) {
	result := (*resultArg).(*TransactionRpcResult)

	fmt.Print(formatTransaction(&result.Result, "", "FULL"))
}
