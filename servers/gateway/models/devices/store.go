package devices

import (
	"errors"

	"gopkg.in/mgo.v2/bson"
)

//ErrUserNotFound is returned when the user can't be found
var ErrDeviceNotFound = errors.New("Device not found")

//Store represents a store for Devices
type Store interface {
	//GetByID returns the Device with the given ObjectID
	GetByID(id bson.ObjectId) (*Device, error)

	//GetByName returns the Device with the given name
	GetByName(name string) (*Device, error)

	//Insert inserts the device into the database, and returns
	//the newly-inserted device, complete with the DBMS-assigned bson.ObjectID
	Insert(device *Device) (*Device, error)

	//Update applies Device updates to the given device ID
	//and returns the newly-updated device
	Update(bson.ObjectId, *Updates) error

	//Delete deletes the device with the given ObjectID
	Delete(id bson.ObjectId) error
}
