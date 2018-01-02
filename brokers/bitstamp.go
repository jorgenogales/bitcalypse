package brokers

import (
	"fmt"

	"encoding/json"
	"log"
	"time"

	"github.com/jorgenogales/bitcalypse/config"
	"github.com/jorgenogales/bitcalypse/httpSched"
)

var cfg config.Config = config.Get()

type Broker string

const (
	BitstampBroker Broker = "Bitstamp"
)

type MarketMsg struct {
	Datetime int64   `json:"timestamp,string"`
	Open     float64 `json:"open,string"`
	High     float64 `json:"high,string"`
	Low      float64 `json:"low,string"`
	Last     float64 `json:"last,string"`
	Volume   float64 `json:"volume,string"`
	Bid      float64 `json:"bid,string"`
	Ask      float64 `json:"ask,string"`
	Vwap     float64 `json:"vwap,string"`
}

type Bitstamp struct {
	baseUrl   string
	tickerUri string
	interval  int
	httpFeeds []httpSched.HttpFeed
}

func (msg MarketMsg) CSVData() string {
	return fmt.Sprintf(
		"%s,%f,%f,%f,%f,%f,%f,%f,%f\n",
		time.Unix(msg.Datetime, 0).Format(time.RFC3339), msg.Open, msg.High, msg.Low, msg.Last, msg.Volume, msg.Bid, msg.Ask, msg.Vwap)
}

func (msg MarketMsg) CSVHeader() string {
	return "Datetime,Open,High,Low,Last,Volume,Bid,Ask,vWap\n"
}
func (msg MarketMsg) String() string {
	return fmt.Sprintf(
		"\nDatetime: %s\nOpen: %f\nHigh: %f\nLow: %f\nLast: %f\nVolume: %f\nBid: %f\nAsk: %f\nvWap: %f\n\n",
		time.Unix(msg.Datetime, 0).Format(time.RFC3339), msg.Open, msg.High, msg.Low, msg.Last, msg.Volume, msg.Bid, msg.Ask, msg.Vwap)
}

func NewBitstamp() Bitstamp {
	return Bitstamp{cfg.Bitstamp.BaseUrl, cfg.Bitstamp.TickerUri, cfg.Bitstamp.Interval,
		[]httpSched.HttpFeed{}}
}

func (bs *Bitstamp) GetTickerUpdates(ticker string, cmsg chan MarketMsg) {
	cstr := make(chan string)
	bs.httpFeeds = append(
		bs.httpFeeds,
		httpSched.NewHttpFeed(
			fmt.Sprintf("%s/%s/%s", bs.baseUrl, bs.tickerUri, ticker),
			bs.interval,
			cstr))
	handler := handleTickerUpdates()
	go handler(cstr, cmsg, ticker, BitstampBroker)
}

func handleTickerUpdates() func(chan string, chan MarketMsg, string, Broker) {
	lastMsg := ""
	lastUpdate := time.Now()
	lastPrice := 0.0
	return func(cstr chan string, cmsg chan MarketMsg, ticker string, broker Broker) {
		for s := range cstr {
			if s != lastMsg {
				lastMsg = s
				var m MarketMsg
				err := json.Unmarshal([]byte(s), &m)
				if err != nil {
					log.Fatalf("Fatal error: %s", err)
				}
				// Checking that the unmarshalling worked as expected
				if m.Last > 0 {
					log.Printf("Update for %s at %s, Last price: %f",
						ticker, time.Unix(m.Datetime, 0).Format(time.RFC3339), m.Last)
					lastUpdate = time.Unix(m.Datetime, 0)
					lastPrice = m.Last
					cmsg <- m
				} else {
					log.Fatal("Wrong HTTP message received: ", s)
				}
			}
			if lastUpdate.Before(time.Now().Add(time.Second * time.Duration(-cfg.DelayThreshold))) {
				log.Printf("DELAY THRESHOLD REACHED FOR %s. Last update: %s",
					ticker, lastUpdate.Format(time.RFC3339))
			}
		}
	}
}
