package domain

import "github.com/group-project/crm_party_service/internal/db"

type CompositeCustomerInput struct {
	// Party
	PartyType   PartyType
	DisplayName string
	LegalName   *string
	Email       *string
	Phone       *string
	Website     *string
	TaxNumber   *string
	Notes       *string

	// Organization profile (when PartyType == ORGANIZATION)
	RegistrationNumber *string
	Industry           *string
	OrganizationType   *string
	FoundedDate        *string
	EmployeeCount      *int
	Description        *string

	// Individual profile (when PartyType == INDIVIDUAL)
	FirstName   *string
	MiddleName  *string
	LastName    *string
	DateOfBirth *string
	Gender      *string
	JobTitle    *string
	NationalID  *string

	// Customer details
	CustomerCode       string
	CustomerGroup      *string
	PriceListID        *string
	CreditLimit        *float64
	PaymentTermsDays   *int
	CreditStatus       *string
	SalesRepID         *string
	TaxExempt          *bool
	TaxExemptionNumber *string

	// Initial Primary Address (optional)
	Address *CreateAddressInput
}

type CompositeSupplierInput struct {
	// Party
	PartyType   PartyType
	DisplayName string
	LegalName   *string
	Email       *string
	Phone       *string
	Website     *string
	TaxNumber   *string
	Notes       *string

	// Organization profile (when PartyType == ORGANIZATION)
	RegistrationNumber *string
	Industry           *string
	OrganizationType   *string
	FoundedDate        *string
	EmployeeCount      *int
	Description        *string

	// Individual profile (when PartyType == INDIVIDUAL)
	FirstName   *string
	MiddleName  *string
	LastName    *string
	DateOfBirth *string
	Gender      *string
	JobTitle    *string
	NationalID  *string

	// Supplier details
	SupplierCode       string
	SupplierGroup      *string
	PaymentTermsDays   *int
	PaymentMethod      *string
	CurrencyCode       *string
	PurchaseRepID      *string
	TaxExempt          *bool
	TaxExemptionNumber *string
	SupplierRating     *float64

	// Initial Primary Address (optional)
	Address *CreateAddressInput
}

type Party360 struct {
	Party          db.Party
	Organization   *db.Organization
	Individual     *db.Individual
	Customer       *db.Customer
	Supplier       *db.Supplier
	Addresses      []db.PartyAddress
	Contacts       []db.PartyContact
	ContactPersons []db.PartyContactPerson
	Notes          []db.PartyNote
}

type CompositeCustomerResult struct {
	Party        db.Party
	Organization *db.Organization
	Individual   *db.Individual
	Customer     db.Customer
	Address      *db.PartyAddress
}

type CompositeSupplierResult struct {
	Party        db.Party
	Organization *db.Organization
	Individual   *db.Individual
	Supplier     db.Supplier
	Address      *db.PartyAddress
}
