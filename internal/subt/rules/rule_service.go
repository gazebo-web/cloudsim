package rules

type Service interface {
	GetRuleByCircuitAndOwner(ruleType Type, circuit, owner string) (*Rule, error)
	GetRemainingSubmissions(owner, circuit string) (*int, error)
}

type service struct {
}
