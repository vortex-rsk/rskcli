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

## Example
`rsk peerList`

## flags
The `rskcli.conf` has a list of target nodes for the CLI to use.
If you ant que query a different target then what's configures as default you can use the name of the target server as a flag:

`rsk -public version`

This will target the public mainnet

## Install

1. Download the binary for your platform in the [releases page](https://github.com/rskcli/rskcli/releases)
put it on your path 
2. Rename the binary to `rsk`
3. Put the file rskcli.conf from the root of the repo on your home folder
4. Be happy
