package repositories

// Filter represents a generic filter used by repositories to filter data by a key-value set.
type Filter interface {
	Template() string
	Value() interface{}
}
