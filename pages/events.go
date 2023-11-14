package pages

import (
	"encoding/json"
	"io"
	"letters/db"
	"letters/models"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func UploadEventsHandler(rw http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		UploadEvents(rw, request)
		return
	}

	return
}
func UploadEvents(rw http.ResponseWriter, request *http.Request) models.DashboardPostResponse {
	var response models.DashboardPostResponse

	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Println("=fa78f5=", "Ошибка чтения данных из запроса", "UploadEvents")
		response.Err = "Ошибка"
		return response
	}

	var eventsData models.Events

	if err := json.Unmarshal(body, &eventsData); err != nil {
		log.Println("=324528f5=", "Ошибка разбора данных JSON", "UploadEvents")
		response.Err = "Ошибка"
		return response
	}
	log.Println("=c6d3a6=", eventsData)
	filter := bson.M{
		"name": eventsData.Name,
	}
	update := bson.M{"$set": bson.M{
		"name":         eventsData.Name,
		"from":         eventsData.From,
		"isDaily":      eventsData.IsDaily,
		"isSent":       eventsData.IsSent,
		"day":          eventsData.Day,
		"month":        eventsData.Month,
		"subject":      eventsData.Subject,
		"sendAt":       eventsData.SendAt,
		"templateName": eventsData.TemplateName,
	}}

	settingInserted := db.InsertIfNotExists(filter, update, "events")
	result := "ok"
	response = models.DashboardPostResponse{
		Err:               result,
		DocumentsInserted: settingInserted.UpsertedCount,
		DocumentsModified: settingInserted.ModifiedCount,
	}
	return response
}
