package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongodb *mongo.Database
var ClientMongo *mongo.Client

// ConnectMongoDB connects to the MongoDB database and sets the DB and Client
func ConnectMongoDB() error {
	// Get MongoDB URI from environment variables
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return fmt.Errorf("MONGODB_URI not set")
	}

	mongoDBName := os.Getenv("MONGODB_DB_NAME")
	if mongoDBName == "" {
		return fmt.Errorf("MONGODB_DB_NAME not set")
	}

	// Connect to MongoDB clientLocal
	clientLocal, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}

	// Assign the client and database to global variables
	ClientMongo = clientLocal
	mongodb = clientLocal.Database(mongoDBName)

	// Log the successful connection to MongoDB
	log.Println("Connected to MongoDB!")

	return nil
}

// DisconnectMongoDB disconnects from MongoDB
func DisconnectMongoDB() error {
	// Disconnect MongoDB client
	if err := ClientMongo.Disconnect(context.TODO()); err != nil {
		return err
	}

	// Log the successful disconnection from MongoDB
	log.Println("Disconnected from MongoDB!")

	return nil
}

// GetCollectionMongo returns a MongoDB collection by name
func GetCollectionMongo(collectionName string) *mongo.Collection {
	return mongodb.Collection(collectionName)
}
