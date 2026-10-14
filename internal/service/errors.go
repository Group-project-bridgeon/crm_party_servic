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
	ErrOrganizationNotFound = errors.New("organization not found")
)
