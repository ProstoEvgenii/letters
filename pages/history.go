package pages

import (
	"context"
	"encoding/json"
	"fmt"
	"letters/db"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func HistoryHandler(rw http.ResponseWriter, request *http.Request) {

	if request.Method == "POST" {

		rw.Write([]byte("Привет"))
	}

	if request.Method == "GET" {
		logsCount := db.CountDocuments("logs")
		today := time.Now().UTC().Truncate(24 * time.Hour)
		yesterday := time.Now().UTC().AddDate(0, 0, -1).Truncate(24 * time.Hour)
		// tomorrow := time.Now().UTC().AddDate(0, 0, 1).Truncate(24 * time.Hour)
		todayLogsNumber := getLogs(today)
		yesterdayLogsNumber := getLogs(yesterday)

		cursor := Find(bson.M{}, "logs")
		var logsSlice []Logs
		if err := cursor.All(context.TODO(), &logsSlice); err != nil {
			log.Println("Cursor All Error Database", err)
			rw.Write([]byte("{}"))
			return
		}
		if len(logsSlice) == 0 {
			rw.Write([]byte("{}"))
			return
		}

		response := GetHistoryResponse{
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
