package pages

import (
	"encoding/json"
	"fmt"
	"io"
	"letters/db"
	"letters/functions"
	"letters/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func UploadTemplateHandler(rw http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		response := UploadTemplate(rw, request)
		eventAdded, err := json.Marshal(response)
		if err != nil {
			fmt.Println("error:", err)
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write(eventAdded)
		return
	}

	return
}
func UploadTemplate(rw http.ResponseWriter, request *http.Request) models.DashboardPostResponse {
	var response models.DashboardPostResponse

	UUID := request.FormValue("UUID")
	_, exists := functions.AuthUsers[UUID]
	if !exists {
		response.Err = "Ошибка Авторизации."
		return response
	}
	file, _, err := request.FormFile("jsonFileTemplate")
	if err != nil {
		response.Err = "Не удалось получить файл."
		return response
	}
	defer file.Close()
	fileContents, err := io.ReadAll(file)
	if err != nil {
		response.Err = "Не удалось получить файл."
		return response
	}
	name := request.FormValue("name")

	indexHTML := string(fileContents)
	if name == "" || indexHTML == "" {
		response.Err = "Полученные данные некорректны."
		return response
	}

	filter := bson.M{
		"name": name,
	}
	update := bson.M{"$set": bson.M{
		"name":      name,
		"indexHTML": indexHTML,
	}}

	templateInserted := db.InsertIfNotExists(filter, update, "templates")
	result := "ok"
	response = models.DashboardPostResponse{
		Err:               result,
		DocumentsInserted: templateInserted.UpsertedCount,
		DocumentsModified: templateInserted.ModifiedCount,
	}
	return response
}
