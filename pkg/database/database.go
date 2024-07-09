package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongodClient(ctx context.Context, connectionString string) *mongo.Client{
    client, err := mongo.Connect(ctx,options.Client().ApplyURI(connectionString))
    if err != nil {
        return nil
    }
    return client
}


func NewMongodCollection(client *mongo.Client,dbName,collectionName string) *mongo.Collection{
   return client.Database(dbName).Collection(collectionName)

}