package sqlite

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"go-final-progect/internal/models"

	_ "modernc.org/sqlite"
)

const (
	searchDateFormat = "02.01.2006"
	dbDateFormat     = "20060102"
)

type Storage struct {
	db *sql.DB
}

func New(dbFile string) (*Storage, error) {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) AddTask(t models.Task) (int64, error) {
	res, err := s.db.Exec(insertTaskQuery,
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat))
	if err != nil {
		return 0, fmt.Errorf("failed to insert task: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return id, nil
}

func (s *Storage) GetTasks(search string, limit int) ([]models.Task, error) {
	var rows *sql.Rows
	var err error

	switch {
	case search != "":
		parsedDate, parseErr := time.Parse(searchDateFormat, search)
		if parseErr == nil {
			dbDate := parsedDate.Format(dbDateFormat)
			rows, err = s.db.Query(selectTasksByDateQuery,
				sql.Named("date", dbDate),
				sql.Named("limit", limit))
		} else {
			pattern := "%" + search + "%"
			rows, err = s.db.Query(selectTasksBySearchQuery,
				sql.Named("title_like", pattern),
				sql.Named("comment_like", pattern),
				sql.Named("limit", limit))
		}
	default:
		rows, err = s.db.Query(selectTasksQuery,
			sql.Named("limit", limit))

	}

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	defer rows.Close()

	tasks := make([]models.Task, 0)

	for rows.Next() {
		var t models.Task
		var idInt int64

		err = rows.Scan(&idInt, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		t.ID = strconv.FormatInt(idInt, 10)
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return tasks, nil

}

func (s *Storage) GetTaskById(id int) (models.Task, error) {
	var t models.Task
	var idInt int64

	err := s.db.QueryRow(selectTaskByIDQuery,
		sql.Named("id", id)).Scan(&idInt, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return t, fmt.Errorf("task not found")
		}
		return t, fmt.Errorf("failed to fetch task by ID: %w", err)

	}

	t.ID = strconv.FormatInt(idInt, 10)
	return t, nil
}

func (s *Storage) UpdateTask(task models.Task) error {
	res, err := s.db.Exec(updateTaskQuery,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (s *Storage) DeleteTask(id int) error {
	res, err := s.db.Exec(deleteTaskQuery, sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
