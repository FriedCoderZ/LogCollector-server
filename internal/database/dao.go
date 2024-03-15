package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCollector(aesKey []byte) (*Collector, error) {
	collection := GetDB().Collection("collectors")
	c := &Collector{
		ID:        primitive.NewObjectID(),
		IP:        "", // 设置空字符串或其他默认值
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

func UpdateCollector(c *Collector) error {
	collection := GetDB().Collection("collectors")
	filter := bson.D{{Key: "_id", Value: c.ID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "ip", Value: c.IP},
		{Key: "aes_key", Value: c.AESKey},
		{Key: "updated_at", Value: time.Now()},
	}}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}

func DeleteCollector(id primitive.ObjectID) error {
	collection := GetDB().Collection("collectors")
	filter := bson.D{{Key: "_id", Value: id}}
	_, err := collection.DeleteOne(context.Background(), filter)
	return err
}

func FindCollectorByID(id primitive.ObjectID) (*Collector, error) {
	collection := GetDB().Collection("collectors")
	filter := bson.D{{Key: "_id", Value: id}}
	var c Collector
	err := collection.FindOne(context.Background(), filter).Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func FindAllCollectors() ([]Collector, error) {
	collection := GetDB().Collection("collectors")
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var collectors []Collector
	for cur.Next(context.Background()) {
		var c Collector
		err := cur.Decode(&c)
		if err != nil {
			return nil, err
		}
		collectors = append(collectors, c)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return collectors, nil
}
