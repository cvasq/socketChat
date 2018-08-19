package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Currency struct {
	Last float64 `json:"last"`
}

type CurrencyList struct {
	USD Currency `json:"USD"`
}

func getPrice() string {

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

	return lastPrice
}

func (h *SocketChat) subscribeLiveTransactions() {

	go func() {
		ticker := time.NewTicker(time.Second * 60)
		for _ = range ticker.C {
			message := Message{}
			message.Type = "bot-message"
			message.Username = "BTC Bot"
			message.Time = currentTime()
			message.Data = getPrice()

			h.broadcast <- message

		}
	}()
}
