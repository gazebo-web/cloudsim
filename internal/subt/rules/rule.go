package rules

import "github.com/jinzhu/gorm"

// CustomRuleType defines the type for circuit custom rules
type RuleType string

// List of rule types
const (
	// Maximum number of submissions allowed for a specific circuit
	MaxSubmissions RuleType = "max_submissions"
)

// Rule holds custom rules for a specific combination of owner
// (user or organization) or circuit. Rules contain arbitrary values that can
// be used to configure specific aspects of a circuit/application (e.g.
// max_submissions - A custom rule for a specific owner to allow for
// extra submissions in a specific circuit). Rules for several owners or
// circuits can be defined by creating rules with NULL values in either fields.
// Rules with NULL values have less priority than rules with values, with the
// following priority: owner, circuit. This means that a rule with NULL circuit
// and owner will apply to ALL circuits for ALL owners, but any rule with either
// circuit or owner will override this general rule.
type Rule struct {
	gorm.Model
	Owner    *string        `json:"owner"`
	Circuit  *string        `json:"circuit" validate:"iscircuit"`
	RuleType RuleType 		`gorm:"not null" json:"rule_type" validate:"isruletype"`
	Value    string         `gorm:"not null" json:"value"`
}


// Rules is a slice of Rule
type Rules []Rule

