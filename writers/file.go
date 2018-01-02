package writers

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/jorgenogales/bitcalypse/brokers"
)

func WriteMsg(path string, msg brokers.MarketMsg) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		ioutil.WriteFile(path, []byte(msg.CSVHeader()), 0644)
	} else {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		if err != nil {
			log.Fatal("Error while opening file: ", err)
		}
		f.Write([]byte(msg.CSVData()))
	}
}
