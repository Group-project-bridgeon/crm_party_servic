package dto

type CreateSupplierRequest struct {
	PartyID            string   `json:"party_id"`
	SupplierCode       string   `json:"supplier_code"`
	SupplierGroup      *string  `json:"supplier_group"`
	PaymentTermsDays   *int     `json:"payment_terms_days"`
	PaymentMethod      *string  `json:"payment_method"`
	CurrencyCode       *string  `json:"currency_code"`
	PurchaseRepID      *string  `json:"purchase_rep_id"`
	TaxExempt          *bool    `json:"tax_exempt"`
	TaxExemptionNumber *string  `json:"tax_exemption_number"`
	SupplierRating     *float64 `json:"supplier_rating"`
}

type UpdateSupplierRequest struct {
	SupplierGroup      *string  `json:"supplier_group"`
	PaymentTermsDays   *int     `json:"payment_terms_days"`
	PaymentMethod      *string  `json:"payment_method"`
	CurrencyCode       *string  `json:"currency_code"`
	PurchaseRepID      *string  `json:"purchase_rep_id"`
	TaxExempt          *bool    `json:"tax_exempt"`
	TaxExemptionNumber *string  `json:"tax_exemption_number"`
	SupplierRating     *float64 `json:"supplier_rating"`
}

type SupplierResponse struct {
	ID                 string   `json:"id"`
	PartyID            string   `json:"party_id"`
	SupplierCode       string   `json:"supplier_code"`
	SupplierGroup      *string  `json:"supplier_group,omitempty"`
	PaymentTermsDays   int32    `json:"payment_terms_days"`
	PaymentMethod      *string  `json:"payment_method,omitempty"`
	CurrencyCode       *string  `json:"currency_code,omitempty"`
	PurchaseRepID      *string  `json:"purchase_rep_id,omitempty"`
	TaxExempt          bool     `json:"tax_exempt"`
	TaxExemptionNumber *string  `json:"tax_exemption_number,omitempty"`
	SupplierRating     *float64 `json:"supplier_rating,omitempty"`
	CreatedAt          string   `json:"created_at"`
	UpdatedAt          string   `json:"updated_at"`
}

type ListSuppliersResponse struct {
	Data   []SupplierResponse `json:"data"`
	Limit  int                `json:"limit"`
	Offset int                `json:"offset"`
}

type CreateSupplierUserAccountRequest struct {
	UserID string `json:"user_id"`
}

type SupplierUserAccountResponse struct {
	ID         string `json:"id"`
	SupplierID string `json:"supplier_id"`
	UserID     string `json:"user_id"`
	IsActive   bool   `json:"is_active"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
