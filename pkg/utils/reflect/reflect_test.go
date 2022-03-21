package reflect

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Tester interface {
	Value() interface{}
}

type Test struct {
	value int
}

func (t *Test) Value() interface{} { return t.value }

func TestAppendToSliceInt(t *testing.T) {
	slice := make([]int, 0)
	assert.NoError(t, AppendToSlice(&slice, 1, 2, 3))

	expected := []int{1, 2, 3}
	assert.Equal(t, expected, slice)
}

func TestAppendToSlicePtr(t *testing.T) {
	slice := make([]*int, 0)
	v1 := 1
	v2 := 2
	v3 := 3
	assert.NoError(t, AppendToSlice(&slice, &v1, &v2, &v3))

	expected := []*int{&v1, &v2, &v3}
	assert.Equal(t, expected, slice)
}

func TestAppendToSliceStruct(t *testing.T) {
	type T struct {
		Value string
	}

	slice := make([]T, 0)
	assert.NoError(t, AppendToSlice(&slice, T{"a"}, T{"b"}, T{"c"}))

	expected := []T{{"a"}, {"b"}, {"c"}}
	assert.Equal(t, expected, slice)
}

func TestAppendToSliceNotPointerError(t *testing.T) {
	slice := make([]int, 0)
	assert.Equal(t, ErrNotPointer, AppendToSlice(slice, 1, 2, 3))
}

func TestAppendToSliceNotSliceError(t *testing.T) {
	value := 0
	assert.Equal(t, ErrNotCollection, AppendToSlice(&value, 1, 2, 3))
}

func TestAppendToSlicePanicsWithInvalidType(t *testing.T) {
	slice := make([]int, 0)
	assert.Panics(t, func() { _ = AppendToSlice(&slice, "a", "b", "c") })
}

func TestSetMapValueInt(t *testing.T) {
	m := make(map[string]int, 0)
	assert.NoError(t, SetMapValue(m, "1", 100))
	assert.NoError(t, SetMapValue(m, "2", 2))
	assert.NoError(t, SetMapValue(m, "3", 3))
	// Replace existing value
	assert.NoError(t, SetMapValue(m, "1", 1))

	expected := map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
	}
	assert.Equal(t, expected, m)
}

func TestSetMapValuePtr(t *testing.T) {
	m := make(map[string]*int, 0)
	v1 := 1
	v2 := 2
	v3 := 3
	assert.NoError(t, SetMapValue(m, "1", nil))
	assert.NoError(t, SetMapValue(m, "2", &v2))
	assert.NoError(t, SetMapValue(m, "3", &v3))
	// Replace existing value
	assert.NoError(t, SetMapValue(m, "1", &v1))

	expected := map[string]*int{
		"1": &v1,
		"2": &v2,
		"3": &v3,
	}
	assert.Equal(t, expected, m)
}

func TestSetMapValueStruct(t *testing.T) {
	type T struct {
		Value string
	}

	m := make(map[string]T, 0)
	assert.NoError(t, SetMapValue(m, "1", T{"a"}))
	assert.NoError(t, SetMapValue(m, "2", T{"b"}))
	assert.NoError(t, SetMapValue(m, "3", T{"c"}))

	expected := map[string]T{
		"1": {"a"},
		"2": {"b"},
		"3": {"c"},
	}
	assert.Equal(t, expected, m)
}

func TestSetMapValueNotSliceError(t *testing.T) {
	value := 0
	assert.Equal(t, ErrNotCollection, SetMapValue(&value, "1", 1))
}

func TestSetMapValuePanicsWithInvalidType(t *testing.T) {
	m := make(map[string]int, 0)
	assert.Panics(t, func() { _ = SetMapValue(m, 1, "a") })
}

func TestNewCollectionValueInstanceIntSlice(t *testing.T) {
	value := make([]int, 0)

	out, err := NewCollectionValueInstance(value)
	assert.NoError(t, err)

	var expected int
	assert.Equal(t, expected, out)
}

func TestNewCollectionValueInstanceIntSlicePtr(t *testing.T) {
	value := make([]*int, 0)

	out, err := NewCollectionValueInstance(value)
	assert.NoError(t, err)

	var expected int
	fmt.Println(&expected, out)
	assert.NotSame(t, &expected, out)
}

func TestNewCollectionValueInstanceStructSlice(t *testing.T) {
	type T struct {
		Value string
	}

	value := make([]T, 0)

	out, err := NewCollectionValueInstance(value)
	assert.NoError(t, err)

	expected := T{}
	assert.Equal(t, expected, out)

	// Validate that the objects are different instances
	expected.Value = "expected"
	result := out.(T)
	result.Value = "out"
	assert.NotEqual(t, expected, result)
}

func TestNewCollectionValueInstanceStructSlicePtr(t *testing.T) {
	type T struct {
		Value string
	}

	value := make([]*T, 0)

	out, err := NewCollectionValueInstance(value)
	assert.NoError(t, err)

	expected := &T{}
	assert.Equal(t, expected, out)

	// Validate that the objects are different instances
	expected.Value = "expected"
	result := out.(*T)
	result.Value = "out"
	assert.NotEqual(t, expected, result)
}

func TestNewCollectionValueInstanceIntMap(t *testing.T) {
	value := make(map[string]int, 0)

	out, err := NewCollectionValueInstance(value)
	assert.NoError(t, err)

	var expected int
	assert.Equal(t, expected, out)
}

func TestNewCollectionValueInstanceIntMapPtr(t *testing.T) {
	value := make(map[string]*int, 0)

	out, err := NewCollectionValueInstance(value)
	assert.NoError(t, err)

	var expected int
	fmt.Println(&expected, out)
	assert.NotSame(t, &expected, out)
}

func TestNewCollectionValueInstanceStructMap(t *testing.T) {
	type T struct {
		Value string
	}

	value := make(map[string]T, 0)

	out, err := NewCollectionValueInstance(value)
	assert.NoError(t, err)

	expected := T{}
	assert.Equal(t, expected, out)

	// Validate that the objects are different instances
	expected.Value = "expected"
	result := out.(T)
	result.Value = "out"
	assert.NotEqual(t, expected, result)
}

func TestNewCollectionValueInstanceStructMapPtr(t *testing.T) {
	type T struct {
		Value string
	}

	value := make(map[string]*T, 0)

	out, err := NewCollectionValueInstance(value)
	assert.NoError(t, err)

	expected := &T{}
	assert.Equal(t, expected, out)

	// Validate that the objects are different instances
	expected.Value = "expected"
	result := out.(*T)
	result.Value = "out"
	assert.NotEqual(t, expected, result)
}

func TestNewCollectionValueInstanceNotSliceError(t *testing.T) {
	value := 0
	_, err := NewCollectionValueInstance(value)
	assert.Equal(t, ErrNotCollection, err)
}

func TestSetValueElementaryValue(t *testing.T) {
	var out int
	expected := 1
	assert.Nil(t, SetValue(&out, 1))
	assert.Equal(t, expected, out)
}

func TestSetValueStructPointer(t *testing.T) {
	out := Test{}
	expected := 1
	assert.Nil(t, SetValue(&out, Test{expected}))
	assert.Equal(t, expected, out.value)
}

func TestSetValueInterface(t *testing.T) {
	var out Tester
	expected := 1
	assert.Nil(t, SetValue(&out, &Test{expected}))
	assert.Equal(t, expected, out.Value())
}

func TestSetValueNotPointerError(t *testing.T) {
	out := 1
	assert.Equal(t, ErrNotPointer, SetValue(out, Test{}))
}

func TestSetValueInvalidOutValueError(t *testing.T) {
	out := 1
	assert.Equal(t, ErrInvalidOutValue, SetValue(&out, Test{}))
}

func TestNewInstanceElementaryValue(t *testing.T) {
	expected := 1
	out := NewInstance(expected)
	assert.IsType(t, expected, out)
	assert.Equal(t, 0, out)
}

func TestNewInstanceElementaryValuePointer(t *testing.T) {
	expected := 1
	out := NewInstance(&expected).(*int)
	assert.IsType(t, &expected, out)
	// Since the new instance will contain the zero value, the values should not match
	assert.NotEqual(t, expected, *out)
	// Verify that the output value contains the zero value
	assert.Equal(t, 0, *out)
}

func TestNewInstanceStructPointer(t *testing.T) {
	expected := &Test{1}
	out := NewInstance(expected).(*Test)
	// Since the new instance will contain the zero value, the values should not match
	assert.NotEqual(t, expected, out)
	// Verify that the output value contains the zero value
	assert.Equal(t, Test{}, *out)
}

func TestNewInstanceInterface(t *testing.T) {
	expected := &Test{1}
	out := NewInstance(expected).(Tester)
	assert.IsType(t, expected, out)
	// Since the new instance will contain the zero value, the values should not match
	assert.NotEqual(t, expected, out)
	// Verify that the output value contains the zero value
	assert.Equal(t, Test{}, *out.(*Test))
}
