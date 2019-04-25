package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
	"sync"
)

// SocketStore contains client connection information
// and a queue channel for sending notifications
type SocketStore struct {
	Connections map[bson.ObjectId]*websocket.Conn
	lock        sync.Mutex
	Chan        *amqp.Channel
}

// NewSocketStore returns a new socket store containing a map of device id's to
// a websocket, a mutex lock for concurrent use and a queue channel for real time
// notifications
func NewSocketStore() *SocketStore {
	return &SocketStore{Connections: map[bson.ObjectId]*websocket.Conn{}}
}

// Control messages for websocket
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

	// name of rabbitmq queue to use for notifications
	qName = "services"
)

// InsertConnection is a Thread-safe method for inserting a connection
func (s *SocketStore) InsertConnection(id bson.ObjectId, conn *websocket.Conn) {
	s.lock.Lock()
	// insert socket connection
	s.Connections[id] = conn
	s.lock.Unlock()
}

// RemoveConnection is a Thread-safe method for removing a connection
func (s *SocketStore) RemoveConnection(id bson.ObjectId) {
	s.lock.Lock()
	_, ok := s.Connections[id]
	if ok {
		delete(s.Connections, id)
	}
	s.lock.Unlock()
}

// WriteToValidConnections sends messages to a subset of connections
// (if the message is intended for a private channel), or to all of them (if the message
// is posted on a public channel
func (s *SocketStore) WriteToValidConnections(deviceIDs []bson.ObjectId, messageType int, data []byte) error {
	fmt.Printf("Number of devices to send to: %d", len(deviceIDs))
	var writeError error
	if len(deviceIDs) > 0 { // private channel
		for _, id := range deviceIDs {
			writeError = s.Connections[id].WriteMessage(messageType, data)
			if writeError != nil {
				return writeError
			}
		}
	} else { // public channel
		for _, conn := range s.Connections {
			writeError = conn.WriteMessage(messageType, data)
			if writeError != nil {
				return writeError
			}
		}
	}

	return nil
}

// Message is a struct to read our message into
type Message struct {
	Type      string                 `json:"type"`
	Channel   map[string]interface{} `json:"channel,omitempty"`
	ChannelID int64                  `json:"channelID,omitempty"`
	Message   map[string]interface{} `json:"message,omitempty"`
	MessageID int64                  `json:"messageID,omitempty"`
	DeviceIDs []int64                `json:"deviceIDs,omitempty"`
}

// upgrader is a variable that stores websocket information and verifies
// the origin of the client request to authenticate
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		orig := r.Header.Get("Origin")
		if strings.Contains(orig, "bfranzen.me") {
			return true
		}
		return false
	},
}

// WebSocketConnectionHandler handles when the client requests an upgrade to a websocket
// if the device is valid (request comes from proper host, device exists) upgrade is performed
// and connection is stored for duration of client session
func (hc *NotifyContext) WebSocketConnectionHandler(w http.ResponseWriter, r *http.Request) {
	r = addWebsocketHeaders(r)
	fmt.Println("UPGRADING TO WEBSOCKET")
	if r.Header.Get("X-Device") == "" {
		http.Error(w, "Unauthorized Access", 401)
		return
	}
	var dest map[string]interface{}
	if err := json.Unmarshal([]byte(r.Header.Get("X-Device")), &dest); err != nil {
		fmt.Printf("error getting message body, %v", err)
		http.Error(w, "Bad Request", 400)
		return
	}

	// handle the websocket handshake
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("ERROR: vvv")
		fmt.Println(err)
		fmt.Printf("ERROR OPENING CONNECTION: %v", err)
		http.Error(w, "Failed to open websocket connection", 401)
		return
	}

	fmt.Println("CONNECTION UPGRADED")
	// Insert our connection onto our datastructure for ongoing usage
	hc.Sockets.InsertConnection(bson.ObjectIdHex(dest["ID"].(string)), conn)
	// Invoke a goroutine for handling control messages from this connection
	fmt.Println("CONNECTION INSERTED")
	fmt.Println(r)
	go (func(conn *websocket.Conn, deviceID bson.ObjectId) {
		defer conn.Close()
		defer hc.Sockets.RemoveConnection(deviceID)

		for {
			messageType, p, err := conn.ReadMessage()
			if len(p) > 0 {
				var j map[string]interface{}
				if err := json.Unmarshal(p, &j); err != nil {
					fmt.Printf("error unmarshaling json: %v", err)
				}
			}

			if messageType == CloseMessage {
				fmt.Println("Close message received...")
				break
			} else if err != nil {
				fmt.Printf("error connecting: %v, closing...", err)
				break
			}
			// ignore ping and pong messages
		}

	})(conn, bson.ObjectIdHex(dest["ID"].(string)))
}

// ConnectQueue connects to the RabbitMQ service at the address defined in the addr variable
// and creates a channel and queue to send/receive messages to. It returns the go channel
// which contains messages living on the RabbitMQ queue. Errors are returned if the
// connection fails
func (s *SocketStore) ConnectQueue(addr string) (<-chan amqp.Delivery, error) {
	con, err := amqp.Dial("amqp://" + addr)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to MQ, %v", err)
	}

	chann, err := con.Channel()
	if err != nil {
		return nil, fmt.Errorf("error creating channel, %v", err)
	}

	s.Chan = chann

	queue, err := chann.QueueDeclare(qName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("error declaring queue, %v", err)
	}

	events, err := chann.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("error retreiving messages, %v", err)
	}
	return events, nil
}

// Read reads events off the passed go channel created by the ConnectQueue method
// and sends the messages to the proper websockets in the SocketStore
func (s *SocketStore) Read(events <-chan amqp.Delivery) {
	for e := range events {
		var event map[string]interface{}
		if err := json.Unmarshal(e.Body, &event); err != nil {
			fmt.Printf("error getting message body, %v", err)
			break
		}
		if event["deviceIDs"] != nil {
			ids := make([]bson.ObjectId, len(event["deviceIDs"].([]interface{})))
			for i, v := range event["deviceIDs"].([]interface{}) {
				ids[i] = v.(bson.ObjectId)
			}
			s.WriteToValidConnections(ids, TextMessage, e.Body)
		} else {
			s.WriteToValidConnections([]bson.ObjectId{}, TextMessage, e.Body)
		}

	}
}

func addWebsocketHeaders(r *http.Request) *http.Request {
	r.Header.Set("Connection", r.Header.Get("X-Connection"))
	r.Header.Set("Upgrade", r.Header.Get("X-Upgrade"))
	r.Header.Set("Sec-Websocket-Key", r.Header.Get("X-Sec-Websocket-Key"))
	return r
}
