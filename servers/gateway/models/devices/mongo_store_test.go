package devices

import (
	//"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	cases := []struct {
		name        string
		column      string
		expectError bool
		needRecord  bool
		record      map[string]interface{}
		coll        string
		DB          string
		val         string
	}{
		{
			"No record",
			"id",
			true,
			false,
			nil,
			"devices",
			"store",
			"1",
		},
		{
			"No record",
			"email",
			true,
			false,
			nil,
			"devices",
			"store",
			"does@not.exist",
		},
		{
			"No record",
			"name",
			true,
			false,
			nil,
			"devices",
			"store",
			"Nonexistent Device",
		},
		{
			"Valid Get By Id",
			"id",
			false,
			true,
			map[string]interface{}{
				"_id":       bson.NewObjectId(),
				"name":      "Seattle School District",
				"latitude":  127.3995,
				"longitude": 193.4564,
				"passHash":  []byte{1, 2, 3, 4},
				"email":     "test@test.com",
				"status":    "up",
			},
			"devices",
			"store",
			"",
		},
		{
			"Valid Get By Email",
			"email",
			false,
			true,
			map[string]interface{}{
				"_id":       bson.NewObjectId(),
				"name":      "Seattle School District",
				"latitude":  127.3995,
				"longitude": 193.4564,
				"passHash":  []byte{1, 2, 3, 4},
				"email":     "test@test.com",
				"status":    "up",
			},
			"devices",
			"store",
			"test@test.com",
		},
		{
			"Valid Get By Name",
			"name",
			false,
			true,
			map[string]interface{}{
				"_id":       bson.NewObjectId(),
				"name":      "Seattle School District",
				"latitude":  127.3995,
				"longitude": 193.4564,
				"passHash":  []byte{1, 2, 3, 4},
				"email":     "test@test.com",
				"status":    "up",
			},
			"devices",
			"store",
			"Seattle School District",
		},
	}

	mongo, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatalf("error dialing mongo: %v", err)
	}

	ms := NewMongoStore(mongo)

	for _, c := range cases {
		if c.needRecord {
			ms.ses.DB(c.DB).C(c.coll).Insert(c.record)
		}
		var dev *Device
		var err error
		switch c.column {
		case "id":
			if c.needRecord {
				dev, err = ms.GetByID(c.record["_id"].(bson.ObjectId))
			} else {
				dev, err = ms.GetByID(bson.ObjectId(c.val))
			}
		case "email":
			dev, err = ms.GetByEmail(c.val)
		case "name":
			dev, err = ms.GetByName(c.val)
		}

		if c.expectError && err == nil {
			t.Errorf("Case: %s, Expected error but received none", c.name)
		}

		if !c.expectError && err != nil {
			t.Errorf("Case: %s, Expected no error but got: %v", c.name, err)
		}

		if dev == nil && !c.expectError {
			t.Errorf("Case: %s, Expecting non-empty struct", c.name)
		}

	}

}

func TestInsert(t *testing.T) {
	cases := []struct {
		name        string
		expectError bool
		device      Device
		coll        string
		DB          string
	}{
		{
			"Invalid Device",
			true,
			Device{},
			"devices",
			"store",
		},
		{
			"Valid Device",
			false,
			Device{
				ID:       bson.NewObjectId(),
				Name:     "Seattle School District",
				Lat:      127.3995,
				Long:     193.4564,
				PassHash: []byte{1, 2, 3, 4},
				Email:    "test@test.com",
				Phone:    "1234567890",
				Status:   "up",
			},
			"devices",
			"store",
		},
	}

	mongo, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatalf("error dialing mongo: %v", err)
	}

	ms := NewMongoStore(mongo)

	for _, c := range cases {
		dev, err := ms.Insert(&c.device)

		if c.expectError && err == nil {
			t.Errorf("Case: %s, Expected error but received none", c.name)
		}

		if !c.expectError && err != nil {
			t.Errorf("Case: %s, Expected no error but got: %v", c.name, err)
		}

		if !reflect.DeepEqual(c.device, Device{}) {
			tmp, err := ms.GetByID(dev.ID)
			if !reflect.DeepEqual(&c.device, tmp) {
				t.Errorf("Case: %s, Expected record to be %v, got %v", c.name, c.device, tmp)
			}
			if c.expectError && err == nil {
				t.Errorf("Case: %s, Expected error but received none", c.name)
			}

			if !c.expectError && err != nil {
				t.Errorf("Case: %s, Expected no error but got: %v", c.name, err)
			}
		}
	}
}

func TestUpdate(t *testing.T) {
	cases := []struct {
		name        string
		expectError bool
		needRecord  bool
		device      Device
		updates     Updates
		coll        string
		DB          string
	}{
		{
			"Invalid Device",
			true,
			false,
			Device{},
			Updates{},
			"devices",
			"store",
		},
		{
			"Invalid Update",
			true,
			false,
			Device{
				ID:       bson.NewObjectId(),
				Name:     "test",
				Lat:      127.3995,
				Long:     193.4564,
				PassHash: []byte{1, 2, 3, 4},
				Email:    "test@test.com",
				Phone:    "1234567890",
				Status:   "up",
			},
			Updates{},
			"devices",
			"store",
		},
		{
			"Valid Update",
			true,
			false,
			Device{
				ID:       bson.NewObjectId(),
				Name:     "test",
				Lat:      127.3995,
				Long:     193.4564,
				PassHash: []byte{1, 2, 3, 4},
				Email:    "test@test.com",
				Phone:    "1234567890",
				Status:   "up",
			},
			Updates{
				Name:   "new name",
				Email:  "upd@email.com",
				Status: "down",
			},
			"devices",
			"store",
		},
	}

	mongo, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatalf("error dialing mongo: %v", err)
	}

	ms := NewMongoStore(mongo)

	for _, c := range cases {
		if c.needRecord {
			ms.Insert(&c.device)

		}
		err := ms.Update(c.device.ID, &c.updates)
		if c.expectError && err == nil {
			t.Errorf("Case: %s, Expected error but received none", c.name)
		}

		if !c.expectError && err != nil {
			t.Errorf("Case: %s, Expected no error but got: %v", c.name, err)
		}
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		name        string
		expectError bool
		needRecord  bool
		device      Device
		coll        string
		DB          string
	}{
		{
			"Invalid Device",
			true,
			false,
			Device{},
			"devices",
			"store",
		},
		{
			"Valid Device",
			true,
			false,
			Device{
				ID:       bson.NewObjectId(),
				Name:     "test",
				Lat:      127.3995,
				Long:     193.4564,
				PassHash: []byte{1, 2, 3, 4},
				Email:    "test@test.com",
				Phone:    "1234567890",
				Status:   "up",
			},
			"devices",
			"store",
		},
	}

	mongo, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatalf("error dialing mongo: %v", err)
	}

	ms := NewMongoStore(mongo)

	for _, c := range cases {
		if c.needRecord {
			ms.Insert(&c.device)
		}
		err := ms.Delete(c.device.ID)
		if c.expectError && err == nil {
			t.Errorf("Case: %s, Expected error but received none", c.name)
		}

		if !c.expectError && err != nil {
			t.Errorf("Case: %s, Expected no error but got: %v", c.name, err)
		}
	}
}
