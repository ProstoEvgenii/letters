package server

import (
	"fmt"
	"letters/db"
	"letters/models"
	"letters/pages"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson"
)

func Start(host string) {
	http.HandleFunc("/", HandleRequest)
	http.HandleFunc("/email/unsubcribe", Unsubcribe)
	http.ListenAndServe(host, nil)
}

var router = map[string]func(http.ResponseWriter, *http.Request){
	"Dashboard": pages.DashboardHandler,
	"Settings":  pages.SettingsHandler,
	"Database":  pages.DatabaseHandler,
	"History":   pages.HistoryHandler,
	"UserAuth":  pages.AuthHandler,
	"Events":    pages.UploadEventsHandler,
	"Templates": pages.UploadTemplateHandler,
}

func HandleRequest(rw http.ResponseWriter, request *http.Request) {

	path := strings.Split(request.URL.Path, "/api/")

	handler, exists := router[path[1]]

	if exists {
		handler(rw, request)
	} else {
		log.Println("Не найден handler => ", path[1])
		// Обработка случая, когда маршрут не найден
		http.NotFound(rw, request)
	}
}

func anyPage(rw http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(rw, "Hello")
}

func Unsubcribe(rw http.ResponseWriter, request *http.Request) {
	params := new(models.Unsubscribe)
	if err := schema.NewDecoder().Decode(params, request.URL.Query()); err != nil {
		log.Println("=Params schema Error News_=", err)
	}

	if params.Email != "" {
		filter := bson.M{
			"E-mail": params.Email,
		}
		update := bson.M{"$set": bson.M{
			"unsubscribe": true,
		}}
		db.UpdateIfExists(filter, update, "users")
	}

	http.Redirect(rw, request, "https://xn----dtbsbdgikgdbazpac.xn--p1ai/", http.StatusSeeOther)

}
