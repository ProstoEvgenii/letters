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
	for _, event := range events {
		if event.IsDaily {
			if event.MustSend != currentDate {
				log.Println("=700beb=", "Обновлено")
				UpdateEvent(event.Name, false)
			} else if event.SendAt == int64(today.Hour()) && today.Minute() == 34 && event.IsSent != true {
				// log.Println("=da64d5=", "Время отправки", event)
				CheckLogsAndSendEmail(event, settings)
				UpdateEvent(event.Name, true)
			}
		} else {
			log.Println("=86f765=", event.Day != int64(currentDate.Day()) && event.Month != int64(currentDate.Month()) && event.IsSent == true)
			if event.Day == int64(currentDate.Day()) && event.Month == int64(currentDate.Month()) && event.IsSent == true {

				log.Println("=69734c=", event.Name)
				UpdateEvent(event.Name, false)
				// log.Println("=1228e4=", "Поздравляю с Днем города")
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
