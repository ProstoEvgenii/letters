package controller

import (
	"github.com/fasthttp/router"
	"mailsender/config"
	"mailsender/internal/controller/handlers"
	"mailsender/internal/events"
)

func NewRouter(cfg *config.Config, u events.UseCase) *router.Router {
	r := router.New()

	// events handler
	eventsHandler := handlers.NewEventsHandler(u)
	r.GET("/events", eventsHandler.GetEvent)
	//r.OPTIONS("/events", output.CORSOptions)

	return r
}
