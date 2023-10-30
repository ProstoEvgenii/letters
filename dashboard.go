package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/gomail.v2"
)

func DashboardHandler(rw http.ResponseWriter, request *http.Request) {
	params := new(Dashboard_Params)
	if err := schema.NewDecoder().Decode(params, request.URL.Query()); err != nil {
		log.Println("=Params schema Error News_=", err)
	}

	var SendEmailResult string
	if params.SendAll == "true" {
		SendEmailResult = checkLogsAndSendEmail()
	}

	usersCount, birthdaysListLen, todayLogsNumber := Dashboard()
	response := Response{
		DocumentsCount: usersCount,
		CountBirtdays:  birthdaysListLen,
		CountLogs:      todayLogsNumber,
		SendEmail:      SendEmailResult,
	}

	itemCountJson, err := json.Marshal(response)
	if err != nil {
		fmt.Println("error:", err)
	}
	// log.Println("=42687c=", string(itemCountJson))
	rw.Write(itemCountJson)
	return
}
func Dashboard() (int64, int, int) {
	usersCount := CountDocuments()
	birthdays_list := CreateBirthdaysSlice()
	todayLogsNumber := getLogs()
	// checkLogsAndSendEmail()
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

func GetTemplate() string {
	filter := bson.M{
		"name": "test1",
	}
	cursor := Find(filter, "templates")
	var template []Templates
	if err := cursor.All(context.TODO(), &template); err != nil {
		log.Println("=8922b7=", err)
	}
	return template[0].IndexHTML
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

	// htmlBytes, err := os.ReadFile("index.html")
	// if err != nil {
	// 	// fmt.Println("Ошибка при чтении файла index.html:", err)
	// 	log.Fatal()
	// 	return
	// }
	html := GetTemplate()
	html = replacer.Replace(html)

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL"))
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", html)

	d := gomail.NewDialer("smtp.mail.ru", 465, os.Getenv("EMAIL"), os.Getenv("EMAIL_PASS"))
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
	// log.Println("=Сегодня поздравлено=", logs)
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
