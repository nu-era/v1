package handlers

import (
	"Capstone/New-Era/servers/gateway/models/devices"
)

// HandlerContext tracks the key that is used to sign and
// validate SessionIDs, the sessions.Store, and the users.Store
type HandlerContext struct {
	DeviceStore devices.Store `json:"mongoStore"`
}

//NewHandlerContext constructs a new HandlerContext,
//ensuring that the dependencies are valid values
func NewHandlerContext(deviceStore devices.Store) *HandlerContext {
	if deviceStore == nil {
		panic("Parameters may not be empty!")
	}
	return &HandlerContext{deviceStore}
}
