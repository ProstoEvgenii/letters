package db

import (
	"context"
	"fmt"
	"mailsender/internal/entity"
)

// Create events
func (e *eventsRepo) Create(ctx context.Context, event *entity.Event) error {
	query := fmt.Sprintf(insert, event.UserID)

	_, err := e.db.Conn.Exec(ctx, query,
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
