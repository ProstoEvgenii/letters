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
	today := time.Now()
	filter := bson.M{}
	cursor := Find(filter)
	var users []CollUsers
	err := cursor.All(context.TODO(), &users)
	if err != nil {
		log.Println("=84ce91=", err)
	}
	// log.Println("=54c6f9=", users)
	// log.Println("=0788d2=", today.Day())
	// log.Println("=5dec2e=",)
	for _, user := range users {
		if user.DateOfBirth.Day() == today.Day() {
			log.Println("=afdf3c=", user.DateOfBirth.Day())
		}
	}

	return itemCount
}
