package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type apiRequest struct {
	Email string `json:"email"`
}
type updateBDRequest struct {
	Email string `json:"email"`
}
type Response struct {
	Count int64 `json:"count"`
}

func Start() {

	http.HandleFunc("/", anyPage)
	http.HandleFunc("/api", ParseRequest)
	http.HandleFunc("/api/updateBD", updateBD)
	http.ListenAndServe(":80", nil)

}

func updateBD(rw http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Max-Age", "15")

		decoder := json.NewDecoder(request.Body)

		var data []Data
		err := decoder.Decode(&data)
		if err != nil {
			log.Println("=b570dd=", err)
		}
		for _, item := range data {
			InsertIfNotExists(item)
		}
		return
	}
}

func ParseRequest(rw http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Max-Age", "15")
		itemCount := CountDocuments()
		response := Response{
			Count: itemCount,
		}

		itemCountJson, err := json.Marshal(response)
		if err != nil {
			fmt.Println("error:", err)
		}
		rw.Write(itemCountJson)
		// fmt.Fprintf(rw, fmt.Sprintf("%v", itemCountJson))
		return
	} else if request.Method == "POST" {
		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Max-Age", "15")

		decoder := json.NewDecoder(request.Body)

		var data apiRequest
		err := decoder.Decode(&data)
		if err != nil {
			panic(err)
			// log.Fatal("Aborting", err)
		}
		log.Println("=90674d=", data.Email)
		SendEmail(data.Email)
		return
	}
}

func anyPage(rw http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(rw, "Hello")
}
