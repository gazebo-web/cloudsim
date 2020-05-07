package robots

import "github.com/jinzhu/gorm"

type IRepository interface {

}

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) IRepository {
	var r IRepository
	r = &Repository{
		DB: db,
	}
	return r
}