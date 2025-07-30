package app

import (
	"context"

	appdto "github.com/Neroframe/sub_crudl/internal/app/dto"
	"github.com/Neroframe/sub_crudl/internal/domain"
	"github.com/google/uuid"
)

// SubscriptionService is the application boundary interface.
type SubscriptionService interface {
	Create(ctx context.Context, input appdto.CreateInput) (*domain.Subscription, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	List(ctx context.Context, userID *uuid.UUID, serviceName *string) ([]*domain.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, input appdto.UpdateInput) (*domain.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Aggregate(ctx context.Context, filter appdto.AggregationFilter) (int32, error)
}
