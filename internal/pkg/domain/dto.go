package domain

type DTO interface {
	ToMap() (map[string]interface{}, error)
}
