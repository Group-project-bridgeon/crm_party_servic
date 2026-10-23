package dto

type CreateCustomerRequest struct {
	PartyID            string   `json:"party_id"`
	CustomerCode       string   `json:"customer_code"`
	CustomerGroup      *string  `json:"customer_group"`
	PriceListID        *string  `json:"price_list_id"`
	CreditLimit        *float64 `json:"credit_limit"`
	PaymentTermsDays   *int     `json:"payment_terms_days"`
	CreditStatus       *string  `json:"credit_status"`
	SalesRepID         *string  `json:"sales_rep_id"`
	TaxExempt          *bool    `json:"tax_exempt"`
	TaxExemptionNumber *string  `json:"tax_exemption_number"`
}

type UpdateCustomerRequest struct {
	CustomerGroup      *string  `json:"customer_group"`
	PriceListID        *string  `json:"price_list_id"`
	CreditLimit        *float64 `json:"credit_limit"`
	PaymentTermsDays   *int     `json:"payment_terms_days"`
	CreditStatus       *string  `json:"credit_status"`
	SalesRepID         *string  `json:"sales_rep_id"`
	TaxExempt          *bool    `json:"tax_exempt"`
	TaxExemptionNumber *string  `json:"tax_exemption_number"`
}

type CustomerResponse struct {
	ID                 string  `json:"id"`
	PartyID            string  `json:"party_id"`
	CustomerCode       string  `json:"customer_code"`
	CustomerGroup      *string `json:"customer_group,omitempty"`
	PriceListID        *string `json:"price_list_id,omitempty"`
	CreditLimit        float64 `json:"credit_limit"`
	PaymentTermsDays   int32   `json:"payment_terms_days"`
	CreditStatus       string  `json:"credit_status"`
	SalesRepID         *string `json:"sales_rep_id,omitempty"`
	TaxExempt          bool    `json:"tax_exempt"`
	TaxExemptionNumber *string `json:"tax_exemption_number,omitempty"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}

type ListCustomersResponse struct {
	Data   []CustomerResponse `json:"data"`
	Limit  int                `json:"limit"`
	Offset int                `json:"offset"`
}

type CreateCustomerUserAccountRequest struct {
	UserID string `json:"user_id"`
}

type UpdateUserAccountStatusRequest struct {
	IsActive bool `json:"is_active"`
}

type CustomerUserAccountResponse struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	UserID     string `json:"user_id"`
	IsActive   bool   `json:"is_active"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
