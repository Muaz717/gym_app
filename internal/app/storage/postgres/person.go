package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gym_app/internal/app/storage"
	"gym_app/internal/models"
)

func (s *Storage) SavePerson(
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

func (s *Storage) FindAllPeople(ctx context.Context) ([]models.Person, error) {
	const op = "postgres.findAllPeople"

	query := `SELECT * FROM person`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.Person])
}

func (s *Storage) FindMemsByPersonName(
	ctx context.Context,
	name string,
) ([]models.Membership, error) {
	const op = "postgres.findPersonByName"

	query := `SELECT m1.number, m1.recording_day FROM membership m1 WHERE m1.person = $1`

	rows, err := s.db.Query(ctx, query, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var mems []models.Membership
	for rows.Next() {
		mem := models.Membership{}
		err := rows.Scan(&mem.Number, &mem.RecordingDay)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		mems = append(mems, mem)
	}

	return mems, nil
}
