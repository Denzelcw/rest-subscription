package usecases

import (
	"context"
	"fmt"
	"log/slog"
	"subscription/internal/domain"
	"subscription/internal/http_server/dto"
	"subscription/internal/lib/logger/sl"
	"subscription/internal/storage/postgres"

	"github.com/google/uuid"
)

type SubscriptionStorage interface {
	AddUserSubscription(ctx context.Context, dto dto.CreateUserSubDTO) (int64, error)
	GetUserSubscriptionById(ctx context.Context, id int) (*domain.UserSubscription, error)
	GetUserSubscriptionsListByUUID(ctx context.Context, userID uuid.UUID) ([]*domain.UserSubscription, error)
	DeleteUserSubscriptionByID(ctx context.Context, id int) error
	UpdateUserSubscription(ctx context.Context, dto dto.UpdateUserSubDTO) (*domain.UserSubscription, error)
	CalculateTotalCost(ctx context.Context, dto dto.TotalCost) (int64, error)
}

type UserSubscriptionService struct {
	log     *slog.Logger
	storage SubscriptionStorage
}

func NewSubscriptionService(storage *postgres.Storage, log *slog.Logger) *UserSubscriptionService {
	return &UserSubscriptionService{storage: storage, log: log}
}

func (s *UserSubscriptionService) Add(ctx context.Context, dto dto.CreateUserSubDTO) (int64, error) {
	const op = "subscription_service.Add"

	id, err := s.storage.AddUserSubscription(ctx, dto)
	if err != nil {
		s.log.Error("can't add subscription", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *UserSubscriptionService) GetById(ctx context.Context, id int) (*domain.UserSubscription, error) {
	const op = "subscription_service.GetById"

	subscription, err := s.storage.GetUserSubscriptionById(ctx, id)
	if err != nil {
		s.log.Error("can't get subscription", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return subscription, nil
}

func (s *UserSubscriptionService) GetListByUUID(ctx context.Context, userId uuid.UUID) ([]*domain.UserSubscription, error) {
	const op = "subscription_service.GetListByUUID"

	subs, err := s.storage.GetUserSubscriptionsListByUUID(ctx, userId)
	if err != nil {
		s.log.Error("can't get subscriptions list", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return subs, nil
}

func (s *UserSubscriptionService) DeleteById(ctx context.Context, id int) error {
	const op = "subscription_service.DeleteById"

	err := s.storage.DeleteUserSubscriptionByID(ctx, id)

	if err != nil {
		s.log.Error("can't delete subscription", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *UserSubscriptionService) UpdateById(ctx context.Context, dto dto.UpdateUserSubDTO) (*domain.UserSubscription, error) {
	const op = "subscription_service.UpdateById"

	sub, err := s.storage.UpdateUserSubscription(ctx, dto)

	if err != nil {
		s.log.Error("can't delete subscription", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return sub, nil
}

func (s *UserSubscriptionService) TotalCost(ctx context.Context, cost dto.TotalCost) (int64, error) {
	const op = "subscription_service.TotalCost"

	totalCost, err := s.storage.CalculateTotalCost(ctx, cost)

	if err != nil {
		s.log.Error("can't get totalCost list", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return totalCost, nil
}
