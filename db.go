package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collectionUsers *mongo.Collection
var collectionLogs *mongo.Collection
var collectionTemplates *mongo.Collection
var collectionSettings *mongo.Collection

func Connect() {
	uri := "mongodb://" + os.Getenv("LOGIN") + ":" + os.Getenv("PASS") + "@" + os.Getenv("SERVER")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatal("Ошибка подключения к базе данный =>", err)
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		log.Fatal("Ping провален =>", err)
	}

	log.Println("База данных подключена упешно!")
	collectionUsers = client.Database(os.Getenv("BASE")).Collection("users")
	collectionLogs = client.Database(os.Getenv("BASE")).Collection("logs")
	collectionTemplates = client.Database(os.Getenv("BASE")).Collection("templates")
	collectionSettings = client.Database(os.Getenv("BASE")).Collection("settings")

	return
}

func InsertIfNotExists(document interface{}, filter, update primitive.M, collName string) *mongo.UpdateResult {

	opts := options.Update().SetUpsert(true)
	ctx := context.TODO()
	switch collName {
	case "users":
		result, err := collectionUsers.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Println("=InsertIfNotExists=", err)
		}
		return result

	case "logs":
		result, err := collectionLogs.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Println("=InsertIfNotExists=", err)
		}
		return result

	case "settings":
		result, err := collectionSettings.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Println("=InsertIfNotExists=", err)
		}
		return result
	default:
		fmt.Println("=InsertIfNotExists=", "Не валидный case")
	}
	// if result.MatchedCount != 0 {
	// 	fmt.Println("matched and replaced an existing document")
	// }
	// if result.UpsertedCount != 0 {
	// 	fmt.Printf("inserted a new document with ID %v\n", result.UpsertedID)
	// }
	return nil
}

func CountDocuments() int64 {

	itemCount, err := collectionUsers.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		log.Println("=2671f1=", err)
	}
	return itemCount
}

func Find(filter primitive.M, collName string) *mongo.Cursor {
	ctx := context.TODO()
	switch collName {
	case "users":
		cursor, err := collectionUsers.Find(ctx, filter)
		if err != nil {
			log.Println("=Find=", err)
		}
		return cursor

	case "logs":
		cursor, err := collectionLogs.Find(ctx, filter)
		if err != nil {
			log.Println("=Find=", err)
		}
		return cursor
	case "templates":
		cursor, err := collectionTemplates.Find(ctx, filter)
		if err != nil {
			log.Println("=Find=", err)
		}
		return cursor
	case "settings":
		result, err := collectionSettings.Find(ctx, filter)
		if err != nil {
			log.Println("=Find=", err)
		}
		return result
	default:
		fmt.Println("=InsertIfNotExists=", "Не валидный case")
	}
	return nil
}

// func Find(filter, sort bson.M, limit int64, collName string) (*mongo.Cursor, error) {
// 	findOptions := options.Find()
// 	findOptions.SetSort(sort)
// 	findOptions.SetLimit(limit)
// 	return DataBase.Collection(collName).Find(context.TODO(), filter, findOptions)
// }

// func FindOneAndUpdate(filter, update bson.M, upsert bool, collName string) *mongo.SingleResult {
// 	after := options.After
// 	opt := options.FindOneAndUpdateOptions{
// 		ReturnDocument: &after,
// 		Upsert:         &upsert,
// 	}
// 	return a := Collection(collName).FindOneAndUpdate(context.TODO(), filter, update, &opt)
// }
