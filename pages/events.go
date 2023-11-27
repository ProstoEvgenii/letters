package pages

import (
	"encoding/json"
	"fmt"
	"io"
	"letters/db"
	"letters/functions"
	"letters/models"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func UploadEventsHandler(rw http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {

		response := UploadEvents(rw, request)
		eventAdded, err := json.Marshal(response)
		if err != nil {
			fmt.Println("error:", err)
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write(eventAdded)
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

	log.Println("=2f49d7=", eventsData.UUID)
	_, exists := functions.AuthUsers[eventsData.UUID]
	if !exists {
		response.Err = "Ошибка Авторизации"
		return response
	}
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
		"active":       eventsData.Active,
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
