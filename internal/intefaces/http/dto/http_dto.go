package dto

type CreateSubscriptionDTO struct {
	ServiceName string `json:"service_name" binding:"required"`
	Price       int    `json:"price" binding:"required,min=0"`
	UserID      string `json:"user_id" binding:"required,uuid"`
	StartDate   string `json:"start_date" binding:"required,datetime=01-2006"` // MM-YYYY
	EndDate     string `json:"end_date" binding:"omitempty,datetime=01-2006"`
}

type UpdateSubscriptionDTO struct {
	ServiceName *string `json:"service_name,omitempty"`
	Price       *int    `json:"price,omitempty"`
	StartDate   *string `json:"start_date,omitempty" binding:"omitempty,datetime=01-2006"`
	EndDate     *string `json:"end_date,omitempty" binding:"omitempty,datetime=01-2006"`
}

type AggregationQueryDTO struct {
	UserID      string `form:"user_id" binding:"omitempty,uuid"`
	ServiceName string `form:"service_name" binding:"omitempty"`
	StartPeriod string `form:"start_period" binding:"required,datetime=01-2006"`
	EndPeriod   string `form:"end_period" binding:"required,datetime=01-2006"`
}

type SubscriptionResponseDTO struct {
	ID          string `json:"id"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"` // or better time.RFC3339
	EndDate     string `json:"end_date,omitempty"`
}
