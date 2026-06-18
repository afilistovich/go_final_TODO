package api

import (
	"fmt"
	"time"

	"github.com/afilistovich/go_final_TODO/internal/calc"
	"github.com/afilistovich/go_final_TODO/internal/db"
)

func normalizeTaskDate(task *db.Task) error {
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var next string

	if task.Date == "" {
		task.Date = now.Format(calc.DateLayout)
		return nil
	}

	t, err := time.Parse(calc.DateLayout, task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYYMMDD")
	}

	if t.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format(calc.DateLayout)
		} else {
			next, err = calc.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("invalid repeat rule: %w", err)
			}
			task.Date = next
		}
	}
	return nil
}
