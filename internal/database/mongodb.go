package database

import (
	"context"
	"fmt"
	"log"

	"github.com/FriedCoderZ/LogCollector-server/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db *mongo.Database
)

func init() {
	// 创建数据库连接
	config := config.GetConfig()
	clientOptions := options.Client().ApplyURI(config.Database.Address)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// 设置数据库实例
	db = client.Database("logCollector")
}

func GetDB() *mongo.Database {
	return db
}

// 其他操作代码...
