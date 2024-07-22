package pages

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"letters/functions"
	"letters/models"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func AuthHandler(rw http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		uploadAuth(rw, request)
		return
	}

	return
}

func uploadAuth(rw http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var params models.Auth
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Error parse post News => ", err)
		fmt.Fprintf(rw, "{\"error\":\" Не верные данные\"}")
		return
	}
	data := []byte(params.Password)
	hash := sha256.Sum256(data)
	hashString := hex.EncodeToString(hash[:])
	var tmp interface{}
	filter := bson.M{
		"login":    params.Login,
		"password": hashString,
	}

	check := functions.CheckInDB(tmp, filter, "auth")
	if !check {
		fmt.Fprintf(rw, "{\"error\":\" Не верные данные\"}")
		return
	} else {
		now := time.Now()
		timestamp := now.Unix()
		functions.AuthUsers[params.UUID] = timestamp
		go functions.CheckAuthUsers()
		fmt.Fprintf(rw, "{\"result\":\"Авторизация успешна\"}")
		return
	}

}
