package app

import (
	"context"

	appdto "github.com/Neroframe/sub_crudl/internal/app/dto"
	"github.com/Neroframe/sub_crudl/internal/domain"
	"github.com/google/uuid"
)



type SubscriptionRepository interface {
	Create(ctx context.Context, sub *domain.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	List(ctx context.Context, userID *uuid.UUID, serviceName *string) ([]*domain.Subscription, error)
	Update(ctx context.Context, sub *domain.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error

	AggregateCost(ctx context.Context, f appdto.AggregationFilter) (int32, error)
}
