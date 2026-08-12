// Package nextdate provides logic for calculating the next execution date
// for repeating tasks based on custom recurrence rules.
package nextdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DateFormat is the standard date format (YYYYMMDD) used across the application.
const DateFormat = "20060102"

// NextDate calculates the next date for a recurring task based on the current time (now),
// the task's start date (dstart, format YYYYMMDD), and the repetition rule (repeat).
//
// Supported repeat rules:
//   - "y": Yearly repetition (e.g. annually on the same month and day).
//   - "d <N>": Repetition every N days, where 1 <= N <= 400 (e.g. "d 5").
//   - "w <days>": Repetition on specific days of the week, where days are 1 (Monday) to 7 (Sunday),
//     comma-separated (e.g. "w 1,3,5").
//   - "m <days> [months]": Repetition on specific days of the month (1..31, -1 for last day,
//     -2 for second-to-last day) and optional specific months (1..12) (e.g. "m 1,15" or "m -1 1,6,12").
//
// Returns the next date in YYYYMMDD format, or an error if any parameter or rule format is invalid.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("empty repeat rule")
	}

	// Normalize 'now' to 00:00:00 UTC for consistent date-only comparison.
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	startDate, err := time.ParseInLocation(DateFormat, dstart, time.UTC)
	if err != nil {
		return "", fmt.Errorf("invalid start date format: %w", err)
	}

	parts := strings.Fields(repeat)

	switch parts[0] {
	case "y":
		return handleYear(startDate, now)
	case "d":
		return handleDays(startDate, now, parts)
	case "w":
		return handleWeeks(startDate, now, parts)
	case "m":
		return handleMonths(startDate, now, parts)
	default:
		return "", fmt.Errorf("unsupported rule format: %s", parts[0])
	}
}

// handleYear calculates the next date for yearly repetition ("y").
// Increments the date by full years until it is strictly strictly strictly strictly after 'now'.
func handleYear(startDate, now time.Time) (string, error) {
	nextDate := startDate
	for {
		nextDate = nextDate.AddDate(1, 0, 0)
		if nextDate.After(now) {
			break
		}
	}
	return nextDate.Format(DateFormat), nil
}

// handleDays calculates the next date for daily repetition ("d <N>").
// N must be an integer between 1 and 400.
func handleDays(startDate, now time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid format for rule 'd'")
	}
	days, err := strconv.Atoi(parts[1])
	if err != nil || days < 1 || days > 400 {
		return "", fmt.Errorf("invalid number of days (expected 1 to 400)")
	}

	nextDate := startDate
	for {
		nextDate = nextDate.AddDate(0, 0, days)
		if nextDate.After(now) {
			break
		}
	}
	return nextDate.Format(DateFormat), nil
}

// handleWeeks calculates the next date for weekly repetition ("w <days>").
// Expects comma-separated weekday numbers (1=Monday ... 7=Sunday).
func handleWeeks(startDate, now time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid format for rule 'w'")
	}

	// Index 0 corresponds to Sunday (7 % 7 = 0), 1..6 to Mon..Sat.
	validWeekday := make([]bool, 7)
	for _, dStr := range strings.Split(parts[1], ",") {
		d, err := strconv.Atoi(dStr)
		if err != nil || d < 1 || d > 7 {
			return "", fmt.Errorf("invalid weekday: %s", dStr)
		}
		validWeekday[d%7] = true
	}

	nextDate := startDate
	if now.After(nextDate) || now.Equal(nextDate) {
		nextDate = now.AddDate(0, 0, 1)
	}

	for {
		if validWeekday[nextDate.Weekday()] {
			return nextDate.Format(DateFormat), nil
		}
		nextDate = nextDate.AddDate(0, 0, 1)
	}
}

// handleMonths calculates the next date for monthly repetition ("m <days> [months]").
// Supports specific days (1..31), relative last days (-1, -2), and optional specific months (1..12).
func handleMonths(startDate, now time.Time, parts []string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", fmt.Errorf("invalid format for rule 'm'")
	}

	validDay := make([]bool, 32)
	var validLastDay, validSecondLastDay bool

	for _, dStr := range strings.Split(parts[1], ",") {
		d, err := strconv.Atoi(dStr)
		if err != nil || d < -2 || d == 0 || d > 31 {
			return "", fmt.Errorf("invalid month day: %s", dStr)
		}
		switch d {
		case -1:
			validLastDay = true
		case -2:
			validSecondLastDay = true
		default:
			if d > 0 {
				validDay[d] = true
			}
		}
	}

	var validMonth []bool
	if len(parts) == 3 {
		validMonth = make([]bool, 13)
		for _, mStr := range strings.Split(parts[2], ",") {
			m, err := strconv.Atoi(mStr)
			if err != nil || m < 1 || m > 12 {
				return "", fmt.Errorf("invalid month: %s", mStr)
			}
			validMonth[m] = true
		}
	}

	nextDate := startDate
	if now.After(nextDate) || now.Equal(nextDate) {
		nextDate = now.AddDate(0, 0, 1)
	}

	for {
		if validMonth == nil || validMonth[nextDate.Month()] {

			// Calculate the last day of the current month.
			lastDay := time.Date(nextDate.Year(), nextDate.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
			day := nextDate.Day()

			if validDay[day] || (validLastDay && day == lastDay) || (validSecondLastDay && day == lastDay-1) {
				return nextDate.Format(DateFormat), nil
			}
		}
		nextDate = nextDate.AddDate(0, 0, 1)
	}
}
