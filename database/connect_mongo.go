package database

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var mongoContext context.Context
var mongoDatabase *mongo.Database
var mongoCollection *mongo.Collection

func connectMongo() {
	mongoClient, _ = mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	mongoContext, _ = context.WithTimeout(context.Background(), 10*time.Second)
	_ = mongoClient.Connect(mongoContext)

	mongoDatabase = mongoClient.Database("pretendo")
	mongoCollection = mongoDatabase.Collection("pnids")
}
