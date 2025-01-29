package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gym_app/internal/app/storage"
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

func (s *Storage) findAllPeople(ctx context.Context) ([]models.Person, error) {
	const op = "postgres.findAllPeople"

	query := `SELECT id, name FROM person`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.Person])
}

func (s *Storage) findPersonByName(
	ctx context.Context,
	name string,
) (models.Person, error) {
	const op = "postgres.findPersonByName"

	query := `SELECT number, recording_day, person FROM membership JOIN person
				ON membership.person = person.name AND person.name = $1`

	row := s.db.QueryRow(ctx, query, name)

	var person models.Person
	err := row.Scan(&person)
	if err != nil {
		return models.Person{}, fmt.Errorf("%s: %w", op, err)
	}

	return person, nil
}
