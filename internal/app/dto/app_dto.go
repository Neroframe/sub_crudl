package appdto

import (
	"time"

	"github.com/google/uuid"
)

type CreateInput struct {
	ServiceName string
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
	Price       int32
}

type UpdateInput struct {
	ServiceName *string
	StartDate   *time.Time
	EndDate     *time.Time
	Price       *int32
}

type AggregationFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	StartPeriod time.Time
	EndPeriod   time.Time
}

