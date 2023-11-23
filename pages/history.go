package pages

import (
	"context"
	"encoding/json"
	"fmt"
	"letters/db"
	"letters/functions"
	"letters/models"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson"
)

func HistoryHandler(rw http.ResponseWriter, request *http.Request) {
	response := models.GetHistoryResponse{}
	if request.Method == "POST" {

		rw.Write([]byte("Привет"))
	}

	if request.Method == "GET" {
		params := new(models.Dashboard_Params)
		if err := schema.NewDecoder().Decode(params, request.URL.Query()); err != nil {
			log.Println("=Params schema Error News_=", err)
		}

		filter := bson.M{}

		if params.Seach != "" {
			filter = bson.M{
				"$or": []bson.M{
					{"Имя": bson.M{"$regex": params.Seach, "$options": "i"}},
					{"Фамилия": bson.M{"$regex": params.Seach, "$options": "i"}},
					{"E-mail": bson.M{"$regex": params.Seach, "$options": "i"}},
				},
			}
		}
		if params.UUID != "" {
			_, exists := functions.AuthUsers[params.UUID]
			if !exists {
				return
			}
		}

		logsCount := db.CountDocuments(bson.M{}, "logs")
		today := time.Now().UTC().Truncate(24 * time.Hour)
		yesterday := time.Now().UTC().AddDate(0, 0, -1).Truncate(24 * time.Hour)
		todayLogsNumber := getLogs(today)
		yesterdayLogsNumber := getLogs(yesterday)

		cursor := db.Find(filter, "logs")
		var logsSlice []models.Logs
		if err := cursor.All(context.TODO(), &logsSlice); err != nil {
			log.Println("Cursor All Error Database", err)
			rw.Write([]byte("{}"))
			return
		}
		if len(logsSlice) == 0 {
			rw.Write([]byte("{}"))
			return
		}

		response = models.GetHistoryResponse{
			Records:        logsSlice,
			LogsCount:      logsCount,
			TodayLogsCount: todayLogsNumber,
			// TommorowLogsCount:  tomorrowLogsNumber,
			YesterdayLogsCount: yesterdayLogsNumber,
		}
		dataBaseJson, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error:", err)
			rw.Write([]byte("{}"))
			return
		}
		rw.Write(dataBaseJson)
		return
	}
}

func getLogs(date time.Time) int {
	filter := bson.M{
		"dateCreate": date,
	}
	cursor := db.Find(filter, "logs")
	var logs []models.Logs
	if err := cursor.All(context.TODO(), &logs); err != nil {
		log.Println("=8922b7=", err)
	}
	return len(logs)

}
