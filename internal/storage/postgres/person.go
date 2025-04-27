package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gym_app/internal/models"
	"gym_app/internal/storage"
)

func (s *Storage) SavePerson(
	ctx context.Context,
	person models.Person,
) (int, error) {
	const op = "postgres.savePerson"

	query := `INSERT INTO person(full_name, phone) VALUES($1, $2) RETURNING id`

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

//func (s *Storage) FindSubsByPersonName(
//	ctx context.Context,
//	name string,
//) ([]models.Subscription, error) {
//	const op = "postgres.findPersonByName"
//
//	query := `SELECT m1.number, m1.recording_day FROM membership m1 WHERE m1.person = $1`
//
//	rows, err := s.db.Query(ctx, query, name)
//	if err != nil {
//		if errors.Is(err, pgx.ErrNoRows) {
//			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
//		}
//		return nil, fmt.Errorf("%s: %w", op, err)
//	}
//
//	var subs []models.Subscription
//	for rows.Next() {
//		sub := models.Subscription{}
//		err := rows.Scan(&sub.Number, &sub.RecordingDay)
//		if err != nil {
//			return nil, fmt.Errorf("unable to scan row: %w", err)
//		}
//		subs = append(subs, sub)
//	}
//
//	return subs, nil
//}
