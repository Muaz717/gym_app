package personService

import (
	"context"
	"errors"
	"fmt"
	"gym_app/internal/app/storage"
	"gym_app/internal/lib/logger/sl"
	"gym_app/internal/models"
	"log/slog"
)

type PersonService struct {
	log           *slog.Logger
	personStorage PersonStorage
}

type PersonStorage interface {
	savePerson(ctx context.Context, person models.Person) (int, error)
	findAllPeople(ctx context.Context) ([]models.Person, error)
	//findPersonByName(ctx context.Context, name string) (models.Person, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppId       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
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

func (p *PersonService) addPerson(ctx context.Context, person models.Person) (int, error) {

	const op = "services.person.addPerson"

	log := p.log.With(
		slog.String("op", op),
	)

	log.Info("Registering new user")

	personId, err := p.personStorage.savePerson(ctx, person)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", sl.Error(err))

			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
	}

	log.Info("person registered", "pid", personId)

	return personId, nil
}

func (p *PersonService) findAllPeople(ctx context.Context) ([]models.Person, error) {
	const op = "services.person.findAllPeople"

	log := p.log.With(
		slog.String("op", op),
	)

	log.Info("Starting to find people")

	allPeople, err := p.personStorage.findAllPeople(ctx)
	if err != nil {
		log.Warn("error", sl.Error(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("People are found")

	return allPeople, nil
}

//func (p *PersonService) findPersonByName(ctx context.Context, name string) (models.Person, error) {
//	const op = "services.person.findPersonByName"
//
//
//}
