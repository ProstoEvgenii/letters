package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
	collectionUsers = client.Database(os.Getenv("BASE")).Collection("users")
	return
}

func InsertIfNotExists(document Data) {
	filter := bson.M{
		"E-mail": document.Email,
	}
	update := bson.M{"$setOnInsert":  bson.M{
		"Имя": document.First_name,
		"Afvbkbz"
		"Дата рождения": dateBirth,
	},

	opts := options.Update().SetUpsert(true)
	log.Println("=08c25b=", document.Date_birth)
	dateBirth, err := time.Parse("01/02/2006", document.Date_birth)
	// log.Println("=449df3=", dt1)
	document.Date_birth = dateBirth
	result, err := collectionUsers.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Println("=InsertIfNotExists=", err)
	}
	// log.Println("=59201e=", result)
	if result.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
		return
	}
	if result.UpsertedCount != 0 {
		fmt.Printf("inserted a new document with ID %v\n", result.UpsertedID)
	}
	return
}

func CountDocuments() int64 {

	itemCount, err := collectionUsers.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		log.Println("=2671f1=", err)
	}
	return itemCount
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
