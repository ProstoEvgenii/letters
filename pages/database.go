package pages

import (
	"context"
	"encoding/json"
	"fmt"
	"letters/db"
	"letters/models"
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson"
)

func DatabaseHandler(rw http.ResponseWriter, request *http.Request) {

	if request.Method == "POST" {

		rw.Write([]byte("Привет"))
	}

	if request.Method == "GET" {

		usersCount := db.CountDocuments(bson.M{}, "users")

		params := new(models.Dashboard_Params)
		if err := schema.NewDecoder().Decode(params, request.URL.Query()); err != nil {
			log.Println("=Params schema Error News_=", err)
		}
		page := 1
		if params.Page != 0 {
			page = params.Page
		}
		limit := 15
		skip := limit * (page - 1)

		cursor := db.FindSkip(bson.M{}, "users", skip, limit)
		var usersSlice []models.Users
		if err := cursor.All(context.TODO(), &usersSlice); err != nil {
			log.Println("Cursor All Error Database", err)
			rw.Write([]byte("{}"))
			return
		}
		// if len(usersSlice) == 0 {
		// 	// rw.Write([]byte("{}"))
		// 	return
		// }

		response := models.GetDataBaseResponse{
			Records:    usersSlice,
			UsersCount: usersCount,
			PageNumber: page,
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
