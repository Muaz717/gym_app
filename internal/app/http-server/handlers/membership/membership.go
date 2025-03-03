package membershipHandler

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "gym_app/internal/lib/api/response"
	"gym_app/internal/lib/logger/sl"
	"gym_app/internal/models"
	"io"
	"log/slog"
	"net/http"
)

type MembershipService interface {
	AddMembership(ctx context.Context, membership models.Membership) (int, error)
	FindAllMemberships(ctx context.Context) ([]models.Membership, error)
}

type MembershipHandler struct {
	ctx               context.Context
	log               *slog.Logger
	membershipService MembershipService
}

func New(
	ctx context.Context,
	log *slog.Logger,
	membershipService MembershipService,
) *MembershipHandler {
	return &MembershipHandler{
		ctx:               ctx,
		log:               log,
		membershipService: membershipService,
	}
}

func (h *MembershipHandler) AddMembership(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.membership.addMembership"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var membership models.Membership

	err := render.DecodeJSON(r.Body, &membership)
	if errors.Is(err, io.EOF) {
		log.Error("request body is empty")

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, resp.Error("empty request"))

		return
	}
	if err != nil {
		log.Error("failed to decode request body", sl.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("failed to decode request"))

		return
	}

	memId, err := h.membershipService.AddMembership(h.ctx, membership)
	if err != nil {
		log.Error("failed to add membership", sl.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, resp.OK("failed to add membership"))

		return
	}

	log.Info("Membership added", slog.Int("person id", memId))
	render.JSON(w, r, resp.OK("Membership added"))
}

func (h *MembershipHandler) FindAllMemberships(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.membership.addMembership"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	memberships, err := h.membershipService.FindAllMemberships(h.ctx)
	if err != nil {
		log.Error("failed to get mems", sl.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("failed to get mems"))

		return
	}

	log.Info("Mems found")

	render.JSON(w, r, memberships)
}
