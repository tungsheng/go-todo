package storage

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "github.com/tursodatabase/go-libsql"

	"github.com/tungsheng/go-todo/internal/model"
)

type Storage struct {
	db *sql.DB
}

func New() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dbDir := filepath.Join(homeDir, ".go-todo")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dbDir, "todos.db")
	db, err := sql.Open("libsql", "file:"+dbPath)
	if err != nil {
		return nil, err
	}

	s := &Storage{db: db}
	if err := s.init(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Storage) init() error {
	// Check if we need to migrate from old DATETIME format
	if err := s.migrate(); err != nil {
		return err
	}

	// Create table with INTEGER timestamps for reliable filtering
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			title      TEXT NOT NULL,
			category   TEXT DEFAULT '',
			detail     TEXT DEFAULT '',
			status     TEXT DEFAULT 'pending',
			time_tag   TEXT DEFAULT '',
			due_date   INTEGER,
			created_at INTEGER,
			updated_at INTEGER
		)
	`)
	if err != nil {
		return err
	}

	// Add time_tag column if it doesn't exist (migration for existing tables)
	s.db.Exec(`ALTER TABLE todos ADD COLUMN time_tag TEXT DEFAULT ''`)

	return nil
}

func (s *Storage) migrate() error {
	// Check if old table exists with DATETIME columns
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM pragma_table_info('todos')
		WHERE name = 'created_at' AND type = 'DATETIME'
	`).Scan(&count)
	if err != nil || count == 0 {
		return nil // No migration needed
	}

	// Migrate: create new table, copy data, replace old table
	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS todos_new (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			title      TEXT NOT NULL,
			category   TEXT DEFAULT '',
			detail     TEXT DEFAULT '',
			status     TEXT DEFAULT 'pending',
			due_date   INTEGER,
			created_at INTEGER,
			updated_at INTEGER
		)
	`)
	if err != nil {
		return err
	}

	// Copy data with timestamp conversion
	_, err = s.db.Exec(`
		INSERT INTO todos_new (id, title, category, detail, status, due_date, created_at, updated_at)
		SELECT id, title, category, detail, status,
			CASE WHEN due_date IS NOT NULL THEN CAST(strftime('%s', due_date) AS INTEGER) ELSE NULL END,
			CAST(strftime('%s', created_at) AS INTEGER),
			CAST(strftime('%s', updated_at) AS INTEGER)
		FROM todos
	`)
	if err != nil {
		return err
	}

	// Drop old table and rename new one
	_, err = s.db.Exec(`DROP TABLE todos`)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`ALTER TABLE todos_new RENAME TO todos`)
	return err
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func scanTodos(rows *sql.Rows) ([]model.Todo, error) {
	var todos []model.Todo
	for rows.Next() {
		var t model.Todo
		var timeTag sql.NullString
		var dueDate, createdAt, updatedAt sql.NullInt64
		err := rows.Scan(&t.ID, &t.Title, &t.Category, &t.Detail, &t.Status, &timeTag, &dueDate, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}
		if timeTag.Valid {
			t.TimeTag = model.TimeTag(timeTag.String)
		}
		if dueDate.Valid {
			d := time.Unix(dueDate.Int64, 0)
			t.DueDate = &d
		}
		if createdAt.Valid {
			t.CreatedAt = time.Unix(createdAt.Int64, 0)
		}
		if updatedAt.Valid {
			t.UpdatedAt = time.Unix(updatedAt.Int64, 0)
		}
		todos = append(todos, t)
	}
	return todos, rows.Err()
}

func (s *Storage) Create(title string, timeTag model.TimeTag) (*model.Todo, error) {
	now := time.Now()
	result, err := s.db.Exec(`
		INSERT INTO todos (title, time_tag, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, title, string(timeTag), now.Unix(), now.Unix())
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &model.Todo{
		ID:        id,
		Title:     title,
		Status:    model.StatusPending,
		TimeTag:   timeTag,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (s *Storage) Update(todo *model.Todo) error {
	todo.UpdatedAt = time.Now()
	var dueDate *int64
	if todo.DueDate != nil {
		d := todo.DueDate.Unix()
		dueDate = &d
	}
	_, err := s.db.Exec(`
		UPDATE todos
		SET title = ?, category = ?, detail = ?, status = ?, time_tag = ?, due_date = ?, updated_at = ?
		WHERE id = ?
	`, todo.Title, todo.Category, todo.Detail, todo.Status, string(todo.TimeTag), dueDate, todo.UpdatedAt.Unix(), todo.ID)
	return err
}

func (s *Storage) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM todos WHERE id = ?`, id)
	return err
}

func (s *Storage) ListFiltered(timeFilter string) ([]model.Todo, error) {
	baseSelect := `SELECT id, title, category, detail, status, time_tag, due_date, created_at, updated_at FROM todos`
	orderBy := ` ORDER BY
		CASE status
			WHEN 'in_progress' THEN 1
			WHEN 'pending' THEN 2
			WHEN 'done' THEN 3
			WHEN 'closed' THEN 4
		END,
		created_at DESC`

	query := baseSelect
	if timeFilter != "" {
		query += " WHERE time_tag = ?"
	}
	query += orderBy

	var rows *sql.Rows
	var err error
	if timeFilter != "" {
		rows, err = s.db.Query(query, timeFilter)
	} else {
		rows, err = s.db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTodos(rows)
}
