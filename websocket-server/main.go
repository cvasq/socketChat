package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {

	// Set custom port by running with --port PORT_NUM
	// Default port is 8000
	httpPort := flag.String("port", "80", "WebSocket Listening Address")
	flag.Parse()

	socketChat := createSocketChat()

	http.HandleFunc("/ws", socketChat.websocketHandler)

	go socketChat.subscribeLiveTransactions()
	go socketChat.trackActiveClients()
	go socketChat.handleMessages()

	log.Println("Starting SocketChat Server")
	log.Println("Listening on port: ", *httpPort)
	err := http.ListenAndServe(":"+*httpPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
