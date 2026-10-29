package domain

type ContactType string

const (
	ContactTypeEmail    ContactType = "EMAIL"
	ContactTypePhone    ContactType = "PHONE"
	ContactTypeMobile   ContactType = "MOBILE"
	ContactTypeWhatsApp ContactType = "WHATSAPP"
	ContactTypeFax      ContactType = "FAX"
	ContactTypeOther    ContactType = "OTHER"
)

func (c ContactType) IsValid() bool {
	switch c {
	case ContactTypeEmail, ContactTypePhone, ContactTypeMobile,
		ContactTypeWhatsApp, ContactTypeFax, ContactTypeOther:
		return true
	default:
		return false
	}
}

type CreateContactInput struct {
	PartyID      string
	ContactType  string
	ContactValue string
	Label        *string
	IsPrimary    bool
	IsActive     *bool
}

type UpdateContactInput struct {
	ID           string
	ContactType  string
	ContactValue string
	Label        *string
	IsPrimary    bool
	IsActive     bool
}

type CreateContactPersonInput struct {
	PartyID     string
	FirstName   string
	MiddleName  *string
	LastName    *string
	Designation *string
	Department  *string
	Email       *string
	Phone       *string
	Mobile      *string
	IsPrimary   bool
	IsActive    *bool
	Notes       *string
}

type UpdateContactPersonInput struct {
	ID          string
	FirstName   string
	MiddleName  *string
	LastName    *string
	Designation *string
	Department  *string
	Email       *string
	Phone       *string
	Mobile      *string
	IsPrimary   bool
	IsActive    bool
	Notes       *string
}
