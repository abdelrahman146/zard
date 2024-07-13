package provider

import (
	"context"
	"github.com/abdelrahman146/zard/shared/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoProvider interface {
	GetDB(dbName string) *mongo.Database
	Close()
}

type mongoProvider struct {
	client *mongo.Client
}

func InitMongoProvider(address string) MongoProvider {
	clientOptions := options.Client().ApplyURI(address)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		logger.GetLogger().Panic("Failed to connect to mongoProvider", logger.Field("address", address), logger.Field("error", err))
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		logger.GetLogger().Panic("Failed to ping to mongoProvider", logger.Field("address", address), logger.Field("error", err))
	}
	logger.GetLogger().Info("Connected to mongoProvider", logger.Field("address", address))
	return &mongoProvider{
		client: client,
	}
}

func (m *mongoProvider) GetDB(dbName string) *mongo.Database {
	return m.client.Database(dbName)
}

func (m *mongoProvider) Close() {
	err := m.client.Disconnect(context.Background())
	if err != nil {
		logger.GetLogger().Warn("Failed to disconnect from mongoProvider", logger.Field("error", err))
		return
	}
	logger.GetLogger().Info("Mongo DB connection closed")
}
