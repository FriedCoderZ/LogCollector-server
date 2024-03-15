package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogRecord struct {
	ID   primitive.ObjectID     `bson:"_id,omitempty"`
	Data map[string]interface{} `bson:"data"`
}

type Collector struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UUID      string             `bson:"uuid"`
	IP        string             `bson:"ip"`
	AESKey    []byte             `bson:"aes_key"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
