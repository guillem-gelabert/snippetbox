package mock

import (
	"time"

	"github.com/guillem-gelabert/snippetbox/pkg/models"
)

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

// SnippetModel is a mock snippet model
type SnippetModel struct{}

// Insert mocks the Insert method of the Snippet Model
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	return 2, nil
}

// Get mocks the Get method of the Snippet Model
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

// Latest mocks the Latest method of the Snippet Model
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
