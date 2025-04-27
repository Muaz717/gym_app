package httpApp

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gym_app/internal/config"
	personHandler "gym_app/internal/http/person"
	membershipHandler "gym_app/internal/http/subscription"
	"gym_app/internal/lib/logger/sl"
	"log/slog"
	"net/http"
	"time"
)

type HttpApp struct {
	HTTPServer *http.Server
	engine     *gin.Engine
	ctx        context.Context
	log        *slog.Logger
	cfg        config.Config
}

func New(
	ctx context.Context,
	log *slog.Logger,
	cfg config.Config,
	personService personHandler.PersonService,
	membershipService membershipHandler.SubscriptionService,
) *HttpApp {

	personHandle := personHandler.New(ctx, log, personService)
	membershipHandle := membershipHandler.New(ctx, log, membershipService)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	engine.Use(gin.Recovery())
	//engine.Use(RequestIDMiddleware(log))
	// engine.Use(CustomLoggerMiddleware(log)) // можно подключить свое логирование запросов

	people := engine.Group("/people")
	{
		people.POST("/add", personHandle.AddPerson)
		people.GET("/:name", personHandle.FindSubsByPersonName)
		people.GET("/", personHandle.FindAllPeople)
	}

	membership := engine.Group("/membership")
	{
		membership.POST("/add", membershipHandle.AddSubscription)
		membership.GET("/", membershipHandle.FindAllSubscriptions)
	}

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      engine,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &HttpApp{
		HTTPServer: srv,
		engine:     engine,
		ctx:        ctx,
		log:        log,
		cfg:        cfg,
	}
}

func (a *HttpApp) Run() error {
	const op = "httpApp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.String("addr", a.cfg.Address),
	)

	log.Info("HTTP server is starting", slog.String("addr", a.cfg.Address))

	if err := a.HTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("failed to run http server", sl.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *HttpApp) Stop() error {
	const op = "httpApp.Stop"

	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()

	a.log.With(slog.String("op", op)).
		Info("stopping HTTP server", slog.String("addr", a.HTTPServer.Addr))

	if err := a.HTTPServer.Shutdown(ctx); err != nil {
		a.log.Error("failed to gracefully shutdown server", sl.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
