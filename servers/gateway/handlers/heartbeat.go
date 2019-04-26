package handlers

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func heartbeat(conn *websocket.Conn) {
	fmt.Println("Beggining to send pings")
	for {
		conn.SetReadLimit(maxMessageSize)
		conn.SetReadDeadline(time.Now().Add(pongWait))

		fmt.Println("Sending: ping.")
		err := conn.WriteMessage(websocket.TextMessage, []byte("ping"))
		if err != nil {
			fmt.Println("Write Error: ", err)
			break
		}

		msgType, bytes, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read Error: ", err)
			Send("+4254229586", trialNum, "Hey Dude")
			break
		}

		// We don't recognize any message that is not "ping".
		if msg := string(bytes[:]); msgType != websocket.TextMessage && msg != "pong" {
			fmt.Println("Unrecognized message received.")
			continue
		} else {
			fmt.Println("Received: pong.")
		}
	}
}
