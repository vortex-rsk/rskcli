# rskcli
CLI for rsk

Available commands:

* clientVersion
* sha3 <arg>
* version
* listening
* peerCount
* peerList
* chainId
* ethProtocolVersion
* protocolVersion
* syncing
* coinbase
* mining
* hashrate
* gasPrice
* accounts
* blockNumber
* balance
* transactionCount
* blockTransactionCountByHash
* blockTransactionCountByNumber
* uncleCountByBlockHash
* uncleCountByBlockNumber
* code
* blockByHash
* blockByNumber
* transactionByHash
* transactionByBlockHashAndIndex
* transactionByBlockNumberAndIndex

## Example
```bash
rsk peerList
```

## Flags

### -serverName

The `rskcli.conf` has a list of target nodes for the CLI to use.
If you want que query a different target then what's configures as default you can use the name of the target server as a flag:

```bash
rsk -public version
```

This will target the public mainnet

### -clean
Return the result without the footer and the final breakline. Useful to put the result inside other commands.

```bash
rsk blockByNumber `rsk -clean blockNumber` false
```
### -jsonreq
Shows the json request to be sent to the RSK node
### -jsonresp
Shows the json response received from the RSK node
### -json
Returns only the json received from the RSK node
```bash
rsk blockNumber 
```
Results in:
```json
{"jsonrpc":"2.0","id":"74","result":"0x457a71"}
```
## Install

1. Download the binary for your platform in the [releases page](https://github.com/vortex-rsk/rskcli/releases)
put it on your path 
2. Rename the binary to `rsk`
3. Put the file rskcli.conf from the root of the repo on your home folder
4. Be happy
