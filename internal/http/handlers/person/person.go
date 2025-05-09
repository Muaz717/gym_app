package personHandler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"gym_app/internal/lib/api/response"
	"gym_app/internal/lib/logger/sl"
	"gym_app/internal/models"
	personService "gym_app/internal/services/person"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

type PersonService interface {
	AddPerson(ctx context.Context, person models.Person) (int, error)
	FindAllPeople(ctx context.Context) ([]models.Person, error)
	UpdatePerson(ctx context.Context, person models.Person, pID int) (int, error)
	DeletePerson(ctx context.Context, pID int) error
	FindPersonByName(ctx context.Context, name string) (models.Person, error)
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

// AddPerson godoc
// @Summary Add a new person
// @Description Add a new person
// @Security BearerAuth
// @Tags person
// @Accept json
// @Produce json
// @Param person body models.Person true "Person"
// @Success 200 {object} response.Response "Person added"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 409 {object} response.Response "Conflict"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /people/add [post]
func (h *PersonHandler) AddPerson(c *gin.Context) {
	const op = "handlers.person.addPerson"

	log := h.log.With(
		slog.String("op", op),
	)

	var person models.Person

	if err := c.ShouldBindJSON(&person); err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			c.JSON(http.StatusBadRequest, response.Error("empty request"))
			return
		}

		log.Error("failed to decode request body", sl.Error(err))
		c.JSON(http.StatusBadRequest, response.Error("failed to decode request"))
		return
	}

	if err := person.Validate(); err != nil {
		log.Error("failed to validate person", err)

		c.JSON(http.StatusBadRequest, err)
		return
	}

	personId, err := h.personService.AddPerson(h.ctx, person)
	if err != nil {
		if errors.Is(err, personService.ErrPersonExists) {
			c.JSON(http.StatusConflict, response.Error("person already exists"))
			return
		}

		log.Error("failed to add person", sl.Error(err))
		c.JSON(http.StatusInternalServerError, response.Error("failed to add person"))
		return
	}

	log.Info("Person added", slog.Int("person_id", personId))
	c.JSON(http.StatusOK, response.OK("Person added, personId: "+strconv.Itoa(personId)))
}

// UpdatePerson godoc
// @Summary Update a person
// @Description Update a person
// @Security BearerAuth
// @Tags person
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Param person body models.Person true "Person"
// @Success 200 {object} response.Response "Person updated"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "Not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /people/update/{id} [put]
func (h *PersonHandler) UpdatePerson(c *gin.Context) {
	const op = "handlers.person.updatePerson"

	log := h.log.With(
		slog.String("op", op),
	)

	pIDStr := c.Param("id")
	if pIDStr == "" {
		log.Error("person id parameter is missing")

		c.JSON(http.StatusBadRequest, response.Error("person id parameter is required"))
		return
	}
	pID, err := strconv.Atoi(pIDStr)
	if err != nil {
		log.Error("failed to parse person id", sl.Error(err))

		c.JSON(http.StatusBadRequest, response.Error("invalid person id"))
		return
	}

	var person models.Person
	if err := c.ShouldBindJSON(&person); err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			c.JSON(http.StatusBadRequest, response.Error("empty request"))
			return
		}

		log.Error("failed to decode request body", sl.Error(err))
		c.JSON(http.StatusBadRequest, response.Error("failed to decode request"))
		return
	}

	if err := person.Validate(); err != nil {
		log.Error("failed to validate person", err)

		c.JSON(http.StatusBadRequest, err)
		return
	}

	personId, err := h.personService.UpdatePerson(h.ctx, person, pID)
	if err != nil {
		if errors.Is(err, personService.ErrPersonNotFound) {
			c.JSON(http.StatusNotFound, response.Error("person not found"))
			return
		}

		if errors.Is(err, personService.ErrPersonExists) {
			c.JSON(http.StatusConflict, response.Error("Person with such name and phone already exists. Set another name or phone"))
			return
		}

		log.Error("failed to update person", sl.Error(err))
		c.JSON(http.StatusInternalServerError, response.Error("failed to update person"))
		return
	}

	log.Info("Person updated", slog.Int("person_id", personId))
	c.JSON(http.StatusOK, response.OK("Person updated"))
}

// DeletePerson godoc
// @Summary Delete a person
// @Description Delete a person
// @Security BearerAuth
// @Tags person
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Success 200 {object} response.Response "Person deleted"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "Not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /people/delete/{id} [delete]
func (h *PersonHandler) DeletePerson(c *gin.Context) {
	const op = "handlers.person.deletePerson"

	log := h.log.With(
		slog.String("op", op),
	)

	pIDStr := c.Param("id")
	if pIDStr == "" {
		log.Error("person id parameter is missing")

		c.JSON(http.StatusBadRequest, response.Error("person id parameter is required"))
		return
	}
	pID, err := strconv.Atoi(pIDStr)
	if err != nil {
		log.Error("failed to parse person id", sl.Error(err))

		c.JSON(http.StatusBadRequest, response.Error("invalid person id"))
		return
	}

	err = h.personService.DeletePerson(h.ctx, pID)
	if err != nil {
		if errors.Is(err, personService.ErrPersonNotFound) {
			c.JSON(http.StatusNotFound, response.Error("person not found"))
			return
		}

		log.Error("failed to delete person", sl.Error(err))
		c.JSON(http.StatusInternalServerError, response.Error("failed to delete person"))
		return
	}

	log.Info("Person deleted", slog.Int("person_id", pID))
	c.JSON(http.StatusOK, response.OK("Person deleted"))
}

// FindPersonByName godoc
// @Summary Find a person by name
// @Description Find a person by name
// @Security BearerAuth
// @Tags person
// @Accept json
// @Produce json
// @Param name query string true "Person name"
// @Success 200 {object} models.Person "Person found"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 404 {object} response.Response "Not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /people/find [get]
func (h *PersonHandler) FindPersonByName(c *gin.Context) {
	const op = "handlers.person.findPersonByName"

	log := h.log.With(
		slog.String("op", op),
	)

	name := c.Query("name")
	if name == "" {
		log.Error("name parameter is missing")

		c.JSON(http.StatusBadRequest, response.Error("name parameter is required"))
		return
	}

	person, err := h.personService.FindPersonByName(h.ctx, name)
	if err != nil {
		if errors.Is(err, personService.ErrPersonNotFound) {
			c.JSON(http.StatusNotFound, response.Error("person not found"))
			return
		}

		log.Error("failed to find person", sl.Error(err))
		c.JSON(http.StatusInternalServerError, response.Error("failed to find person"))
		return
	}

	log.Info("Person found", slog.String("name", name))
	c.JSON(http.StatusOK, person)
}

// FindAllPeople godoc
// @Summary Find all people
// @Description Find all people
// @Security BearerAuth
// @Tags person
// @Accept json
// @Produce json
// @Success 200 {array} models.Person "People found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /people [get]
func (h *PersonHandler) FindAllPeople(c *gin.Context) {
	const op = "handlers.person.findAllPeople"

	log := h.log.With(
		slog.String("op", op),
	)

	people, err := h.personService.FindAllPeople(h.ctx)
	if err != nil {
		log.Error("failed to get people", sl.Error(err))

		c.JSON(http.StatusInternalServerError, response.Error("failed to get people"))
		return
	}

	log.Info("People found")

	c.JSON(http.StatusOK, people)
}
