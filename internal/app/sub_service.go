package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Neroframe/sub_crudl/internal/domain"
	"github.com/Neroframe/sub_crudl/pkg/logger"
	"github.com/google/uuid"
)

type CreateInput struct {
	ServiceName string
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
	Price       int32
}

type UpdateInput struct {
	ServiceName *string
	StartDate   *time.Time
	EndDate     *time.Time
	Price       *int32
}

type service struct {
	repo SubscriptionRepository
	log  *logger.Logger
}

func NewSubscriptionService(repo SubscriptionRepository, logger *logger.Logger) SubscriptionService {
	return &service{repo: repo, log: logger}
}

var (
	ErrNotFound     = errors.New("subscription not found")
	ErrInvalidInput = errors.New("invalid input")
)

func (s *service) Create(ctx context.Context, input CreateInput) (*domain.Subscription, error) {
	log := s.log.With("service", "Create")
	log.Debug("creating subscription", "input", input)

	// Input validation
	if input.ServiceName == "" {
		log.Error("service_name is required")
		return nil, fmt.Errorf("%w: service_name", ErrInvalidInput)
	}
	if input.Price < 0 {
		log.Error("price must be non-negative", "price", input.Price)
		return nil, fmt.Errorf("%w: price", ErrInvalidInput)
	}
	if input.EndDate != nil && input.StartDate.After(*input.EndDate) {
		log.Error("start_date cannot be after end_date", "start", input.StartDate, "end", *input.EndDate)
		return nil, fmt.Errorf("%w: date range", ErrInvalidInput)
	}
	sub := &domain.Subscription{
		ID:          uuid.New(),
		ServiceName: input.ServiceName,
		UserID:      input.UserID,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		Price:       input.Price,
	}

	if err := s.repo.Create(ctx, sub); err != nil {
		s.log.Error("repo.Create failed", "err", err)
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	s.log.Info("service.Create success", "id", sub.ID)
	return sub, nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	log := s.log.With("service", "Get", "id", id)
	log.Debug("fetching subscription")

	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Error("repo.GetByID failed", "error", err)
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	if sub == nil {
		log.Info("subscription not found")
		return nil, ErrNotFound
	}

	log.Info("subscription fetched", "id", id)
	return sub, nil
}

func (s *service) List(ctx context.Context, userID *uuid.UUID, serviceName *string) ([]*domain.Subscription, error) {
	log := s.log.With("service", "List", "user_id", userID, "service_name", serviceName)
	log.Debug("listing subscriptions")

	subs, err := s.repo.List(ctx, userID, serviceName)
	if err != nil {
		log.Error("repo.List failed", "error", err)
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	log.Info("subscriptions listed", "count", len(subs))
	return subs, nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*domain.Subscription, error) {
	log := s.log.With("service", "Update", "id", id)
	log.Debug("updating subscription", "input", input)

	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Error("repo.GetByID failed", "error", err)
		return nil, fmt.Errorf("failed to fetch subscription: %w", err)
	}
	if sub == nil {
		log.Info("subscription not found")
		return nil, ErrNotFound
	}

	// Validate input
	if input.ServiceName != nil {
		if *input.ServiceName == "" {
			log.Error("service_name cannot be empty")
			return nil, fmt.Errorf("%w: service_name", ErrInvalidInput)
		}
		sub.ServiceName = *input.ServiceName
	}
	if input.Price != nil {
		if *input.Price < 0 {
			log.Error("price must be non-negative", "price", *input.Price)
			return nil, fmt.Errorf("%w: price", ErrInvalidInput)
		}
		sub.Price = *input.Price
	}
	if input.StartDate != nil {
		sub.StartDate = *input.StartDate
	}
	if input.EndDate != nil {
		sub.EndDate = input.EndDate
	}
	if sub.EndDate != nil && sub.StartDate.After(*sub.EndDate) {
		log.Error("start_date cannot be after end_date", "start", sub.StartDate, "end", *sub.EndDate)
		return nil, fmt.Errorf("%w: date range", ErrInvalidInput)
	}

	if err := s.repo.Update(ctx, sub); err != nil {
		log.Error("repo.Update failed", "error", err)
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	log.Info("subscription updated", "id", id)
	return sub, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	log := s.log.With("service", "Delete", "id", id)
	log.Debug("deleting subscription")

	if err := s.repo.Delete(ctx, id); err != nil {
		log.Error("repo.Delete failed", "error", err)
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	log.Info("subscription deleted", "id", id)
	return nil
}

func (s *service) Aggregate(ctx context.Context, filter AggregationFilter) (int32, error) {
	log := s.log.With("service", "Aggregate", "filter", filter)
	log.Debug("aggregating subscriptions")

	// Validate date range
	if filter.StartPeriod.After(filter.EndPeriod) {
		log.Error("start period cannot be after end period", "start", filter.StartPeriod, "end", filter.EndPeriod)
		return 0, fmt.Errorf("%w: start_period > end_period", ErrInvalidInput)
	}

	total, err := s.repo.AggregateCost(ctx, filter)
	if err != nil {
		log.Error("repo.AggregateCost failed", "error", err)
		return 0, fmt.Errorf("failed to aggregate subscription cost: %w", err)
	}

	log.Info("subscription cost aggregated", "total", total)
	return int32(total), nil
}
