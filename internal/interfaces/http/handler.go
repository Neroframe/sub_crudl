package httpapi

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Neroframe/sub_crudl/internal/app"
	"github.com/Neroframe/sub_crudl/internal/interfaces/http/dto"
	"github.com/Neroframe/sub_crudl/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}

type AggregateResponse struct {
	Total int `json:"total" example:"123"`
}

type Handler struct {
	SubService app.SubscriptionService
	log        *logger.Logger
}

func NewHandler(subService app.SubscriptionService, logger *logger.Logger) *Handler {
	return &Handler{SubService: subService, log: logger}
}

// CreateSubscription godoc
// @Summary     Create a new subscription
// @Description Create subscription with service name, price, user ID, start and optional end date
// @Tags        subscriptions
// @Accept      json
// @Produce     json
// @Param       subscription body dto.CreateSubscriptionDTO true "Subscription data"
// @Success     201 {object} dto.SubscriptionDTO
// @Failure     400 {object} httpapi.ErrorResponse
// @Failure     500 {object} httpapi.ErrorResponse
// @Router      /subscriptions [post]
func (h *Handler) CreateSubscription(c *gin.Context) {
	log := h.log.With("handler", "CreateSubscription")
	log.Debug("parsing request")

	var req dto.CreateSubscriptionDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("invalid request body", "error", err)
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fieldErr := ve[0]
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s failed %s validation", fieldErr.Field(), fieldErr.Tag())})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		}
		return
	}

	log.Debug("parsed request", "service_name", req.ServiceName, "user_id", req.UserID)

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		log.Error("invalid user_id format", "user_id", req.UserID, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		log.Error("invalid start_date format", "start_date", req.StartDate, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date"})
		return
	}

	var endDate *time.Time
	if req.EndDate != "" {
		t, err2 := time.Parse("01-2006", req.EndDate)
		if err2 != nil {
			log.Error("invalid end_date format", "end_date", req.EndDate, "error", err2)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date"})
			return
		}
		endDate = &t
	}

	input := app.CreateInput{
		ServiceName: req.ServiceName,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	sub, err := h.SubService.Create(c.Request.Context(), input)
	if err != nil {
		log.Error("failed to create subscription", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	log.Info("subscription created", "id", sub.ID, "user", sub.UserID)
	c.JSON(http.StatusCreated, sub)
}

// GetSubscription godoc
// @Summary     Get subscription by ID
// @Description Retrieve subscription details by subscription ID
// @Tags        subscriptions
// @Produce     json
// @Param       id path string true "Subscription ID"
// @Success     200 {object} dto.SubscriptionDTO
// @Failure     400 {object} httpapi.ErrorResponse
// @Failure     404 {object} httpapi.ErrorResponse
// @Failure     500 {object} httpapi.ErrorResponse
// @Router      /subscriptions/{id} [get]
func (h *Handler) GetSubscription(c *gin.Context) {
	log := h.log.With("handler", "GetSubscription")

	idStr := c.Param("id")
	log.Debug("received get request", "id", idStr)

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error("invalid subscription ID format", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	sub, err := h.SubService.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info("subscription not found", "id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		} else {
			log.Error("failed to retrieve subscription", "id", id, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subscription"})
		}
		return
	}

	log.Info("subscription retrieved", "id", sub.ID)
	c.JSON(http.StatusOK, sub)
}

// ListSubscriptions godoc
// @Summary     List subscriptions
// @Description Get all subscriptions, optionally filter by user_id and service_name
// @Tags        subscriptions
// @Produce     json
// @Param       user_id      query string false "User ID"
// @Param       service_name query string false "Service Name"
// @Success     200 {array}  dto.SubscriptionDTO
// @Failure     500 {object} httpapi.ErrorResponse
// @Router      /subscriptions [get]
func (h *Handler) ListSubscriptions(c *gin.Context) {
	log := h.log.With("handler", "ListSubscriptions")

	userIDStr := c.Query("user_id")
	serviceNameStr := c.Query("service_name")
	log.Debug("received list request", "user_id", userIDStr, "service_name", serviceNameStr)

	var userID *uuid.UUID
	if userIDStr != "" {
		parsed, err := uuid.Parse(userIDStr)
		if err != nil {
			log.Error("invalid user_id format", "user_id", userIDStr, "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}
		userID = &parsed
	}

	var serviceName *string
	if serviceNameStr != "" {
		serviceName = &serviceNameStr
	}

	subs, err := h.SubService.List(c.Request.Context(), userID, serviceName)
	if err != nil {
		log.Error("failed to list subscriptions", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscriptions"})
		return
	}

	log.Info("subscriptions listed", "count", len(subs))
	c.JSON(http.StatusOK, subs)
}

// UpdateSubscription godoc
// @Summary     Update a subscription
// @Description Update subscription fields by ID
// @Tags        subscriptions
// @Accept      json
// @Produce     json
// @Param       id           path string               true  "Subscription ID"
// @Param       subscription body dto.UpdateSubscriptionDTO true "Updated subscription data"
// @Success     200 {object} dto.SubscriptionDTO
// @Failure     400 {object} httpapi.ErrorResponse
// @Failure     500 {object} httpapi.ErrorResponse
// @Router      /subscriptions/{id} [put]
func (h *Handler) UpdateSubscription(c *gin.Context) {
	log := h.log.With("handler", "UpdateSubscription")
	idStr := c.Param("id")
	log.Debug("received update request", "id", idStr)

	subID, err := uuid.Parse(idStr)
	if err != nil {
		log.Error("invalid subscription ID format", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	var req dto.UpdateSubscriptionDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("invalid update payload", "error", err)
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fieldErr := ve[0]
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s failed %s validation", fieldErr.Field(), fieldErr.Tag())})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		}
		return
	}

	var startDate *time.Time
	if req.StartDate != nil {
		t, err2 := time.Parse("01-2006", *req.StartDate)
		if err2 != nil {
			log.Error("invalid start_date format", "start_date", *req.StartDate, "error", err2)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date"})
			return
		}
		startDate = &t
	}

	var endDate *time.Time
	if req.EndDate != nil {
		t, err2 := time.Parse("01-2006", *req.EndDate)
		if err2 != nil {
			log.Error("invalid end_date format", "end_date", *req.EndDate, "error", err2)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date"})
			return
		}
		endDate = &t
	}

	input := app.UpdateInput{
		ServiceName: req.ServiceName,
		StartDate:   startDate,
		EndDate:     endDate,
		Price:       req.Price,
	}

	sub, err := h.SubService.Update(c.Request.Context(), subID, input)
	if err != nil {
		log.Error("failed to update subscription", "id", subID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
		return
	}

	log.Info("subscription updated", "id", sub.ID)
	c.JSON(http.StatusOK, sub)
}

// DeleteSubscription godoc
// @Summary     Delete a subscription
// @Description Delete subscription by ID
// @Tags        subscriptions
// @Param       id   path   string true "Subscription ID"
// @Success     204 {object} nil
// @Failure     400 {object} httpapi.ErrorResponse
// @Failure     500 {object} httpapi.ErrorResponse
// @Router      /subscriptions/{id} [delete]
func (h *Handler) DeleteSubscription(c *gin.Context) {
	log := h.log.With("handler", "DeleteSubscription")
	idStr := c.Param("id")
	log.Debug("received delete request", "id", idStr)

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error("invalid subscription ID format", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	if err := h.SubService.Delete(c.Request.Context(), id); err != nil {
		log.Error("failed to delete subscription", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subscription"})
		return
	}

	log.Info("subscription deleted", "id", id)
	c.Status(http.StatusNoContent)
}

// AggregateSubscriptions godoc
// @Summary     Aggregate subscription costs
// @Description Calculate total cost over period with optional filters
// @Tags        subscriptions
// @Produce     json
// @Param       user_id      query string false "User ID"
// @Param       service_name query string false "Service Name"
// @Param       start_period query string true  "Start period (MM-YYYY)"
// @Param       end_period   query string true  "End period (MM-YYYY)"
// @Success     200 {object} httpapi.AggregateResponse
// @Failure     400 {object} httpapi.ErrorResponse
// @Failure     500 {object} httpapi.ErrorResponse
// @Router      /subscriptions/aggregate [get]
func (h *Handler) AggregateSubscriptions(c *gin.Context) {
	log := h.log.With("handler", "AggregateSubscriptions")

	userIDStr := c.Query("user_id")
	serviceNameStr := c.Query("service_name")
	startStr := c.Query("start_period")
	endStr := c.Query("end_period")
	log.Debug("received aggregate request", "user_id", userIDStr, "service_name", serviceNameStr, "start", startStr, "end", endStr)

	var userID *uuid.UUID
	if userIDStr != "" {
		parsed, err := uuid.Parse(userIDStr)
		if err != nil {
			log.Error("invalid user_id format", "user_id", userIDStr, "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}
		userID = &parsed
	}

	var serviceName *string
	if serviceNameStr != "" {
		serviceName = &serviceNameStr
	}

	start, err := time.Parse("01-2006", startStr)
	if err != nil {
		log.Error("invalid start_period format", "start_period", startStr, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_period"})
		return
	}

	end, err := time.Parse("01-2006", endStr)
	if err != nil {
		log.Error("invalid end_period format", "end_period", endStr, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_period"})
		return
	}

	filter := app.AggregationFilter{
		UserID:      userID,
		ServiceName: serviceName,
		StartPeriod: start,
		EndPeriod:   end,
	}

	sum, err := h.SubService.Aggregate(c.Request.Context(), filter)
	if err != nil {
		log.Error("aggregate calculation failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate aggregate"})
		return
	}

	log.Info("aggregate calculated", "total", sum)
	c.JSON(http.StatusOK, gin.H{"total": sum})
}
