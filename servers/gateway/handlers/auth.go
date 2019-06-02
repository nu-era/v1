package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/New-Era/servers/gateway/models/devices"
	"github.com/New-Era/servers/gateway/sessions"
	//webpush "github.com/SherClockHolmes/webpush-go"
)

// DevicesHandler handles the creation (POST) of new devices. Primary key is stored using the name
// of the device. A standard should be observed.
func (ctx *HandlerContext) DevicesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" && r.Method != "PATCH" {
		http.Error(w, "method must be Post or PATCH", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get(headerContentType) != contentTypeJSON {
		http.Error(w, "content type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	if r.Method == "POST" {
		newDevice := &devices.NewDevice{}
		if err := json.NewDecoder(r.Body).Decode(newDevice); err != nil {
			http.Error(w, fmt.Sprintf("error decoding JSON: %v", err), http.StatusBadRequest)
			return
		}

		device, err := newDevice.ToDevice()
		if err != nil {
			http.Error(w, fmt.Sprintf("error creating new device: %v", err), http.StatusBadRequest)
			return
		}

		if _, err := ctx.deviceStore.GetByName(device.Name); err == nil {
			http.Error(w, fmt.Sprintf("device name already exists"), http.StatusBadRequest)
			return
		}

		device, err = ctx.deviceStore.Insert(device)
		if err != nil {
			http.Error(w, fmt.Sprintf("error adding device: %v", err), http.StatusInternalServerError)
			return
		}

		newSession := &SessionState{StartTime: time.Now(), Device: device}
		if _, err := sessions.BeginSession(ctx.SigningKey, ctx.SessStore, newSession, w); err != nil {
			http.Error(w, fmt.Sprintf("error creating new session: %v", err), http.StatusInternalServerError)
			return
		}
		Verify("+1"+device.Phone, trialNum, "sup")
		// send public key for push notifications as header
		w.Header().Add("X-VapidKey", ctx.PubVapid)
		respond(w, device, http.StatusCreated)
	} else {
		type VerificationCheck struct {
			Code  string `json:"code"`
			Phone string `json:"phone"`
		}

		newVerificationCheck := &VerificationCheck{}
		if err := json.NewDecoder(r.Body).Decode(newVerificationCheck); err != nil {
			http.Error(w, fmt.Sprintf("error decoding JSON into VerificationCheck: %v", err), http.StatusBadRequest)
			return
		}

		err := CheckVerification(newVerificationCheck.Code, newVerificationCheck.Phone)
		if err != nil {
			http.Error(w, fmt.Sprintf("error sending VerificationCheck: %v", err), 500)
		}

	}

}

// SpecificDeviceHandler handles request for a specific device, requiring a prexisting sessions.
// GET returns device info with given device name
// PATCH updates the device name
func (ctx *HandlerContext) SpecificDeviceHandler(w http.ResponseWriter, r *http.Request) {
	sessionState := &SessionState{}
	sessID, err := sessions.GetState(r, ctx.SigningKey, ctx.SessStore, sessionState)
	if err != nil {
		http.Error(w, fmt.Sprintf("problem with session %v", err), http.StatusUnauthorized)
		return
	}
	deviceID := sessionState.Device.ID
	switch r.Method {
	case "GET":
		device, err := ctx.deviceStore.GetByID(deviceID)
		if err != nil {
			http.Error(w, fmt.Sprintf("device not found: %v", err), http.StatusNotFound)
			return
		}
		w.Header().Add("X-VapidKey", ctx.PubVapid)
		respond(w, device, http.StatusOK)
	case "PATCH":
		// segment := path.Base(r.URL.Path)
		// if segment != "me" || segment != deviceID.Hex() {
		// 	http.Error(w, "unathorized", http.StatusUnauthorized)
		// 	return
		// }
		if r.Header.Get(headerContentType) != contentTypeJSON {
			http.Error(w, "content type must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		device, err := ctx.deviceStore.GetByID(deviceID)
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting device: %v", err), http.StatusNotFound)
		}

		deviceUpdates := &devices.Updates{}
		if err := json.NewDecoder(r.Body).Decode(deviceUpdates); err != nil {
			http.Error(w, fmt.Sprintf("error decoding JSON: %v", err), http.StatusBadRequest)
			return
		}

		if err = device.ApplyUpdates(deviceUpdates); err != nil {
			http.Error(w, fmt.Sprintf("error updating device: %v", err), http.StatusBadRequest)
			return
		}

		if err = ctx.deviceStore.Update(deviceID, device); err != nil {
			http.Error(w, fmt.Sprintf("error updating device: %v", err), http.StatusBadRequest)
			return
		}

		sessionState.Device = device
		if err := ctx.SessStore.Save(sessID, sessionState); err != nil {
			http.Error(w, fmt.Sprintf("error updating session state: %v", err), http.StatusInternalServerError)
			return
		}
		respond(w, device, http.StatusOK)

	default:
		http.Error(w, "method must be GET or PATCH", http.StatusMethodNotAllowed)
		return
	}
}

//SessionsHandler handels requests for the sessions resource. POST allows for logging in with creditials.
func (ctx *HandlerContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method must be POST", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get(headerContentType) != contentTypeJSON {
		http.Error(w, "content type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	deviceCredentials := &devices.Credentials{}
	if err := json.NewDecoder(r.Body).Decode(deviceCredentials); err != nil {
		http.Error(w, fmt.Sprintf("error decoding JSON: %v", err), http.StatusBadRequest)
		return
	}
	device, err := ctx.deviceStore.GetByName(deviceCredentials.Name)
	if device == nil { //Set dummy device
		device = &devices.Device{}
	}
	err = device.Authenticate(deviceCredentials.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid credentials"), http.StatusUnauthorized)
		return
	}
	newSession := &SessionState{StartTime: time.Now(), Device: device}
	if _, err := sessions.BeginSession(ctx.SigningKey, ctx.SessStore, newSession, w); err != nil {
		http.Error(w, fmt.Sprintf("error creating new session: %v", err), http.StatusInternalServerError)
		return
	}

	// ip := r.RemoteAddr
	// ips := strings.Split(r.Header.Get(headerXForwarded), ", ")
	// if len(ips) > 0 {
	// 	ip = ips[0]
	// }
	// if err := ctx.usersStore.Log(user.ID, ip); err != nil {
	// 	http.Error(w, fmt.Sprintf("error recording ip of user: %v", err), http.StatusBadRequest)
	// 	return
	// }
	respond(w, device, http.StatusCreated)
}

//SpecificSessionHandler handles requests related to a specific authenticated session
func (ctx *HandlerContext) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "method must be DELETE", http.StatusMethodNotAllowed)
		return
	}
	// if segment := path.Base(r.URL.Path); segment != "mine" {
	// 	http.Error(w, "unknown path segment", http.StatusForbidden)
	// 	return
	// }
	if _, err := sessions.EndSession(r, ctx.SigningKey, ctx.SessStore); err != nil {
		http.Error(w, fmt.Sprintf("session not found: %v", err), http.StatusBadRequest)
		return
	}
	fmt.Println("signed out")
	respond(w, nil, http.StatusOK)
	return
}

//SubscriptionHandler handles requests for users wanting to receive push notifications
func (ctx *HandlerContext) SubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// requires auth header to get device trying to subscribe
		sessionState := &SessionState{}
		_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessStore, sessionState)
		if err != nil {
			http.Error(w, fmt.Sprintf("problem with session %v", err), http.StatusUnauthorized)
			return
		}
		dev := sessionState.Device
		if err := json.NewDecoder(r.Body).Decode(dev.Subscription); err != nil {
			http.Error(w, fmt.Sprintf("error decoding JSON: %v", err), http.StatusBadRequest)
			return
		}
		// user subsribed successfully
		respond(w, nil, http.StatusCreated)
	} else {
		http.Error(w, "Method Not Allowed", 405)
	}

}

//respond responds with the status, content type of JSON and encoded value
func respond(w http.ResponseWriter, value interface{}, status int) {
	w.WriteHeader(status)
	if value != nil {
		w.Header().Add(headerContentType, contentTypeJSON)
		if err := json.NewEncoder(w).Encode(value); err != nil {
			http.Error(w, fmt.Sprintf("error encoding response value to JSON: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
