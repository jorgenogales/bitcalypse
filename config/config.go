package config

import (
	"encoding/json"
	"log"
	"os"
)

const configFilePath = "/Users/jorgenogales/go/src/github.com/jorgenogales/bitcalypse/config/bitcalypse.json"

type Config struct {
	Bitstamp struct {
		BaseUrl   string
		TickerUri string
		Interval  int
	}
	Tickers        []string
	DumpFilePath   string
	DelayThreshold int
	HttpTimeout    int
}

func Get() Config {
	var config Config
	configFile, err := os.Open(configFilePath)
	defer configFile.Close()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
