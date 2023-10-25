package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
	return
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
