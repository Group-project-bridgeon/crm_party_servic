package dto

type CreateOrganizationRequest struct {
	PartyID            string  `json:"party_id"`
	RegistrationNumber *string `json:"registration_number"`
	Industry           *string `json:"industry"`
	OrganizationType   *string `json:"organization_type"`
	FoundedDate        *string `json:"founded_date"`
	EmployeeCount      *int    `json:"employee_count"`
	Description        *string `json:"description"`
}

type OrganizationResponse struct {
	ID                 string  `json:"id"`
	PartyID            string  `json:"party_id"`
	RegistrationNumber *string `json:"registration_number,omitempty"`
	Industry           *string `json:"industry,omitempty"`
	OrganizationType   *string `json:"organization_type,omitempty"`
	FoundedDate        *string `json:"founded_date,omitempty"`
	EmployeeCount      *int32  `json:"employee_count,omitempty"`
	Description        *string `json:"description,omitempty"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}
