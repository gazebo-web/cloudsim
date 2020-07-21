package repositories

// gormFilter is a filter to be used with gorm.
type gormFilter struct {
	template string
	values    []interface{}
}

// Template returns the filter's query template.
func (w gormFilter) Template() string {
	return w.template
}

// Value returns the filter's value.
func (w gormFilter) Values() []interface{} {
	return w.values
}

// NewGormFilter initializes a new filter with the given template and value to be used with gorm.
func NewGormFilter(template string, values ...interface{}) Filter {
	return &gormFilter{
		template: template,
		values:    values,
	}
}
