package handlers

import (
	"sync"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/websocket"
)

// Connections is a map of userIds to websocket connections.
// When a new websocket connection is made, it stores that new
// connection for that user, until the connection is ended,
// in which case the connection is deleted from Connections
type Connections struct {
	Conns map[bson.ObjectId]*websocket.Conn
	mx    sync.RWMutex
}

// NewConnections creates a new Connections map
func NewConnections() *Connections {
	connections := Connections{
		Conns: map[bson.ObjectId]*websocket.Conn{},
	}
	return &connections
}

// Add adds a new connection for the given userID and connection
func (c *Connections) Add(deviceID bson.ObjectId, conn *websocket.Conn) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.Conns[deviceID] = conn
}

// Remove deletes the connection for the given userID
func (c *Connections) Remove(deviceID bson.ObjectId) {
	c.mx.Lock()
	defer c.mx.Unlock()
	delete(c.Conns, deviceID)
}
