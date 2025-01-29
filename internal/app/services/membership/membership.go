package membershipService

import (
	"context"
	"errors"
	"fmt"
	"gym_app/internal/lib/logger/sl"
	"gym_app/internal/models"
	"log/slog"
)

type MembershipService struct {
	log               *slog.Logger
	membershipStorage MembershipStorage
}

type MembershipStorage interface {
	saveMembership(ctx context.Context, membership models.Membership) (int, error)
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

func (m *MembershipService) addMembership(ctx context.Context, membership models.Membership) (int, error) {
	const op = "services.membership.addMembership"

	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("Adding new membership")

	memId, err := m.membershipStorage.saveMembership(ctx, membership)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, sl.Error(err))
	}

	log.Info("membership registered", "mid", memId)

	return memId, nil
}
