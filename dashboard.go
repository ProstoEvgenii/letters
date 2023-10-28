package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func Dashboard() (int64, int) {
	usersCount := CountDocuments()
	birthdays_list := CreateBirthdaysSlice()
	return usersCount, len(birthdays_list)
}
func CreateBirthdaysSlice() []Users {
	today := time.Now()
	filter := bson.M{}
	cursor := Find(filter, "users")
	var users []Users
	err := cursor.All(context.TODO(), &users)
	if err != nil {
		log.Println("=84ce91=", err)
	}

	var birthdays_list []Users
	for _, user := range users {
		if user.DateOfBirth.Day() == today.Day() && user.DateOfBirth.Month() == today.Month() {
			result := CreateLog(user)
			if result != 0 {
				log.Println("=56eccc=", result)
				////Отправить письмо тут .Вынести этот цикл в отдельную функцию
			}

			birthdays_list = append(birthdays_list, user)
		}
	}

	return birthdays_list
}

func GetTemplate() {
	filter := bson.M{
		"name": "test1",
	}
	cursor := Find(filter, "templates")
	var template []Templates
	if err := cursor.All(context.TODO(), &template); err != nil {
		log.Println("=8922b7=", err)
	}
	log.Println("=a1e37e=", (template))

}

func CreateLog(user Users) int64 {
	currentDate := time.Now().UTC().Truncate(24 * time.Hour)
	filter := bson.M{
		"E-mail":     user.Email,
		"dateCreate": currentDate,
	}
	update := bson.M{"$setOnInsert": bson.M{
		"Имя":           user.FirstName,
		"Фамилия":       user.LastName,
		"Отчество":      user.MiddleName,
		"Дата рождения": user.DateOfBirth,
		"E-mail":        user.Email,
		"dateCreate":    currentDate,
	}}
	result := InsertIfNotExists(user, filter, update, "logs").UpsertedCount
	return result
}

// func FindLogs() []Users {
// 	today := time.Now()
// 	twentyFourHoursAgo := today.Add(-24 * time.Hour)
// 	// filter := bson.M{}
// 	filter := bson.M{"Дата рождения": bson.M{"$gte": twentyFourHoursAgo}}
// 	cursor := Find(filter, "logs")
// 	var users_logs []Users
// 	err := cursor.All(context.TODO(), &users_logs)
// 	if err != nil {
// 		log.Println("=84ce91=", err)
// 	}

// 	// var logs_list []Users
// 	// for _, user := range users_logs {
// 	// 	if user.DateOfBirth.Day() == today.Day() && user.DateOfBirth.Month() == today.Month() {
// 	// 		birthdays_list = append(birthdays_list, user)
// 	// 	}
// 	// }

// 	return
// }
