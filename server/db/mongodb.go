package db

import (
	"context"
	"example/one-page/server/models"
	"fmt"
	// "reflect"
	"strconv"
	"time"

	// "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitMongoDB() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mdb, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	if err = mdb.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return mdb, nil
}

func InsertPostInMongoDB(collection *mongo.Collection, username string, content string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createdAt := time.Now()
	_, ierr := collection.InsertOne(ctx, bson.D{{Key: "username", Value: username}, {Key: "content", Value: content}, {Key: "createdAt", Value: createdAt}})
	if ierr != nil {
		return ierr
	}

	return nil
}

func InsertCommentInMongoDB(collection *mongo.Collection, postObjId string, username string, commentContent string) error{

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createdAt := time.Now()

	objectID, err := primitive.ObjectIDFromHex(postObjId)
	if err != nil {
		fmt.Println(" > Error in Convertin string to primitive.ObjectID ")
		fmt.Println(err.Error())
		return err
	}
	filter := bson.D{{Key: "_id", Value: objectID}}
	newComment := bson.D{
		{Key: "commentUsername", Value: username},
		{Key: "commmentContent", Value: commentContent},
		{Key: "commentCreatedAt", Value: createdAt},
	}

	update := bson.D{{Key: "$push", Value: bson.D{{ Key: "comments", Value: newComment }}}}
	_, eerr := collection.UpdateOne(ctx, filter, update)
	if eerr != nil {
    	fmt.Println(" > Error Occured in Executing the MongoDB code : ")
		fmt.Println(eerr.Error())
		return eerr
	}
	return nil
}

func GetPostsFromMongoDB(collection *mongo.Collection, offset int, limit int) ([]models.Post, bool) {
	var resultSet []models.Post
	var emptyList bool
	emptyList = false

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		fmt.Println(" Error Occured at Getting Cursor [GetPostsFromMongoDB] : ")
		fmt.Println(err.Error())
		return nil, emptyList
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result bson.D
		var post models.Post

		err := cursor.Decode(&result)
		if err != nil {
			fmt.Println(" Error Occured in Accessing Cursor : ")
			fmt.Println(err)
		}

		for _, elem := range result {
			if elem.Key == "username" {
				post.Username, _ = elem.Value.(string)
				// fmt.Println(" >> Username : ", post.Username)
			} else if elem.Key == "content" {
				post.Data, _ = elem.Value.(string)
			} else if elem.Key == "createdAt" {

				// Convert primitive.DateTime -> time.Time -> day, month, year
				year, month, day := (elem.Value.(primitive.DateTime)).Time().Date()
				date := strconv.Itoa(day) + "." + strconv.Itoa(int(month)) + "." + strconv.Itoa(year)
				post.CreatedAt = date
				// fmt.Println(date)
			} else if elem.Key == "_id" {
				id := elem.Value.(primitive.ObjectID).Hex()
				post.ObjectID = id
				// fmt.Println(id, " - ", reflect.TypeOf(id))
				// fmt.Println("Object ID : ", id.(primitive.ObjectID).Hex())
			}
		}

		resultSet = append(resultSet, post)
	}

	if err := cursor.Err(); err != nil {
		fmt.Println(" Error Occured at cursor : ")
		fmt.Println(err.Error())
	}

	if len(resultSet) == 0 {
		// fmt.Println(" > Last Post Reached ")
		emptyList = true
	}

	return resultSet, emptyList
}

func GetUsernameThroughObjectID(collection *mongo.Collection, objectHexID string) string {
	objID, err := primitive.ObjectIDFromHex(objectHexID)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	curr, cerr := collection.Find(ctx, bson.D{{Key: "_id", Value: objID}})
	if cerr != nil {
		fmt.Println(" Error Occured at Getting Cursor [GetUsernameThroughObjectID] : ")
		fmt.Println(err.Error())
		return ""
	}
	defer curr.Close(ctx)

	for curr.Next(ctx) {
		var result bson.D

		derr := curr.Decode(&result)
		if derr != nil {
			fmt.Println(" Error Occured in Accessing Cursor : ")
			fmt.Println(derr.Error())
		}

		for _, elem := range result {
			if elem.Key == "username" {
				return elem.Value.(string)
			}
		}
	}

	return ""
}
