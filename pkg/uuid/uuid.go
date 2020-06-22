package uuid

import (
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

type UUID interface {
	Generate() string
}

// NewUUID
func NewUUID() UUID {
	return &generator{}
}

// generator is a UUID implementation.
type generator struct{}

// Generate returns a UUID in string format.
func (*generator) Generate() string {
	return uuid.NewV4().String()
}

// MockUUID is a UUID implementation for testing purposes.
type MockUUID struct {
	*mock.Mock
}

// Generate returns a UUID in string format.
func (m *MockUUID) Generate() string {
	return m.Called().String(0)
}

// NewTestUUID creates a Mock to test the UUID generator.
// It should only be used for testing.
func NewTestUUID() *MockUUID {
	return &MockUUID{
		Mock: new(mock.Mock),
	}
}
