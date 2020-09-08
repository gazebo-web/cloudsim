package fake

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestNewStore(t *testing.T) {
	s := NewFakeStore(new(int))
	assert.NotNil(t, s)
}

type fakeStoreSuite struct {
	suite.Suite
	store *store
}

func TestStore(t *testing.T) {
	suite.Run(t, new(fakeStoreSuite))
}

func (s *fakeStoreSuite) SetupTest() {
	st := NewFakeStore(new(int))
	s.store = st.(*store)
}

func (s *fakeStoreSuite) TestInnerValueIsZeroValueAtBeginning() {
	expected := 0
	s.Equal(&expected, s.store.value)
}

func (s *fakeStoreSuite) TestInnerValueIsNotDifferentThanZero() {
	expected := 1
	s.NotEqual(&expected, s.store.value)
}

func (s *fakeStoreSuite) TestGetValueIsSameAsInnerValue() {
	s.Equal(s.store.value, s.store.Get())
}

func (s *fakeStoreSuite) TestSetChangesInnerValue() {
	before := 0
	s.Equal(&before, s.store.value)

	var value int
	value = 5
	err := s.store.Set(&value)

	s.NoError(err)
	s.Equal(&value, s.store.value)
}

func (s *fakeStoreSuite) TestSetThrowsErrorWhenValueTypeIsDifferent() {
	before := 0
	s.Equal(&before, s.store.value)

	var value string
	value = "test"
	err := s.store.Set(&value)

	s.Equal(&before, s.store.value)
	s.Error(err)
}

func (s *fakeStoreSuite) TestLoadReadsChangesInnerValue() {
	before := 5
	err := s.store.Set(&before)
	s.NoError(err)

	actual := s.store.value.(*int)
	s.Equal(before, *actual)

	err = s.store.Load()
	s.NoError(err)

	after := 0
	actual = s.store.value.(*int)
	s.Equal(after, *actual)
}

func (s *fakeStoreSuite) TestSavePersistsInnerValue() {
	value := 5
	err := s.store.Set(&value)
	s.NoError(err)

	actual := s.store.value.(*int)
	s.Equal(value, *actual)

	err = s.store.Save()
	s.NoError(err)

	actual = s.store.data.(*int)
	s.Equal(value, *actual)
	s.True(s.store.persisted)
}
