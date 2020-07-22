package mysql

import (
	"database/sql"
	"errors"

	"github.com/guillem-gelabert/snippetbox/pkg/models"
)

// SnippetModel wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert inserts a new snippet into the database.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `
		INSERT INTO snippets (title, content, created, expires)
		VALUES(?,?,UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(),INTERVAL ? DAY))
	`
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get returns a specific snipped based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `
		SELECT id, title, content, created, expires FROM snippets
		WHERE expires > UTC_TIMESTAMP() AND id = ?	
	`

	s := &models.Snippet{}
	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return s, nil
}

// Latest returns the 10 most recently  created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
