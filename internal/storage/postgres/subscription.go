package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gym_app/internal/models"
	"gym_app/internal/storage"
)

func (s *Storage) AddSubscription(
	ctx context.Context,
	subscription models.Subscription,
) (int, error) {
	const op = "postgres.addSubscription"

	query := `INSERT INTO subscriptions(title, price, duration_days, freeze_days) VALUES($1, $2, $3, $4) RETURNING id`

	row := s.db.QueryRow(ctx, query, subscription.Title, subscription.Price, subscription.DurationDays, subscription.FreezeDays)

	var subId int
	if err := row.Scan(&subId); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return subId, nil
}

func (s *Storage) FindAllSubscriptions(ctx context.Context) ([]models.Subscription, error) {
	const op = "postgres.FindAllSubscriptions"

	query := `SELECT * FROM subscriptions`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var subs []models.Subscription
	for rows.Next() {
		sub := models.Subscription{}
		err := rows.Scan(&sub.ID, &sub.Title, &sub.Price, &sub.DurationDays, &sub.FreezeDays)

		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		subs = append(subs, sub)
	}

	return subs, nil
}
