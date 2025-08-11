package handler

import (
	"context"
	"log/slog"
	"subscription/internal/domain"
	"subscription/internal/http_server/dto"
	"subscription/internal/usecases"
	"time"

	"github.com/google/uuid"
)

type UserSubUseCases interface {
	Add(ctx context.Context, dto dto.CreateUserSubDTO) (int64, error)
	GetById(ctx context.Context, id int) (*domain.UserSubscription, error)
	GetListByUUID(ctx context.Context, userId uuid.UUID) ([]*domain.UserSubscription, error)
	DeleteById(ctx context.Context, id int) error
	UpdateById(ctx context.Context, dto dto.UpdateUserSubDTO) (*domain.UserSubscription, error)
	TotalCost(ctx context.Context, cost dto.TotalCost) (int64, error)
}

type UserSubscriptionHandler struct {
	log     *slog.Logger
	service UserSubUseCases
	timeOut time.Duration
}

func NewUserSubscriptionHandler(
	service *usecases.UserSubscriptionService,
	l *slog.Logger,
	timeOut time.Duration,
) *UserSubscriptionHandler {
	return &UserSubscriptionHandler{service: service, log: l, timeOut: timeOut}
}
