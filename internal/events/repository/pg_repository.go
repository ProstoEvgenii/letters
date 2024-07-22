package repository

import (
	"context"
	"mailsender/internal/entity"
)

type Events interface {
	Create(ctx context.Context, event *entity.Event) error
	Update()
	Delete()
	GetByID()
}
