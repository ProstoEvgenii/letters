package pages

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

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
			log.Println("=Params schema Error =", err)
		}
		_, exists := functions.AuthUsers[params.UUID]
		if !exists {
			return
		}
		var templates []models.Templates
		if params.Templates {
			templates = GetTemplates()
		}

		settings := GetSettings()
		events := GetEvents()

		response := models.SettingsUpload{
			Records:    events,
			Templates:  templates,
			EmailLogin: settings.EmailLogin,
			Smtp:       settings.Smtp,
			Port:       settings.Port,
		}
		settingsJson, err := json.Marshal(response)
		if err != nil {
			fmt.Println("error:", err)
		}
		rw.Write(settingsJson)
		return
	}
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

func GetTemplates() []models.Templates {
	filter := bson.M{}

	var templates []models.Templates

	cursor := db.Find(filter, "templates")
	if err := cursor.All(context.TODO(), &templates); err != nil {
		log.Println("=8922b7=", err)
	}
	var names []string
	for _, template := range templates {
		names = append(names, template.Name)
	}
	return templates
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
		log.Println("", "Ошибка разбора данных JSON", "uploadSettings")
		response.Err = "Ошибка"
		return response
	}
	// log.Println("=bf12b7=", settingsData.UUID)
	_, exists := functions.AuthUsers[settingsData.UUID]
	if !exists {
		response.Err = "Ошибка Авторизации."
		return response
	}
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
		"emailLogin": settingsData.EmailLogin,
		"emailPass":  settingsData.EmailPass,
		"smtp":       settingsData.Smtp,
		"port":       settingsData.Port,
		"dateCreate": currentDate,
	}}

	settingInserted := db.InsertIfNotExists(filter, update, "settings", true)

	response = models.DashboardPostResponse{
		Err:               result,
		DocumentsInserted: settingInserted.UpsertedCount,
		DocumentsModified: settingInserted.ModifiedCount,
	}
	return response
}
