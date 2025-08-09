package handler

import (
	"log/slog"
	"task_manager/internal/domain"
	"task_manager/internal/http_server/dto"
	"task_manager/internal/usecases"

	"github.com/google/uuid"
)

type SubUseCases interface {
	Add(dto dto.CreateSubDTO) (int64, error)
	GetById(id int) (*domain.Subscription, error)
	GetListByUUID(userId uuid.UUID) ([]*domain.Subscription, error)
	DeleteById(id int) error
	UpdateById(dto dto.UpdateSubDTO) (*domain.Subscription, error)
	TotalCost(cost dto.TotalCost) (int64, error)
}

type SubscriptionHandler struct {
	log     *slog.Logger
	service SubUseCases
}

func NewSubscriptionHandler(service *usecases.SubscriptionService, l *slog.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{service: service, log: l}
}
