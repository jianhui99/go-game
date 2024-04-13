package database

import (
	"common/config"
	"common/logs"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type MongoManager struct {
	Client *mongo.Client
	Db     *mongo.Database
}

func NewMongo() *MongoManager {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Conf.Database.MongoConf.Url))
	if err != nil {
		logs.Fatal("mongo connect err: %v", err)
		return nil
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		logs.Fatal("mongo ping err: %v", err)
		return nil
	}

	m := &MongoManager{
		Client: client,
	}
	m.Db = m.Client.Database(config.Conf.Database.MongoConf.Db)
	return m
}

func (m *MongoManager) Close() {
	err := m.Client.Disconnect(context.TODO())
	if err != nil {
		logs.Error("close mongodb err:", err)
	}
}
