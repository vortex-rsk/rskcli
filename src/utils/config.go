package utils

import (
	"github.com/gurkankaymak/hocon"
	"log"
	"os"
)

var Config *hocon.Config

func LoadConfig() {
	var err error
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	Config, err = hocon.ParseResource(dirname + "/rskcli.conf")
	if err != nil {
		log.Fatal("error while parsing configuration: ", err)
	}
}
