package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	appdto "github.com/Neroframe/sub_crudl/internal/app/dto"
	"github.com/Neroframe/sub_crudl/internal/domain"
	queries "github.com/Neroframe/sub_crudl/internal/infra/postgres/queries/generated"
	"github.com/Neroframe/sub_crudl/pkg/logger"
	"github.com/google/uuid"
)

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

func (s *service) Create(ctx context.Context, input appdto.CreateInput) (*domain.Subscription, error) {
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

	sub := queries.CreateSubscriptionParams{
		ID:          uuid.New(),
		ServiceName: input.ServiceName,
		UserID:      input.UserID,
		StartDate:   input.StartDate,
		EndDate:     sql.NullTime{Time: *input.EndDate, Valid: input.EndDate != nil},
		Price:       input.Price,
	}

	if err := s.repo.Create(ctx, sub); err != nil {
		s.log.Error("repo.Create failed", "err", err)
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	var endDate *time.Time
	if sub.EndDate.Valid {
		endDate = &sub.EndDate.Time
	}

	s.log.Info("service.Create success", "id", sub.ID)
	return &domain.Subscription{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate,
		EndDate:     endDate,
		Price:       sub.Price,
	}, nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	log := s.log.With("service", "Get", "id", id)
	log.Debug("fetching subscription")

	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Info("subscription not found")
			return nil, ErrNotFound
		}
		log.Error("repo.GetByID failed", "error", err)
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return mapToDomain(sub), nil
}

func (s *service) List(ctx context.Context, userID *uuid.UUID, serviceName *string) ([]*domain.Subscription, error) {
	log := s.log.With("service", "List", "user_id", userID, "service_name", serviceName)
	log.Debug("listing subscriptions")

	subs, err := s.repo.List(ctx, userID, serviceName)
	if err != nil {
		log.Error("repo.List failed", "error", err)
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	var result []*domain.Subscription
	for _, sub := range subs {
		result = append(result, mapToDomain(sub))
	}

	log.Info("subscriptions listed", "count", len(result))
	return result, nil

}

func (s *service) Update(
	ctx context.Context,
	id uuid.UUID,
	input appdto.UpdateInput,
) (*domain.Subscription, error) {
	log := s.log.With("service", "Update", "id", id)
	log.Debug("updating subscription", "input", input)

	// 1) Fetch existing record from repo (SQLC type)
	qsub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Error("repo.GetByID failed", "error", err)
		return nil, fmt.Errorf("failed to fetch subscription: %w", err)
	}

	// 2) Map SQLC type → domain model
	var dom = &domain.Subscription{
		ID:          qsub.ID,
		ServiceName: qsub.ServiceName,
		Price:       qsub.Price,
		UserID:      qsub.UserID,
		StartDate:   qsub.StartDate,
	}
	if qsub.EndDate.Valid {
		dom.EndDate = &qsub.EndDate.Time
	}

	// 3) Apply updates + validate
	if input.ServiceName != nil {
		if *input.ServiceName == "" {
			log.Error("service_name cannot be empty")
			return nil, fmt.Errorf("%w: service_name", ErrInvalidInput)
		}
		dom.ServiceName = *input.ServiceName
	}
	if input.Price != nil {
		if *input.Price < 0 {
			log.Error("price must be non-negative", "price", *input.Price)
			return nil, fmt.Errorf("%w: price", ErrInvalidInput)
		}
		dom.Price = *input.Price
	}
	if input.StartDate != nil {
		dom.StartDate = *input.StartDate
	}
	if input.EndDate != nil {
		dom.EndDate = input.EndDate
	}
	if dom.EndDate != nil && dom.StartDate.After(*dom.EndDate) {
		log.Error("start_date cannot be after end_date",
			"start", dom.StartDate, "end", *dom.EndDate)
		return nil, fmt.Errorf("%w: date range", ErrInvalidInput)
	}

	// 4) Map domain → SQLC params
	params := queries.UpdateSubscriptionParams{
		ID:          dom.ID,
		ServiceName: dom.ServiceName,
		Price:       dom.Price,
		StartDate:   dom.StartDate,
		EndDate: sql.NullTime{Time: func() time.Time {
			if dom.EndDate != nil {
				return *dom.EndDate
			}
			return time.Time{}
		}(), Valid: dom.EndDate != nil},
	}

	// 5) Call repo.Update
	if err := s.repo.Update(ctx, params); err != nil {
		log.Error("repo.Update failed", "error", err)
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	log.Info("subscription updated", "id", id)
	return dom, nil
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

func (s *service) Aggregate(
	ctx context.Context,
	filter appdto.AggregationFilter,
) (int32, error) {
	log := s.log.With("service", "Aggregate", "filter", filter)
	log.Debug("aggregating subscriptions")

	// Validate date range
	if filter.StartPeriod.After(filter.EndPeriod) {
		log.Error("start_period cannot be after end_period",
			"start", filter.StartPeriod, "end", filter.EndPeriod)
		return 0, fmt.Errorf("%w: date range", ErrInvalidInput)
	}

	// Map app DTO → SQLC params
	params := queries.AggregateCostParams{
		Column1: func() uuid.UUID {
			if filter.UserID != nil {
				return *filter.UserID
			}
			return uuid.Nil
		}(),
		Column2: func() string {
			if filter.ServiceName != nil {
				return *filter.ServiceName
			}
			return ""
		}(),
		StartDate: filter.StartPeriod,
		EndDate:   sql.NullTime{Time: filter.EndPeriod, Valid: true},
	}

	// Call repo.AggregateCost (returns interface{})
	raw, err := s.repo.AggregateCost(ctx, params)
	if err != nil {
		log.Error("repo.AggregateCost failed", "error", err)
		return 0, fmt.Errorf("failed to aggregate subscription cost: %w", err)
	}

	// Cast to int64 (Postgres COALESCE SUM gives int64), then to int32
	total64, ok := raw.(int64)
	if !ok {
		return 0, fmt.Errorf("invalid aggregate type: %T", raw)
	}
	total := int32(total64)

	log.Info("subscription cost aggregated", "total", total)
	return total, nil
}

func mapToDomain(sub queries.Subscription) *domain.Subscription {
	var endDate *time.Time
	if sub.EndDate.Valid {
		endDate = &sub.EndDate.Time
	}
	return &domain.Subscription{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate,
		EndDate:     endDate,
		Price:       sub.Price,
	}
}
