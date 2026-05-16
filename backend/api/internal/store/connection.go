package store

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Database struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func Connect() (*Database, error) {

	url := os.Getenv("MONGO_URL")
	if url == "" {
		return nil, fmt.Errorf("MONGO_URL environment variable not set")
	}

	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		return nil, fmt.Errorf("MONGO_DB environment variable not set")
	}

	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(url).
		SetServerAPIOptions(serverApi)

	client, err := mongo.Connect(opts)

	if err != nil {
		return nil, fmt.Errorf("Error connecting to db: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("Error pinging db: %w", err)
	}

	db := client.Database(dbName)

	return &Database{Client: client, Database: db}, nil
}
