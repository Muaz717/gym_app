package membershipService

import (
	"context"
	"errors"
	"fmt"
	"gym_app/internal/lib/logger/sl"
	"gym_app/internal/models"
	"gym_app/internal/services"
	"log/slog"
	"time"
)

type MembershipService struct {
	log               *slog.Logger
	membershipStorage MembershipStorage
}

type MembershipStorage interface {
	SaveMembership(ctx context.Context, membership models.Membership) (int, error)
	FindAllMemberships(ctx context.Context) ([]models.Membership, error)
}

var (
	ErrMemExists = errors.New("membership with that number already exists")
)

func New(
	log *slog.Logger,
	membershipStorage MembershipStorage,
) *MembershipService {
	return &MembershipService{
		log:               log,
		membershipStorage: membershipStorage,
	}
}

func (m *MembershipService) AddMembership(ctx context.Context, membership models.Membership) (int, error) {
	const op = "services.membership.AddMembership"

	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("Adding new membership")

	membership.RecordingDay = time.Now()

	memId, err := m.membershipStorage.SaveMembership(ctx, membership)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, sl.Error(err))
	}

	log.Info("membership registered", "mid", memId)

	return memId, nil
}

func (m *MembershipService) FindAllMemberships(ctx context.Context) ([]models.Membership, error) {
	const op = "services.membership.FindAllMemberships"

	log := m.log.With(
		slog.String("op", op),
	)

	memberships, err := m.membershipStorage.FindAllMemberships(ctx)
	if err != nil {
		log.Warn("error", sl.Error(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}
	var mems []models.Membership
	for _, membership := range memberships {
		mems = append(mems, services.EnrichMembership(membership))
	}

	log.Info("Memberships are found")

	return mems, nil
}
