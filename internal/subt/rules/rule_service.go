package rules

type IService interface {
	GetRuleByCircuitAndOwner(ruleType RuleType, circuit, owner string) (*Rule, error)
	GetRemainingSubmissions(owner, circuit string) (*int, error)
}

type Service struct {
}
