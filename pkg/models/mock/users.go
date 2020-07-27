package mock

import (
	"time"

	"github.com/guillem-gelabert/snippetbox/pkg/models"
)

var mockUser = &models.User{
	ID:      1,
	Name:    "Alice",
	Email:   "alice@example.com",
	Created: time.Now(),
	Active:  true,
}

// UserModel is a mock snippet model
type UserModel struct{}

// Insert mocks the Insert method of the User Model
func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

// Authenticate mocks the Authenticate method of the User Model
func (m *UserModel) Authenticate(email, password string) (int, error) {
	switch email {
	case "alice@example.com":
		return 1, nil
	default:
		return 0, models.ErrInvalidCredentials
	}
}

// Get mocks the Get method of the User Model
func (m *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}

// ChangePassword mocks the ChangePassword method of the User Model
func (m *UserModel) ChangePassword(id int, currentPassword, newPassword string) error {
	if id != 1 {
		return models.ErrNoRecord
	}

	if currentPassword != "123123123123" {
		return models.ErrInvalidCredentials
	}

	return nil
}
