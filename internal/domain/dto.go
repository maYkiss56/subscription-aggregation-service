package domain

import (
	"github.com/google/uuid"
	"github.com/maYkiss56/subscription-aggregation-service/internal/utils"
)

// CreateSubRequest represents request to create subscription
type CreateSubRequest struct {
	ServiceName string    `json:"service_name" example:"Netflix"`
	Price       int       `json:"price" example:"1000"`
	UserID      uuid.UUID `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	StartDate   string    `json:"start_date" example:"07-2025"`
	EndDate     string    `json:"end_date" example:"07-2025"`
}

// UpdateSubRequest represents request to update subscription
type UpdateSubRequest struct {
	ServiceName *string `json:"service_name,omitempty" example:"Netflix Premium"`
	Price       *int    `json:"price,omitempty" example:"1500"`
	StartDate   *string `json:"start_date,omitempty" example:"07-2025"`
	EndDate     *string `json:"end_date,omitempty" example:"07-2025"`
}

// TotalCostFilter represents filter for total cost calculation
type TotalCostFilter struct {
	UserID      *uuid.UUID `json:"user_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	ServiceName *string    `json:"service_name,omitempty" example:"Netflix"`
	StartPeriod string     `json:"start_period" example:"07-2025"`
	EndPeriod   string     `json:"end_period" example:"07-2025"`
}

// SubResponse представляет ответ с датами в формате MM-YYYY
type SubResponse struct {
	ID          uuid.UUID `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"` // MM-YYYY
	EndDate     string    `json:"end_date"`   // MM-YYYY
}

// convertSubToResponse преобразует доменную Sub в SubResponse
func ConvertSubToResponse(sub *Sub) *SubResponse {
	return &SubResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   utils.ToMonthYearString(sub.StartDate),
		EndDate:     utils.ToMonthYearString(sub.EndDate),
	}
}

// convertSubsToResponse преобразует список доменных Sub в список SubResponse
func ConvertSubsToResponse(subs []*Sub) []*SubResponse {
	result := make([]*SubResponse, len(subs))
	for i, sub := range subs {
		result[i] = ConvertSubToResponse(sub)
	}
	return result
}
