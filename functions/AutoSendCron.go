package functions

import (
	"letters/db"
	"letters/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func AutoSend() {
	var settings models.SettingsUpload
	cursor := db.FindOne(bson.M{}, "settings")
	cursor.Decode(&settings)

	now := time.Now()
	log.Println("=Вызвалась=", settings.SendAutoAt)
	if now.Hour() == settings.SendAutoAt && now.Minute() == 20 {
		result := GetStatusToday()
		log.Println("=f84318=", result)
		if !result.IsSent {
			var info models.IsSent
			//Отправляю письмо и CreateStatusToday() = true
			CheckLogsAndSendEmail()
			CreateStatusToday(info, true)
		}

	}

	time.AfterFunc(time.Duration(3)*time.Second, func() {
		AutoSend()
	})

}

func CreateStatusToday(info models.IsSent, isSent bool) {
	currentDate := time.Now().UTC().Truncate(24 * time.Hour)
	filter := bson.M{
		"date": currentDate,
	}

	update := bson.M{"$set": bson.M{
		"date":   currentDate,
		"isSent": isSent,
	}}
	db.InsertIfNotExists(info, filter, update, "isSentToday")

}

func GetStatusToday() models.IsSent {
	currentDate := time.Now().UTC().Truncate(24 * time.Hour)
	filter := bson.M{
		"date": currentDate,
	}
	cursor := db.FindOne(filter, "isSentToday")
	var info models.IsSent
	err := cursor.Decode(&info)
	if err != nil {
		log.Println("=5dd75c=", err)
		return models.IsSent{}
	}
	// if cursor.Err().Error() == "no documents in result"
	// log.Println("=2bf842=", cursor.Err().Error())
	return (info)
}
