package membershipService

import (
	"context"
	"errors"
	"fmt"
	"gym_app/internal/lib/logger/sl"
	"gym_app/internal/models"
	"log/slog"
)

type SubscriptionService struct {
	log                 *slog.Logger
	subscriptionStorage SubscriptionStorage
}

type SubscriptionStorage interface {
	SaveSubscription(ctx context.Context, subscription models.Subscription) (int, error)
	FindAllSubscription(ctx context.Context) ([]models.Subscription, error)
}

var (
	ErrSubExists = errors.New("subscription with that number already exists")
)

func New(
	log *slog.Logger,
	subscriptionStorage SubscriptionStorage,
) *SubscriptionService {
	return &SubscriptionService{
		log:                 log,
		subscriptionStorage: subscriptionStorage,
	}
}

func (m *SubscriptionService) AddSubscription(ctx context.Context, subscription models.Subscription) (int, error) {
	const op = "services.subscription.AddSubscription"

	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("Adding new membership")

	subId, err := m.subscriptionStorage.SaveSubscription(ctx, subscription)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, sl.Error(err))
	}

	log.Info("subscription registered", "mid", subId)

	return subId, nil
}

func (m *SubscriptionService) FindAllSubscriptions(ctx context.Context) ([]models.Subscription, error) {
	const op = "services.subscription.FindAllSubscriptions"

	log := m.log.With(
		slog.String("op", op),
	)

	subscriptions, err := m.subscriptionStorage.FindAllSubscription(ctx)
	if err != nil {
		log.Warn("error", sl.Error(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Subscriptions are found")

	return subscriptions, nil
}
