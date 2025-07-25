package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Neroframe/sub_crudl/internal/app"
	"github.com/Neroframe/sub_crudl/internal/domain"
	"github.com/Neroframe/sub_crudl/pkg/logger"
	"github.com/google/uuid"

	"github.com/jmoiron/sqlx"
)

type dbSub struct {
	ID          uuid.UUID  `db:"id"`
	ServiceName string     `db:"service_name"`
	Price       int32      `db:"price"`
	UserID      uuid.UUID  `db:"user_id"`
	StartDate   time.Time  `db:"start_date"`
	EndDate     *time.Time `db:"end_date"`
}

type repo struct {
	db  *sqlx.DB
	log *logger.Logger
}

func NewSubscriptionRepo(db *sqlx.DB, logger *logger.Logger) app.SubscriptionRepository {
	return &repo{db: db, log: logger}
}

func (r *repo) Create(ctx context.Context, sub *domain.Subscription) error {
	log := r.log.With("repo", "Create")
	log.Debug("inserting subscription", "id", sub.ID)

	query := `INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
				VALUES (:id, :service_name, :price, :user_id, :start_date, :end_date)`

	args := dbSub{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate,
		EndDate:     sub.EndDate,
	}

	_, err := r.db.NamedExecContext(ctx, query, args)
	if err != nil {
		log.Error("failed to insert subscription", "error", err)
		return fmt.Errorf("repo.Create: %w", err)
	}

	log.Info("subscription inserted", "id", sub.ID)
	return err
}

func (r *repo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	log := r.log.With("repo", "GetByID")
	log.Debug("fetching subscription by id", "id", id)

	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1`

	var sub domain.Subscription
	err := r.db.GetContext(ctx, &sub, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info("subscription not found", "id", id)
			return nil, nil
		}
		log.Error("failed to fetch subscription", "id", id, "error", err)
		return nil, fmt.Errorf("repo.GetByID: %w", err)
	}

	log.Info("subscription fetched", "id", id)
	return &sub, nil
}

func (r *repo) List(ctx context.Context, userID *uuid.UUID, serviceName *string) ([]*domain.Subscription, error) {
	log := r.log.With("repo", "List")
	log.Debug("listing subscriptions", "user_id", userID, "service_name", serviceName)

	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions`
	var args []interface{}
	var conditions []string

	if userID != nil {
		conditions = append(conditions, "user_id=$"+strconv.Itoa(len(args)+1))
		args = append(args, *userID)
	}

	if serviceName != nil {
		conditions = append(conditions, "service_name ILIKE $"+strconv.Itoa(len(args)+1))
		args = append(args, "%"+*serviceName+"%")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var subs []*domain.Subscription
	if err := r.db.SelectContext(ctx, &subs, query, args...); err != nil {
		log.Error("failed to list subscriptions", "error", err)
		return nil, fmt.Errorf("repo.List: %w", err)
	}

	log.Info("subscriptions listed", "count", len(subs))
	return subs, nil
}

func (r *repo) Update(ctx context.Context, sub *domain.Subscription) error {
	log := r.log.With("repo", "Update")
	log.Debug("updating subscription", "id", sub.ID)

	query := `UPDATE subscriptions
				SET service_name = :service_name,
					price = :price,
					start_date = :start_date,
					end_date = :end_date,
				WHERE id=:id;`

	args := dbSub{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		StartDate:   sub.StartDate,
		EndDate:     sub.EndDate,
	}

	if _, err := r.db.NamedExecContext(ctx, query, args); err != nil {
		log.Error("failed to update subscription", "id", sub.ID, "error", err)
		return fmt.Errorf("repo.Update: %w", err)
	}

	log.Info("subscription updated", "id", sub.ID)
	return nil
}

func (r *repo) Delete(ctx context.Context, id uuid.UUID) error {
	log := r.log.With("repo", "Delete")
	log.Debug("deleting subscription", "id", id)

	query := `DELETE FROM subscriptions WHERE id = $1`
	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		log.Error("failed to delete subscription", "id", id, "error", err)
		return fmt.Errorf("repo.Delete: %w", err)
	}

	log.Info("subscription deleted", "id", id)
	return nil
}

// AggregateCost returns the sum of Price for all subscriptions
// in the given period and matching filters
func (r *repo) AggregateCost(ctx context.Context, f app.AggregationFilter) (int32, error) {
	log := r.log.With("repo", "AggregateCost")
	log.Debug("aggregating subscription cost", "filter", f)

	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions` // ensure returns 0 if no match
	var args []interface{}
	var conditions []string

	if f.UserID != nil {
		conditions = append(conditions, "user_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, *f.UserID)
	}
	if f.ServiceName != nil {
		conditions = append(conditions, "service_name = $"+strconv.Itoa(len(args)+1))
		args = append(args, *f.ServiceName)
	}
	conditions = append(conditions, "start_date >= $"+strconv.Itoa(len(args)+1))
	args = append(args, f.StartPeriod)
	conditions = append(conditions, "start_date <= $"+strconv.Itoa(len(args)+1))
	args = append(args, f.EndPeriod)

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var sum int32
	err := r.db.GetContext(ctx, &sum, query, args...)
	if err != nil {
		log.Error("failed to aggregate subscription cost", "error", err)
		return 0, fmt.Errorf("repo.AggregateCost: %w", err)
	}

	log.Info("subscription cost aggregated", "total", sum)
	return sum, nil
}
