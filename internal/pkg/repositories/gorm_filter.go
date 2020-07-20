package repositories

import "fmt"

// gormFilter is a filter to be used with gorm.
type gormFilter struct {
	key   string
	value interface{}
}

// Key returns the filter's key.
func (w gormFilter) Key() string {
	return w.key
}

// Value returns the filter's value.
func (w gormFilter) Value() interface{} {
	return w.value
}

// NewGormFilter initializes a new filter with the given key and value to be used with gorm.
func NewGormFilter(key string, value interface{}) Filter {
	return &gormFilter{
		key:   fmt.Sprintf("%s = ?", key),
		value: value,
	}
}
