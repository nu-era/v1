package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/New-Era/servers/gateway/models/devices"
	"net/http"
	"strings"
)

// Check request for 'Content-Type' header equal to the passed content type. If
// not the correct content-type, returns error and writes 415 status response
// to writer.
func contentTypeCheck(w http.ResponseWriter, r *http.Request, contentT string) error {
	contentType := r.Header.Get(headerContentType)
	typeList := strings.Split(contentType, ",") // Get the first content type
	if len(typeList) == 0 || typeList[0] != contentT {
		http.Error(w, "Incorrect Content Type", http.StatusUnsupportedMediaType)
		return fmt.Errorf("Incorrect Content-Type")
	}
	return nil
}

// Post: registering a new device, create new session with device
// Get: Create a new session for a registered device

func (ctx *HandlerContext) DevicesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" && r.Method != "GET" {
		http.Error(w, "method must be Post or Get", http.StatusMethodNotAllowed)
		return
	}

	if r.Method == http.MethodPost {
		// Check to make sure content-type is application/json
		err := contentTypeCheck(w, r, contentTypeJSON)
		if err != nil {
			return
		}

		newDevice := devices.NewDevice{} // Create an empty NewDevice struct to be filled by request body
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&newDevice); err != nil {
			http.Error(w, "Request body unable to be decoded to new device", 400)
			return
		}

		validDevice, err := newDevice.ToDevice() // Turn new device to a Device type
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		// Insert the new device to db, and get the newly inserted device
		// with the new database-assigned primary key value
		validDevice, err = ctx.DeviceStore.Insert(validDevice)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		err = ctx.WebSocketConnectionHandler(w, r, validDevice.ID)
		if err != nil {
			return
		}
		conn := ctx.WsConnections.Conns[validDevice.ID]
		go heartbeat(conn)

		// Respond to the client with an http.StatusCreated code, and
		// the json encoded new user profile
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(201)
		err = json.NewEncoder(w).Encode(validDevice)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}
