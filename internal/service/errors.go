package service

import "errors"

var (
	ErrInvalidPartyType          = errors.New("invalid party type")
	ErrEmptyDisplayName          = errors.New("display name is required")
	ErrInvalidPagination         = errors.New("invalid pagination")
	ErrPartyNotFound             = errors.New("party not found")
	ErrInvalidIsActive           = errors.New("invalid is_active")
	ErrInvalidPartyID            = errors.New("invalid party id")
	ErrPartyNotOrganization      = errors.New("party is not an organization")
	ErrOrganizationAlreadyExists = errors.New("organization already exists")
	ErrInvalidEmployeeCount      = errors.New("invalid employee count")
	ErrInvalidFoundedDate        = errors.New("invalid founded date")
	ErrOrganizationNotFound      = errors.New("organization not found")
	ErrPartyNotIndividual        = errors.New("party is not an individual")
	ErrIndividualAlreadyExists   = errors.New("individual already exists")
	ErrIndividualNotFound        = errors.New("individual not found")
	ErrEmptyFirstName            = errors.New("first name is required")
	ErrInvalidDateOfBirth        = errors.New("invalid date of birth")

	// Customer errors
	ErrCustomerNotFound                 = errors.New("customer not found")
	ErrCustomerAlreadyExists            = errors.New("customer already exists for this party")
	ErrDuplicateCustomerCode            = errors.New("customer code already exists")
	ErrInvalidCreditLimit               = errors.New("credit limit must be zero or greater")
	ErrInvalidPaymentTerms              = errors.New("payment terms days must be zero or greater")
	ErrInvalidCreditStatus              = errors.New("invalid credit status, must be ACTIVE, BLOCKED, or HOLD")
	ErrEmptyCustomerCode                = errors.New("customer code is required")
	ErrInvalidCustomerID                = errors.New("invalid customer id")
	ErrInvalidUserID                    = errors.New("invalid user id")
	ErrCustomerUserAccountAlreadyExists = errors.New("user account already linked to this customer")
	ErrCustomerUserAccountNotFound      = errors.New("customer user account not found")

	// Supplier errors
	ErrSupplierNotFound                 = errors.New("supplier not found")
	ErrSupplierAlreadyExists            = errors.New("supplier already exists for this party")
	ErrDuplicateSupplierCode            = errors.New("supplier code already exists")
	ErrEmptySupplierCode                = errors.New("supplier code is required")
	ErrInvalidSupplierRating            = errors.New("supplier rating must be between 0 and 5")
	ErrInvalidCurrencyCode              = errors.New("invalid currency code, must be 3 uppercase letters")
	ErrInvalidSupplierID                = errors.New("invalid supplier id")
	ErrSupplierUserAccountAlreadyExists = errors.New("user account already linked to this supplier")
	ErrSupplierUserAccountNotFound      = errors.New("supplier user account not found")

	// Address errors
	ErrAddressNotFound    = errors.New("party address not found")
	ErrInvalidAddressType = errors.New("invalid address type")
	ErrEmptyAddressLine1  = errors.New("address line 1 is required")
	ErrEmptyCity          = errors.New("city is required")
	ErrInvalidCountryCode = errors.New("country code must be a 2-letter ISO code")
	ErrInvalidCoordinates = errors.New("invalid latitude or longitude")
	ErrInvalidAddressID   = errors.New("invalid address id")

	// Contact errors
	ErrContactNotFound    = errors.New("party contact not found")
	ErrInvalidContactType = errors.New("invalid contact type")
	ErrEmptyContactValue  = errors.New("contact value is required")
	ErrInvalidContactID   = errors.New("invalid contact id")

	// Contact Person errors
	ErrContactPersonNotFound  = errors.New("party contact person not found")
	ErrInvalidContactPersonID = errors.New("invalid contact person id")

	// Note errors
	ErrNoteNotFound     = errors.New("party note not found")
	ErrEmptyNoteContent = errors.New("note content cannot be empty")
	ErrInvalidNoteID    = errors.New("invalid note id")
)
