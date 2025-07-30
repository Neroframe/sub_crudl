package postgres

import (
	"context"
	"database/sql"

	"github.com/Neroframe/sub_crudl/internal/app"
	generated "github.com/Neroframe/sub_crudl/internal/infra/postgres/queries/generated"
	queries "github.com/Neroframe/sub_crudl/internal/infra/postgres/queries/generated"
	"github.com/google/uuid"
)

type repo struct {
	q *generated.Queries
}

func NewSubscriptionRepo(db *sql.DB) app.SubscriptionRepository {
	return &repo{
		q: generated.New(db),
	}
}

func (r *repo) Create(ctx context.Context, arg queries.CreateSubscriptionParams) error {
	return r.q.CreateSubscription(ctx, arg)
}

func (r *repo) GetByID(ctx context.Context, id uuid.UUID) (queries.Subscription, error) {
	return r.q.GetSubscriptionByID(ctx, id)
}

func (r *repo) List(ctx context.Context, userID *uuid.UUID, serviceName *string) ([]queries.Subscription, error) {
	// sqlc expects concrete values, not pointers
	var uid uuid.UUID
	var svc string
	if userID != nil {
		uid = *userID
	}
	if serviceName != nil {
		svc = *serviceName
	}
	params := queries.ListSubscriptionsPaginatedParams{
		Column1: uid,
		Column2: svc,
	}
	return r.q.ListSubscriptionsPaginated(ctx, params)
}

func (r *repo) Update(ctx context.Context, arg queries.UpdateSubscriptionParams) error {
	return r.q.UpdateSubscription(ctx, arg)
}

func (r *repo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteSubscription(ctx, id)
}

// AggregateCost computes the total subscription cost based on filters
func (r *repo) AggregateCost(ctx context.Context, arg queries.AggregateCostParams) (interface{}, error) {
	return r.q.AggregateCost(ctx, arg)
}
