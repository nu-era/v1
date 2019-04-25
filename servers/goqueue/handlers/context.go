package handlers

import (
	"github.com/streadway/amqp"
	mgo "gopkg.in/mgo.v2"
)

const (
	// name of rabbitmq queue to use for services
	ShakeAlert = "receive"
	NewEra     = "send"
)

// QueueContext is a receiver that stores
// information for the queue microservice
type QueueContext struct {
	// Database for accessing devices
	MongoDB *mgo.Session

	// channel for publishing messages
	Channel *amqp.Channel
}
