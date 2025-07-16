package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Workflow represents a stored workflow definition.
type Workflow struct {
	ID         string
	Name       string
	Definition string // YAML workflow definition
	CreatedAt  time.Time
}

// WorkflowStore defines the interface for storing and retrieving workflows.
type WorkflowStore interface {
	Init() error
	SaveWorkflow(workflow *Workflow) error
	GetWorkflow(id string) (*Workflow, error)
	GetWorkflowByName(name string) (*Workflow, error)
	ListWorkflows() ([]*Workflow, error)
}

// SQLiteStore implements WorkflowStore for SQLite.
type SQLiteStore struct {
	db *sql.DB
	dbPath string
}

// NewSQLiteStore creates a new SQLiteStore.
func NewSQLiteStore(dbPath string) *SQLiteStore {
	return &SQLiteStore{dbPath: dbPath}
}

// Init initializes the SQLite database and creates the workflows table.
func (s *SQLiteStore) Init() error {
	var err error
	s.db, err = sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS workflows (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		definition TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = s.db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create workflows table: %w", err)
	}

	return nil
}

// SaveWorkflow saves a workflow definition to the database.
func (s *SQLiteStore) SaveWorkflow(workflow *Workflow) error {
	stmt, err := s.db.Prepare("INSERT INTO workflows(id, name, definition) VALUES(?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(workflow.ID, workflow.Name, workflow.Definition)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

// GetWorkflow retrieves a workflow definition by ID.
func (s *SQLiteStore) GetWorkflow(id string) (*Workflow, error) {
	row := s.db.QueryRow("SELECT id, name, definition, created_at FROM workflows WHERE id = ?", id)

	workflow := &Workflow{}
	err := row.Scan(&workflow.ID, &workflow.Name, &workflow.Definition, &workflow.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to scan workflow: %w", err)
	}

	return workflow, nil
}

// GetWorkflowByName retrieves a workflow definition by name.
func (s *SQLiteStore) GetWorkflowByName(name string) (*Workflow, error) {
	row := s.db.QueryRow("SELECT id, name, definition, created_at FROM workflows WHERE name = ?", name)

	workflow := &Workflow{}
	err := row.Scan(&workflow.ID, &workflow.Name, &workflow.Definition, &workflow.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to scan workflow: %w", err)
	}

	return workflow, nil
}

// ListWorkflows lists all stored workflow definitions.
func (s *SQLiteStore) ListWorkflows() ([]*Workflow, error) {
	rows, err := s.db.Query("SELECT id, name, definition, created_at FROM workflows")
	if err != nil {
		return nil, fmt.Errorf("failed to query workflows: %w", err)
	}
	defer rows.Close()

	var workflows []*Workflow
	for rows.Next() {
		workflow := &Workflow{}
		err := rows.Scan(&workflow.ID, &workflow.Name, &workflow.Definition, &workflow.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workflow row: %w", err)
		}
		workflows = append(workflows, workflow)
	}

	return workflows, nil
}
