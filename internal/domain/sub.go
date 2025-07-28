package domain

import (
	"time"

	"github.com/google/uuid"
)

// Sub represents subscription model
type Sub struct {
	ID          uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ServiceName string    `json:"service_name" example:"Netflix"`
	Price       int       `json:"price" example:"1000"`
	UserID      uuid.UUID `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	StartDate   time.Time `json:"start_date" example:"01-2023"`
	EndDate     time.Time `json:"end_date" example:"01-2023"`
}

func New(serviceName string, price int, userID uuid.UUID, startDate time.Time, endDate time.Time) (*Sub, error) {
	return &Sub{
		ID:          uuid.New(),
		ServiceName: serviceName,
		Price:       price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}
