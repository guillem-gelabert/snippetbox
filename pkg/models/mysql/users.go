package mysql

import (
	"database/sql"

	"github.com/guillem-gelabert/snippetbox/pkg/models"
)

// UserModel wraps a sql.DB connection pool
type UserModel struct {
	DB *sql.DB
}

// Insert adds a new record to the users table.
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// Authenticate verifies whether a user exists with the provided email and password.
// Returns the relevant user ID if it exists.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Get fetches details for a specific user based on their user ID.
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
