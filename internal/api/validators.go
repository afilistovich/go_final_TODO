package api

import (
	"fmt"
	"time"

	"github.com/afilistovich/go_final_TODO/internal/calc"
	"github.com/afilistovich/go_final_TODO/internal/db"
)

// normalizeTaskDate normalizes task date: sets today if empty or in the past
func normalizeTaskDate(task *db.Task) error {
	// Get today's date at midnight
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var next string

	// If date is empty, use today
	if task.Date == "" {
		task.Date = now.Format(calc.DateLayout)
		return nil
	}

	t, err := time.Parse(calc.DateLayout, task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYYMMDD")
	}

	// If date is in the past, adjust it
	if t.Before(now) {
		if task.Repeat == "" {
			// No repeat rule - set to today
			task.Date = now.Format(calc.DateLayout)
		} else {
			// Has repeat rule - calculate next valid date
			next, err = calc.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("invalid repeat rule: %w", err)
			}
			task.Date = next
		}
	}
	return nil
}
