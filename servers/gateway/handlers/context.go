package handlers

import (
<<<<<<< HEAD
	"Capstone/New-Era/servers/gateway/models/devices"
=======
	"github.com/New-Era/servers/gateway/models/devices"
>>>>>>> 47550d9798a8f2b1d80770d26df1b466dbc525b1
)

// HandlerContext tracks the key that is used to sign and
// validate SessionIDs, the sessions.Store, and the users.Store
type HandlerContext struct {
<<<<<<< HEAD
	DeviceStore devices.Store `json:"mongoStore"`
=======
	DeviceStore   devices.Store `json:"mongoStore"`
	WsConnections *Connections  `json:"connections"`
>>>>>>> 47550d9798a8f2b1d80770d26df1b466dbc525b1
}

//NewHandlerContext constructs a new HandlerContext,
//ensuring that the dependencies are valid values
<<<<<<< HEAD
func NewHandlerContext(deviceStore devices.Store) *HandlerContext {
	if deviceStore == nil {
		panic("Parameters may not be empty!")
	}
	return &HandlerContext{deviceStore}
=======
func NewHandlerContext(deviceStore devices.Store, connections *Connections) *HandlerContext {
	if deviceStore == nil || connections == nil {
		panic("Parameters may not be empty!")
	}
	return &HandlerContext{deviceStore, connections}
>>>>>>> 47550d9798a8f2b1d80770d26df1b466dbc525b1
}
