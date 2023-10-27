package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type CollUsers struct {
	FirstName   string    `bson:"Имя"`
	LastName    string    `bson:"Фамилия"`
	MiddleName  string    `bson:"Отчество"`
	DateOfBirth time.Time `bson:"Дата рождения"`
	Email       string    `bson:"E-mail"`
}

func Dashboard() int64 {
	itemCount := CountDocuments()

	filter := bson.M{}
	cursor := Find(filter)
	var users CollUsers
	err := cursor.All(context.TODO(), &users)
	if err != nil {
		log.Println("=84ce91=", err)
	}
	log.Println("=54c6f9=",users)

	return itemCount
}
