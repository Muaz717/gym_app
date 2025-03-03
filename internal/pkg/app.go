package app

import (
	"context"
	membershipService "gym_app/internal/app/services/membership"
	personService "gym_app/internal/app/services/person"
	"gym_app/internal/app/storage/postgres"
	"gym_app/internal/config"
	"gym_app/internal/lib/logger/sl"
	httpApp "gym_app/internal/pkg/http"
	"log/slog"
)

type App struct {
	HTTPSrv *httpApp.HttpApp
}

func New(
	ctx context.Context,
	log *slog.Logger,
	cfg *config.Config,
) *App {
	storage, err := postgres.New(ctx, cfg.DB)
	if err != nil {
		log.Error("failed to init storage", sl.Error(err))
		panic(err)
	}

	personSrv := personService.New(log, storage)
	membershipSrv := membershipService.New(log, storage)

	httpApplication := httpApp.New(ctx, log, *cfg, personSrv, membershipSrv)

	return &App{
		HTTPSrv: httpApplication,
	}
}
