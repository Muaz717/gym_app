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
	FindSubsByPersonName(ctx context.Context, name string) ([]models.Subscription, error)
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

			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
	}

	log.Info("person registered", "pid", personId)

	return personId, nil
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

func (p *PersonService) FindSubsByPersonName(ctx context.Context, name string) ([]models.Subscription, error) {
	const op = "services.person.findPersonByName"

	log := p.log.With(
		slog.String("op", op),
	)

	subscriptions, err := p.personStorage.FindSubsByPersonName(ctx, name)
	if err != nil {
		log.Warn("error", sl.Error(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//var subss []models.Subscription
	//for _, subscription := range subscriptions {
	//	subss = append(subss, services.EnrichSubscription(subscription))
	//}

	log.Info("Person are found")

	return subscriptions, nil
}
