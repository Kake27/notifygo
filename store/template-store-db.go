package store

import (
	"database/sql"
	"time"
)

type DBTemplate struct {
	ID string
	Name string
	Content string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DBTemplateStore struct {
	db *sql.DB
}

func NewDBTemplateStore(db *sql.DB) *DBTemplateStore {
	return &DBTemplateStore{db: db}
}

// retrieve template by name

func (s *DBTemplateStore) GetByName(name string) (*DBTemplate, error) {
	query := `SELECT id, name, content, created_at, updated_at FROM templates WHERE name = $1`

	row := s.db.QueryRow(query, name)

	var tmpl DBTemplate
	err := row.Scan(&tmpl.ID, &tmpl.Name, &tmpl.Content, &tmpl.CreatedAt, &tmpl.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &tmpl, nil
}


// create and delete
func (s *DBTemplateStore) Create(name, content string) error {
	query := `INSERT INTO templates (name, content) VALUES ($1, $2)`
	_, err := s.db.Exec(query, name, content)
	return err
}

func (s *DBTemplateStore) Delete(name string) error {
	query := `DELETE FROM templates WHERE name = $1`
	_, err := s.db.Exec(query, name)
	return err
}