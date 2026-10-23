package domain

type CreateIndividualInput struct {
	PartyID     string
	FirstName   string
	MiddleName  *string
	LastName    *string
	DateOfBirth *string
	Gender      *string
	JobTitle    *string
	NationalID  *string
}

type UpdateIndividualInput struct {
	PartyID     string
	FirstName   string
	MiddleName  *string
	LastName    *string
	DateOfBirth *string
	Gender      *string
	JobTitle    *string
	NationalID  *string
}
