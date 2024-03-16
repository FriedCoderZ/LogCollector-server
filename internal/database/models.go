package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Collector struct {
	ID        primitive.ObjectID       `bson:"_id,omitempty"`
	UUID      string                   `bson:"uuid"`
	AESKey    []byte                   `bson:"aes_key"`
	Logs      []map[string]interface{} `bson:"logs"`
	CreatedAt time.Time                `bson:"created_at"`
	UpdatedAt time.Time                `bson:"updated_at"`
}
