package dto

type CreateIndividualRequest struct {
	PartyID     string  `json:"party_id"`
	FirstName   string  `json:"first_name"`
	MiddleName  *string `json:"middle_name"`
	LastName    *string `json:"last_name"`
	DateOfBirth *string `json:"date_of_birth"`
	Gender      *string `json:"gender"`
	JobTitle    *string `json:"job_title"`
	NationalID  *string `json:"national_id"`
}

type UpdateIndividualRequest struct {
	FirstName   string  `json:"first_name"`
	MiddleName  *string `json:"middle_name"`
	LastName    *string `json:"last_name"`
	DateOfBirth *string `json:"date_of_birth"`
	Gender      *string `json:"gender"`
	JobTitle    *string `json:"job_title"`
	NationalID  *string `json:"national_id"`
}

type IndividualResponse struct {
	ID          string  `json:"id"`
	PartyID     string  `json:"party_id"`
	FirstName   string  `json:"first_name"`
	MiddleName  *string `json:"middle_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	DateOfBirth *string `json:"date_of_birth,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	JobTitle    *string `json:"job_title,omitempty"`
	NationalID  *string `json:"national_id,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
