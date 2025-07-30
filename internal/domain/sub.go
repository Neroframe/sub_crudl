package domain

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID  `db:"id"`
	ServiceName string     `db:"service_name"`
	Price       int32      `db:"price"`
	UserID      uuid.UUID  `db:"user_id"`
	StartDate   time.Time  `db:"start_date"`
	EndDate     *time.Time `db:"end_date"`
}
