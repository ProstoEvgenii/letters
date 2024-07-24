package handlers

import (
	"github.com/valyala/fasthttp"
	"mailsender/internal/events"
)

type eventsHandler struct {
	uc events.UseCase
}

func NewEventsHandler(uc events.UseCase) *eventsHandler {
	return &eventsHandler{uc}
}

func (h *eventsHandler) GetEvent(ctx *fasthttp.RequestCtx) {
	eventId := int64(ctx.QueryArgs().GetUintOrZero("event_id"))
	_ = eventId
}
