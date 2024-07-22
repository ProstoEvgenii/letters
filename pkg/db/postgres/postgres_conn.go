package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"mailsender/config"
	"time"
)

type Postgres struct {
	Conn *pgx.Conn
}

func NewPSQL(ctx context.Context, conf *config.Config) (*Postgres, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, conf.PG.DSN)
	if err != nil {
		return nil, err
	}

	return &Postgres{Conn: conn}, nil
}
