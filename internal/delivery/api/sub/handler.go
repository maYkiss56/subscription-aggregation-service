package sub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/maYkiss56/subscription-aggregation-service/internal/domain"
	"github.com/maYkiss56/subscription-aggregation-service/internal/utils"
)

const (
	ErrInvalidBody      = "invalid request body"
	ErrInvalidSubData   = "invalid subscription data"
	ErrInternalServer   = "internal server error"
	ErrInvalidUserID    = "invalid user id"
	ErrInvalidSubID     = "invalid subscription id"
	ErrInvalidDateRange = "invalid date range"
)

type SubService interface {
	CreateSub(ctx context.Context, sub *domain.Sub) (id uuid.UUID, err error)
	GetAllSubs(ctx context.Context) ([]*domain.Sub, error)
	GetSubByUserID(ctx context.Context, userUID uuid.UUID) ([]*domain.Sub, error)
	UpdateSub(ctx context.Context, id uuid.UUID, req *domain.UpdateSubRequest) (*domain.Sub, error)
	DeleteSub(ctx context.Context, id uuid.UUID) error
	CalculateTotalCost(ctx context.Context, filter domain.TotalCostFilter) (int, error)
}

type HandlerSub struct {
	service SubService
}

func New(service SubService) *HandlerSub {
	return &HandlerSub{
		service: service,
	}
}

// CreateSub godoc
// @Summary Create a new subscription
// @Description Create a new subscription with the input payload
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Param input body domain.CreateSubRequest true "Create subscription"
// @Success 201 {object} map[string]interface{} "Subscription created"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /create [post]
func (h *HandlerSub) CreateSub(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateSubRequest
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", ErrInvalidBody, err), http.StatusBadRequest)
		return
	}

	// Парсим start_date как первый день месяца
	startDate, err := utils.ParseMonthYear(req.StartDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid start date: %v", err), http.StatusBadRequest)
		return
	}

	// Парсим end_date как последний день месяца
	endDate, err := utils.ParseMonthYearToEndOfMonth(req.EndDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid end date: %v", err), http.StatusBadRequest)
		return
	}

	newSub, err := domain.New(req.ServiceName, req.Price, req.UserID, startDate, endDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", ErrInvalidSubData, err), http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateSub(r.Context(), newSub)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", ErrInternalServer, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "subscription created successfully",
		"id":      id,
	}); err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", ErrInternalServer, err), http.StatusInternalServerError)
	}
}

// GetAllSubs godoc
// @Summary Get all subscriptions
// @Description Get list of all subscriptions
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Success 200 {array} domain.Sub "List of subscriptions"
// @Failure 500 {string} string "Internal server error"
// @Router / [get]
func (h *HandlerSub) GetAllSubs(w http.ResponseWriter, r *http.Request) {
	subs, err := h.service.GetAllSubs(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get subscriptions: %v", err), http.StatusInternalServerError)
		return
	}

	response := domain.ConvertSubsToResponse(subs)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", ErrInternalServer, err), http.StatusInternalServerError)
	}
}

// GetSubByUserID godoc
// @Summary Get subscriptions by user ID
// @Description Get list of subscriptions for specific user
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Param user_id path string true "User ID"
// @Success 200 {array} domain.Sub "List of user subscriptions"
// @Failure 400 {string} string "Invalid user ID"
// @Failure 500 {string} string "Internal server error"
// @Router /{user_id} [get]
func (h *HandlerSub) GetSubByUserID(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", ErrInvalidUserID, err), http.StatusBadRequest)
		return
	}

	subs, err := h.service.GetSubByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get user subscriptions: %v", err), http.StatusInternalServerError)
		return
	}

	response := domain.ConvertSubsToResponse(subs)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", ErrInternalServer, err), http.StatusInternalServerError)
	}
}

// UpdateSub godoc
// @Summary Update subscription
// @Description Update existing subscription
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Param id path string true "Subscription ID"
// @Param input body domain.UpdateSubRequest true "Update data"
// @Success 200 {object} domain.Sub "Updated subscription"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Subscription not found"
// @Failure 500 {string} string "Internal server error"
// @Router /update/{id} [patch]
func (h *HandlerSub) UpdateSub(w http.ResponseWriter, r *http.Request) {
	subIDStr := chi.URLParam(r, "id")
	subID, err := uuid.Parse(subIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", ErrInvalidSubID, err), http.StatusBadRequest)
		return
	}

	var req domain.UpdateSubRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", ErrInvalidBody, err), http.StatusBadRequest)
		return
	}

	// Преобразуем даты перед передачей в сервис
	if req.StartDate != nil {
		startDate, err := utils.ParseMonthYear(*req.StartDate)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid start date: %v", err), http.StatusBadRequest)
			return
		}
		*req.StartDate = startDate.Format("2006-01-02") // Преобразуем в YYYY-MM-DD
	}

	if req.EndDate != nil {
		endDate, err := utils.ParseMonthYearToEndOfMonth(*req.EndDate)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid end date: %v", err), http.StatusBadRequest)
			return
		}
		*req.EndDate = endDate.Format("2006-01-02") // Преобразуем в YYYY-MM-DD
	}

	updatedSub, err := h.service.UpdateSub(r.Context(), subID, &req)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to update subscription: %v", err), http.StatusInternalServerError)
		return
	}

	response := domain.ConvertSubToResponse(updatedSub)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", ErrInternalServer, err), http.StatusInternalServerError)
	}
}

// DeleteSub godoc
// @Summary Delete subscription
// @Description Delete existing subscription
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Param id path string true "Subscription ID"
// @Success 204 "No content"
// @Failure 400 {string} string "Invalid subscription ID"
// @Failure 404 {string} string "Subscription not found"
// @Failure 500 {string} string "Internal server error"
// @Router /delete/{id} [delete]
func (h *HandlerSub) DeleteSub(w http.ResponseWriter, r *http.Request) {
	subIDStr := chi.URLParam(r, "id")
	subID, err := uuid.Parse(subIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", ErrInvalidSubID, err), http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteSub(r.Context(), subID); err != nil {
		http.Error(w, fmt.Sprintf("failed to delete subscription: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CalculateTotalCost godoc
// @Summary Calculate total cost
// @Description Calculate total cost of subscriptions for given period
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Param input body domain.TotalCostFilter true "Filter parameters"
// @Success 200 {object} map[string]int "Total cost"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /total [post]
func (h *HandlerSub) CalculateTotalCost(w http.ResponseWriter, r *http.Request) {
	var filter domain.TotalCostFilter
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", ErrInvalidBody, err), http.StatusBadRequest)
		return
	}

	// Преобразуем даты в формат, понятный БД
	if filter.StartPeriod != "" {
		startDate, err := utils.ParseMonthYear(filter.StartPeriod)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid start period: %v", err), http.StatusBadRequest)
			return
		}
		filter.StartPeriod = startDate.Format("2006-01-02") // Преобразуем в YYYY-MM-DD
	}

	if filter.EndPeriod != "" {
		endDate, err := utils.ParseMonthYearToEndOfMonth(filter.EndPeriod)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid end period: %v", err), http.StatusBadRequest)
			return
		}
		filter.EndPeriod = endDate.Format("2006-01-02") // Преобразуем в YYYY-MM-DD
	}

	total, err := h.service.CalculateTotalCost(r.Context(), filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to calculate total cost: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]int{
		"total_cost": total,
	}); err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", ErrInternalServer, err), http.StatusInternalServerError)
	}
}
