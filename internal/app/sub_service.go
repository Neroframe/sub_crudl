package app

import (
	"context"

	"github.com/Neroframe/sub_crudl/internal/domain"
	"github.com/google/uuid"
)

type SubscriptionUsecase interface {
	Create(ctx context.Context, input CreateInput) (*domain.Subscription, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	List(ctx context.Context, userID *uuid.UUID, serviceName *string) ([]domain.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*domain.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Aggregate(ctx context.Context, filter AggregationFilter) (int, error)
}
