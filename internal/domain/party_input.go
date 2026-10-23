package domain

type CreatePartyInput struct {
	PartyType   PartyType
	DisplayName string
	LegalName   *string
	Email       *string
	Phone       *string
	Website     *string
	TaxNumber   *string
	Notes       *string
}

type UpdatePartyInput struct {
	ID          string
	DisplayName string
	LegalName   *string
	Email       *string
	Phone       *string
	Website     *string
	TaxNumber   *string
	Notes       *string
}

