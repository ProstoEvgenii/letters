package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson"
)

// type apiRequest struct {
// 	Email string `json:"email"`
// }

func Start() {
	http.HandleFunc("/", HandleRequest)
	http.HandleFunc("/api/Dashboard/upload", uploadHandler)
	// http.HandleFunc("/api/Dashboard/upload", uploadHandler)
	http.ListenAndServe(":80", nil)
}

var router = map[string]func(http.ResponseWriter, *http.Request){
	"Dashboard": DashboardHandler,
	// "/api/Dashboard/upload": UploadHandler,
	// "/api/Users": UsersHandler,
}

func HandleRequest(rw http.ResponseWriter, request *http.Request) {
	//разбиваем полученный запрос на массив строк по разделителю /
	path := strings.Split(request.URL.Path, "/api/")
	//берем первую по индексу строку и проверяем существует ли такой маршрут в router

	// Ищем обработчик в карте по URL-пути
	handler, exists := router[path[1]]
	if exists {
		// Вызываем соответствующую фнкцию из роутера
		handler(rw, request)
	} else {
		log.Println("Не найден event => ", path[1])
		// Обработка случая, когда маршрут не найден
		http.NotFound(rw, request)
	}
}

// func main() {
// Установите сервер и маршрутизацию на HandleRequest
// 	http.HandleFunc("/", HandleRequest)
// 	http.ListenAndServe(":3000", nil)
// }

// func updateBD(rw http.ResponseWriter, request *http.Request) {
// 	if request.Method == "POST" {
// 		rw.Header().Set("Content-Type", "application/json")
// 		rw.Header().Set("Access-Control-Allow-Origin", "*")
// 		rw.Header().Set("Access-Control-Max-Age", "15")

// 		decoder := json.NewDecoder(request.Body)

// 		var data []Data
// 		err := decoder.Decode(&data)
// 		if err != nil {
// 			log.Println("=b570dd=", err)
// 		}
// 		for _, item := range data {
// 			InsertIfNotExists(item)
// 		}
// 		return
// 	}
// }

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
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

	response := Response{
		DocumentsInserted: documentsInserted,
	}

	usersAdded, err := json.Marshal(response)
	if err != nil {
		fmt.Println("error:", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(usersAdded)
}
func GetInfo(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Max-Age", "15")
	if request.Method == "GET" {
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
	} else if request.Method == "POST" {
		// r.FormFile("userFile")
		// rw.Header().Set("Content-Type", "application/json")
		// rw.Header().Set("Access-Control-Allow-Origin", "*")
		// rw.Header().Set("Access-Control-Max-Age", "15")

		// decoder := json.NewDecoder(request.Body)

		// var data apiRequest
		// err := decoder.Decode(&data)
		// if err != nil {
		// 	panic(err)
		// 	// log.Fatal("Aborting", err)
		// }
		// log.Println("=90674d=", data.Email)
		// SendEmail(data.Email)
		return
	}
}

func anyPage(rw http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(rw, "Hello")
}
