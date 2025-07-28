package domain

import (
	"time"

	"github.com/google/uuid"
)

type Sub struct {
	ID          uuid.UUID
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
}

func New(serviceName string, price int, userID uuid.UUID, startDate time.Time, endDate *time.Time) (*Sub, error) {
	return &Sub{
		ID:          uuid.New(),
		ServiceName: serviceName,
		Price:       price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}

func (s *Sub) MonthYear() string {
	return s.StartDate.Format("01-2006")
}
