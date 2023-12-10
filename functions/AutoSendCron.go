package functions

import (
	"context"
	"letters/db"
	"letters/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var activeEvents map[models.Events]bool

func init() {
	activeEvents = make(map[models.Events]bool)
}
func AutoSend() {
	events := GetEvents()
	unhappenedDailyEvents := make(map[models.Events]bool)
	currentDate := time.Now().UTC().Truncate(24 * time.Hour)
	today := time.Now()
	currentDay := int64(today.Day())
	currentMonth := int64(today.Month())
	currentHour := int64(today.Hour())
	currentMinute := today.Minute()

	for _, event := range events {
		if event.Active {
			activeEvents[event] = event.IsDaily
		}
		if !event.IsDaily && event.Active && event.Day == currentDay && event.Month == currentMonth {
			if event.IsSent {
				UpdateEvent(event.Name, false)
			} else if !event.IsSent && event.SendAt == currentHour && currentMinute == 00 {
				log.Println("=72e334=", "Отправлено", event.Name)
				go SendToEverybody(event)
				UpdateEvent(event.Name, true)
			}
		}
	}
	for activeEvent, isDaily := range activeEvents {
		if isDaily {
			// fmt.Printf("%+v", event)
			if activeEvent.MustSend != currentDate {
				UpdateEvent(activeEvent.Name, false)
				unhappenedDailyEvents[activeEvent] = true
			} else if !activeEvent.IsSent {
				unhappenedDailyEvents[activeEvent] = true
			}
		}

	}
	for event, ok := range unhappenedDailyEvents {
		if ok && event.SendAt == currentHour && currentMinute == 00 {
			birthdays_list, anniversary_list := CreateBirthdaysSlice()
			if event.Name == "День рождения" {
				CheckLogsAndSendEmail(event, birthdays_list)
			} else if event.Name == "Юбилей" {
				CheckLogsAndSendEmail(event, anniversary_list)
			}
			UpdateEvent(event.Name, true)
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
