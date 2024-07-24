package app

import (
	"context"
	"github.com/rs/zerolog/log"
	"mailsender/config"
	"mailsender/internal/controller"
	db "mailsender/internal/events/repository/pg"
	"mailsender/internal/events/usecase"
	"mailsender/pkg/cache/rediscache"
	"mailsender/pkg/db/postgres"
)

func Run(cfg *config.Config) {
	ctx := context.Background()

	log.Log().Msg("Connection to psql")
	psqlDB, err := postgres.NewPSQL(ctx, cfg)
	if err != nil {
		log.Err(err).Send()
	}
	defer psqlDB.Conn.Close(ctx)

	log.Log().Msg("Connection to redis")
	rdb, err := rediscache.NewRedisClient(cfg)
	if err != nil {
		log.Err(err).Send()
	}
	_ = rdb

	eventsRepo := db.NewEventsRepository(psqlDB)

	eventsUC := usecase.NewEventsUseCase(cfg, eventsRepo)

	router := controller.NewRouter(cfg, eventsUC)

	_ = router

}
