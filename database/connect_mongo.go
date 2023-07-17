package database

import (
	"context"
	"os"
	"time"

	"github.com/PretendoNetwork/friends-secure/globals"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var mongoContext context.Context
var mongoDatabase *mongo.Database
var MongoCollection *mongo.Collection

func connectMongo() {
	mongoClient, _ = mongo.NewClient(options.Client().ApplyURI(os.Getenv("PN_FRIENDS_CONFIG_MONGO_URI")))
	mongoContext, _ = context.WithTimeout(context.Background(), 10*time.Second)
	_ = mongoClient.Connect(mongoContext)

	mongoDatabase = mongoClient.Database("pretendo")
	MongoCollection = mongoDatabase.Collection("pnids")

	globals.Logger.Success("Connected to Mongo!")
}
