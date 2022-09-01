package database

import (
	"context"
	"os"
	"sidewarslobby/app/queries"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Queries struct {
	*queries.UserQueries
}

func MongoDBConnection() (*Queries, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_CONNECTION_URL")))

	if err != nil {
		return nil, err
	}

	return &Queries{
		UserQueries: &queries.UserQueries{Client: client},
	}, nil
}
