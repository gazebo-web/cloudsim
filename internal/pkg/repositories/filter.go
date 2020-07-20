package repositories

// Filter represents a generic filter used by repositories to filter data by a key-value set.
type Filter interface {
	Key() string
	Value() interface{}
}
