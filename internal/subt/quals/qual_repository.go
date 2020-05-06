package quals

import "github.com/jinzhu/gorm"

type IRepository interface {
	GetByOwnerAndCircuit(owner, circuit string) (*Qualification, error)
}

type Repository struct {
	DB *gorm.DB
}

// NewRepository initializes a new IRepository implementation.
func NewRepository(db *gorm.DB) IRepository {
	var repository IRepository
	repository = &Repository{
		DB: db,
	}
	return repository
}

// GetByOwnerAndCircuit returns a qualification from the given owner and circuit.
// In any other case, it will return an error.
func (r Repository) GetByOwnerAndCircuit(owner, circuit string) (*Qualification, error) {
	var qualification Qualification
	if err := r.DB.Model(&Qualification{}).
		Where("owner = ? AND circuit = ?", owner, circuit).
		First(&qualification).Error; err != nil {
		return nil, err
	}
	return &qualification, nil
}
