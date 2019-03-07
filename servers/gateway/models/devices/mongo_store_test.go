package devices

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"testing"
)

func TestGetByID(t *testing.T) {
	cases := []struct {
		name        string
		column      string
		expectError bool
		needRecord  bool
		record      map[interface{}]string
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
			"Valid Get By Id",
			"id",
			false,
			true,
			nil,
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
			dev, err = ms.GetByID(bson.ObjectId(c.val))
		case "email":
			dev, err = ms.GetByEmail(c.val)
		case "name":
			dev, err = ms.GetByName(c.val)
		}

		if c.expectError && err == nil {
			fmt.Println(dev)
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

func TestGetByEmail(t *testing.T) {
	return
}

func TestGetByName(t *testing.T) {
	return
}

func TestInsert(t *testing.T) {
	return
}

func TestUpdate(t *testing.T) {
	return
}

func TestDelete(t *testing.T) {
	return
}
