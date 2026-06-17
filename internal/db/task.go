package db

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrTaskNotFound = errors.New("task not found")

type Task struct {
	ID      int64  `json:"id,string"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(t *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat)

	var id int64
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

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

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler
              SET date = ?, title = ?, comment = ?, repeat = ?
              WHERE id = ?`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
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
