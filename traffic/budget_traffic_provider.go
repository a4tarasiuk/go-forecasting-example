package traffic

type budgetTrafficProvider struct {
}

func NewBudgetTrafficProvider() *budgetTrafficProvider {
	return &budgetTrafficProvider{}
}

func (p *budgetTrafficProvider) Get(options BudgetTrafficOptions) []BudgetTrafficRecord {
	return nil
}
