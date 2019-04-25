package handlers

import (
	"github.com/streadway/amqp"
	mgo "gopkg.in/mgo.v2"
)

const (
	// name of rabbitmq queue to use for services
	ShakeAlert = "api"     // comes from shakealert
	NewEra     = "devices" // want to go to devices
)

// QueueContext is a receiver that stores
// information for the queue microservice
type QueueContext struct {
	// Database for accessing devices
	MongoDB *mgo.Session

	// channel for publishing messages
	Channel *amqp.Channel
}
