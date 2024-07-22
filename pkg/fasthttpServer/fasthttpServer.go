package server

import (
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"mailsender/config"
)

func Start(cfg *config.Config, r fasthttp.RequestHandler) {
	server := &fasthttp.Server{
		Handler:            r,
		MaxRequestBodySize: 1024 * 1024,
	}

	addr := cfg.Server.ServerAddr
	if addr == "" {
		log.Fatal().Msg("failed to get 'SERVER_ADDR' from .env")
	}
	log.Info().Str("addr", addr).Msg("HTTP listening")
	log.Fatal().Err(server.ListenAndServe(addr)).Send()
}
