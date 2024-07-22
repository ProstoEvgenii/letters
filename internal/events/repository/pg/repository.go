package db

import (
	"mailsender/internal/events/repository"
	"mailsender/pkg/db/postgres"
)

// eventsRepo - Events Repository
type eventsRepo struct {
	db *postgres.Postgres
}

// NewEventsRepository - Events repository constructor
func NewEventsRepository(db *postgres.Postgres) repository.Events {
	return &eventsRepo{db: db}
}
