package main

import (
	"encoding/json"
	"fmt"
	"github.com/New-Era/servers/goqueue/handlers"
	"github.com/streadway/amqp"
	mgo "gopkg.in/mgo.v2"
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

	ch, err := connectQueue(os.Getenv("RABBITMQ"))
	if err != nil {
		log.Fatalf("error connecting to queue, %v", err)
	}

	mongoAddr := os.Getenv("MONGO_ADDR")
	mongoSess, err := mgo.Dial(mongoAddr)
	if err != nil {
		log.Fatalf("error dialing mongo: %v", err)
	}

	hc := handlers.QueueContext{
		MongoDB: mongoSess,
		Channel: ch,
	}

	go hc.Routine()
	log.Printf("Server is listening on port %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// connectQueue connects to the RabbitMQ service at the address defined in the addr variable
// and creates a channel and queue to send messages to. It returns the go channel
// which contains messages living on the RabbitMQ queue. Errors are returned if the
// connection fails
func connectQueue(addr string) (*amqp.Channel, error) {
	con, err := amqp.Dial("amqp://" + addr)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to MQ, %v", err)
	}

	chann, err := con.Channel()
	if err != nil {
		return nil, fmt.Errorf("error creating channel, %v", err)
	}
	return chann, nil
}
