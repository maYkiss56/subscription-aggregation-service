package utils

import (
	"errors"
	"time"
)

const (
	monthYearLayout = "01-2006"
)

func ParseMonthYear(dateStr string) (time.Time, error) {
	if len(dateStr) != 7 || dateStr[2] != '-' {
		return time.Time{}, errors.New("invalid date format, expected MM-YYYY")
	}

	// Парсим дату как первый день месяца
	parsed, err := time.Parse(monthYearLayout, dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return parsed, nil
}

func ParseMonthYearToEndOfMonth(dateStr string) (time.Time, error) {
	parsed, err := ParseMonthYear(dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return parsed.AddDate(0, 1, -1), nil
}

func ToMonthYearString(date time.Time) string {
	return date.Format(monthYearLayout)
}
