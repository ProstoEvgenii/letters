package pages

import (
	"context"
	"encoding/json"
	"fmt"
	"letters/db"
	"letters/functions"
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
		filter := bson.M{}
		limitPerPage := 15
		page := 1
		skip := limitPerPage * (page - 1)

		if err := schema.NewDecoder().Decode(params, request.URL.Query()); err != nil {
			log.Println("=Params schema Error Database=", err)
		}

		if params.UUID != "" {
			_, exists := functions.AuthUsers[params.UUID]
			if !exists {
				return
			}
		}

		if params.Page != 0 {
			page = params.Page
		}

		if params.Seach != "" {
			filter = bson.M{
				"$or": []bson.M{
					{"Имя": bson.M{"$regex": params.Seach, "$options": "i"}},
					{"Фамилия": bson.M{"$regex": params.Seach, "$options": "i"}},
					{"E-mail": bson.M{"$regex": params.Seach, "$options": "i"}},
				},
			}
		}

		cursor := db.FindSkip(filter, "users", skip, limitPerPage)
		var usersSlice []models.Users
		if err := cursor.All(context.TODO(), &usersSlice); err != nil {
			log.Println("Cursor All Error Database", err)
			rw.Write([]byte("{}"))
			return
		}
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
