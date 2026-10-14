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
