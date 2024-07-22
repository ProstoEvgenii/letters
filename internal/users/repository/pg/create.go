package db

import (
	"context"
	"mailsender/internal/entity"
)

// Create - creates new user
func (e *usersRepo) Create(ctx context.Context, event *entity.Event) error {

	_, err := e.db.Conn.Exec(ctx, insert,
		event.ID,
		event.UserID,
		event.Active,
		event.Title,
		event.Daily,
		event.IsSent,
		event.Subject,
		event.Author,
		event.Letter.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
