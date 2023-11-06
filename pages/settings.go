package pages

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"letters/db"
	"letters/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/gomail.v2"
)

func SettingsHandler(rw http.ResponseWriter, request *http.Request) {

	if request.Method == "POST" {
		response := uploadSettings(rw, request)
		log.Println("=ad4c75=", response)

		settingsAdded, err := json.Marshal(response)
		if err != nil {
			fmt.Println("error:", err)
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write(settingsAdded)
	}

	if request.Method == "GET" {
		params := new(models.Dashboard_Params)
		if err := schema.NewDecoder().Decode(params, request.URL.Query()); err != nil {
			log.Println("=Params schema Error News_=", err)
		}
		if params.SendTo != "" {

		}

		settings := GetSettings()
		events := GetEvents()

		settings_response := models.SettingsUpload{
			Template:   settings.Template,
			EmailLogin: settings.EmailLogin,
			Smtp:       settings.Smtp,
			Port:       settings.Port,
		}

		response := models.GetSettingsResponse{
			EventsRecords: events,
			Settings:      settings_response,
		}
		settingsJson, err := json.Marshal(response)
		if err != nil {
			fmt.Println("error:", err)
		}
		rw.Write(settingsJson)
		return
	}
}

func GetSettings() models.SettingsUpload {
	filter := bson.M{}

	var settings models.SettingsUpload

	cursor := db.FindOne(filter, "settings")
	cursor.Decode(&settings)

	return settings
}
func CheckConnectionToEmail(settingsData models.SettingsUpload) string {
	port, err := strconv.Atoi(settingsData.Port)
	if err != nil {
		fmt.Println("SendEmail Ошибка форматирования строки в int")
	}

	d := gomail.NewDialer(settingsData.Smtp, port, settingsData.EmailLogin, settingsData.EmailPass)
	if err := d.DialAndSend(); err != nil {
		log.Println("=51d73a=", err)
		return "Ошибка при подключении к почте."
	}
	log.Println("Соединение с почтовым ящиком установлено.")
	return "Соединение с почтовым ящиком установлено."
}

func uploadSettings(rw http.ResponseWriter, request *http.Request) models.DashboardPostResponse {
	var response models.DashboardPostResponse

	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Println("=fa78f5=", "Ошибка чтения данных из запроса", "uploadSettings")
		response.Err = "Ошибка"
		return response
	}

	var settingsData models.SettingsUpload

	if err := json.Unmarshal(body, &settingsData); err != nil {
		log.Println("=324528f5=", "Ошибка разбора данных JSON", "uploadSettings")
		response.Err = "Ошибка"
		return response
	}
	// log.Println("=09c43c=", settingsData)
	result := CheckConnectionToEmail(settingsData)

	if result != "Соединение с почтовым ящиком установлено." {
		response.Err = result
		return response
	}

	currentDate := time.Now().UTC().Truncate(24 * time.Hour)
	objectId, _ := primitive.ObjectIDFromHex("6540ff760fc1b4b7a36a287b")
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

	settingInserted := db.InsertIfNotExists(settingsData, filter, update, "settings")

	response = models.DashboardPostResponse{
		Err:               result,
		DocumentsInserted: settingInserted.UpsertedCount,
		DocumentsModified: settingInserted.ModifiedCount,
	}
	return response

}

func EventsHandler(rw http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		response := uploadEvents(rw, request)

		settingsAdded, err := json.Marshal(response)
		if err != nil {
			fmt.Println("error:", err)
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write(settingsAdded)
	}
}
func GetEvents() []models.Events {
	cursor := db.Find(bson.M{}, "events")
	var eventsSlice []models.Events
	if err := cursor.All(context.TODO(), &eventsSlice); err != nil {
		log.Println("Cursor All Error Events", err)
	}
	return eventsSlice
}

func uploadEvents(rw http.ResponseWriter, request *http.Request) models.DashboardPostResponse {
	var response models.DashboardPostResponse

	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Println("=fa78f5=", "Ошибка чтения данных из запроса", "uploadEvents", err)
		response.Err = "Ошибка"
		return response
	}

	var eventsData models.EventUpload

	if err := json.Unmarshal(body, &eventsData); err != nil {
		log.Println("=324528f5=", "Ошибка разбора данных JSON", "uploadEvents", err)
		response.Err = "Ошибка"
		return response
	}
	date, errDay := strconv.ParseInt(eventsData.Date, 10, 64)
	month, errMonth := strconv.ParseInt(eventsData.Month, 10, 64)
	SendAt, errSendAt := strconv.ParseInt(eventsData.SendAt, 10, 64)
	if errDay != nil || errMonth != nil || errSendAt != nil {
		log.Println("=0d19ba=", "Ошибка форматирования в int даты и месяца", errDay, errMonth)
	}

	currentDate := time.Now().UTC().Truncate(24 * time.Hour)
	filter := bson.M{
		"name": eventsData.Name,
	}
	update := bson.M{"$set": bson.M{
		"templateName": eventsData.TemplateName,
		"sendAt":       SendAt,
		"isDaily":      eventsData.IsDaily,
		"date":         date,
		"month":        month,
		"dateCreate":   currentDate,
		"isSent":       false,
	}}

	settingInserted := db.InsertIfNotExists(eventsData, filter, update, "events")

	response = models.DashboardPostResponse{
		Err:               "ok",
		DocumentsInserted: settingInserted.UpsertedCount,
		DocumentsModified: settingInserted.ModifiedCount,
	}
	return response
}



// response := map[string]string{"message": "Дата получена успешно"}
//             json.NewEncoder(w).Encode(response)
