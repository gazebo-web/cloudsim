package domain

// DTO represents a generic Data Transfer Object. It has a method to convert a struct into a map.
type DTO interface {
	ToMap() (map[string]interface{}, error)
}
