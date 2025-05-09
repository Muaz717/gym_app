package personService

import (
	"context"
	"errors"
	"fmt"
	"gym_app/internal/lib/logger/sl"
	"gym_app/internal/models"
	"gym_app/internal/storage"
	"log/slog"
)

type PersonService struct {
	log           *slog.Logger
	personStorage PersonStorage
}

type PersonStorage interface {
	SavePerson(ctx context.Context, person models.Person) (int, error)
	FindAllPeople(ctx context.Context) ([]models.Person, error)
	UpdatePerson(ctx context.Context, person models.Person, pID int) (int, error)
	DeletePerson(ctx context.Context, pID int) error
	FindPersonByName(ctx context.Context, name string) (models.Person, error)
}

var (
	ErrPersonExists   = errors.New("person already exists")
	ErrPersonNotFound = errors.New("person not found")
)

func New(
	log *slog.Logger,
	personStorage PersonStorage,
) *PersonService {
	return &PersonService{
		log:           log,
		personStorage: personStorage,
	}
}

func (p *PersonService) AddPerson(ctx context.Context, person models.Person) (int, error) {

	const op = "services.person.addPerson"

	log := p.log.With(
		slog.String("op", op),
	)

	log.Info("Registering new user")

	personId, err := p.personStorage.SavePerson(ctx, person)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", sl.Error(err))

			return 0, fmt.Errorf("%s: %w", op, ErrPersonExists)
		}

		return 0, err
	}

	log.Info("person registered", "pid", personId)

	return personId, nil
}

func (p *PersonService) UpdatePerson(ctx context.Context, person models.Person, pID int) (int, error) {
	const op = "services.person.UpdatePerson"

	log := p.log.With(
		slog.String("op", op),
	)

	log.Info("Updating user")

	personId, err := p.personStorage.UpdatePerson(ctx, person, pID)
	if err != nil {
		if errors.Is(err, storage.ErrPersonNotFound) {
			log.Warn("user not found", sl.Error(err))

			return 0, fmt.Errorf("%s: %w", op, ErrPersonNotFound)
		}

		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", sl.Error(err))

			return 0, fmt.Errorf("%s: %w", op, ErrPersonExists)
		}

		return 0, err
	}

	log.Info("person updated", "pid", personId)

	return personId, nil
}

func (p *PersonService) DeletePerson(ctx context.Context, pID int) error {
	const op = "services.person.DeletePerson"

	log := p.log.With(
		slog.String("op", op),
	)

	log.Info("Deleting user")

	err := p.personStorage.DeletePerson(ctx, pID)
	if err != nil {
		if errors.Is(err, storage.ErrPersonNotFound) {
			log.Warn("user not found", sl.Error(err))

			return fmt.Errorf("%s: %w", op, ErrPersonNotFound)
		}
	}

	log.Info("person deleted")

	return nil
}

func (p *PersonService) FindPersonByName(ctx context.Context, name string) (models.Person, error) {
	const op = "services.person.FindPersonByName"

	log := p.log.With(
		slog.String("op", op),
	)

	log.Info("Finding user by name")

	person, err := p.personStorage.FindPersonByName(ctx, name)
	if err != nil {
		if errors.Is(err, storage.ErrPersonNotFound) {
			log.Warn("user not found", sl.Error(err))

			return models.Person{}, fmt.Errorf("%s: %w", op, ErrPersonNotFound)
		}
	}

	log.Info("person found")

	return person, nil
}

func (p *PersonService) FindAllPeople(ctx context.Context) ([]models.Person, error) {
	const op = "services.person.FindAllPeople"

	log := p.log.With(
		slog.String("op", op),
	)

	log.Info("Starting to find people")

	allPeople, err := p.personStorage.FindAllPeople(ctx)
	if err != nil {
		log.Warn("error", sl.Error(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("People are found")

	return allPeople, nil
}
