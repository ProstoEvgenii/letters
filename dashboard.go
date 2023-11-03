package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	if params.SendTo != "" {
		var userTest Users
		userTest.FirstName, userTest.LastName, userTest.Email = "Иван", "Иванов", params.SendTo
		SendEmailResult = SendTest(userTest)
	}
	if params.SendAutoAt != 0 {
		var settingsData SettingsUpload
		objectId, _ := primitive.ObjectIDFromHex("6540ff760fc1b4b7a36a287b")
		filter := bson.M{
			"_id": objectId,
		}
		update := bson.M{"$set": bson.M{
			"sendAutoAt": params.SendAutoAt,
		}}

		InsertIfNotExists(settingsData, filter, update, "settings")

	}

	usersCount, logsCount, birthdaysListLen, todayLogsNumber, sendAutoAt := Dashboard()
	response := DashboardGetResponse{
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

func Dashboard() (int64, int64, int, int, int) {
	settings := GetSettings()
	usersCount := CountDocuments("users")
	logsCount := CountDocuments("logs")

	birthdays_list := CreateBirthdaysSlice()
	// currentDate := today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	todayLogsNumber := getLogs(today)
	// GetTemplate("test1")
	return usersCount, logsCount, len(birthdays_list), todayLogsNumber, settings.SendAutoAt
}

// func CheckSettingsAndEmail() string{
// 	settings := GetSettings()
// 	if settings.EmailLogin == "" || settings.EmailPass == "" || settings.Smtp == "" || settings.Port == "" || settings.Template == "" {
// 		log.Println("=82842e=", "Настройки не верны либо отсутствуют.")
// 		return "Настройки не верны либо отсутствуют."
// 	}

// 	html := GetTemplate(settings.Template)
// 	if html == "" {
// 		return fmt.Sprintf("Шаблона %s не существует", settings.Template)
// 	}
// 	return "ok"
// }

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

	cursor := FindOne(filter, "templates")

	if cursor.Err() != nil {
		log.Println("=ce7969=", cursor.Err())
		return ""
	}
	var template Templates
	cursor.Decode(&template)

	return template.IndexHTML

}

func SendTest(user Users) string {
	settings := GetSettings()
	if settings.EmailLogin == "" || settings.EmailPass == "" || settings.Smtp == "" || settings.Port == "" || settings.Template == "" {
		log.Println("=82842e=", "Настройки не верны либо отсутствуют.")
		return "Настройки не верны либо отсутствуют."
	}

	html := GetTemplate(settings.Template)
	if html == "" {
		return fmt.Sprintf("Шаблона %s не существует", settings.Template)
	}

	err := SendEmail(user, settings, html)
	if err != "ok" {
		return err
	}

	return fmt.Sprintf("Пользователь %s поздравлен", user.Email)
}
func checkLogsAndSendEmail() string {
	birthdays_list := CreateBirthdaysSlice()
	if len(birthdays_list) == 0 {
		return "Нет Дней рождений сегодня"
	}

	emailSent := 0

	settings := GetSettings()
	if settings.EmailLogin == "" || settings.EmailPass == "" || settings.Smtp == "" || settings.Port == "" || settings.Template == "" {
		log.Println("=82842e=", "Настройки не верны либо отсутствуют.")
		return "Настройки не верны либо отсутствуют."
	}

	html := GetTemplate(settings.Template)
	if html == "" {
		return fmt.Sprintf("Шаблона %s не существует", settings.Template)
	}

	for _, user := range birthdays_list {
		result := CreateLog(user)
		if result != 0 {
			//Если результат создания лога == 0 ,значит лог с таким email существует и поздравлять его не нужно
			err := SendEmail(user, settings, html)
			if err != "ok" {
				return err
			}

			emailSent += 1
		}
	}
	if emailSent == 0 {
		return "Сегодня все поздравлены"
	} else {
		log.Printf("Поздравлено %d пользователей", emailSent)
		return fmt.Sprintf("Поздравлено %d пользователей", emailSent)
	}

}

func SendEmail(user Users, settings SettingsUpload, html string) string {
	first_name := user.FirstName
	last_name := user.LastName
	subject := "C днем рождения! От главы администрации."

	replacer := strings.NewReplacer("${first_name}", first_name, "${last_name}", last_name)

	port, err := strconv.Atoi(settings.Port)
	if err != nil {
		fmt.Println("SendEmail Ошибка форматирования строки в int")
		return "Ошибка"
	}

	html = replacer.Replace(html)

	m := gomail.NewMessage()
	m.SetHeader("From", settings.EmailLogin)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", html)

	d := gomail.NewDialer(settings.Smtp, port, settings.EmailLogin, settings.EmailPass)
	if err := d.DialAndSend(m); err != nil {
		log.Println("=SendEmail79fc04 Отправка письма=", err)
		return "Ошибка при отправкe сообщения"
	}
	fmt.Printf("Поздравление отправлено:%s", user.Email)
	return "ok"
}

func getLogs(date time.Time) int {
	
	filter := bson.M{
		"dateCreate": date,
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
		documentsInserted += InsertIfNotExists(document, filter, update, "users").UpsertedCount
		documentsModified += InsertIfNotExists(document, filter, update, "users").ModifiedCount
	}

	if documentsInserted != 0 {
		result := getStatusToday()
		if result.IsSent {
			checkLogsAndSendEmail()
		}

	}
	response := DashboardPostResponse{
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
