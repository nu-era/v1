package main

import (
	"github.com/New-Era/servers/notify/handlers"
	"log"
	"net/http"
	"os"
)

// microservice that handles websockets for an api
func main() {
	tlscert := os.Getenv("TLSCERT")
	tlskey := os.Getenv("TLSKEY")

	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":8000"
	}

	if len(tlscert) == 0 {
		log.Fatal("No TLSCERT variable specified, exiting...")
	}
	if len(tlskey) == 0 {
		log.Fatal("No TLSKEY variable specified, exiting...")
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

	mux := http.NewServeMux()
	//wrappedMux := handlers.NewWsMiddleware(mux)
	mux.HandleFunc("/ws", hc.WebSocketConnectionHandler)

	log.Printf("https Server is listening on port %s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlscert, tlskey, mux))
}
