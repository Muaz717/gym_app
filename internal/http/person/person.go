package personHandler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
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
	FindMemsByPersonName(ctx context.Context, name string) ([]models.Subscription, error)
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

func (h *PersonHandler) AddPerson(c *gin.Context) {
	const op = "handlers.person.addPerson"

	log := h.log.With(
		slog.String("op", op),
	)

	var person models.Person

	if err := c.ShouldBindJSON(&person); err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			c.JSON(http.StatusBadRequest, resp.Error("empty request"))
			return
		}

		log.Error("failed to decode request body", sl.Error(err))
		c.JSON(http.StatusBadRequest, resp.Error("failed to decode request"))
		return
	}

	personId, err := h.personService.AddPerson(h.ctx, person)
	if err != nil {
		log.Error("failed to add person", sl.Error(err))

		c.JSON(http.StatusInternalServerError, resp.Error("failed to add person"))
		return
	}

	log.Info("Person added", slog.Int("person_id", personId))
	c.JSON(http.StatusOK, resp.OK("Person added"))
}

func (h *PersonHandler) FindAllPeople(c *gin.Context) {
	const op = "handlers.person.findAllPeople"

	log := h.log.With(
		slog.String("op", op),
	)

	people, err := h.personService.FindAllPeople(h.ctx)
	if err != nil {
		log.Error("failed to get people", sl.Error(err))

		c.JSON(http.StatusInternalServerError, resp.Error("failed to get people"))
		return
	}

	log.Info("People found")

	c.JSON(http.StatusOK, people)
}

func (h *PersonHandler) FindMemsByPersonName(c *gin.Context) {
	type Response struct {
		Name        string                `json:"name"`
		Memberships []models.Subscription `json:"memberships"`
	}

	const op = "handlers.person.findPersonByName"

	log := h.log.With(
		slog.String("op", op),
	)

	name := c.Param("name")
	if name == "" {
		log.Error("name parameter is missing")

		c.JSON(http.StatusBadRequest, resp.Error("name parameter is required"))
		return
	}

	memberships, err := h.personService.FindMemsByPersonName(h.ctx, name)
	if err != nil {
		log.Error("failed to get person memberships", sl.Error(err))

		c.JSON(http.StatusInternalServerError, resp.Error("failed to get person memberships"))
		return
	}

	personToShow := Response{
		Name:        name,
		Memberships: memberships,
	}

	log.Info("Person found")

	c.JSON(http.StatusOK, personToShow)
}
