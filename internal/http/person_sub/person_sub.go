package person_sub

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

type PersonSubService interface {
	AddPersonSub(ctx context.Context, personSub models.PersonSubscription) (string, error)
	GetPersonSubByNumber(ctx context.Context, number string) (models.PersonSubscription, error)
	GetAllPersonSubs(ctx context.Context) ([]models.PersonSubscription, error)
}

type PersonSubHandler struct {
	ctx              context.Context
	log              *slog.Logger
	personSubService PersonSubService
}

func New(ctx context.Context, log *slog.Logger, personSubService PersonSubService) *PersonSubHandler {
	return &PersonSubHandler{
		ctx:              ctx,
		log:              log,
		personSubService: personSubService,
	}
}

func (h *PersonSubHandler) AddPersonSub(c *gin.Context) {

	const op = "handlers.personSub.addPersonSub"

	log := h.log.With(
		slog.String("op", op),
	)

	var personSub models.PersonSubscription

	if err := c.ShouldBindJSON(&personSub); err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			c.JSON(http.StatusBadRequest, resp.Error("empty request"))
			return
		}

		log.Error("failed to decode request body", sl.Error(err))
		c.JSON(http.StatusBadRequest, resp.Error("failed to decode request"))
		return
	}

}
