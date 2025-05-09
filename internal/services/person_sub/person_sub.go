package personSubService

import (
	"context"
	"errors"
	"fmt"
	"gym_app/internal/lib/logger/sl"
	"gym_app/internal/models"
	"gym_app/internal/storage"
	"log/slog"
	"time"
)

const (
	activeStatus  = "active"
	frozenStatus  = "frozen"
	expiredStatus = "expired"
)

type PersonSubStorage interface {
	AddPersonSub(ctx context.Context, personSub models.PersonSubscription) (string, error)
	GetPersonSubByNumber(ctx context.Context, number string) (models.PersonSubscription, error)
	GetAllPersonSubs(ctx context.Context) ([]models.PersonSubscription, error)
	DeletePersonSub(ctx context.Context, number string) error
	FindPersonSubByPersonName(ctx context.Context, name string) ([]models.PersonSubscription, error)
	UpdatePersonSubStatus(ctx context.Context, number string, status string) error
}

var (
	ErrSubExists      = errors.New("subscription with that number already exists")
	ErrSubNotFound    = errors.New("subscription not found")
	ErrPersonNotFound = errors.New("person not found")
)

type PersonSubService struct {
	log              *slog.Logger
	personSubStorage PersonSubStorage
}

func New(log *slog.Logger, personSubStorage PersonSubStorage) *PersonSubService {
	return &PersonSubService{
		log:              log,
		personSubStorage: personSubStorage,
	}
}

func (p *PersonSubService) AddPersonSub(ctx context.Context, personSubStrDate models.PersonSubStrDate) (string, error) {
	const op = "services.personSub.AddPersonSub"

	log := p.log.With(
		slog.String("op", op),
	)

	log.Info("Adding new person subscription")

	personSub := convertToPersonSub(personSubStrDate)

	personSubNumber, err := p.personSubStorage.AddPersonSub(ctx, personSub)
	if err != nil {

		if errors.Is(err, storage.ErrSubscriptionExists) {
			log.Warn("subscription already exists", slog.String("number", personSub.Number), sl.Error(err))

			return "", fmt.Errorf("%s: %w", op, ErrSubExists)
		} else if errors.Is(err, storage.ErrPersonNotFound) {

			log.Warn("person not found", slog.String("number", personSub.Number), sl.Error(err))

			return "", fmt.Errorf("%s: %w", op, ErrPersonNotFound)
		}

		log.Error("failed to add person subscription", sl.Error(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("person subscription added", "number", personSubNumber)

	return personSubNumber, nil
}

func (p *PersonSubService) DeletePersonSub(ctx context.Context, number string) error {
	const op = "services.personSub.DeletePersonSub"

	log := p.log.With(
		slog.String("op", op),
	)

	log.Info("Deleting person subscription")

	err := p.personSubStorage.DeletePersonSub(ctx, number)
	if err != nil {

		if errors.Is(err, storage.ErrSubscriptionNotFound) {
			log.Warn("subscription not found", slog.String("number", number), sl.Error(err))

			return fmt.Errorf("%s: %w", op, ErrSubNotFound)
		}

		log.Error("failed to delete person subscription", sl.Error(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("person subscription deleted", "number", number)

	return nil
}

func (p *PersonSubService) GetPersonSubByNumber(ctx context.Context, number string) (models.PersonSubStrDate, error) {
	const op = "services.personSub.FindPersonSubByNumber"

	log := p.log.With(
		slog.String("op", op),
	)

	log.Info("Getting person subscription by number")

	personSub, err := p.personSubStorage.GetPersonSubByNumber(ctx, number)
	if err != nil {
		return models.PersonSubStrDate{}, err
	}

	log.Info("dsd", personSub.SubscriptionID)

	log.Info("person subscription found", "number", number)

	personSubStrDate := convertToPersonSubStrDate(personSub)

	return personSubStrDate, nil
}

func (p *PersonSubService) GetAllPersonSubs(ctx context.Context) ([]models.PersonSubStrDate, error) {
	const op = "services.personSub.GetAllPersonSubs"

	log := p.log.With(
		slog.String("op", op),
	)

	log.Info("Getting all person subscriptions")

	personSubs, err := p.personSubStorage.GetAllPersonSubs(ctx)
	if err != nil {
		return nil, err
	}

	log.Info("all person subscriptions found")

	var personSubsStrDate []models.PersonSubStrDate

	for _, personSub := range personSubs {
		personSubsStrDate = append(personSubsStrDate, convertToPersonSubStrDate(personSub))
	}

	return personSubsStrDate, nil
}

func (p *PersonSubService) FindPersonSubByPersonName(ctx context.Context, name string) ([]models.PersonSubStrDate, error) {
	const op = "services.personSub.FindPersonSubByPersonName"

	log := p.log.With(
		slog.String("op", op),
	)

	log.Info("Finding person subscription by person name")

	personSubs, err := p.personSubStorage.FindPersonSubByPersonName(ctx, name)
	if err != nil {

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var personSubsStrDate []models.PersonSubStrDate
	for _, personSub := range personSubs {
		personSubsStrDate = append(personSubsStrDate, convertToPersonSubStrDate(personSub))
	}

	log.Info("person subscriptions found")

	return personSubsStrDate, nil
}

func (p *PersonSubService) UpdateStatuses(ctx context.Context) error {
	const op = "services.personSub.UpdateStatuses"

	log := p.log.With(
		slog.String("op", op),
	)

	log.Info("Updating person subscription statuses")

	subs, err := p.personSubStorage.GetAllPersonSubs(ctx)
	if err != nil {
		log.Error("failed to get all person subscriptions", sl.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	today := time.Now().Truncate(24 * time.Hour)

	for _, sub := range subs {
		newStatus := ""

		if sub.StartDate.After(today) {
			newStatus = frozenStatus
		} else if sub.EndDate.Before(today) {
			newStatus = expiredStatus
		} else {
			newStatus = activeStatus
		}

		if sub.Status != newStatus {
			err := p.personSubStorage.UpdatePersonSubStatus(ctx, sub.Number, newStatus)
			if err != nil {
				log.Error("failed to update person subscription status", sl.Error(err))
				return fmt.Errorf("%s: %w", op, err)
			}
		}
	}

	log.Info("person subscription statuses updated")
	return nil
}

func convertToPersonSubStrDate(personSub models.PersonSubscription) models.PersonSubStrDate {
	return models.PersonSubStrDate{
		PersonID:       personSub.PersonID,
		SubscriptionID: personSub.SubscriptionID,
		Number:         personSub.Number,
		StartDate:      personSub.StartDate.Format("02-01-2006"),
		EndDate:        personSub.EndDate.Format("02-01-2006"),
		Status:         personSub.Status,
	}
}

func convertToPersonSub(personSubStrDate models.PersonSubStrDate) models.PersonSubscription {
	startDate, _ := time.Parse("02-01-2006", personSubStrDate.StartDate)
	endDate, _ := time.Parse("02-01-2006", personSubStrDate.EndDate)

	if startDate.IsZero() {
		startDate = time.Now()
	}
	if endDate.IsZero() {
		endDate = startDate.AddDate(0, 0, 30) // Default to one month subscription
	}

	if personSubStrDate.Status == "" {

		today := time.Now().Truncate(24 * time.Hour)
		start := startDate.Truncate(24 * time.Hour)

		if !today.Equal(start) {
			personSubStrDate.Status = frozenStatus
		} else {
			personSubStrDate.Status = activeStatus
		}
	}

	return models.PersonSubscription{
		PersonID:       personSubStrDate.PersonID,
		SubscriptionID: personSubStrDate.SubscriptionID,
		Number:         personSubStrDate.Number,
		StartDate:      startDate,
		EndDate:        endDate,
		Status:         personSubStrDate.Status,
	}
}
