package rules

type Repository interface {
	GetRuleByCircuitAndOwner(ruleType Type, circuit, owner string) (*Rule, error)
	GetRemainingSubmissions(owner, circuit string) (*int, error)
}

type repository struct {
}
