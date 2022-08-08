package rpc

import (
	"errors"
	"golang.org/x/exp/slices"
	"rskcli/src/utils"
	"strconv"
	"strings"
)

type DefParam struct {
	Name        string
	Mandatory   bool
	Description string
	Encoded     bool
}

type Method struct {
	Names         []string
	Params        []*DefParam
	ProcessResult func(result *interface{})
	ReturnType    string
	Description   string
}

var methods []Method

func init() {

	var NOPARAMS []*DefParam

	blockNumberTag := &DefParam{"blockNumber|tag", true, " integer block number, or the string \"latest\", \"earliest\" or \"pending\"", true}

	SIMPLE := "SimpleRpcResult"
	BLOCK := "BlockRpcResult"

	methods = []Method{
		Method{[]string{"clientVersion", "web3ClientVersion", "web3_clientVersion"}, NOPARAMS, processSimple, SIMPLE, "The current node version."},
		Method{[]string{"sha3", "web3Sha3", "web3_sha3"}, []*DefParam{{"data", true, "the data to convert into a SHA3 hash.", true}}, processSimple, SIMPLE, "Returns Keccak-256 (not the standardized SHA3-256) of the given data."},
		Method{[]string{"version", "netVersion", "net_version"}, NOPARAMS, processSimple, SIMPLE, "Returns the current network id."},
		Method{[]string{"listening", "netListening", "net_listening"}, NOPARAMS, processBool, SIMPLE, "Returns true if client is actively listening for network connections."},
		Method{[]string{"peerCount", "netPeerCount", "net_peerCount"}, NOPARAMS, processHexInt, SIMPLE, "Returns number of peers currently connected to the client."},
		Method{[]string{"peerList", "netPeerList", "net_peerList"}, NOPARAMS, processList, SIMPLE, "Returns list of peers known to the client."},
		Method{[]string{"chainId", "ethChainId", "eth_chainId"}, NOPARAMS, processHexInt, SIMPLE, "Returns the currently configured chain id"},
		Method{[]string{"ethProtocolVersion", "eth_protocolVersion"}, NOPARAMS, processHexInt, SIMPLE, "Returns the current node protocol version."},
		Method{[]string{"protocolVersion", "rskProtocolVersion", "rsk_protocolVersion"}, NOPARAMS, processHexInt, SIMPLE, "Returns the current node protocol version."},
		Method{[]string{"syncing", "ethSyncing", "eth_syncing"}, NOPARAMS, processSyncing, SIMPLE, "Returns an object with data about the sync status or false."},
		Method{[]string{"coinbase", "ethCoinbase", "eth_coinbase"}, NOPARAMS, processSimple, SIMPLE, "Returns the client coinbase address."},
		Method{[]string{"mining", "ethMining", "eth_mining"}, NOPARAMS, processBool, SIMPLE, "Returns true if client is actively mining new blocks."},
		Method{[]string{"hashrate", "ethHashrate", "eth_hashrate"}, NOPARAMS, processFloat64String, SIMPLE, "Returns the number of hashes per second that the node is mining with."},
		Method{[]string{"gasPrice", "ethGasPrice", "eth_gasPrice"}, NOPARAMS, processHexInt, SIMPLE, "Returns the current price per gas in wei."},
		Method{[]string{"accounts", "ethAccounts", "eth_accounts"}, NOPARAMS, processList, SIMPLE, "Returns a list of addresses owned by node."},
		Method{[]string{"blockNumber", "ethBlockNumber", "eth_blockNumber"}, NOPARAMS, processHexInt, SIMPLE, "Returns the number of most recent block."},
		Method{
			[]string{"balance", "getBalance", "ethGetBalance", "eth_getBalance"},
			[]*DefParam{
				{"data", true, "address to check for balance.", false},
				// getBalance is the only one that has the optional blockNumber the nodes defaults to 'latest'
				{"blockNumber|tag", false, " integer block number, or the string \"latest\", \"earliest\" or \"pending\"", false},
			},
			processHexInt, SIMPLE,
			"Returns the balance of the account of given address in wei.",
		},
		// TODO: TESTAR
		Method{
			[]string{"storageAt", "getStorageAt", "ethGetStorageAt", "eth_getStorageAt"},
			[]*DefParam{
				{"contractAddress", true, "contract address", false},
				{"position", true, "integer of the position in the storage", false},
				{"blockNumber|tag", false, " integer block number, or the string \"latest\", \"earliest\" or \"pending\"", false},
			},
			processSimple, SIMPLE,
			"Returns the value from a storage position at a given contract address.",
		},
		Method{
			[]string{"transactionCount", "getTransactionCount", "ethGetTransactionCount", "eth_getTransactionCount"},
			[]*DefParam{{"address", true, "the address", false}, blockNumberTag},
			processHexInt, SIMPLE,
			"Returns the number of transactions sent from an address.",
		},
		Method{
			[]string{"blockTransactionCountByHash", "getBlockTransactionCountByHash", "ethGetBlockTransactionCountByHash", "eth_getBlockTransactionCountByHash"},
			[]*DefParam{{"blockHash", true, "hash of a block", false}},
			processHexInt, SIMPLE,
			"Returns the number of transactions in a block from a block matching the given block hash.",
		},
		Method{
			[]string{"blockTransactionCountByNumber", "getBlockTransactionCountByNumber", "ethGetBlockTransactionCountByNumber", "eth_getBlockTransactionCountByNumber"},
			[]*DefParam{blockNumberTag},
			processHexInt, SIMPLE,
			"Returns the number of transactions in a block matching the given block number.",
		},
		Method{
			[]string{"uncleCountByBlockHash", "getUncleCountByBlockHash", "ethGetUncleCountByBlockHash", "eth_getUncleCountByBlockHash"},
			[]*DefParam{{"blockHash", true, "hash of a block", false}},
			processHexInt, SIMPLE,
			"Returns the number of uncles in a block from a block matching the given block hash.",
		},
		Method{
			[]string{"uncleCountByBlockNumber", "getUncleCountByBlockNumber", "ethGetUncleCountByBlockNumber", "eth_getUncleCountByBlockNumber"},
			[]*DefParam{blockNumberTag},
			processHexInt, SIMPLE,
			"Returns the number of uncles in a block from a block matching the given block number.",
		},
		Method{
			[]string{"code", "getCode", "ethGetCode", "eth_getCode"},
			[]*DefParam{{"address", true, "the contract address", false}, blockNumberTag},
			processSimple, SIMPLE,
			"String value of the compiled bytecode or 0 if it's not a contract.",
		},
		// TODO TESTAR
		Method{
			[]string{"sign", "ethSign", "eth_sign"},
			[]*DefParam{{"address", true, "the address", false}, {"message", true, "the message to be signed", false}},
			processSimple, SIMPLE,
			"The sign method calculates an Ethereum specific signature with: sign(keccak256(\"\\x19Ethereum Signed Message:\\n\" + len(message) + message))).\n\nBy adding a prefix to the message makes the calculated signature recognisable as an Ethereum specific signature. This prevents misuse where a malicious DApp can sign arbitrary data (e.g. transaction) and use the signature to impersonate the victim.\n\nNote the address to sign with must be unlocked.",
		},
		// TODO: TESTAR
		Method{
			[]string{"sendTransaction", "ethSendTransaction", "eth_sendTransaction"},
			[]*DefParam{
				{"from:<VALUE>", true, "Origin address", false},
				{"to:<VALUE>", false, "The address the transaction is directed to. (optional when creating new contract)", false},
				{"gas:<VALUE>", false, "Integer of the gas provided for the transaction execution. It will return unused gas. (optional, default: 90000)", true},
				{"gasPrice:<VALUE>", false, "Integer of the gasPrice used for each paid gas", true},
				{"value", false, "Integer of the value sent with this transaction", true},
				{"data:<VALUE>", false, "The compiled code of a contract OR the hash of the invoked method signature and encoded parameters.", false},
				{"nonce:<VALUE>", false, " Integer of a nonce. This allows to overwrite your own pending transactions that use the same nonce.", true},
			},
			processSimple, SIMPLE,
			"Creates new message call transaction or a contract creation, if the data field contains code. Returns the transaction hash, or the zero hash if the transaction is not yet available.",
		},
		// TODO: TESTAR
		Method{
			[]string{"sendRawTransaction", "ethSendRawTransaction", "eth_sendRawTransaction"},
			[]*DefParam{{"data", true, "The signed transaction data", false}},
			processSimple, SIMPLE,
			"Creates new message call transaction or a contract creation for signed transactions. Returns the transaction hash, or the zero hash if the transaction is not yet available.",
		},
		// TODO: TESTAR
		Method{
			[]string{"call", "ethCall", "eth_call"},
			[]*DefParam{
				{"from:<VALUE>", false, "Origin address", false},
				{"to:<VALUE>", false, "The address the transaction is directed to. (optional when creating new contract)", false},
				{"gas:<VALUE>", false, "Integer of the gas provided for the transaction execution. eth_call consumes zero gas, but this parameter may be needed by some executions.", true},
				{"gasPrice:<VALUE>", false, "Integer of the gasPrice used for each paid gas", true},
				{"value:<VALUE>", false, "Integer of the value sent with this transaction", true},
				{"data:<VALUE>", false, "Hash of the method signature and encoded parameters", false},
				{"nonce:<VALUE>", false, " Integer of a nonce. This allows to overwrite your own pending transactions that use the same nonce.", true},
				{"blockNumber:<VALUE>", true, " integer block number, or the string \"latest\", \"earliest\" or \"pending\"", false},
			},
			processSimple, SIMPLE,
			"Executes a new message call immediately without creating a transaction on the block chain. Returns the value of executed contract method.",
		},
		// TODO: TESTAR
		Method{
			[]string{"estimateGas", "ethEstimateGas", "eth_estimateGas"},
			[]*DefParam{
				{"from:<VALUE>", false, "Origin address", false},
				{"to:<VALUE>", false, "The address the transaction is directed to. (optional when creating new contract)", false},
				{"gas:<VALUE>", false, "Integer of the gas provided for the transaction execution. eth_call consumes zero gas, but this parameter may be needed by some executions.", true},
				{"gasPrice:<VALUE>", false, "Integer of the gasPrice used for each paid gas", true},
				{"value:<VALUE>", false, "Integer of the value sent with this transaction", true},
				{"data:<VALUE>", false, "Hash of the method signature and encoded parameters", false},
				{"nonce:<VALUE>", false, " Integer of a nonce. This allows to overwrite your own pending transactions that use the same nonce.", true},
				{"blockNumber:<VALUE>", true, " integer block number, or the string \"latest\", \"earliest\" or \"pending\"", false},
			},
			processHexInt, SIMPLE,
			"Generates and returns an estimate of how much gas is necessary to allow the transaction to complete. The transaction will not be added to the blockchain. Note that the estimate may be significantly more than the amount of gas actually used by the transaction, for a variety of reasons including EVM mechanics and node performance.",
		},

		// BLOCKS

		Method{
			[]string{"blockByHash", "getBlockByHash", "ethGetBlockByHash", "eth_getBlockByHash"},
			[]*DefParam{
				{"blockHash", true, "Hash of a block", false},
				{"completeTransactions", true, "If true it returns the full transaction objects, if false only the hashes of the transactions.", false}},
			processBlock, BLOCK,
			"Returns information about a block by hash.",
		},
		Method{
			[]string{"blockByNumber", "getBlockByNumber", "ethGetBlockByNumber", "eth_getBlockByNumber"},
			[]*DefParam{
				blockNumberTag,
				{"completeTransactions", true, "If true it returns the full transaction objects, if false only the hashes of the transactions.", false},
			},
			processBlock, BLOCK,
			"Returns information about a block by hash.",
		},
	}

}

func (method *Method) GetRPCName() string {
	return method.Names[len(method.Names)-1]
}

func parseMethodParams(params []*DefParam, rpcName string, name string, args []string) ([]interface{}, error) {

	if rpcName == "eth_call" || rpcName == "eth_sendTransaction" || rpcName == "eth_estimateGas" {
		retSize := 1
		var blockNumber string
		blockNumberIdx := utils.IndexOfContain(args, "blockNumber")
		//fmt.Println(args)
		if blockNumberIdx >= 0 {
			retSize = 2
			blockNumber = args[blockNumberIdx]
			args = slices.Delete(args, 1, 2)
		}

		ret := make([]interface{}, retSize)

		callArgs := make(map[string]string)
		for _, val := range args {
			parts := strings.Split(val, ":")
			if len(parts) == 1 {
				return nil, errors.New(name + " parameters must be prefixed by its name and colon. Like in \"param:value\"")
			}
			callArgs[parts[0]] = parts[1]
		}
		ret[0] = callArgs

		if len(blockNumber) > 0 {
			parts := strings.Split(blockNumber, ":")
			ret[1] = parts[1]
		}

		return ret, nil
	} else {
		return stringArrayToInterface(params, args), nil
	}

}
func stringArrayToInterface(params []*DefParam, values []string) []interface{} {
	inter := make([]interface{}, len(values))
	for idx, val := range values {
		if params[idx].Encoded && !strings.HasPrefix(val, "0x") {
			intVar, _ := strconv.ParseInt(val, 0, 32)
			inter[idx] = "0x" + strconv.FormatInt(intVar, 16)
		} else {
			inter[idx] = val
		}
	}
	return inter
}
