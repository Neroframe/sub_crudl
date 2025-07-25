package domain

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID
	ServiceName string
	Price       int32
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
}
