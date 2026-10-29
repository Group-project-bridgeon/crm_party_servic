package domain

type CreditStatus string

const (
	CreditStatusActive  CreditStatus = "ACTIVE"
	CreditStatusBlocked CreditStatus = "BLOCKED"
	CreditStatusHold    CreditStatus = "HOLD"
)

func (c CreditStatus) IsValid() bool {
	switch c {
	case CreditStatusActive, CreditStatusBlocked, CreditStatusHold:
		return true
	default:
		return false
	}
}

type CreateCustomerInput struct {
	PartyID            string
	CustomerCode       string
	CustomerGroup      *string
	PriceListID        *string
	CreditLimit        *float64
	PaymentTermsDays   *int
	CreditStatus       *string
	SalesRepID         *string
	TaxExempt          *bool
	TaxExemptionNumber *string
}

type UpdateCustomerInput struct {
	ID                 string
	CustomerGroup      *string
	PriceListID        *string
	CreditLimit        *float64
	PaymentTermsDays   *int
	CreditStatus       *string
	SalesRepID         *string
	TaxExempt          *bool
	TaxExemptionNumber *string
}

type CreateCustomerUserAccountInput struct {
	CustomerID string
	UserID     string
}
