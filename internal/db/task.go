package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/afilistovich/go_final_TODO/internal/calc"
)

// ErrTaskNotFound returned when task with given ID doesn't exist
var ErrTaskNotFound = errors.New("task not found")

// Task represents a single task in the scheduler
type Task struct {
	ID      int64  `json:"id,string"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask inserts new task into database and returns its ID
func AddTask(t *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return 0, fmt.Errorf("add task: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id: %w", err)
	}

	return id, nil
}

// scanTasks scans multiple rows into Task slice
func scanTasks(rows *sql.Rows) ([]*Task, error) {
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		t := &Task{}
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return tasks, nil
}

// Tasks returns all tasks sorted by date with limit
func Tasks(limit int) ([]*Task, error) {
	if limit <= 0 {
		limit = 50
	}

	query := `SELECT * FROM scheduler 
              ORDER BY date ASC 
              LIMIT ?`

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}

	return scanTasks(rows)
}

// TasksBySearch returns tasks matching search pattern in title or comment
func TasksBySearch(search string, limit int) ([]*Task, error) {
	if limit <= 0 {
		limit = 50
	}

	query := `SELECT * FROM scheduler
              WHERE title LIKE ? OR comment LIKE ?
              ORDER BY date ASC
              LIMIT ?`

	searchPattern := "%" + search + "%"

	rows, err := db.Query(query, searchPattern, searchPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("query tasks by search: %w", err)
	}

	return scanTasks(rows)
}

// TasksByDate returns tasks for specific date
func TasksByDate(date string, limit int) ([]*Task, error) {
	if limit <= 0 {
		limit = 50
	}

	query := `SELECT * FROM scheduler
              WHERE date = ?
              ORDER BY date ASC
              LIMIT ?`

	rows, err := db.Query(query, date, limit)
	if err != nil {
		return nil, fmt.Errorf("query tasks by date: %w", err)
	}

	return scanTasks(rows)
}

// GetTask returns single task by ID
func GetTask(id int64) (*Task, error) {

	query := `SELECT * FROM scheduler
              WHERE id = ?`

	row := db.QueryRow(query, id)
	var t Task

	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("query task: %w", err)
	}

	return &t, nil
}

// UpdateTask updates existing task fields
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler
              SET date = ?, title = ?, comment = ?, repeat = ?
              WHERE id = ?`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	// Check if task was actually updated
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if count == 0 {
		return ErrTaskNotFound
	}
	return nil
}

// DeleteTask removes task by ID
func DeleteTask(id int64) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	res, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	// Check if task was actually deleted
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if count == 0 {
		return ErrTaskNotFound
	}
	return nil
}

// UpdateDate changes task date
func UpdateDate(id int64, newDate string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := db.Exec(query, newDate, id)
	if err != nil {
		return fmt.Errorf("update date: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if count == 0 {
		return ErrTaskNotFound
	}
	return nil
}

// DoneTask marks task as done: deletes if no repeat, or updates date to next occurrence
func DoneTask(id int64) error {
	task, err := GetTask(id)
	if err != nil {
		return err
	}

	// No repeat rule - delete task
	if task.Repeat == "" {
		return DeleteTask(id)
	}

	// Has repeat rule - calculate and set next date
	now := time.Now()
	nextDate, err := calc.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return fmt.Errorf("calculate next date: %w", err)
	}
	return UpdateDate(id, nextDate)
}
