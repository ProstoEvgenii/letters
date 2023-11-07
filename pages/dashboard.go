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
	"time"

	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// var Ico = Router_struct{
// 	Name:     "Ico",
// 	CollName: Ico_CollName,
// 	Event:    Ico_Event,
// 	Action:   Ico_Action,
// }

func DashboardHandler(rw http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		uploadUsers(rw, request)
		return
	}
	params := new(models.Dashboard_Params)
	if err := schema.NewDecoder().Decode(params, request.URL.Query()); err != nil {
		log.Println("=Params schema Error News_=", err)
	}
	var SendEmailResult string
	if params.SendTo != "" {
		var userTest models.Users
		userTest.FirstName, userTest.LastName, userTest.Email = "Иван", "Иванов", params.SendTo
		SendEmailResult = functions.SendTest(userTest)
	}

	if params.SendAutoAt != 0 {
		var settingsData models.SettingsUpload
		objectId, _ := primitive.ObjectIDFromHex("6540ff760fc1b4b7a36a287b")
		filter := bson.M{
			"_id": objectId,
		}
		update := bson.M{"$set": bson.M{
			"sendAutoAt": params.SendAutoAt,
		}}

		db.InsertIfNotExists(settingsData, filter, update, "settings")

	}

	usersCount, logsCount, birthdaysListLen, todayLogsNumber, sendAutoAt := Dashboard()
	response := models.DashboardGetResponse{
		UsersCount:    usersCount,
		LogsCount:     logsCount,
		CountBirtdays: birthdaysListLen,
		CountLogs:     todayLogsNumber,
		SendEmail:     SendEmailResult,
		SendAutoAt:    sendAutoAt,
	}

	itemCountJson, err := json.Marshal(response)

	if err != nil {
		fmt.Println("error:", err)
	}
	rw.Write(itemCountJson)

	return
}

func Dashboard() (int64, int64, int, int64, int) {
	usersCount := db.CountDocuments(bson.M{}, "users")
	logsCount := db.CountDocuments(bson.M{}, "logs")
	today := time.Now().UTC().Truncate(24 * time.Hour)
	logsLogsToday := db.CountDocuments(bson.M{"dateCreate": today}, "logs")
	sendAutoAt := GetSettings().SendAutoAt
	birthdays_list := functions.CreateBirthdaysSlice()

	return usersCount, logsCount, len(birthdays_list), logsLogsToday, sendAutoAt
}

func uploadUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Max-Age", "15")
	file, _, err := r.FormFile("jsonFile")
	if err != nil {
		http.Error(w, "Не удалось получить файл", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file) // Читаем содержимое файла в срез байтов
	if err != nil {
		http.Error(w, "Не удалось прочитать файл", http.StatusInternalServerError)
		return
	}
	var users []models.UsersUpload //  Форматитирую срез байтов в структуру
	if err := json.Unmarshal(fileBytes, &users); err != nil {
		fmt.Println("Ошибка при разборе JSON:", err)
		return
	}
	var documentsInserted int64
	var documentsModified int64

	for _, document := range users {
		filter := bson.M{
			"E-mail": document.Email,
		}
		dateBirth, _ := time.Parse("01/02/2006", document.Date_birth)
		update := bson.M{"$set": bson.M{
			"Имя":           document.First_name,
			"Фамилия":       document.Last_name,
			"Отчество":      document.Middle_name,
			"Дата рождения": dateBirth,
			"E-mail":        document.Email,
		}}
		documentsInserted += db.InsertIfNotExists(document, filter, update, "users").UpsertedCount
		documentsModified += db.InsertIfNotExists(document, filter, update, "users").ModifiedCount
	}

	response := models.DashboardPostResponse{
		Err:               "Ok",
		DocumentsInserted: documentsInserted,
		DocumentsModified: documentsModified,
	}

	usersAdded, err := json.Marshal(response)
	if err != nil {
		fmt.Println("error:", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(usersAdded)
}

// func FindLogs() []Users {
// 	today := time.Now()
// 	twentyFourHoursAgo := today.Add(-24 * time.Hour)
// 	// filter := bson.M{}
// 	filter := bson.M{"Дата рождения": bson.M{"$gte": twentyFourHoursAgo}}
// 	cursor := Find(filter, "logs")
// 	var users_logs []Users
// 	err := cursor.All(context.TODO(), &users_logs)
// 	if err != nil {
// 		log.Println("=84ce91=", err)
// 	}

// 	// var logs_list []Users
// 	// for _, user := range users_logs {
// 	// 	if user.DateOfBirth.Day() == today.Day() && user.DateOfBirth.Month() == today.Month() {
// 	// 		birthdays_list = append(birthdays_list, user)
// 	// 	}
// 	// }

// 	return
// }
