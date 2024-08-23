package traffic

import "fmt"

type budgetTrafficProvider struct {
}

func NewBudgetTrafficProvider() *budgetTrafficProvider {
	return &budgetTrafficProvider{}
}

func (p *budgetTrafficProvider) Get(options BudgetTrafficOptions) []BudgetTrafficRecord {
	fmt.Printf("%v\n", options)
	return nil
}
