package functions

import (
	"context"
	"fmt"
	"letters/db"
	"letters/models"
	"log"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/gomail.v2"
)

func GetSettings() models.SettingsUpload {
	filter := bson.M{}

	var settings models.SettingsUpload

	cursor := db.FindOne(filter, "settings")
	cursor.Decode(&settings)

	return settings
}
func CreateBirthdaysSlice() []models.Users {
	today := time.Now()
	filter := bson.M{}
	cursor := db.Find(filter, "users")
	var users []models.Users
	err := cursor.All(context.TODO(), &users)
	if err != nil {
		log.Println("=84ce91=", err)
	}
	var birthdays_list []models.Users
	for _, user := range users {
		if user.DateOfBirth.Day() == today.Day() && user.DateOfBirth.Month() == today.Month() && !user.Unsubscribe {
			birthdays_list = append(birthdays_list, user)
		}
	}

	return birthdays_list
}

func GetTemplate(templateName string) string {
	filter := bson.M{
		"name": templateName,
	}

	cursor := db.FindOne(filter, "templates")

	if cursor.Err() != nil {
		log.Println("=ce7969=", cursor.Err())
		return ""
	}
	var template models.Templates
	cursor.Decode(&template)

	return template.IndexHTML

}

func SendTest(user models.Users, templateName string) string {
	settings := GetSettings()
	if settings.EmailLogin == "" || settings.EmailPass == "" || settings.Smtp == "" || settings.Port == "" {
		log.Println("=82842e=", "Настройки не верны либо отсутствуют.")
		return "Настройки не верны либо отсутствуют."
	}

	html := GetTemplate(templateName)
	if html == "" {
		return fmt.Sprintf("Шаблона %s не существует", templateName)
	}

	err := SendEmail(user, "Teстовое со", html, settings)
	if err != "ok" {
		return err
	}

	return fmt.Sprintf("Пользователь %s поздравлен", user.Email)
}

// Пока работает только с днем рождения.
func CheckLogsAndSendEmail(event models.Events) string {
	birthdays_list := CreateBirthdaysSlice() //Получил выборку пользователей с днем рождения
	if len(birthdays_list) == 0 {
		log.Println("=91c8c4=", "Нет Дней рождений сегодня")
		return "Нет Дней рождений сегодня"
	}
	html := GetTemplate(event.TemplateName)
	settings := GetSettings()
	emailSent := 0
	for _, user := range birthdays_list {
		result := CreateLog(user, event.Name)
		if result != 0 {
			//Если результат создания лога == 0 ,значит лог с таким email существует и поздравлять его не нужно
			err := SendEmail(user, event.Subject, html, settings)
			if err != "ok" {
				return err
			}
			emailSent += 1
		}
	}
	if emailSent == 0 {
		log.Println("=5c58cc=", "Сегодня все поздравлены")
		return "Сегодня все поздравлены"
	} else {
		log.Printf("Поздравлено %d пользователей", emailSent)
		return fmt.Sprintf("Поздравлено %d пользователей", emailSent)
	}
}

func SendToEverybody(event models.Events) {
	filter := bson.M{}
	cursor := db.Find(filter, "users")
	var users []models.Users
	err := cursor.All(context.TODO(), &users)
	if err != nil {
		log.Println("=84ce91=", err)
	}

	html := GetTemplate(event.TemplateName)
	settings := GetSettings()
	for _, user := range users {
		CreateLog(user, event.Name)
		SendEmail(user, event.Subject, html, settings)
	}
}
func SendEmail(user models.Users, subject string, html string, settings models.SettingsUpload) string {

	first_name := user.FirstName
	last_name := user.LastName

	replacer := strings.NewReplacer("${first_name}", first_name, "${last_name}", last_name, "${email}", user.Email)

	port, err := strconv.Atoi(settings.Port)
	if err != nil {
		fmt.Println("SendEmail Ошибка форматирования строки в int")
		return "Ошибка"
	}

	html = replacer.Replace(html)

	m := gomail.NewMessage()
	from := m.FormatAddress(settings.EmailLogin, "Глава администрации")
	m.SetHeader("From", from)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", html)

	d := gomail.NewDialer(settings.Smtp, port, settings.EmailLogin, settings.EmailPass)
	if err := d.DialAndSend(m); err != nil {
		log.Println("=SendEmail79fc04 Отправка письма=", err)
		return "Ошибка при отправкe сообщения"
	}
	fmt.Printf("Поздравление отправлено:%s\n", user.Email)
	return "ok"
}
func CreateLog(user models.Users, eventName string) int64 {
	currentDate := time.Now().UTC().Truncate(24 * time.Hour)
	filter := bson.M{
		"event":      eventName,
		"E-mail":     user.Email,
		"dateCreate": currentDate,
	}
	update := bson.M{"$setOnInsert": bson.M{
		"event":         eventName,
		"Имя":           user.FirstName,
		"Фамилия":       user.LastName,
		"Отчество":      user.MiddleName,
		"Дата рождения": user.DateOfBirth,
		"E-mail":        user.Email,
		"dateCreate":    currentDate,
	}}
	result := db.InsertIfNotExists(filter, update, "logs").UpsertedCount
	return result
}
