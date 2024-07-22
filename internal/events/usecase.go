package events

import (
	"context"
	"mailsender/internal/entity"
)

// Events use case
type UseCase interface {
	Create(ctx context.Context, event entity.Event) error
	//GetEventByID()
	//GetEvents()
	//Update()
	//Delete()
}
