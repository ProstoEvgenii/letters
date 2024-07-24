package usecase

import (
	"context"
	"mailsender/config"
	"mailsender/internal/entity"
	"mailsender/internal/events"
	"mailsender/internal/events/repository"
)

type eventsUC struct {
	eventsRepo repository.Events
}

// Events UseCase constructor
func NewEventsUseCase(cfg *config.Config, eventsRepo repository.Events) events.UseCase {
	return &eventsUC{eventsRepo: eventsRepo}
}

// Create event
func (r *eventsUC) Create(ctx context.Context, event entity.Event) error {
	err := r.eventsRepo.Create(ctx, &event)
	if err != nil {
		return err
	}
	return nil
}
