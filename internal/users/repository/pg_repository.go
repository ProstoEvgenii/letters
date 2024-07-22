package repository

import (
	"context"
	"mailsender/internal/entity"
)

type Users interface {
	Create(ctx context.Context, event *entity.Event) error
	Update()
	Delete()
	GetByID()
}
