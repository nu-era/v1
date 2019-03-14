package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/New-Era/servers/gateway/models/devices"
)

// Check request for 'Content-Type' header equal to 'application/json'. If
// not correct content-type, returns error and writes 415 status response
// to writer.
func contentTypeCheck(w http.ResponseWriter, r *http.Request) error {
	contentType := r.Header.Get("Content-Type")
	typeList := strings.Split(contentType, ",") // Get the first content type
	if len(typeList) == 0 || typeList[0] != "application/json" {
		http.Error(w, "Request body must be of type JSON", 415)
		return fmt.Errorf("Incorrect Content Type")
	}
	return nil
}

func (ctx *HandlerContext) DevicesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Check to make sure content-type is application/json
		err := contentTypeCheck(w, r)
		if err != nil {
			return
		}

		newDevice := devices.NewDevice{} // Create a empty NewDevice struct to be filled by request body
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
