package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gym_app/internal/models"
	"gym_app/internal/storage"
)

func (s *Storage) SaveMembership(
	ctx context.Context,
	membership models.Membership,
) (int, error) {
	const op = "postgres.saveMembership"

	query := `INSERT INTO membership(number, recording_day, person) VALUES($1, $2, $3) RETURNING id`

	row := s.db.QueryRow(ctx, query, membership.Number, membership.RecordingDay, membership.Person.Name)

	var memId int
	if err := row.Scan(&memId); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return memId, nil
}

func (s *Storage) FindAllMemberships(ctx context.Context) ([]models.Membership, error) {
	const op = "postgres.FindAllMemberships"

	query := `SELECT * FROM membership`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var mems []models.Membership
	for rows.Next() {
		mem := models.Membership{}
		err := rows.Scan(&mem.Id, &mem.Number, &mem.RecordingDay, &mem.Person.Name)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		mems = append(mems, mem)
	}

	return mems, nil
}
