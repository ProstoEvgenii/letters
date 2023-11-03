package main

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func autoSend() {
	settings := GetSettings()

	now := time.Now()
	log.Println("=Вызвалась=", settings.SendAutoAt)
	if now.Hour() == settings.SendAutoAt && now.Minute() == 20 {
		result := getStatusToday()
		log.Println("=f84318=", result)
		if !result.IsSent {
			var info IsSent
			//Отправляю письмо и CreateStatusToday() = true
			checkLogsAndSendEmail()
			CreateStatusToday(info, true)
		}

	}

	time.AfterFunc(time.Duration(3)*time.Second, func() {
		autoSend()
	})

}

func CreateStatusToday(info IsSent, isSent bool) {
	currentDate := time.Now().UTC().Truncate(24 * time.Hour)
	filter := bson.M{
		"date": currentDate,
	}

	update := bson.M{"$set": bson.M{
		"date":   currentDate,
		"isSent": isSent,
	}}
	InsertIfNotExists(info, filter, update, "isSentToday")

}

func getStatusToday() IsSent {
	currentDate := time.Now().UTC().Truncate(24 * time.Hour)
	filter := bson.M{
		"date": currentDate,
	}
	cursor := FindOne(filter, "isSentToday")
	var info IsSent
	err := cursor.Decode(&info)
	if err != nil {
		log.Println("=5dd75c=", err)
		return IsSent{}
	}
	// if cursor.Err().Error() == "no documents in result"
	// log.Println("=2bf842=", cursor.Err().Error())
	return (info)
}
