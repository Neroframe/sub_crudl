package dto

type CreateSubscriptionDTO struct {
	ServiceName string `json:"service_name" binding:"required"`
	UserID      string `json:"user_id" binding:"required,uuid"`
	StartDate   string `json:"start_date" binding:"required"` // format: MM-YYYY, validated manually
	EndDate     string `json:"end_date,omitempty"`            // optional, same format
	Price       int32  `json:"price" binding:"required,min=0"`
}

type UpdateSubscriptionDTO struct {
	ServiceName *string `json:"service_name,omitempty"`
	StartDate   *string `json:"start_date,omitempty"` // validated manually
	EndDate     *string `json:"end_date,omitempty"`   // validated manually
	Price       *int32  `json:"price,omitempty"`
}

type SubscriptionDTO struct {
	ID          string  `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ServiceName string  `json:"service_name" example:"Netflix"`
	Price       int32   `json:"price" example:"999"`
	UserID      string  `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	StartDate   string  `json:"start_date" example:"01-2025"` // MM-YYYY
	EndDate     *string `json:"end_date,omitempty" example:"12-2025"`
}
