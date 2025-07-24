package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateInput struct {
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
}

type UpdateInput struct {
	ServiceName *string
	Price       *int
	StartDate   *time.Time
	EndDate     *time.Time
}

type AggregationFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	StartPeriod time.Time
	EndPeriod   time.Time
}
