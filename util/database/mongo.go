package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupMongoConnection(ctx context.Context, dbUri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(dbUri)
	conn, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
