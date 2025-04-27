package httpApp

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gym_app/internal/config"
	membershipHandler "gym_app/internal/http/membership"
	personHandler "gym_app/internal/http/person"
	"gym_app/internal/lib/logger/sl"
	"log/slog"
	"net/http"
)

type HttpApp struct {
	HTTPServer *http.Server
	router     chi.Router
	ctx        context.Context
	log        *slog.Logger
	cfg        config.Config
}

func New(
	ctx context.Context,
	log *slog.Logger,
	cfg config.Config,
	personService personHandler.PersonService,
	membershipService membershipHandler.MembershipService,
) *HttpApp {

	personHandle := personHandler.New(ctx, log, personService)
	membershipHandle := membershipHandler.New(ctx, log, membershipService)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	//router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/people", func(people chi.Router) {
		people.Post("/add", personHandle.AddPerson)
		people.Get("/{name}", personHandle.FindMemsByPersonName)
		people.Get("/", personHandle.FindAllPeople)
	})

	router.Route("/membership", func(membership chi.Router) {
		membership.Post("/add", membershipHandle.AddMembership)
		membership.Get("/", membershipHandle.FindAllMemberships)
	})

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &HttpApp{
		HTTPServer: srv,
		router:     router,
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

	if err := a.HTTPServer.ListenAndServe(); err != nil {
		log.Error("failed to run http server", sl.Error(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("HTTP server is running", slog.String("addr", a.HTTPServer.Addr))

	return nil
}

func (a *HttpApp) Stop() error {
	const op = "httpApp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping Http server", slog.String("addr", a.HTTPServer.Addr))

	if err := a.HTTPServer.Shutdown(a.ctx); err != nil {
		a.log.Error("failed to stop server", sl.Error(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
