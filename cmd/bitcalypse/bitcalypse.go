package main

import (
	"fmt"

	"github.com/jorgenogales/bitcalypse/brokers"
	"github.com/jorgenogales/bitcalypse/config"
	"github.com/jorgenogales/bitcalypse/writers"
)

func main() {
	bs := brokers.NewBitstamp()
	cfg := config.Get()
	tab := make(map[chan brokers.MarketMsg]string, len(cfg.Tickers))
	for i := range cfg.Tickers {
		c := make(chan brokers.MarketMsg)
		tab[c] = cfg.Tickers[i]
		bs.GetTickerUpdates(cfg.Tickers[i], c)
		go writeUpdates(fmt.Sprintf("%s/%s.csv", cfg.DumpFilePath, cfg.Tickers[i]), c)
	}

	checkExit()
}

func writeUpdates(filepath string, c chan brokers.MarketMsg) {
	for m := range c {
		writers.WriteMsg(filepath, m)
	}
}

func checkExit() {
	var input string
	fmt.Scanln(&input)
	for input != ".-." {
		fmt.Println("Input \".-.\"  if exit is what you want")
		fmt.Scanln(&input)
	}
	fmt.Println("Thanks for using Bitcalypse")
}
