package usecases

import (
	"fmt"
	"log/slog"
	"task_manager/internal/domain"
	"task_manager/internal/http_server/dto"
	"task_manager/internal/lib/logger/sl"
	"task_manager/internal/storage/postgres"

	"github.com/google/uuid"
)

type SubscriptionStorage interface {
	AddSubscription(dto dto.CreateSubDTO) (int64, error)
	GetSubscriptionById(id int) (*domain.Subscription, error)
	GetSubscriptionsListByUUID(userID uuid.UUID) ([]*domain.Subscription, error)
	DeleteSubscriptionByID(id int) error
	UpdateSubscription(dto dto.UpdateSubDTO) (*domain.Subscription, error)
	CalculateTotalCost(dto dto.TotalCost) (int64, error)
}

type SubscriptionService struct {
	log     *slog.Logger
	storage SubscriptionStorage
}

func NewSubscriptionService(storage *postgres.Storage, log *slog.Logger) *SubscriptionService {
	return &SubscriptionService{storage: storage, log: log}
}

func (s *SubscriptionService) Add(dto dto.CreateSubDTO) (int64, error) {
	const op = "subscription_service.Add"

	id, err := s.storage.AddSubscription(dto)
	if err != nil {
		s.log.Error("can't add subscription", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *SubscriptionService) GetById(id int) (*domain.Subscription, error) {
	const op = "subscription_service.GetById"

	subscription, err := s.storage.GetSubscriptionById(id)
	if err != nil {
		s.log.Error("can't get subscription", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return subscription, nil
}

func (s *SubscriptionService) GetListByUUID(userId uuid.UUID) ([]*domain.Subscription, error) {
	const op = "subscription_service.GetListByUUID"

	subs, err := s.storage.GetSubscriptionsListByUUID(userId)
	if err != nil {
		s.log.Error("can't get subscriptions list", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return subs, nil
}

func (s *SubscriptionService) DeleteById(id int) error {
	const op = "subscription_service.DeleteById"

	err := s.storage.DeleteSubscriptionByID(id)

	if err != nil {
		s.log.Error("can't delete subscription", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *SubscriptionService) UpdateById(dto dto.UpdateSubDTO) (*domain.Subscription, error) {
	const op = "subscription_service.UpdateById"

	sub, err := s.storage.UpdateSubscription(dto)

	if err != nil {
		s.log.Error("can't delete subscription", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return sub, nil
}

func (s *SubscriptionService) TotalCost(cost dto.TotalCost) (int64, error) {
	const op = "subscription_service.TotalCost"

	totalCost, err := s.storage.CalculateTotalCost(cost)

	if err != nil {
		s.log.Error("can't get totalCost list", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return totalCost, nil
}
