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
	events := GetEvents()

	currentDate := time.Now().UTC().Truncate(24 * time.Hour)
	today := time.Now()
	currentDay := int64(today.Day())
	currentMonth := int64(today.Month())
	currentHour := int64(today.Hour())
	currentMinute := today.Minute()

	for _, event := range events {
		if event.IsDaily && event.Active {
			if event.MustSend != currentDate {
				UpdateEvent(event.Name, false)
			} else if event.SendAt == currentHour && currentMinute == 00 && !event.IsSent {
				birthdays_list, anniversary_list := CreateBirthdaysSlice()
				if event.Type == "anniversary" && len(anniversary_list) != 0 {
					CheckLogsAndSendEmail(event, anniversary_list)
				} else if event.Type == "birthday" && len(birthdays_list) != 0 {
					CheckLogsAndSendEmail(event, birthdays_list)
				}
				UpdateEvent(event.Name, true)
			}
		}
		if !event.IsDaily && event.Active && event.Day == currentDay && event.Month == currentMonth {
			if event.IsSent {
				UpdateEvent(event.Name, false)
			} else if !event.IsSent && event.SendAt == currentHour && currentMinute == 00 {
				log.Println("=72e334=", "Отправлено", event.Name)
				SendToEverybody(event)
				UpdateEvent(event.Name, true)
			}
		}
	}

	time.AfterFunc(time.Duration(60)*time.Second, func() {
		AutoSend()
	})
}

func GetEvents() []models.Events {
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
