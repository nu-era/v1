package main

import (
	"github.com/TriviaRoulette/servers/notify/handlers"
	"log"
	"net/http"
	"os"
)

// microservice that handles websockets for an api
func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":8080"
	}

	rmq := os.Getenv("RABBITMQ")

	hc := handlers.NotifyContext{
		Sockets: handlers.NewSocketStore(),
	}

	// connect to RabbitMQ
	events, err := hc.Sockets.ConnectQueue(rmq)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ, %v", err)
	}

	// start go routine to read/send event/message notifications
	// to sockets
	go hc.Sockets.Read(events)

	http.HandleFunc("/v1/ws", hc.WebSocketConnectionHandler)

	log.Printf("Server is listening at http:/trivia/%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
