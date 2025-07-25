package app

import (
	"context"
	"time"

	"github.com/Neroframe/sub_crudl/internal/domain"
	"github.com/google/uuid"
)

type AggregationFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	StartPeriod time.Time
	EndPeriod   time.Time
}

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *domain.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	List(ctx context.Context, userID *uuid.UUID, serviceName *string) ([]*domain.Subscription, error)
	Update(ctx context.Context, sub *domain.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error

	AggregateCost(ctx context.Context, f AggregationFilter) (int32, error)
}
