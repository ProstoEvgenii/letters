package functions

import (
	"context"
	"letters/db"
	"letters/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func AutoSend() {

	cursor := db.Find(bson.M{}, "events")
	var eventsSlice []models.Events
	if err := cursor.All(context.TODO(), &eventsSlice); err != nil {
		log.Println("Cursor All Error Events", err)
	}
	// log.Println("=dc1ead=", "Вызвалась")
	now := time.Now()
	for _, event := range eventsSlice {
		if event.IsDaily == "daily" {
			log.Println("=event=", event)
			if now.Hour() == event.SendAt && now.Minute() == 21 {
				// log.Println("=7a16de=", event)
				result := GetEventLogToday(event.Name)
				if !result.IsSent {
					info := models.IsSent{
						Name: event.Name,
					}
					//Отправляю письмо и CreateEventLogToday() = true
					CheckLogsAndSendEmail()
					CreateEventLogToday(info, true)
				}

			}
		}

	}

	time.AfterFunc(time.Duration(3)*time.Second, func() {
		AutoSend()
	})

}

func CreateEventLogToday(info models.IsSent, isSent bool) {
	currentDate := time.Now().UTC().Truncate(24 * time.Hour)
	filter := bson.M{
		"event": "День рождения",
		"date":  currentDate,
	}

	update := bson.M{"$set": bson.M{
		"event":  "День рождения",
		"date":   currentDate,
		"isSent": isSent,
	}}
	db.InsertIfNotExists(info, filter, update, "isSentToday")

}

func GetEventLogToday(event_name string) models.IsSent {
	currentDate := time.Now().UTC().Truncate(24 * time.Hour)
	filter := bson.M{
		"event": event_name,
		"date":  currentDate,
	}
	cursor := db.FindOne(filter, "isSentToday")
	var info models.IsSent
	err := cursor.Decode(&info)
	if err != nil {
		log.Println("=5dd75c=", err)
		return models.IsSent{}
	}
	return info
}
