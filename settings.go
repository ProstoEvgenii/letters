package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SettingsHandler(rw http.ResponseWriter, request *http.Request) {

	if request.Method == "POST" {
		uploadSettings(rw, request)
	}
	if request.Method == "GET" {
		// log.Println("=e1c06c=", "GET")
		settings := GetSettings()
		response := SettingsUpload{
			Template:   settings.Template,
			EmailLogin: settings.EmailLogin,
			Smtp:       settings.Smtp,
			Port:       settings.Port,
		}
		GetSettings()
		settingsJson, err := json.Marshal(response)
		if err != nil {
			fmt.Println("error:", err)
		}
		rw.Write(settingsJson)
		return
	}
}
func GetSettings() SettingsUpload {
	objectId, err := primitive.ObjectIDFromHex("6540ff760fc1b4b7a36a287b")
	if err != nil {
		fmt.Println("=getSettings Ошибка преобразования ID=", err)
	}
	filter := bson.M{
		"_id": objectId,
	}
	cursor := Find(filter, "settings")

	var settings []SettingsUpload

	if err := cursor.All(context.TODO(), &settings); err != nil {
		log.Println("=8922b7=", err)
	}
	log.Println("=499684=", settings)
	return settings[0]

}
func uploadSettings(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Max-Age", "15")
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(rw, "Ошибка чтения данных из запроса", http.StatusBadRequest)
		return
	}

	var settingsData SettingsUpload
	if err := json.Unmarshal(body, &settingsData); err != nil {
		http.Error(rw, "Ошибка разбора данных JSON", http.StatusBadRequest)
		return
	}

	log.Println("=5cd991=", settingsData)
	currentDate := time.Now().UTC().Truncate(24 * time.Hour)
	objectId, err := primitive.ObjectIDFromHex("6540ff760fc1b4b7a36a287b")
	filter := bson.M{
		"_id": objectId,
	}
	update := bson.M{"$set": bson.M{
		"template":   settingsData.Template,
		"emailLogin": settingsData.EmailLogin,
		"emailPass":  settingsData.EmailPass,
		"smtp":       settingsData.Smtp,
		"port":       settingsData.Port,
		"dateCreate": currentDate,
	}}
	settingInserted := InsertIfNotExists(settingsData, filter, update, "settings")
	log.Println("=ModifiedCount=", settingInserted.ModifiedCount, "=UpsertedCount=", settingInserted.UpsertedCount, "=UpsertedCount=", settingInserted.UpsertedID)
}
