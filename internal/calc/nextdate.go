package calc

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DateLayout represents dates in YYYYMMDD format
const DateLayout = "20060102"

// NextDate calculates next date for task, based on repeat rule
// now: current time, dstart: task creation date, repeat: repeat rule string.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	date, err := time.Parse(DateLayout, dstart)
	if err != nil {
		return "", fmt.Errorf("repeat rule format is wrong: invalid start date format")
	}

	if repeat == "" {
		return "", fmt.Errorf("repeat rule is empty")
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {

	// case "d": repeat every N days
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("repeat rule format is wrong")
		}

		var interval int
		interval, err = strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("repeat rule format is wrong")
		}

		if interval <= 0 || interval > 400 {
			return "", fmt.Errorf("interval must be positive and less than 400")
		}

		for {
			date = date.AddDate(0, 0, interval)
			if date.After(now) {
				break
			}
		}
		return date.Format(DateLayout), nil

	// case "y": repeat every year
	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("repeat rule format is wrong")
		}

		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				break
			}
		}
		return date.Format(DateLayout), nil

	// case "w": repeat on specific days of the week (e.g. "w 1,3,5" for Mon, Wed, Fri)
	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("repeat rule format is wrong")
		}

		days := strings.Split(parts[1], ",")
		if len(days) > 7 {
			return "", fmt.Errorf("the number of days exceeded")
		}

		var targetDays []int
		var intDay int

		// Convert user input (1-7, Mon-Sun) to Go's Weekday format (0-6, Sun-Sat).
		// intDay % 7 transforms 7 (Sun) to 0, and 1-6 remain 1-6.
		for _, day := range days {
			intDay, err = strconv.Atoi(day)
			if err != nil {
				return "", fmt.Errorf("repeat rule format is wrong")
			}
			if intDay < 1 || intDay > 7 {
				return "", fmt.Errorf("weekday must be between 1 and 7, got %d", intDay)
			}
			goDay := intDay % 7
			targetDays = append(targetDays, goDay)
		}

		// Calculate the minimum difference from 'now', not from 'dstart'.
		currentNow := int(now.Weekday())
		minDiff := 7
		for _, targetDay := range targetDays {
			diff := targetDay - currentNow
			if diff <= 0 {
				diff += 7
			}
			if diff < minDiff {
				minDiff = diff
			}
		}

		// Add the calculated difference to 'now' to get the next valid date.
		date = now.AddDate(0, 0, minDiff)
		return date.Format(DateLayout), nil

	// case "m": repeat on specific days of the month (e.g., "m 15", "m -1", "m 15 6,12")
	case "m":

		var targetDays []int
		var targetMonths []int

		if len(parts) < 2 || len(parts) > 3 {
			return "", fmt.Errorf("repeat rule format is wrong")
		}

		// Parse days (always present in parts[1]).
		days := strings.Split(parts[1], ",")
		for _, value := range days {
			var day int
			day, err = strconv.Atoi(value)
			if err != nil {
				return "", fmt.Errorf("repeat rule format is wrong")
			}
			if day >= -2 && day <= 31 && day != 0 {
				targetDays = append(targetDays, day)
			} else {
				return "", fmt.Errorf("repeat rule format is wrong, want from -2 to 31, got %d", day)
			}
		}

		// Parse months (optional, present only if len(parts) == 3).
		if len(parts) == 3 {
			months := strings.Split(parts[2], ",")

			for _, value := range months {
				var month int
				month, err = strconv.Atoi(value)
				if err != nil {
					return "", fmt.Errorf("repeat rule format is wrong")
				}

				if month >= 1 && month <= 12 {
					targetMonths = append(targetMonths, month)
				} else {
					return "", fmt.Errorf("repeat rule format is wrong, want from 1 to 12, got %d", month)
				}
			}
		}

		currentDate := date

		for {
			// Step 1: Check if the current month is allowed by the rule.
			var monthIsValid bool
			if len(parts) == 3 {
				for _, monthRule := range targetMonths {
					if int(currentDate.Month()) == monthRule {
						monthIsValid = true
						break
					}
				}
			} else {
				// If no months are specified in the rule, any month is valid.
				// This prevents an infinite loop where the algorithm skips every month.
				monthIsValid = true
			}

			// If the month is not allowed, advance to the 1st of the next month and retry.
			if !monthIsValid {
				currentDate = time.Date(currentDate.Year(), currentDate.Month()+1, 1, 0, 0, 0, 0, currentDate.Location())
				continue
			}

			// Step 2: The month is valid. Check all specified days in this month.
			var candidate time.Time
			var bestCandidate time.Time

			for _, dayRule := range targetDays {
				actualDay := dayRule
				lastDay := lastDayOfMonth(currentDate)

				// Handle negative days (-1 = last day, -2 = second to last)
				if dayRule < 0 {
					actualDay = lastDay + dayRule + 1
				}

				// Safety check: ensure the day exists in this month
				if actualDay > lastDay {
					continue
				}

				candidate = time.Date(currentDate.Year(), currentDate.Month(), actualDay, 0, 0, 0, 0, currentDate.Location())
				if candidate.After(now) {
					if bestCandidate.IsZero() || candidate.Before(bestCandidate) {
						bestCandidate = candidate
					}
				}
			}

			// Step 3: If a valid date was found in this month, return it immediately.
			if !bestCandidate.IsZero() {
				return bestCandidate.Format(DateLayout), nil
			}

			// Step 4: If no valid days were found in this month (all are <= now),
			// advance to the 1st day of the next month and continue the loop.
			currentDate = time.Date(currentDate.Year(), currentDate.Month()+1, 1, 0, 0, 0, 0, currentDate.Location())

		}

	default:
		return "", fmt.Errorf("repeat rule format is wrong")
	}
}

// lastDayOfMonth returns the last day of the month for a given date.
func lastDayOfMonth(t time.Time) int {
	firstOfNextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	return firstOfNextMonth.AddDate(0, 0, -1).Day()
}
