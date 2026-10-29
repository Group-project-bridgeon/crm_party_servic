package domain

type CreateSupplierInput struct {
	PartyID            string
	SupplierCode       string
	SupplierGroup      *string
	PaymentTermsDays   *int
	PaymentMethod      *string
	CurrencyCode       *string
	PurchaseRepID      *string
	TaxExempt          *bool
	TaxExemptionNumber *string
	SupplierRating     *float64
}

type UpdateSupplierInput struct {
	ID                 string
	SupplierGroup      *string
	PaymentTermsDays   *int
	PaymentMethod      *string
	CurrencyCode       *string
	PurchaseRepID      *string
	TaxExempt          *bool
	TaxExemptionNumber *string
	SupplierRating     *float64
}

type CreateSupplierUserAccountInput struct {
	SupplierID string
	UserID     string
}
