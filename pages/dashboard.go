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

	_, exists := functions.AuthUsers[params.UUID]
	if !exists {
		return
	}

	if params.SendTo != "" {
		var userTest models.Users
		userTest.FirstName, userTest.LastName, userTest.Email = "Иван", "Иванов", params.SendTo
		SendEmailResult = functions.SendTest(userTest, params.SendTemplate)
	}

	usersCount, logsCount, birthdaysListLen, todayLogsNumber, templates := Dashboard()
	response := models.DashboardGetResponse{
		Templates:     templates,
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

func Dashboard() (int64, int64, int, int64, []string) {
	usersCount := db.CountDocuments(bson.M{}, "users")
	logsCount := db.CountDocuments(bson.M{}, "logs")

	today := time.Now().UTC().Truncate(24 * time.Hour)
	logsLogsToday := db.CountDocuments(bson.M{"dateCreate": today}, "logs")
	birthdays_list, anniversary_list := functions.CreateBirthdaysSlice()
	birthdaysNum := len(birthdays_list) + len(anniversary_list)
	templates := GetTemplates()

	return usersCount, logsCount, birthdaysNum, logsLogsToday, templates
}

func uploadUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Max-Age", "15")

	UUID := r.FormValue("UUID")
	_, exists := functions.AuthUsers[UUID]
	if !exists {
		return
	}
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
		documentsInserted += db.InsertIfNotExists(filter, update, "users").UpsertedCount
		documentsModified += db.InsertIfNotExists(filter, update, "users").ModifiedCount
	}
	if documentsInserted != 0 {
		objectId, _ := primitive.ObjectIDFromHex("6548eb240fc1b4b7a3800f31")
		filter := bson.M{
			"_id": objectId,
		}
		var eventBirth models.Events
		result := db.FindOne(filter, "events")
		result.Decode(&eventBirth)

		objectId, _ = primitive.ObjectIDFromHex("65647ad3a62203657bf27b62")
		filter = bson.M{
			"_id": objectId,
		}
		var eventAnniversery models.Events
		result = db.FindOne(filter, "events")
		result.Decode(&eventAnniversery)
		birthdays_list, anniversary_list := functions.CreateBirthdaysSlice()
		if eventBirth.Name == "День рождения" && eventBirth.IsSent == true && eventBirth.Active {
			functions.CheckLogsAndSendEmail(eventBirth, birthdays_list)
		} else if eventAnniversery.Name == "День рождения" && eventAnniversery.IsSent == true && eventAnniversery.Active {
			functions.CheckLogsAndSendEmail(eventAnniversery, anniversary_list)
		}

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
