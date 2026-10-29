package dto

type CreateContactRequest struct {
	ContactType  string  `json:"contact_type"`
	ContactValue string  `json:"contact_value"`
	Label        *string `json:"label"`
	IsPrimary    bool    `json:"is_primary"`
	IsActive     *bool   `json:"is_active"`
}

type UpdateContactRequest struct {
	ContactType  string  `json:"contact_type"`
	ContactValue string  `json:"contact_value"`
	Label        *string `json:"label"`
	IsPrimary    bool    `json:"is_primary"`
	IsActive     bool    `json:"is_active"`
}

type ContactResponse struct {
	ID           string  `json:"id"`
	PartyID      string  `json:"party_id"`
	ContactType  string  `json:"contact_type"`
	ContactValue string  `json:"contact_value"`
	Label        *string `json:"label,omitempty"`
	IsPrimary    bool    `json:"is_primary"`
	IsActive     bool    `json:"is_active"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type CreateContactPersonRequest struct {
	FirstName   string  `json:"first_name"`
	MiddleName  *string `json:"middle_name"`
	LastName    *string `json:"last_name"`
	Designation *string `json:"designation"`
	Department  *string `json:"department"`
	Email       *string `json:"email"`
	Phone       *string `json:"phone"`
	Mobile      *string `json:"mobile"`
	IsPrimary   bool    `json:"is_primary"`
	IsActive    *bool   `json:"is_active"`
	Notes       *string `json:"notes"`
}

type UpdateContactPersonRequest struct {
	FirstName   string  `json:"first_name"`
	MiddleName  *string `json:"middle_name"`
	LastName    *string `json:"last_name"`
	Designation *string `json:"designation"`
	Department  *string `json:"department"`
	Email       *string `json:"email"`
	Phone       *string `json:"phone"`
	Mobile      *string `json:"mobile"`
	IsPrimary   bool    `json:"is_primary"`
	IsActive    bool    `json:"is_active"`
	Notes       *string `json:"notes"`
}

type ContactPersonResponse struct {
	ID          string  `json:"id"`
	PartyID     string  `json:"party_id"`
	FirstName   string  `json:"first_name"`
	MiddleName  *string `json:"middle_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	Designation *string `json:"designation,omitempty"`
	Department  *string `json:"department,omitempty"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Mobile      *string `json:"mobile,omitempty"`
	IsPrimary   bool    `json:"is_primary"`
	IsActive    bool    `json:"is_active"`
	Notes       *string `json:"notes,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
