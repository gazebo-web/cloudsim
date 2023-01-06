package structs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TagTest struct {
	Name  string `validate:"required" default:"test"`
	Value int    `validate:"required" default:"100"`
}

func TestGetFieldTagValueStruct(t *testing.T) {
	value := TagTest{}

	tagValue, err := GetFieldTagValue(value, "Name", "default")
	assert.Equal(t, nil, err)
	assert.Equal(t, "test", tagValue)
}

func TestGetNameFieldDefaultTagStructPointer(t *testing.T) {
	value := &TagTest{}

	tagValue, err := GetFieldTagValue(value, "Name", "default")
	assert.Equal(t, nil, err)
	assert.Equal(t, "test", tagValue)
}

func TestGetNameFieldDefaultTagNiPointer(t *testing.T) {
	var value *TagTest

	assert.Panics(t, func() {
		_, _ = GetFieldTagValue(value, "Name", "default")
	})
}

func TestGetInvalidFieldDefaultTag(t *testing.T) {
	value := TagTest{}

	_, err := GetFieldTagValue(value, "invalid", "default")
	assert.Equal(t, ErrFieldNotFound, err)
}

func TestGetNameFieldInvalidTag(t *testing.T) {
	value := TagTest{}

	_, err := GetFieldTagValue(value, "Name", "invalid")
	assert.NotEqual(t, nil, err)
}
