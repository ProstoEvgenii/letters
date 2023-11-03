package pages

import (
	"context"
	"encoding/json"
	"fmt"
	"letters/db"
	"letters/models"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func DatabaseHandler(rw http.ResponseWriter, request *http.Request) {

	if request.Method == "POST" {

		rw.Write([]byte("Привет"))
	}

	if request.Method == "GET" {
		usersCount := db.CountDocuments(bson.M{}, "users")
		cursor := db.Find(bson.M{}, "users")
		var usersSlice []models.Users
		if err := cursor.All(context.TODO(), &usersSlice); err != nil {
			log.Println("Cursor All Error Database", err)
			rw.Write([]byte("{}"))
			return
		}
		if len(usersSlice) == 0 {
			rw.Write([]byte("{}"))
			return
		}

		response := models.GetDataBaseResponse{
			Records:    usersSlice,
			UsersCount: usersCount,
		}
		dataBaseJson, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error:", err)
			rw.Write([]byte("{}"))
			return
		}
		rw.Write(dataBaseJson)
		return
	}
}
