package rpc

import (
	"fmt"
	"rskcli/src/utils"
	"rskcli/src/utils/color"
	"time"
)

var chains map[string]string = map[string]string{"30": "mainnet", "31": "testnet", "33": "regtest"}

func PrintFooter(elapsed time.Duration, serverName string, ctx *Context) {

	if utils.Config.GetBoolean("printFooter") {
		result, _, _ := CallInternal("SimpleRpcResult", "eth_chainId", []interface{}{}, ctx)
		chain := "----"
		if result != nil && result.(*SimpleRpcResult).Result != nil {
			chain = chains[utils.HexInt(result.(*SimpleRpcResult).Result)]
		}
		var info = time.Now().Format("2006-01-02 15:04:05") + " | " + serverName + " | " + ctx.Get("serverUrl") + " | " + chain + " | took " + elapsed.String()
		fmt.Println(color.LightGrey(info))
	}

}
