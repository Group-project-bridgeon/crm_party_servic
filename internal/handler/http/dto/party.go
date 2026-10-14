package dto

type CreatePartyRequest struct {
	PartyType   string  `json:"party_type"`
	DisplayName string  `json:"display_name"`
	LegalName   *string `json:"legal_name"`
	Email       *string `json:"email"`
	Phone       *string `json:"phone"`
	Website     *string `json:"website"`
	TaxNumber   *string `json:"tax_number"`
	Notes       *string `json:"notes"`
}

type PartyResponse struct {
	ID          string  `json:"id"`
	PartyType   string  `json:"party_type"`
	DisplayName string  `json:"display_name"`
	LegalName   *string `json:"legal_name,omitempty"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Website     *string `json:"website,omitempty"`
	TaxNumber   *string `json:"tax_number,omitempty"`
	Notes       *string `json:"notes,omitempty"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type ListPartiesResponse struct {
	Data   []PartyResponse `json:"data"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}
