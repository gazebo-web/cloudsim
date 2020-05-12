package uuid

import (
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

type UUID interface {
	Generate() string
}

func NewUUID() UUID {
	return &generator{}
}

type generator struct {}

func (*generator) Generate() string {
	return uuid.NewV4().String()
}

type MockUUID struct {
	*mock.Mock
}

func (m *MockUUID) Generate() string {
	return m.Called().String(0)
}

func NewTestUUID() *MockUUID {
	return &MockUUID{
		Mock: new(mock.Mock),
	}
}