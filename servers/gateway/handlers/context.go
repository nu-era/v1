package handlers

import (
	"github.com/New-Era/servers/gateway/models/devices"
)

// HandlerContext tracks the key that is used to sign and
// validate SessionIDs, the sessions.Store, and the users.Store
type HandlerContext struct {
	DeviceStore   devices.Store
	WsConnections *Connections
}

//NewHandlerContext constructs a new HandlerContext,
//ensuring that the dependencies are valid values
func NewHandlerContext(deviceStore devices.Store, connections *Connections) *HandlerContext {
	if deviceStore == nil || connections == nil {
		panic("Parameters may not be empty!")
	}
	return &HandlerContext{deviceStore, connections}
}
