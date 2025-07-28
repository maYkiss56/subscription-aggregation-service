package sub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/maYkiss56/subscription-aggregation-service/internal/domain"
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

func (h *HandlerSub) CreateSub(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateSubRequest
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", ErrInvalidBody, err), http.StatusBadRequest)
		return
	}

	newSub, err := domain.New(req.ServiceName, req.Price, req.UserID, req.StartDate, req.EndDate)
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

func (h *HandlerSub) GetAllSubs(w http.ResponseWriter, r *http.Request) {
	subs, err := h.service.GetAllSubs(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get subscriptions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(subs); err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", ErrInternalServer, err), http.StatusInternalServerError)
	}
}

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(subs); err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", ErrInternalServer, err), http.StatusInternalServerError)
	}
}

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

	updatedSub, err := h.service.UpdateSub(r.Context(), subID, &req)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to update subscription: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedSub); err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", ErrInternalServer, err), http.StatusInternalServerError)
	}
}

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

func (h *HandlerSub) CalculateTotalCost(w http.ResponseWriter, r *http.Request) {
	var filter domain.TotalCostFilter
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", ErrInvalidBody, err), http.StatusBadRequest)
		return
	}

	if filter.StartPeriod.IsZero() || filter.EndPeriod.IsZero() {
		http.Error(w, ErrInvalidDateRange, http.StatusBadRequest)
		return
	}

	if filter.EndPeriod.Before(filter.StartPeriod) {
		http.Error(w, "end period must be after start period", http.StatusBadRequest)
		return
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
