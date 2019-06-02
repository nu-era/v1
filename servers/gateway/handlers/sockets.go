package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/New-Era/servers/gateway/sessions"
	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
	"gopkg.in/mgo.v2/bson"
	"reflect"
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
	qName = "devices"
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
	var writeError error
	if len(deviceIDs) > 0 && len(s.Connections) > 0 { // send to necessary users
		for _, id := range deviceIDs {
			if _, ok := s.Connections[id]; ok { // if connection exists
				writeError = s.Connections[id].WriteMessage(messageType, data)
				if writeError != nil {
					return writeError
				}
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
		if strings.Contains(orig, "bfranzen.me") || strings.Contains(orig, "127.0.0.1") {
			return true
		}
		return false
	},
}

// WebSocketConnectionHandler handles when the client requests an upgrade to a websocket
// if the device is valid (request comes from proper host, device exists) upgrade is performed
// and connection is stored for duration of client session
func (hc *HandlerContext) WebSocketConnectionHandler(w http.ResponseWriter, r *http.Request) {
	var sess SessionState
	if _, err := sessions.GetState(r, hc.SigningKey, hc.SessStore, &sess); err != nil {
		http.Error(w, "Not authorized to access resource", 401)
		return
	}

	// handle the websocket handshake
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to open websocket connection", 401)
		return
	}
	//4256148818
	// Insert our connection onto our datastructure for ongoing usage
	hc.Sockets.InsertConnection(sess.Device.ID, conn)
	// Invoke a goroutine for handling control messages from this connection
	fmt.Println("CONNECTION INSERTED")

	// Get phone number to send twilio messages to
	dev, _ := hc.deviceStore.GetByID(sess.Device.ID)
	heartbeat(conn, dev.Email)
	go (func(conn *websocket.Conn, deviceID bson.ObjectId) {
		defer conn.Close()
		defer hc.Sockets.RemoveConnection(deviceID)

		for {
			messageType, p, conErr := conn.ReadMessage()
			if len(p) > 0 {
				var j map[string]interface{}
				if err := json.Unmarshal(p, &j); err != nil {
					fmt.Printf("error unmarshaling json: %v", err)
				}
			}

			if messageType == CloseMessage {
				fmt.Println("Close message received...")
				break
			} else if conErr != nil {
				fmt.Printf("error connecting: %v, closing...", err)
				break
			}
			// ignore ping and pong messages
		}

	})(conn, sess.Device.ID)
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
func (s *SocketStore) Read(events <-chan amqp.Delivery, ctx *HandlerContext) {
	for e := range events {
		var event map[string]interface{}
		if err := json.Unmarshal(e.Body, &event); err != nil {
			fmt.Printf("error getting message body, %v", err)
			break
		}
		if event["deviceIDs"] != nil {
			ids := make([]bson.ObjectId, len(event["deviceIDs"].([]interface{})))
			for i, v := range event["deviceIDs"].([]interface{}) {
				//ids[i] = v.(bson.ObjectId)
				ids[i] = bson.ObjectIdHex(v.(string))
			}
			s.WriteToValidConnections(ids, TextMessage, e.Body)
			ctx.Push(ids, e.Body)
		} else {
			//s.WriteToValidConnections([]bson.ObjectId{}, TextMessage, e.Body)
		}
	}
}

func (ctx *HandlerContext) Push(ids []bson.ObjectId, data []byte) {
	for _, id := range ids {
		// get user device
		dev, err := ctx.deviceStore.GetByID(id)
		if err != nil {
			fmt.Errorf("error retrieving device, %v", err)
		}
		// check that user is subscribed
		if reflect.ValueOf(dev.Subscription).IsNil() {
			// Send Notification
			_, err := webpush.SendNotification(data, &dev.Subscription, &webpush.Options{
				Subscriber:      dev.Email,
				VAPIDPublicKey:  ctx.PubVapid,
				VAPIDPrivateKey: ctx.PriVapid,
				TTL:             30,
			})
			if err != nil {
				fmt.Errorf("error sending notification, %v", err)
			}
		}
	}
}
