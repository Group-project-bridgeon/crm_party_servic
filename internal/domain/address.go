package domain

type AddressType string

const (
	AddressTypeBilling    AddressType = "BILLING"
	AddressTypeShipping   AddressType = "SHIPPING"
	AddressTypeOffice     AddressType = "OFFICE"
	AddressTypeRegistered AddressType = "REGISTERED"
	AddressTypeWarehouse  AddressType = "WAREHOUSE"
	AddressTypeHome       AddressType = "HOME"
	AddressTypeOther      AddressType = "OTHER"
)

func (a AddressType) IsValid() bool {
	switch a {
	case AddressTypeBilling, AddressTypeShipping, AddressTypeOffice,
		AddressTypeRegistered, AddressTypeWarehouse, AddressTypeHome, AddressTypeOther:
		return true
	default:
		return false
	}
}

type CreateAddressInput struct {
	PartyID      string
	AddressType  string
	AddressName  *string
	AddressLine1 string
	AddressLine2 *string
	Landmark     *string
	City         string
	State        *string
	PostalCode   *string
	CountryCode  string
	Latitude     *float64
	Longitude    *float64
	IsPrimary    bool
	IsActive     *bool
}

type UpdateAddressInput struct {
	ID           string
	AddressType  string
	AddressName  *string
	AddressLine1 string
	AddressLine2 *string
	Landmark     *string
	City         string
	State        *string
	PostalCode   *string
	CountryCode  string
	Latitude     *float64
	Longitude    *float64
	IsPrimary    bool
	IsActive     bool
}
