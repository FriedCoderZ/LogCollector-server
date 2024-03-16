package database

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collector
func CreateCollector(aesKey []byte) (*Collector, error) {
	collection := GetDB().Collection("collectors")
	c := &Collector{
		ID:        primitive.NewObjectID(),
		UUID:      uuid.New().String(),
		AESKey:    aesKey,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := collection.InsertOne(context.Background(), c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
func GetCollectorByUUID(uuid string) (*Collector, error) {
	collection := GetDB().Collection("collectors")
	var collector Collector
	filter := bson.M{"uuid": uuid}

	err := collection.FindOne(context.Background(), filter).Decode(&collector)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("Collector not found")
		}
		return nil, err
	}

	return &collector, nil
}

func (c *Collector) GetAESKey() ([]byte, error) {
	if c.AESKey == nil {
		return nil, fmt.Errorf("AESKey is not set")
	}
	return c.AESKey, nil
}

func convertToNumber(m map[string]interface{}) {
	for key, value := range m {
		switch v := value.(type) {
		case string:
			if num, err := strconv.Atoi(v); err == nil {
				m[key] = num
			} else if num64, err := strconv.ParseInt(v, 10, 64); err == nil {
				m[key] = num64
			}
		case float32, float64:
			m[key] = int64(v.(float64))
		}
	}
}

func (c *Collector) AppendLogs(logs []map[string]interface{}) error {
	for _, log := range logs {
		convertToNumber(log)
	}
	c.Logs = append(c.Logs, logs...)
	filter := bson.M{"_id": c.ID}
	update := bson.M{"$set": bson.M{"logs": c.Logs}}
	opts := options.Update().SetUpsert(false)
	db := GetDB()
	collection := db.Collection("collectors")
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}
