package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *Handler) {
	api := r.Group("/subscriptions")
	{
		api.POST("", h.CreateSubscription)
		api.GET("", h.ListSubscriptions)
		api.GET("/:id", h.GetSubscription)
		api.PUT("/:id", h.UpdateSubscription)
		api.DELETE("/:id", h.DeleteSubscription)
	}

	r.Get("/subscriptions/aggregate", h.AggregateSubscriptions)
}
