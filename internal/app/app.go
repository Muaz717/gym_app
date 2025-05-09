package app

import (
	"context"
	"gym_app/internal/app/http"
	"gym_app/internal/clients/sso/grpc"
	"gym_app/internal/config"
	"gym_app/internal/cron"
	"gym_app/internal/lib/logger/sl"
	authService "gym_app/internal/services/auth"
	"gym_app/internal/services/person"
	personSubService "gym_app/internal/services/person_sub"
	"gym_app/internal/services/subscription"
	"gym_app/internal/storage/postgres"
	"log/slog"
)

type App struct {
	HTTPSrv *httpApp.HttpApp
	Cron    *cron.CronJobs
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

	ssoClient, err := grpc.NewSSOClient(
		log,
		cfg.Clients.SSO.Address,
		cfg.Clients.SSO.Timeout,
		cfg.Clients.SSO.RetriesCount,
	)
	if err != nil {
		log.Error("failed to init sso client", sl.Error(err))
		panic(err)
	}

	personSrv := personService.New(log, storage)
	subscriptionSrv := subscriptionService.New(log, storage)
	personSubSrv := personSubService.New(log, storage)
	authSrv := authService.New(log, ssoClient, cfg.AppID)

	cr := cron.New(personSubSrv)

	httpApplication := httpApp.New(ctx, log, *cfg, ssoClient, authSrv, personSrv, subscriptionSrv, personSubSrv)

	return &App{
		HTTPSrv: httpApplication,
		Cron:    cr,
	}
}
