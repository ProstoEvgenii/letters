package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// type apiRequest struct {
// 	Email string `json:"email"`
// }

type Response struct {
	Count int64 `json:"count"`
}
type UsersUpload struct {
	Last_name   string `json:"Фамилия" bson:"Фамилия"`
	First_name  string `json:"Имя" bson:"Имя"`
	Middle_name string `json:"Отчество" bson:"Отчество"`
	Date_birth  string `json:"Дата рождения" bson:"Дата рождения"`
	Email       string `json:"E-mail" bson:"E-mail"`
}

func Start() {

	http.HandleFunc("/", anyPage)
	http.HandleFunc("/api", ParseRequest)
	// http.HandleFunc("/api/updateBD", updateBD)
	http.HandleFunc("/api/upload", uploadHandler)
	http.ListenAndServe(":80", nil)

}

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
		// http.Error(w, "Не удалось прочитать файл", http.StatusInternalServerError)
		return
	}
	var users []UsersUpload //  Форматитирую срез байтов в структуру
	if err := json.Unmarshal(fileBytes, &users); err != nil {
		fmt.Println("Ошибка при разборе JSON:", err)
		return
	}
	var result int64
	for _, item := range users {
		result += InsertIfNotExists(item)
	}

	response := Response{
		Count: result,
	}

	usersAdded, err := json.Marshal(response)
	if err != nil {
		fmt.Println("error:", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(usersAdded)
}
func ParseRequest(rw http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Max-Age", "15")

		itemCount := Dashboard()
		response := Response{
			Count: itemCount,
		}

		itemCountJson, err := json.Marshal(response)
		if err != nil {
			fmt.Println("error:", err)
		}
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
