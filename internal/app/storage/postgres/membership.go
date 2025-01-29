package postgres

import (
	"context"
	"fmt"
	"gym_app/internal/models"
)

func (s *Storage) saveMembership(
	ctx context.Context,
	membership models.Membership,
) (int, error) {
	const op = "postgres.saveMembership"

	query := `INSERT INTO membership(number, recording_day, person) VALUES($1, $2, $3) RETURNING id`

	row := s.db.QueryRow(ctx, query, membership.Number, membership.RecordingDay, membership.Person)

	var memId int
	if err := row.Scan(&memId); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return memId, nil
}
