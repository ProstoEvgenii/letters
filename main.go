package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

type Data struct {
	Last_name   string `json:"Фамилия"`
	First_name  string `json:"Имя"`
	Middle_name string `json:"Отчество"`
	Date_birth  string `json:"Дата рождения"`
	Email       string `json:"E-mail"`
}

func main() {
	files_name := []string{"index.html", ".env", "users.json"}

	CheckFilesAndConnectToEmail(files_name)

	records := readJson("users.json")
	todayMonthDate := time.Now().Format("01/02")
	foundBirthday := false
	for _, item := range records {
		if item.Date_birth[:5] == todayMonthDate {
			checkAndLog(item)
			foundBirthday = true
		}
	}

	if !foundBirthday {
		log.Println("Сегодня нет дней рождений среди пользователей.")
		time.Sleep(10 * time.Second)
		log.Fatal()
	}
}

func CheckFilesAndConnectToEmail(files_name []string) {
	fmt.Println("Проверяю файлы в папке...")

	for _, item := range files_name {
		if _, err := os.Stat(item); os.IsNotExist(err) {
			fmt.Printf("Файл %s не найден в корне проекта.\n", item)
			time.Sleep(10 * time.Second)
			log.Fatal()
		}
	}

	log.Println("Все файлы на присутствуют.")

	if err := godotenv.Load(".env"); err != nil {
		log.Println("Файл .env не найден")
		time.Sleep(10 * time.Second)
		log.Fatal()
	}

	// d := gomail.NewDialer("smtp.yandex.ru", 465, "support@crypto-emergency.com", os.Getenv("EMAIL_PASS"))
	// if err := d.DialAndSend(); err != nil {
	// 	log.Printf("Не удалось отправить установить соединение с почтовым ящиком. Убедитесь ,что E-mail и пароль в файле .env указаны верно \n%v", err)
	// 	time.Sleep(10 * time.Second)
	// 	log.Fatal(err)
	// } else {
	// 	log.Println("Соединение с почтовым ящиком установлено.")
	// }
}

func readJson(jsonName string) []Data {
	file, err := os.Open(jsonName)
	if err != nil {
		fmt.Println("Не найден файл БД users.json")
		time.Sleep(10 * time.Second)
		log.Fatal()
	}
	defer file.Close()

	var records []Data
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&records)
	if err != nil {
		if err.Error() == "EOF" && len(records) == 0 {
			fmt.Println("Файл users.json пуст")
			time.Sleep(10 * time.Second)
			log.Fatal()
		} else {
			fmt.Println("Убедитесь, что вы используете верную базу данных", err)
			time.Sleep(10 * time.Second)
			log.Fatal()
		}
	}
	return records
}

func checkAndLog(item Data) {
	log_name := time.Now().Format("01.02.2006")

	existingLogs, err := os.Open("./logs/" + log_name + ".json")
	if os.IsNotExist(err) { //Если файла log не существует создаю и записываю в него item
		create_log(item)
		SendEmailReg(item)
		return
	}
	defer existingLogs.Close()

	var logs []Data
	decoder := json.NewDecoder(existingLogs)
	err = decoder.Decode(&logs)
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return
	}

	itemAlreadyExists := false
	for _, log := range logs {
		if log == item {
			itemAlreadyExists = true
			break
		}
	}

	if !itemAlreadyExists { //Если нет в логах - поздравить и записать в логи
		SendEmailReg(item)
		logs = append(logs, item)
		logJson, err := json.Marshal(logs)
		overwriteLogs, err := os.Create("./logs/" + log_name + ".json")
		if err != nil {
			fmt.Println("Unable to create file:", err)
			os.Exit(1)
		}
		_, err = overwriteLogs.Write(logJson)
		if err != nil {
			log.Println("=Ошибка записи в json=", err)
		}
		defer overwriteLogs.Close()
	}

}
func create_log(item Data) {
	log_name := time.Now().Format("01.02.2006")

	newLog := []Data{item} //Форматирую item в Json
	logJson, err := json.Marshal(newLog)
	if err != nil {

		// log.Fatal("=Ошибка форматирования лога в json=", err)
	}

	if err := os.MkdirAll("./logs", os.ModePerm); err != nil {
		log.Fatal("Ошибка при создании директории logs:", err)
	}
	newLogs, err := os.Create("./logs/" + log_name + ".json")
	if err != nil {
		fmt.Println("Unable to create file:", err)
	}
	_, err = newLogs.Write(logJson)
	if err != nil {
		log.Println("=Ошибка записи в json=", err)
	}
	defer newLogs.Close()
}

func SendEmailReg(item Data) {
	subject := "C др!"
	first_name := item.First_name
	last_name := item.Last_name

	replacer := strings.NewReplacer("${first_name}", first_name, "${last_name}", last_name)

	htmlBytes, err := os.ReadFile("index.html")
	if err != nil {
		fmt.Println("Ошибка при чтении файла index.html:", err)
		return
	}
	html := string(htmlBytes)
	html = replacer.Replace(html)

	m := gomail.NewMessage()
	m.SetHeader("From", "support@crypto-emergency.com")
	log.Println("=30747d=", item.Email)
	m.SetHeader("To", item.Email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", html)

	// html = strings.ReplaceAll(html, "${s_uvazheniem_glava_munici.png}", s_uvazheniem_glava_munici)
	// html = strings.ReplaceAll(html, "${image2}", image2)
	// html = strings.ReplaceAll(html, "${image3}", image3)
	m.Embed("./images/")
	m.Embed("image1.png")
	m.Embed("image2.png", "/images/img_4424.png")

	related.Embed("image1", "https://example.com/image1.jpg")

	d := gomail.NewDialer("smtp.yandex.ru", 465, "support@crypto-emergency.com", os.Getenv("EMAIL_PASS"))
	if err := d.DialAndSend(m); err != nil {
		// log.Println("Error SendEmailReg", err)
		time.Sleep(10 * time.Second)
		log.Fatal(err)
	}

}

func getImageBase64(imagePath string) string {
	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		log.Println("Ошибка при чтении изображения:", err)
		return ""
	}
	base64Image := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(imageBytes)
	return base64Image
}

// log.Println("=fba203=", reflect.TypeOf(today))
// date := "11/06/1969"
// t, err := time.Parse("02/01/2006", date)
// fmt.Println("=008c37=", "День", t.Day(), "Месяц", t.Month())
// fmt.Println("=008c37=", "День", today.Day(), "Месяц", today.Month())

// reflect.TypeOf(tst)
// \n
// today := time.Now()
// today.Day() == birthdate.Day() && today.Month()
