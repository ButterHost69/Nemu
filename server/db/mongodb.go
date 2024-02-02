package db

import (
	"context"
	
	"time"

	// "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitMongoDB()(*mongo.Client, error){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mdb, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	
	if err = mdb.Ping(ctx, readpref.Primary()) ; err != nil{
		return nil, err
	}

	return mdb, nil
}

func InsertPostInMongoDB(collection *mongo.Collection, username string, content string) error{
	ctx, cancel := context.WithTimeout(context.Background(), 5* time.Second)
	defer cancel()

	_, ierr := collection.InsertOne(ctx, bson.D{{Key: "username", Value: username}, {Key: "content", Value: content}})
	if ierr != nil {
		return	ierr
	}
	
	return nil
}
