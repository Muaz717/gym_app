package membershipHandler

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"gym_app/internal/lib/api/response"
	"gym_app/internal/lib/logger/sl"
	"gym_app/internal/models"
	"io"
	"log/slog"
	"net/http"
)

type SubscriptionService interface {
	AddSubscription(ctx context.Context, subscription models.Subscription) (int, error)
	FindAllSubscription(ctx context.Context) ([]models.Subscription, error)
}

type SubscriptionHandler struct {
	ctx                 context.Context
	log                 *slog.Logger
	subscriptionService SubscriptionService
}

func New(
	ctx context.Context,
	log *slog.Logger,
	membershipService SubscriptionService,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		ctx:                 ctx,
		log:                 log,
		subscriptionService: membershipService,
	}
}

func (h *SubscriptionHandler) AddSubscription(c *gin.Context) {
	const op = "handlers.subscription.addSubscription"

	log := h.log.With(
		slog.String("op", op),
	)

	var subscription models.Subscription

	if err := c.ShouldBindJSON(&subscription); err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			c.JSON(http.StatusBadRequest, response.Error("empty request"))
			return
		}

		log.Error("failed to decode request body", sl.Error(err))

		c.JSON(http.StatusBadRequest, response.Error("failed to decode request"))
		return
	}

	subId, err := h.subscriptionService.AddSubscription(h.ctx, subscription)
	if err != nil {
		log.Error("failed to add subscription", sl.Error(err))

		c.JSON(http.StatusInternalServerError, response.Error("failed to add subscription"))
		return
	}

	log.Info("Subscription added", slog.Int("Subscription_id", subId))
	c.JSON(http.StatusOK, response.OK("Subscription added"))
}

func (h *SubscriptionHandler) FindAllSubscriptions(c *gin.Context) {
	const op = "handlers.subscription.findAllSubscriptions"

	log := h.log.With(
		slog.String("op", op),
	)

	subscriptions, err := h.subscriptionService.FindAllSubscription(h.ctx)
	if err != nil {
		log.Error("failed to get Subscriptions", sl.Error(err))

		c.JSON(http.StatusInternalServerError, response.Error("failed to get Subscriptions"))
		return
	}

	log.Info("Subscription found")

	c.JSON(http.StatusOK, subscriptions)
}
