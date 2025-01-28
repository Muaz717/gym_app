package postgres

import (
	"context"
	"fmt"
	"gym_app/internal/models"
)

func (s *Storage) savePerson(
	ctx context.Context,
	person models.Person,
) (int, error) {
	const op = "postgres.savePerson"

	query := `INSERT INTO person(name) VALUES($1) RETURNING id`

	row := s.db.QueryRow(ctx, query, person.Name)

	var personId int
	if err := row.Scan(&personId); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return personId, nil
}
