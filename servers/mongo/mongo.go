package mongo

import (
	"github.com/mongodb/mongo-go-driver/bson"
)

type Device struct {
	ID	bson.ObjectId `bson:"_id"`,
	lat float64,
	long float64,
	passHash []byte,
	emailHash []byte,
	Status string,
	messages map[string]interface{}
}

// insert methods below


