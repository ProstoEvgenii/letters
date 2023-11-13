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
	events := getEvents()
	settings := GetSettings()

	currentDate := time.Now().UTC().Truncate(24 * time.Hour)
	today := time.Now()
	currentDay := int64(today.Day())
	currentMonth := int64(today.Month())
	currentHour := int64(today.Hour())
	currentMinute := today.Minute()
	for _, event := range events {
		if event.IsDaily {
			if event.MustSend != currentDate {
				log.Println("=700beb=", "Обновлено")
				UpdateEvent(event.Name, false)
			} else if event.SendAt == currentDay && currentMinute == 34 && !event.IsSent {
				// log.Println("=da64d5=", "Время отправки", event)
				CheckLogsAndSendEmail(event, settings)
				UpdateEvent(event.Name, true)
			}
		}

		if !event.IsDaily && event.Day == currentDay && event.Month == currentMonth {
			if event.IsSent && event.SendAt != currentHour {
				log.Println("=69734c=", event.Name)
				UpdateEvent(event.Name, false)
				// log.Println("=1228e4=", "Поздравляю с Днем города")
			} else if !event.IsSent && event.SendAt == currentHour && currentMinute == 00 {
				log.Println("=72e334=")
				UpdateEvent(event.Name, true)
			}
		}
	}

	time.AfterFunc(time.Duration(3)*time.Second, func() {
		AutoSend()
	})
}
func getEvents() []models.Events {
	filter := bson.M{}
	cursor := db.Find(filter, "events")
	var events []models.Events
	if err := cursor.All(context.TODO(), &events); err != nil {
		log.Println("=8922b7=", err)
	}
	return events
}

func UpdateEvent(eventName string, isSent bool) {
	currentDate := time.Now().UTC().Truncate(24 * time.Hour)
	filter := bson.M{
		"name": eventName,
	}
	update := bson.M{"$set": bson.M{
		"mustSend": currentDate,
		"isSent":   isSent,
	}}
	db.UpdateIfExists(filter, update, "events")
}
