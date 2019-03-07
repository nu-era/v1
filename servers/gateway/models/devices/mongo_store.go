package devices

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// insert methods below

type MongoStore struct {
	ses *mgo.Session
}

func NewMySqlStore(ses *mgo.Session) *MongoStore {
	if ses == nil {
		return nil
	}
	return &MongoStore{
		ses: ses,
	}
}

//GetByID returns the Device with the given ID
func (ms *MongoStore) GetByID(id bson.ObjectId) (*Device, error) {
	return ms.get("_id", string(id))
}

//GetByEmail returns the Device with the given email
func (ms *MongoStore) GetByEmail(email string) (*Device, error) {
	return ms.get("email", email)
}

//GetByName returns the Device with the given Name
func (ms *MongoStore) GetByName(name string) (*Device, error) {
	return ms.get("name", name)
}

func (ms *MongoStore) get(col string, val string) (*Device, error) {
	coll := ms.ses.DB("store").C("devices")
	dev := Device{}
	if col == "val" {
		coll.Find(bson.M{col: bson.ObjectId(val)}).One(&dev)
	}
	coll.Find(bson.M{col: val}).One(&dev)
	return &dev, nil
}

//Insert inserts the device into the database, and returns
//the newly-inserted Device, complete with the DBMS-assigned ID
func (ms *MongoStore) Insert(dev *Device) (*Device, error) {
	coll := ms.ses.DB("store").C("devices")
	//insert struct into collection
	if err := coll.Insert(dev); err != nil {
		return nil, fmt.Errorf("error inserting document: %v\n", err)
	} else {
		fmt.Printf("inserted document with ID %s\n", dev.ID.Hex())
		return dev, nil
	}
}

//Update applies DeviceUpdates to the given device ID
//and returns the newly-updated device. Only applies valid updates to Device
func (ms *MongoStore) Update(id bson.ObjectId, updates *Updates) (*Device, error) {
	coll := ms.ses.DB("store").C("devices")
	//send an update document with a `$set` property set to the updates map
	if err := coll.UpdateId(id, bson.M{"$set": updates}); err != nil {
		return nil, err
	}
	dev, err := ms.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err = dev.ApplyUpdates(updates); err != nil {
		return nil, err
	}
	return dev, nil
}

//Delete deletes the device with the given ID
func (ms *MongoStore) Delete(id bson.ObjectId) error {
	coll := ms.ses.DB("store").C("devices")
	if err := coll.RemoveId(id); err != nil {
		return err
	}

	return nil
}

//InsertMessage inserts a message into the database for the Device with given id
func (ms *MongoStore) InsertMessage(id bson.ObjectId, msg *Message) error {
	coll := ms.ses.DB("store").C("devices")
	dev := bson.M{"_id": id}
	push := bson.M{"$push": bson.M{"messages": msg}}
	if err := coll.Update(dev, push); err != nil {
		return err
	}
	return nil
}
