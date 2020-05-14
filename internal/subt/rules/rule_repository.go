package rules

type IRepository interface {
	GetRuleByCircuitAndOwner(ruleType Type, circuit, owner string) (*Rule, error)
	GetRemainingSubmissions(owner, circuit string) (*int, error)
}

type Repository struct {
}
