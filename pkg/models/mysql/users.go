package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/guillem-gelabert/snippetbox/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// UserModel wraps a sql.DB connection pool
type UserModel struct {
	DB *sql.DB
}

// Insert adds a new record to the users table.
func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `
		INSERT INTO users (name, email, hashed_password, created)
		VALUES(?,?,?,UTC_TIMESTAMP())
	`

	_, err = m.DB.Exec(stmt, name, email, hashedPassword)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}

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
