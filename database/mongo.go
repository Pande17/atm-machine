package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongoclient *mongo.Client

func ConnectMongoDB() error {
	clientOptions := options.Client().ApplyURI(os.Getenv("mongodb"))

	var err error
	Mongoclient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = Mongoclient.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB!")

	return nil
}

func ConnectCollection(collectionName string) *mongo.Collection {
	if Mongoclient == nil {
		if err := ConnectMongoDB(); err != nil {
			log.Fatal("can't connect to MongoDB!")
		}
	}

	return Mongoclient.Database("atm-machine").Collection(collectionName)
}
