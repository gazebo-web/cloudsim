package rules

import "github.com/jinzhu/gorm"

type Repository interface {
	GetByCircuitAndOwner(ruleType Type, circuit, owner string) (*Rule, error)
}

type repository struct {
	Db *gorm.DB
}

// GetByCircuitAndOwner returns the rule value for a specific circuit and owner.
func (r *repository) GetByCircuitAndOwner(ruleType Type, circuit, owner string) (*Rule, error) {
	var rule Rule
	err := r.Db.Model(&Rule{}).Where("rule_type = ?", ruleType).
		Where("circuit = ? OR circuit IS NULL", circuit).
		Where("owner = ? OR owner IS NULL", owner).
		Order("owner DESC, circuit DESC").Limit(1).First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		Db: db,
	}
}
