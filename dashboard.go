package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/gomail.v2"
)

func DashboardHandler(rw http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		uploadUsers(rw, request)
		return
	}
	params := new(Dashboard_Params)
	if err := schema.NewDecoder().Decode(params, request.URL.Query()); err != nil {
		log.Println("=Params schema Error News_=", err)
	}
	var SendEmailResult string
	if params.SendAll == "true" {
		SendEmailResult = checkLogsAndSendEmail()
		log.Println("=cd010b=", "тут")
	}

	usersCount, birthdaysListLen, todayLogsNumber := Dashboard()
	response := DashboardGetResponse{
		DocumentsCount: usersCount,
		CountBirtdays:  birthdaysListLen,
		CountLogs:      todayLogsNumber,
		SendEmail:      SendEmailResult,
	}

	itemCountJson, err := json.Marshal(response)
	if err != nil {
		fmt.Println("error:", err)
	}
	rw.Write(itemCountJson)

	return
}

func Dashboard() (int64, int, int) {
	usersCount := CountDocuments()
	birthdays_list := CreateBirthdaysSlice()
	todayLogsNumber := getLogs()
	// GetTemplate("test1")
	return usersCount, len(birthdays_list), todayLogsNumber
}

func CreateBirthdaysSlice() []Users {
	today := time.Now()
	filter := bson.M{}
	cursor := Find(filter, "users")
	var users []Users
	err := cursor.All(context.TODO(), &users)
	if err != nil {
		log.Println("=84ce91=", err)
	}
	var birthdays_list []Users
	for _, user := range users {
		if user.DateOfBirth.Day() == today.Day() && user.DateOfBirth.Month() == today.Month() {
			birthdays_list = append(birthdays_list, user)
		}
	}

	return birthdays_list
}

func GetTemplate(templateName string) string {
	filter := bson.M{
		"name": templateName,
	}
	cursor := Find(filter, "templates")
	var template []Templates
	if err := cursor.All(context.TODO(), &template); err != nil {
		log.Println("=8922b7=", err)
	}
	if len(template) > 0 {
		return template[0].Name
	} else {
		log.Println("GetTemplate: template[0].Name не обнаружен")
		return ""
	}

}

func checkLogsAndSendEmail() string {
	var result string
	birthdays_list := CreateBirthdaysSlice()
	if len(birthdays_list) == 0 {
		result = "Нет Дней рождений сегодня"
		return result
	}

	emailSent := 0

	for _, user := range birthdays_list {
		result := CreateLog(user)
		if result != 0 {
			//Если результат создания лога == 0 ,значит лог с таким email существует и поздравлять его не нужно
			SendEmail(user)
			emailSent += 1
		}
	}
	if emailSent == 0 {
		result = "Сегодня все поздравлены"
		return result
	} else {
		result = fmt.Sprintf("Поздравлено %d пользователей", emailSent)
		return result
	}

}

func SendEmail(user Users) {
	first_name := user.FirstName
	last_name := user.LastName
	subject := "C днем рождения!"

	replacer := strings.NewReplacer("${first_name}", first_name, "${last_name}", last_name)

	settings := GetSettings()
	port, err := strconv.Atoi(settings.Port)
	if err != nil {
		fmt.Println("=SendEmail Ошибка форматирования строки в int=")
	}
	html := GetTemplate(settings.Template)
	html = replacer.Replace(html)

	fmt.Println("=950d99=", html)
	fmt.Println("=ad2a17=", settings)

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL"))
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", html)

	d := gomail.NewDialer(settings.Smtp, port, settings.EmailLogin, settings.EmailPass)
	if err := d.DialAndSend(m); err != nil {

		time.Sleep(10 * time.Second)
		log.Fatal()
	}
	fmt.Printf("Поздравление отправлено:%s", user.Email)
	time.Sleep(10 * time.Second)
}

func getLogs() int {
	currentDate := time.Now().UTC().Truncate(24 * time.Hour)
	filter := bson.M{
		"dateCreate": currentDate,
	}
	cursor := Find(filter, "logs")
	var logs []Logs
	if err := cursor.All(context.TODO(), &logs); err != nil {
		log.Println("=8922b7=", err)
	}
	return len(logs)

}
func CreateLog(user Users) int64 {
	currentDate := time.Now().UTC().Truncate(24 * time.Hour)
	filter := bson.M{
		"E-mail":     user.Email,
		"dateCreate": currentDate,
	}
	update := bson.M{"$setOnInsert": bson.M{
		"Имя":           user.FirstName,
		"Фамилия":       user.LastName,
		"Отчество":      user.MiddleName,
		"Дата рождения": user.DateOfBirth,
		"E-mail":        user.Email,
		"dateCreate":    currentDate,
	}}
	result := InsertIfNotExists(user, filter, update, "logs").UpsertedCount
	return result
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
	var users []UsersUpload //  Форматитирую срез байтов в структуру
	if err := json.Unmarshal(fileBytes, &users); err != nil {
		fmt.Println("Ошибка при разборе JSON:", err)
		return
	}
	var documentsInserted int64

	for _, document := range users {
		filter := bson.M{
			"E-mail": document.Email,
		}
		dateBirth, _ := time.Parse("01/02/2006", document.Date_birth)
		update := bson.M{"$setOnInsert": bson.M{
			"Имя":           document.First_name,
			"Фамилия":       document.Last_name,
			"Отчество":      document.Middle_name,
			"Дата рождения": dateBirth,
			"E-mail":        document.Email,
		}}
		documentsInserted += InsertIfNotExists(document, filter, update, "users").UpsertedCount
	}

	response := DashboardPostResponse{
		DocumentsInserted: documentsInserted,
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
