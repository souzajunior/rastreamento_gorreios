package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoGorreiosENV = "MONGO_GORREIOS"

func OpenMongo(ctx context.Context) (client *mongo.Client, err error) {
	var (
		mongoURI string
		found    bool
	)

	// Looking for env variable
	if mongoURI, found = os.LookupEnv(mongoGorreiosENV); !found {
		log.Println(mongoGorreiosENV, "environment variable wasn't detected. Using default mongo URI {mongodb://localhost:27017/} to create a new client...")
		mongoURI = "mongodb://localhost:27017/"
	}

	// Instantiating a new client
	if client, err = mongo.NewClient(
		options.Client().ApplyURI(mongoURI),
	); err != nil {
		return nil, err
	}

	// Connecting
	if err = client.Connect(ctx); err != nil {
		return
	}

	return
}
