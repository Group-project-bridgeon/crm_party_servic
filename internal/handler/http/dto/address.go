package dto

type CreateAddressRequest struct {
	AddressType  string   `json:"address_type"`
	AddressName  *string  `json:"address_name"`
	AddressLine1 string   `json:"address_line1"`
	AddressLine2 *string  `json:"address_line2"`
	Landmark     *string  `json:"landmark"`
	City         string   `json:"city"`
	State        *string  `json:"state"`
	PostalCode   *string  `json:"postal_code"`
	CountryCode  string   `json:"country_code"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	IsPrimary    bool     `json:"is_primary"`
	IsActive     *bool    `json:"is_active"`
}

type UpdateAddressRequest struct {
	AddressType  string   `json:"address_type"`
	AddressName  *string  `json:"address_name"`
	AddressLine1 string   `json:"address_line1"`
	AddressLine2 *string  `json:"address_line2"`
	Landmark     *string  `json:"landmark"`
	City         string   `json:"city"`
	State        *string  `json:"state"`
	PostalCode   *string  `json:"postal_code"`
	CountryCode  string   `json:"country_code"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	IsPrimary    bool     `json:"is_primary"`
	IsActive     bool     `json:"is_active"`
}

type AddressResponse struct {
	ID           string   `json:"id"`
	PartyID      string   `json:"party_id"`
	AddressType  string   `json:"address_type"`
	AddressName  *string  `json:"address_name,omitempty"`
	AddressLine1 string   `json:"address_line1"`
	AddressLine2 *string  `json:"address_line2,omitempty"`
	Landmark     *string  `json:"landmark,omitempty"`
	City         string   `json:"city"`
	State        *string  `json:"state,omitempty"`
	PostalCode   *string  `json:"postal_code,omitempty"`
	CountryCode  string   `json:"country_code"`
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`
	IsPrimary    bool     `json:"is_primary"`
	IsActive     bool     `json:"is_active"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}
