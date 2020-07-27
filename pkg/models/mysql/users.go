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
	var id int
	var hashedPassword []byte

	stmt := `
		SELECT id, hashed_password FROM users
		WHERE email = ? AND active = TRUE
	`

	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}

	if err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

// Get fetches details for a specific user based on their user ID.
func (m *UserModel) Get(id int) (*models.User, error) {
	u := &models.User{}

	stmt := `
		SELECT id, name, email, created, active FROM users WHERE id = ?
	`

	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.Created, &u.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}

		return nil, err
	}
	return u, nil
}

// ChangePassword updates user password
func (m *UserModel) ChangePassword(id int, currentPassword, newPassword string) error {
	u, err := m.Get(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNoRecord
		}

		return err
	}

	_, err = m.Authenticate(u.Email, currentPassword)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	stmt := `
		UPDATE users
		SET hashed_password = ?
		WHERE id = ?
	`

	_, err = m.DB.Exec(stmt, hashedPassword, u.ID)
	if err != nil {
		return err
	}

	return nil
}
