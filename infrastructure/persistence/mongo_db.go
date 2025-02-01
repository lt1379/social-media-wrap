package persistence

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"my-project/infrastructure/logger"
)

func NewMongoDb(host string, port string, username string, password string, database string) (*mongo.Client, error) {
	// connect to mongodb
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://" + username + ":" + password + "@" + host + ":" + port + "/" + database))
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Failed to connect to MongoDB")
		return nil, err
	}

	return client, nil
}
