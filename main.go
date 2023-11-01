package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Файл .env не найден")
		time.Sleep(10 * time.Second)
		log.Fatal()
	}
	Connect()
	Start()

}

// func main() {
// files_name := []string{"index.html", ".env", "users.json"}
// CheckFilesAndConnectToEmail(files_name)
// records := readJson("users.json")
// findBirthdays(records)

// Connect()
// GetTemplate()
// Dashboard()
// Start()
// }
// func findBirthdays(records []Data) {
// 	todayMonthDate := time.Now().Format("01/02")
// 	foundBirthday := false
// 	for _, item := range records {
// 		if item.Date_birth[:5] == todayMonthDate {
// 			checkAndLog(item)
// 			foundBirthday = true
// 		}
// 	}

// 	if !foundBirthday {
// 		log.Println("Сегодня нет дней рождений среди пользователей.")
// 		time.Sleep(10 * time.Second)
// 		log.Fatal()
// 	}
// }

// func CheckFilesAndConnectToEmail(files_name []string) {
// fmt.Println("Проверяю файлы в папке...")

// for _, item := range files_name {
// 	if _, err := os.Stat(item); os.IsNotExist(err) {
// 		fmt.Printf("Файл %s не найден в корне проекта.\n", item)
// 		time.Sleep(10 * time.Second)
// 		log.Fatal()
// 	}
// }

// log.Println("Все файлы присутствуют.")

// if err := godotenv.Load(".env"); err != nil {
// 	log.Println("Файл .env не найден")
// 	time.Sleep(10 * time.Second)
// 	log.Fatal()
// }

// d := gomail.NewDialer("smtp.mail.ru", 465, os.Getenv("EMAIL"), os.Getenv("EMAIL_PASS"))
// if err := d.DialAndSend(); err != nil {
// 	log.Println("Не удалось отправить установить соединение с почтовым ящиком. Убедитесь ,что E-mail и пароль в файле .env указаны верно")
// 	time.Sleep(10 * time.Second)
// 	log.Fatal()
// } else {
// 	log.Println("Соединение с почтовым ящиком установлено.")
// }
// }

// func readJson(jsonName string) []Data {
// 	file, err := os.Open(jsonName)
// 	if err != nil {
// 		fmt.Println("Не найден файл БД users.json")
// 		time.Sleep(10 * time.Second)
// 		log.Fatal()
// 	}
// 	defer file.Close()

// 	var records []Data
// 	decoder := json.NewDecoder(file)
// 	err = decoder.Decode(&records)
// 	if err != nil {
// 		if err.Error() == "EOF" && len(records) == 0 {
// 			fmt.Println("Файл users.json пуст")
// 			time.Sleep(10 * time.Second)
// 			log.Fatal()
// 		} else {
// 			fmt.Println("Убедитесь, что вы используете верную базу данных", err)
// 			time.Sleep(10 * time.Second)
// 			log.Fatal()
// 		}
// 	}
// 	return records
// }

// func checkAndLog(item Data) {
// 	log_name := time.Now().Format("01.02.2006")

// 	existingLogs, err := os.Open("./logs/" + log_name + ".json")
// 	if os.IsNotExist(err) { //Если файла log не существует создаю и записываю в него item
// 		create_log(item)
// 		SendEmailReg(item)
// 		return
// 	}
// 	defer existingLogs.Close()

// 	var logs []Data
// 	decoder := json.NewDecoder(existingLogs)
// 	err = decoder.Decode(&logs)
// 	if err != nil {
// 		time.Sleep(10 * time.Second)
// 		log.Fatal()
// 		// fmt.Println("Ошибка при декодировании JSON:", err)
// 		return
// 	}

// itemAlreadyExists := false
// for _, log := range logs {
// 	if log == item {
// 		itemAlreadyExists = true
// 		break
// 	}
// }
// if itemAlreadyExists {
// 	log.Println("Сегодня все поздравлены.")
// 	time.Sleep(10 * time.Second)
// 	log.Fatal()
// }
// if !itemAlreadyExists { //Если нет в логах - поздравить и записать в логи
// 	SendEmailReg(item)
// 	logs = append(logs, item)
// 	logJson, _ := json.Marshal(logs)
// 	overwriteLogs, err := os.Create("./logs/" + log_name + ".json")
// 	if err != nil {
// 		// fmt.Println("Unable to create file:", err)
// 		log.Fatal()
// 	}
// 	_, err = overwriteLogs.Write(logJson)
// 	if err != nil {
// 		// log.Println("=Ошибка записи в json=", err)
// 		log.Fatal()
// 	}
// 	defer overwriteLogs.Close()
// }

// }
// func create_log(item Data) {
// 	log_name := time.Now().Format("01.02.2006")

// 	newLog := []Data{item}
// 	logJson, err := json.Marshal(newLog)
// 	if err != nil {
// 		log.Fatal()
// 		// log.Fatal("=Ошибка форматирования лога в json=", err)
// 	}

// 	if err := os.MkdirAll("./logs", os.ModePerm); err != nil {
// 		// log.Fatal("Ошибка при создании директории logs:", err)
// 		log.Fatal()
// 	}
// 	newLogs, err := os.Create("./logs/" + log_name + ".json")
// 	if err != nil {
// 		// fmt.Println("Unable to create file:", err)
// 		log.Fatal()
// 	}
// 	_, err = newLogs.Write(logJson)
// 	if err != nil {
// 		// log.Println("=Ошибка записи в json=", err)
// 		log.Fatal()
// 	}
// 	defer newLogs.Close()
// }

// func SendEmail(item Data) {
// 	subject := "C днем рождения!"
// 	first_name := item.First_name
// 	last_name := item.Last_name

// 	replacer := strings.NewReplacer("${first_name}", first_name, "${last_name}", last_name)

// 	htmlBytes, err := os.ReadFile("index.html")
// 	if err != nil {
// 		// fmt.Println("Ошибка при чтении файла index.html:", err)
// 		log.Fatal()
// 		return
// 	}
// 	html := string(htmlBytes)
// 	html = replacer.Replace(html)

// 	m := gomail.NewMessage()
// 	m.SetHeader("From", os.Getenv("EMAIL"))
// 	m.SetHeader("To", item.Email)
// 	m.SetHeader("Subject", subject)
// 	m.SetBody("text/html", html)

// 	d := gomail.NewDialer("smtp.mail.ru", 465, os.Getenv("EMAIL"), os.Getenv("EMAIL_PASS"))
// 	if err := d.DialAndSend(m); err != nil {

// 		time.Sleep(10 * time.Second)
// 		log.Fatal()
// 	}
// 	fmt.Printf("Поздравление отправлено:%s", item.Email)
// 	time.Sleep(10 * time.Second)

// }

// func SendEmail(email string) {

// 	subject := "C днем рождения!"
// 	// first_name := item.First_name
// 	// last_name := item.Last_name

// 	// replacer := strings.NewReplacer("${first_name}", first_name, "${last_name}", last_name)

// 	htmlBytes, err := os.ReadFile("index.html")
// 	if err != nil {
// 		// fmt.Println("Ошибка при чтении файла index.html:", err)
// 		log.Fatal()
// 		return
// 	}
// 	html := string(htmlBytes)
// 	// html = replacer.Replace(html)
// 	log.Println("=b771ac=", reflect.TypeOf(email))
// 	m := gomail.NewMessage()
// 	m.SetHeader("From", os.Getenv("EMAIL"))
// 	m.SetHeader("To", email)
// 	m.SetHeader("Subject", subject)
// 	m.SetBody("text/html", html)

// 	d := gomail.NewDialer("smtp.mail.ru", 465, os.Getenv("EMAIL"), os.Getenv("EMAIL_PASS"))
// 	if err := d.DialAndSend(m); err != nil {
// 		log.Println("=8b66a8=", err)
// 		time.Sleep(10 * time.Second)
// 		log.Fatal()
// 	}
// 	fmt.Printf("Поздравление отправлено:%s", email)
// 	time.Sleep(10 * time.Second)

// }

// log.Println("=fba203=", reflect.TypeOf(today))
// date := "11/06/1969"
// t, err := time.Parse("02/01/2006", date)
// fmt.Println("=008c37=", "День", t.Day(), "Месяц", t.Month())
// fmt.Println("=008c37=", "День", today.Day(), "Месяц", today.Month())

// reflect.TypeOf(tst)
// \n
// today := time.Now()
// today.Day() == birthdate.Day() && today.Month()
