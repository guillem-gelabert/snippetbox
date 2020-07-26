package mysql

import (
	"reflect"
	"testing"
	"time"

	"github.com/guillem-gelabert/snippetbox/pkg/models"
)

func TestUserModelGet(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	testCases := []struct {
		desc          string
		userID        int
		expectedUser  *models.User
		expectedError error
	}{
		{
			desc:   "Valid ID",
			userID: 1,
			expectedUser: &models.User{
				ID:      1,
				Name:    "Alice Jones",
				Email:   "alice@example.com",
				Created: time.Date(2018, 12, 23, 17, 25, 22, 0, time.UTC),
				Active:  true,
			},
			expectedError: nil,
		},
		{
			desc:          "Zero ID",
			userID:        0,
			expectedUser:  nil,
			expectedError: models.ErrNoRecord,
		},
		{
			desc:          "Non existent ID",
			userID:        2,
			expectedUser:  nil,
			expectedError: models.ErrNoRecord,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := UserModel{db}

			user, err := m.Get(tC.userID)
			if err != tC.expectedError {
				t.Errorf("expected %v, got %s", tC.expectedError, err)
			}
			if !reflect.DeepEqual(user, tC.expectedUser) {
				t.Errorf("expected %v, got %v", tC.expectedUser, user)
			}
		})
	}
}
