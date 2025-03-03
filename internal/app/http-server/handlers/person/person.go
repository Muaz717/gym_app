package personHandler

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "gym_app/internal/lib/api/response"
	"gym_app/internal/lib/logger/sl"
	"gym_app/internal/models"
	"io"
	"log/slog"
	"net/http"
)

type PersonService interface {
	AddPerson(ctx context.Context, person models.Person) (int, error)
	FindAllPeople(ctx context.Context) ([]models.Person, error)
	FindMemsByPersonName(ctx context.Context, name string) ([]models.Membership, error)
}

type PersonHandler struct {
	ctx           context.Context
	log           *slog.Logger
	personService PersonService
}

func New(
	ctx context.Context,
	log *slog.Logger,
	personService PersonService,
) *PersonHandler {
	return &PersonHandler{
		ctx:           ctx,
		log:           log,
		personService: personService,
	}
}

func (h *PersonHandler) AddPerson(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.person.addPerson"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var person models.Person

	err := render.DecodeJSON(r.Body, &person)
	if errors.Is(err, io.EOF) {
		log.Error("request body is empty")

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, resp.Error("empty request"))

		return
	}
	if err != nil {
		log.Error("failed to decode request body", sl.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("failed to decode request"))

		return
	}

	personId, err := h.personService.AddPerson(h.ctx, person)
	if err != nil {
		log.Error("failed to addPerson", sl.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, resp.OK("failed to addPerson"))

		return
	}

	log.Info("Person added", slog.Int("person id", personId))
	render.JSON(w, r, resp.OK("Person added"))
}

func (h *PersonHandler) FindAllPeople(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.person.findAllPeople"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	people, err := h.personService.FindAllPeople(h.ctx)
	if err != nil {
		log.Error("failed to get people", sl.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("failed to get people"))

		return
	}

	log.Info("People found")

	render.JSON(w, r, people)
}

func (h *PersonHandler) FindMemsByPersonName(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Name        string
		Memberships []models.Membership
	}

	const op = "handlers.person.findPersonByName"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	name := chi.URLParam(r, "name")

	memberships, err := h.personService.FindMemsByPersonName(h.ctx, name)
	if err != nil {
		log.Error("failed to get person", sl.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("failed to get person"))

		return
	}

	personToShow := Response{
		Name:        name,
		Memberships: memberships,
	}

	log.Info("Person found")

	render.JSON(w, r, personToShow)
}
