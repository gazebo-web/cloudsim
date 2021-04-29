package structs

import (
	"errors"
	"github.com/fatih/structtag"
	"reflect"
)

var (
	// ErrFieldNotFound is returned when an operation is unable to find a target field.
	ErrFieldNotFound = errors.New("struct field not found")
)

// GetFieldTagValue gets a tag value from a struct field.
// The tag is always returned as a string, independent of field type.
//
// Example
// ```
// s := Struct{
//     Value int `default:"100" validate:"required"`
// }{}
// // Get the Value field default value
// tag := GetFieldTagValue(s, "Value", "default")
//
// // Prints 100
// fmt.Println(tag)
// ```
func GetFieldTagValue(targetStruct interface{}, fieldName string, tagName string) (string, error) {
	// If target struct is a pointer, get pointer value
	if reflect.TypeOf(targetStruct).Kind() == reflect.Ptr {
		targetStruct = reflect.ValueOf(targetStruct).Elem().Interface()
	}

	// Get field
	field, ok := reflect.TypeOf(targetStruct).FieldByName(fieldName)
	if !ok {
		return "", ErrFieldNotFound
	}

	// Parse field tags
	tags, err := structtag.Parse(string(field.Tag))
	if err != nil {
		return "", err
	}

	// Get tag value
	tag, err := tags.Get(tagName)
	if err != nil {
		return "", err
	}

	return tag.Value(), nil
}
