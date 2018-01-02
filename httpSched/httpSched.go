package httpSched

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jorgenogales/bitcalypse/config"
)

var cfg config.Config = config.Get()

type HttpFeed struct {
	Url      string
	Interval int
	Ch       chan string
	ticker   *time.Ticker
}

func NewHttpFeed(url string, interval int, c chan string) HttpFeed {
	log.Printf("New HttpFeed setup in %s with %d interval", url, interval)
	h := HttpFeed{
		url,
		interval,
		c,
		time.NewTicker(time.Second * time.Duration(interval))}
	h.Start()
	return h
}

func (h *HttpFeed) Start() {
	go func() {
		for range h.ticker.C {
			h.Ch <- httpGetRequest(h.Url)
		}
	}()

}

func (h *HttpFeed) Stop() {
	h.ticker.Stop()
}

func httpGetRequest(url string) string {
	client := http.Client{
		Timeout: time.Duration(cfg.HttpTimeout) * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Error in HTTP call: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error in reading HTTP response body: %s. Body: %s", err, body)
	}
	return string(body)
}
