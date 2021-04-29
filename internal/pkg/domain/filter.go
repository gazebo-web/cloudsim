package domain

// Filter represents a generic filter used by repositories to filter data by a key-value set.
type Filter interface {
	// Template represents a SQL Syntax.
	// Example: `name = ? AND age = ?`
	Template() string
	// Values returns the values used by the SQL Syntax.
	// Example: `["Test", 33]`
	Values() []interface{}
}
