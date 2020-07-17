package repositories

// Repository represents a generic repository layer interface.
type Repository interface {
	GetAll(offset, limit *int) ([]interface{}, error)
	Get(uuids []string) ([]interface{}, error)
	Update(uuids []string, in interface{}) ([]interface{}, error)
	Delete(uuids []string) ([]interface{}, error)
}
