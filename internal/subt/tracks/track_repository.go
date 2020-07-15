package tracks

import "github.com/jinzhu/gorm"

// Repository groups a set of methods to perform CRUD operations in the database for a certain Track.
type Repository interface {
	repositoryCreate
	repositoryRead
	repositoryUpdate
	repositoryDelete
}

// repositoryCreate has a method to Create a track in the database.
type repositoryCreate interface {
	Create(track Track) (*Track, error)
}

// repositoryRead has a method to get one or multiple tracks from the database.
type repositoryRead interface {
	Get(name string) (*Track, error)
	GetAll() ([]Track, error)
}

// repositoryUpdate has a method to update a track in the database.
type repositoryUpdate interface {
	Update(name string, track Track) (*Track, error)
}

// repositoryDelete has a method to delete a track from the database.
type repositoryDelete interface {
	Delete(name string) (*Track, error)
}

// repository is a Repository implementation.
type repository struct {
	db *gorm.DB
}

// Create adds the given track into the database.
// It returns the created track.
// It will return an error if the track creation failed.
func (r repository) Create(track Track) (*Track, error) {
	err := r.db.Model(&Track{}).Create(&track).Error
	if err != nil {
		return nil, err
	}
	return &track, nil
}

// Get returns the track with the given name.
// If the track wasn't found, it will return an error.
func (r repository) Get(name string) (*Track, error) {
	var t Track
	err := r.db.Model(&Track{}).First(&t).Where("name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// GetAll returns the list of tracks.
func (r repository) GetAll() ([]Track, error) {
	var t []Track
	err := r.db.Model(&Track{}).Find(&t).Error
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Update sets the given track values to the track that matches the given name.
// It returns the updated track.
// It will return an error if the update could not be processed.
func (r repository) Update(name string, track Track) (*Track, error) {
	err := r.db.Model(&Track{}).Save(&track).Where("name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &track, nil
}

// Delete removes a record with the given name.
// It returns the deleted record.
// It will return an error if the record could not be deleted.
func (r repository) Delete(name string) (*Track, error) {
	t, err := r.Get(name)
	if err != nil {
		return nil, err
	}
	err = r.db.Model(&Track{}).Delete(t, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return t, nil
}

// NewRepository initializes a new Repository implementation using gorm.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
