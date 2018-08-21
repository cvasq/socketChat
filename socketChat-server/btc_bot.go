package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const priceCheckInterval = 60

type Currency struct {
	Last float64 `json:"last"`
}

type CurrencyList struct {
	USD Currency `json:"USD"`
}

func getPrice() Message {

	priceMessage := Message{}
	priceMessage.Type = "bot-message"
	priceMessage.Username = "BTC Bot"
	priceMessage.Time = currentTime()

	resp, err := http.Get("https://blockchain.info/ticker")
	if err != nil {
		log.Println(err)
		return priceMessage
	}

	currentPrice := CurrencyList{}
	err = json.NewDecoder(resp.Body).Decode(&currentPrice)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
	}

	lastPrice := fmt.Sprintf("The current Bitcoin price is $%.2f USD", currentPrice.USD.Last)

	priceMessage.Data = lastPrice

	return priceMessage
}

func (h *SocketChat) runBtcBot() {

	go func() {
		time.Sleep(time.Second * 1)
		botListing := Message{
			Type:     "bot-listing",
			Username: "system",
			Time:     currentTime(),
			Data:     "btc-bot",
		}
		log.Println("Sending bot listing message")
		h.broadcast <- botListing

	}()
	// Sends Time (tick) to channel every X seconds
	tickChan := time.NewTicker(time.Second * priceCheckInterval).C

	for {
		select {
		case <-tickChan:
			h.broadcast <- getPrice()
		case b := <-h.triggerBot:
			if b == "btc-bot" {
				h.broadcast <- getPrice()
			}
		}
	}
}
