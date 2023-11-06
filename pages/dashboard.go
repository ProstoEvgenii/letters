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
)

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
		log.Println("=52316b=", userTest)
		SendEmailResult = functions.SendTest(userTest)
	}

	usersCount, logsCount, birthdaysListLen, todayLogsNumber := Dashboard()
	response := models.DashboardGetResponse{
		UsersCount:    usersCount,
		LogsCount:     logsCount,
		CountBirtdays: birthdaysListLen,
		CountLogs:     todayLogsNumber,
		SendEmail:     SendEmailResult,
	}

	itemCountJson, err := json.Marshal(response)

	if err != nil {
		fmt.Println("error:", err)
	}
	rw.Write(itemCountJson)

	return
}

func Dashboard() (int64, int64, int, int64) {
	usersCount := db.CountDocuments(bson.M{}, "users")
	logsCount := db.CountDocuments(bson.M{}, "logs")
	today := time.Now().UTC().Truncate(24 * time.Hour)
	logsLogsToday := db.CountDocuments(bson.M{"dateCreate": today}, "logs")

	birthdays_list := functions.CreateBirthdaysSlice()

	return usersCount, logsCount, len(birthdays_list), logsLogsToday
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
