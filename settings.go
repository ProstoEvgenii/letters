package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

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
		settings := GetSettings()
		response := SettingsUpload{
			Template:   settings.Template,
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

//	func CheckAll(rw http.ResponseWriter, request *http.Request) DashboardPostResponse {
//		var response DashboardPostResponse
//		body, err := io.ReadAll(request.Body)
//		if err != nil {
//			log.Println("=fa78f5=", "Ошибка чтения данных из запроса", "uploadSettings")
//			response.Err = "Ошибка"
//			return response
//		}
//		var settingsData SettingsUpload
//		if err := json.Unmarshal(body, &settingsData); err != nil {
//			log.Println("=324528f5=", "Ошибка разбора данных JSON", "uploadSettings")
//			response.Err = "Ошибка"
//			return response
//		}
//		result := CheckConnectionToEmail(settingsData)
//		if result != "Соединение с почтовым ящиком установлено." {
//			response.Err = "Ошибка"
//			return response
//		}
//		response = uploadSettings(settingsData, rw)
//		return response
//
//	if result != "ok" {
//		response := DashboardPostResponse{
//			Err: result,
//		}
//		errBody, err := json.Marshal(response)
//		if err != nil {
//			fmt.Println("error:", err)
//		}
//		// rw.WriteHeader(http.StatusOK)
//		rw.Write(errBody)
//		return
//	}
//
// }
func GetSettings() SettingsUpload {
	filter := bson.M{}

	var settings SettingsUpload

	cursor := FindOne(filter, "settings")
	cursor.Decode(&settings)

	return settings
}
func CheckConnectionToEmail(settingsData SettingsUpload) string {
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

func uploadSettings(rw http.ResponseWriter, request *http.Request) DashboardPostResponse {
	var response DashboardPostResponse

	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Println("=fa78f5=", "Ошибка чтения данных из запроса", "uploadSettings")
		response.Err = "Ошибка"
		return response
	}

	var settingsData SettingsUpload

	if err := json.Unmarshal(body, &settingsData); err != nil {
		log.Println("=324528f5=", "Ошибка разбора данных JSON", "uploadSettings")
		response.Err = "Ошибка"
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
		"template":   settingsData.Template,
		"emailLogin": settingsData.EmailLogin,
		"emailPass":  settingsData.EmailPass,
		"smtp":       settingsData.Smtp,
		"port":       settingsData.Port,
		"dateCreate": currentDate,
	}}

	settingInserted := InsertIfNotExists(settingsData, filter, update, "settings")

	response = DashboardPostResponse{
		Err:               result,
		DocumentsInserted: settingInserted.UpsertedCount,
		DocumentsModified: settingInserted.ModifiedCount,
	}
	return response

}
