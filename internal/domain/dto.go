package domain

import (
	"time"

	"github.com/google/uuid"
)

type CreateSubRequest struct {
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type CreateSubResponse struct {
	ID uuid.UUID `json:"id"`
}

type UpdateSubRequest struct {
	ServiceName *string    `json:"service_name"`
	Price       *int       `json:"price"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type TotalCostFilter struct {
	UserID      *uuid.UUID `json:"user_id"`
	ServiceName *string    `json:"service_name"`
	StartPeriod time.Time  `json:"start_period"`
	EndPeriod   time.Time  `json:"end_period"`
}
