package dto

type CreateCompositeCustomerRequest struct {
	PartyType   string  `json:"party_type"`
	DisplayName string  `json:"display_name"`
	LegalName   *string `json:"legal_name"`
	Email       *string `json:"email"`
	Phone       *string `json:"phone"`
	Website     *string `json:"website"`
	TaxNumber   *string `json:"tax_number"`
	Notes       *string `json:"notes"`

	// Organization fields
	RegistrationNumber *string `json:"registration_number"`
	Industry           *string `json:"industry"`
	OrganizationType   *string `json:"organization_type"`
	FoundedDate        *string `json:"founded_date"`
	EmployeeCount      *int    `json:"employee_count"`
	Description        *string `json:"description"`

	// Individual fields
	FirstName   *string `json:"first_name"`
	MiddleName  *string `json:"middle_name"`
	LastName    *string `json:"last_name"`
	DateOfBirth *string `json:"date_of_birth"`
	Gender      *string `json:"gender"`
	JobTitle    *string `json:"job_title"`
	NationalID  *string `json:"national_id"`

	// Customer fields
	CustomerCode       string   `json:"customer_code"`
	CustomerGroup      *string  `json:"customer_group"`
	PriceListID        *string  `json:"price_list_id"`
	CreditLimit        *float64 `json:"credit_limit"`
	PaymentTermsDays   *int     `json:"payment_terms_days"`
	CreditStatus       *string  `json:"credit_status"`
	SalesRepID         *string  `json:"sales_rep_id"`
	TaxExempt          *bool    `json:"tax_exempt"`
	TaxExemptionNumber *string  `json:"tax_exemption_number"`

	// Optional initial address
	Address *CreateAddressRequest `json:"address"`
}

type CreateCompositeSupplierRequest struct {
	PartyType   string  `json:"party_type"`
	DisplayName string  `json:"display_name"`
	LegalName   *string `json:"legal_name"`
	Email       *string `json:"email"`
	Phone       *string `json:"phone"`
	Website     *string `json:"website"`
	TaxNumber   *string `json:"tax_number"`
	Notes       *string `json:"notes"`

	// Organization fields
	RegistrationNumber *string `json:"registration_number"`
	Industry           *string `json:"industry"`
	OrganizationType   *string `json:"organization_type"`
	FoundedDate        *string `json:"founded_date"`
	EmployeeCount      *int    `json:"employee_count"`
	Description        *string `json:"description"`

	// Individual fields
	FirstName   *string `json:"first_name"`
	MiddleName  *string `json:"middle_name"`
	LastName    *string `json:"last_name"`
	DateOfBirth *string `json:"date_of_birth"`
	Gender      *string `json:"gender"`
	JobTitle    *string `json:"job_title"`
	NationalID  *string `json:"national_id"`

	// Supplier fields
	SupplierCode       string   `json:"supplier_code"`
	SupplierGroup      *string  `json:"supplier_group"`
	PaymentTermsDays   *int     `json:"payment_terms_days"`
	PaymentMethod      *string  `json:"payment_method"`
	CurrencyCode       *string  `json:"currency_code"`
	PurchaseRepID      *string  `json:"purchase_rep_id"`
	TaxExempt          *bool    `json:"tax_exempt"`
	TaxExemptionNumber *string  `json:"tax_exemption_number"`
	SupplierRating     *float64 `json:"supplier_rating"`

	// Optional initial address
	Address *CreateAddressRequest `json:"address"`
}

type Party360Response struct {
	Party          PartyResponse           `json:"party"`
	Organization   *OrganizationResponse   `json:"organization,omitempty"`
	Individual     *IndividualResponse     `json:"individual,omitempty"`
	Customer       *CustomerResponse       `json:"customer,omitempty"`
	Supplier       *SupplierResponse       `json:"supplier,omitempty"`
	Addresses      []AddressResponse       `json:"addresses"`
	Contacts       []ContactResponse       `json:"contacts"`
	ContactPersons []ContactPersonResponse `json:"contact_persons"`
	Notes          []NoteResponse          `json:"notes"`
}

type CompositeCustomerResponse struct {
	Party        PartyResponse         `json:"party"`
	Organization *OrganizationResponse `json:"organization,omitempty"`
	Individual   *IndividualResponse   `json:"individual,omitempty"`
	Customer     CustomerResponse      `json:"customer"`
	Address      *AddressResponse      `json:"address,omitempty"`
}

type CompositeSupplierResponse struct {
	Party        PartyResponse         `json:"party"`
	Organization *OrganizationResponse `json:"organization,omitempty"`
	Individual   *IndividualResponse   `json:"individual,omitempty"`
	Supplier     SupplierResponse      `json:"supplier"`
	Address      *AddressResponse      `json:"address,omitempty"`
}
