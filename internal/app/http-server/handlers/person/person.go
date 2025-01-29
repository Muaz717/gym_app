package personHandler

import (
	"context"
	"errors"
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
	addPerson(ctx context.Context, person models.Person) (int, error)
	findAllPeople(ctx context.Context) ([]models.Person, error)
	findPersonByName(ctx context.Context, name string) (models.Person, error)
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

func (h *PersonHandler) addPerson(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.person.addPerson"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req models.Person

	err := render.DecodeJSON(r.Body, &req)
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
}
