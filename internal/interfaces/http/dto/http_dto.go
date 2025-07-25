package dto

type CreateSubscriptionDTO struct {
	ServiceName string `json:"service_name" binding:"required"`
	UserID      string `json:"user_id" binding:"required,uuid"`
	StartDate   string `json:"start_date" binding:"required,datetime=01-2006"` // MM-YYYY
	EndDate     string `json:"end_date" binding:"omitempty,datetime=01-2006"`
	Price       int32  `json:"price" binding:"required,min=0"`
}

type UpdateSubscriptionDTO struct {
	ServiceName *string `json:"service_name,omitempty"`
	StartDate   *string `json:"start_date,omitempty" binding:"omitempty,datetime=01-2006"`
	EndDate     *string `json:"end_date,omitempty" binding:"omitempty,datetime=01-2006"`
	Price       *int32  `json:"price,omitempty"`
}
