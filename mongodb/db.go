package mongodb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Config struct {
	URI  string
	Name string
}

func GetNoSQLDatabaseConnection(conf Config) *mongo.Database {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(conf.URI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	return client.Database(conf.Name)
}
