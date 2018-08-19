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

	resp, err := http.Get("https://blockchain.info/ticker")
	if err != nil {
		log.Println(err)
	}

	currentPrice := CurrencyList{}
	err = json.NewDecoder(resp.Body).Decode(&currentPrice)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	lastPrice := fmt.Sprintf("The current Bitcoin price is $%.2f USD", currentPrice.USD.Last)

	priceMessage := Message{}
	priceMessage.Type = "bot-message"
	priceMessage.Username = "BTC Bot"
	priceMessage.Time = currentTime()
	priceMessage.Data = lastPrice

	return priceMessage
}

func (h *SocketChat) subscribeLiveTransactions() {

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
