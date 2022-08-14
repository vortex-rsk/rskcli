package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"rskcli/src/cmd/rpc"
	"rskcli/src/utils"
	"strconv"
	"strings"
	"time"

	_ "rskcli/src/cmd/rpc"

	"github.com/decred/dcrd/dcrec/secp256k1"
	"golang.org/x/crypto/sha3"
)

var availableTargets map[string]*bool

func main() {

	args := os.Args[1:]
	validateCliArgs(args)

	utils.LoadConfig()

	clean := flag.Bool("clean", false, "-clean")
	jsonreq := flag.Bool("jsonreq", false, "-jsonreq")
	jsonresp := flag.Bool("jsonresp", false, "-jsonresp")
	json := flag.Bool("json", false, "-json")
	nofooter := flag.Bool("nofooter", false, "-nofooter")

	setUpServerFlags()
	flag.Parse()

	// needs to be done after the flag.parse
	globalFlags := &map[string]bool{"clean": *clean, "jsonreq": *jsonreq, "jsonresp": *jsonresp, "json": *json, "nofooter": *nofooter}

	serverName := evaluateTargetServer()

	var serverAddr = utils.Config.GetString("servers." + serverName + ".address")
	var serverPort = utils.Config.GetString("servers." + serverName + ".port")
	var serverUrl = serverAddr + ":" + serverPort

	var commandArgs []string
	for idx := range args {
		if !strings.HasPrefix(args[idx], "-") {
			commandArgs = args[idx:]
			break
		}
	}

	rpc.AddToContext("serverName", serverName)
	rpc.AddToContext("serverUrl", serverUrl)
	rpc.AddArgsToContext(commandArgs)
	rpc.AddFlags(globalFlags)

	rpc.Handle(commandArgs)

	//generateRandomPrivateKey()

}

func validateCliArgs(args []string) {
	if len(args) == 0 {
		rpc.PrintHelp()
		os.Exit(0)
	} else if len(args) == 1 && strings.HasPrefix(args[0], "-") {
		var hasCmd bool = false
		for idx := range args {
			if strings.HasPrefix(args[idx], "-") {
				hasCmd = true
				break
			}
		}
		if !hasCmd {
			fmt.Printf("Please provide command\n")
			fmt.Printf("Example: rsk help\n")
			os.Exit(1)
		}
	}
}

func evaluateTargetServer() string {
	var server string
	for key, value := range availableTargets {
		if *value {
			server = key
			break
		}
	}

	if len(server) == 0 {
		server = utils.Config.GetString("default")
	}
	if len(server) == 0 {
		panic("No server provided and no default defined in the config file.")
	}

	return server
}

func setUpFlags() *map[string]*bool {
	ret := make(map[string]*bool)
	ret["clean"] = flag.Bool("clean", false, "-clean")
	ret["jsonreq"] = flag.Bool("jsonreq", false, "-jsonreq")
	ret["jsonresp"] = flag.Bool("jsonresp", false, "-jsonresp")
	ret["json"] = flag.Bool("json", false, "-json")

	return &ret
}

func setUpServerFlags() {
	servers := utils.Config.GetStringMap("servers")

	availableTargets = make(map[string]*bool)

	for key, _ := range servers {
		availableTargets[key] = flag.Bool(key, false, "-"+key)
	}
}

func generateRandomHexaByte() string {

	rand.Seed(time.Now().UnixNano())

	hex := strconv.FormatInt(rand.Int63n(255), 16)

	for len(hex) < 2 {
		hex = "0" + hex
	}

	return hex
}

func generateRandomPrivateKey() /*(string, string)*/ {

	keyStr := ""

	for k := 0; k < 32; k++ {
		keyStr += generateRandomHexaByte()
	}

	pkBytes, err := hex.DecodeString(keyStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	privKey, pubKey := secp256k1.PrivKeyFromBytes(pkBytes)

	fmt.Println(privKey)
	fmt.Println(pubKey)

	var buf []byte

	hash := sha3.NewLegacyKeccak256()
	hash.Write((*pubKey).Serialize())
	buf = hash.Sum(nil)

	address := hex.EncodeToString(buf[len(buf)-20:])

	fmt.Println(address)

	//return
	//	hex.EncodeToString((*pubKey).Serialize()),
	//	address

}
