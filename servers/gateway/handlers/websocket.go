package handlers

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

//TODO: add a handler that upgrades clients to a WebSocket connection
//and adds that to a list of WebSockets to notify when events are
//read from the RabbitMQ server. Remember to synchronize changes
//to this list, as handlers are called concurrently from multiple
//goroutines.

//TODO: start a goroutine that connects to the RabbitMQ server,
//reads events off the queue, and broadcasts them to all of
//the existing WebSocket connections that should hear about
//that event. If you get an error writing to the WebSocket,
//just close it and remove it from the list
//(client went away without closing from
//their end). Also make sure you start a read pump that
//reads incoming control messages, as described in the
//Gorilla WebSocket API documentation:
//http://godoc.org/github.com/gorilla/websocket

// WebSocketConnectionHandler upgrades the connection to a WebSocket if the the
// user is authenticated
func (ctx *HandlerContext) WebSocketConnectionHandler(w http.ResponseWriter, r *http.Request, deviceID bson.ObjectId) error {
	// TODO: Check origin before upgrade

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to open websocket connection", 401)
		return err
	}
	fmt.Println("Opened Websocket Connection")

	ctx.WsConnections.Add(deviceID, conn)
	go (func(conn *websocket.Conn, deviceID bson.ObjectId) {
		defer conn.Close()
		defer ctx.WsConnections.Remove(deviceID)

		for {
			messageType, _, err := conn.ReadMessage()

			if messageType == CloseMessage {
				fmt.Println("Close message received.")
				// TODO: SEND MSG TO DEVICE OWNER VIA OTHER MEDIUM
				// TWILIO OR EMAIL
				break
			} else if err != nil {
				fmt.Println("Error reading message.")
				// TODO: SEND MSG TO DEVICE OWNER VIA OTHER MEDIUM
				// TWILIO OR EMAIL
				break
			}

		}
	})(conn, deviceID)
	return nil
}
