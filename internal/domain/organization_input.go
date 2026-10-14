package domain

type CreateOrganizationInput struct {
	PartyID            string
	RegistrationNumber *string
	Industry           *string
	OrganizationType   *string
	FoundedDate        *string
	EmployeeCount      *int
	Description        *string
}
