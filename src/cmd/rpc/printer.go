package rpc

import (
	"fmt"
	"rskcli/src/utils"
	"rskcli/src/utils/color"
	"time"
)

var chains map[string]string = map[string]string{"30": "mainnet", "31": "testnet", "33": "regtest"}

func PrintFooter(elapsed time.Duration, serverName string, serverUrl string) {

	if utils.Config.GetBoolean("printFooter") {
		result := CallInternal("eth_chainId", []interface{}{}, serverUrl)
		chain := "----"
		if result != nil && result.Result != nil {
			chain = chains[utils.HexInt(result.Result)]
		}
		var info = time.Now().Format("2006-01-02 15:04:05") + " | " + serverName + " | " + serverUrl + " | " + chain + " | took " + elapsed.String()
		fmt.Println(color.LightGrey(info))
	}

}
