package db

import (
	"mailsender/internal/users/repository"
	"mailsender/pkg/db/postgres"
)

// usersRepo - Events Repository
type usersRepo struct {
	db *postgres.Postgres
}

// NewUsersRepository - Events repository constructor
func NewUsersRepository(db *postgres.Postgres) repository.Users {
	return &usersRepo{db: db}
}
