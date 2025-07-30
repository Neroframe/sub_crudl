package app

import (
	"context"

	queries "github.com/Neroframe/sub_crudl/internal/infra/postgres/queries/generated"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, arg queries.CreateSubscriptionParams) error
	GetByID(ctx context.Context, id uuid.UUID) (queries.Subscription, error)
	List(ctx context.Context, userID *uuid.UUID, serviceName *string) ([]queries.Subscription, error)
	Update(ctx context.Context, arg queries.UpdateSubscriptionParams) error
	Delete(ctx context.Context, id uuid.UUID) error
	AggregateCost(ctx context.Context, arg queries.AggregateCostParams) (interface{}, error)
}
