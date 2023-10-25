package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collectionUsers *mongo.Collection

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
	collectionUsers = client.Database("crem").Collection("users")
	return
}

func InsertIfNotExists(newDocument Data) *mongo.UpdateResult {
	filter := bson.M{
		"uniqueField": "uniqueValue",
		// Фильтр для проверки уникальности
	}
	update := bson.M{
		"$setOnInsert": newDocument,
	}

	opts := options.Update().SetUpsert(true)

	cursor, err := collectionUsers.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Println("=440c2b=", err)
	}
	return cursor
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
