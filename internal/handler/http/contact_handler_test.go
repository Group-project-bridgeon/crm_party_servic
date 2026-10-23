package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
	"github.com/group-project/crm_party_service/internal/domain"
	"github.com/group-project/crm_party_service/internal/handler/http/dto"
)

type mockContactService struct {
	createContactFunc       func(ctx context.Context, input domain.CreateContactInput) (db.PartyContact, error)
	getContactByIDFunc      func(ctx context.Context, id pgtype.UUID) (db.PartyContact, error)
	listContactsByPartyIDFunc func(ctx context.Context, partyID pgtype.UUID, contactType *string) ([]db.PartyContact, error)
	updateContactFunc       func(ctx context.Context, input domain.UpdateContactInput) (db.PartyContact, error)
	deleteContactFunc       func(ctx context.Context, id pgtype.UUID) error

	createContactPersonFunc       func(ctx context.Context, input domain.CreateContactPersonInput) (db.PartyContactPerson, error)
	getContactPersonByIDFunc      func(ctx context.Context, id pgtype.UUID) (db.PartyContactPerson, error)
	listContactPersonsByPartyIDFunc func(ctx context.Context, partyID pgtype.UUID) ([]db.PartyContactPerson, error)
	updateContactPersonFunc       func(ctx context.Context, input domain.UpdateContactPersonInput) (db.PartyContactPerson, error)
	deleteContactPersonFunc       func(ctx context.Context, id pgtype.UUID) error
}

func (m *mockContactService) CreateContact(ctx context.Context, input domain.CreateContactInput) (db.PartyContact, error) {
	if m.createContactFunc != nil {
		return m.createContactFunc(ctx, input)
	}
	return db.PartyContact{}, nil
}

func (m *mockContactService) GetContactByID(ctx context.Context, id pgtype.UUID) (db.PartyContact, error) {
	if m.getContactByIDFunc != nil {
		return m.getContactByIDFunc(ctx, id)
	}
	return db.PartyContact{}, nil
}

func (m *mockContactService) ListContactsByPartyID(ctx context.Context, partyID pgtype.UUID, contactType *string) ([]db.PartyContact, error) {
	if m.listContactsByPartyIDFunc != nil {
		return m.listContactsByPartyIDFunc(ctx, partyID, contactType)
	}
	return []db.PartyContact{}, nil
}

func (m *mockContactService) UpdateContact(ctx context.Context, input domain.UpdateContactInput) (db.PartyContact, error) {
	if m.updateContactFunc != nil {
		return m.updateContactFunc(ctx, input)
	}
	return db.PartyContact{}, nil
}

func (m *mockContactService) DeleteContact(ctx context.Context, id pgtype.UUID) error {
	if m.deleteContactFunc != nil {
		return m.deleteContactFunc(ctx, id)
	}
	return nil
}

func (m *mockContactService) CreateContactPerson(ctx context.Context, input domain.CreateContactPersonInput) (db.PartyContactPerson, error) {
	if m.createContactPersonFunc != nil {
		return m.createContactPersonFunc(ctx, input)
	}
	return db.PartyContactPerson{}, nil
}

func (m *mockContactService) GetContactPersonByID(ctx context.Context, id pgtype.UUID) (db.PartyContactPerson, error) {
	if m.getContactPersonByIDFunc != nil {
		return m.getContactPersonByIDFunc(ctx, id)
	}
	return db.PartyContactPerson{}, nil
}

func (m *mockContactService) ListContactPersonsByPartyID(ctx context.Context, partyID pgtype.UUID) ([]db.PartyContactPerson, error) {
	if m.listContactPersonsByPartyIDFunc != nil {
		return m.listContactPersonsByPartyIDFunc(ctx, partyID)
	}
	return []db.PartyContactPerson{}, nil
}

func (m *mockContactService) UpdateContactPerson(ctx context.Context, input domain.UpdateContactPersonInput) (db.PartyContactPerson, error) {
	if m.updateContactPersonFunc != nil {
		return m.updateContactPersonFunc(ctx, input)
	}
	return db.PartyContactPerson{}, nil
}

func (m *mockContactService) DeleteContactPerson(ctx context.Context, id pgtype.UUID) error {
	if m.deleteContactPersonFunc != nil {
		return m.deleteContactPersonFunc(ctx, id)
	}
	return nil
}

func TestContactHandler_Flow(t *testing.T) {
	app := fiber.New()
	var partyID pgtype.UUID
	partyID.Scan("550e8400-e29b-41d4-a716-446655440000")

	mockSvc := &mockContactService{
		createContactFunc: func(ctx context.Context, input domain.CreateContactInput) (db.PartyContact, error) {
			return db.PartyContact{
				ID:           partyID,
				PartyID:      partyID,
				ContactType:  input.ContactType,
				ContactValue: input.ContactValue,
			}, nil
		},
		createContactPersonFunc: func(ctx context.Context, input domain.CreateContactPersonInput) (db.PartyContactPerson, error) {
			return db.PartyContactPerson{
				ID:        partyID,
				PartyID:   partyID,
				FirstName: input.FirstName,
			}, nil
		},
	}
	handler := NewContactHandler(mockSvc)
	app.Post("/api/v1/parties/:party_id/contacts", handler.CreateContact)
	app.Post("/api/v1/parties/:party_id/contact-persons", handler.CreateContactPerson)

	// Create contact
	bodyContact, _ := json.Marshal(dto.CreateContactRequest{
		ContactType:  "EMAIL",
		ContactValue: "test@example.com",
	})
	req := httptest.NewRequest("POST", "/api/v1/parties/550e8400-e29b-41d4-a716-446655440000/contacts", bytes.NewReader(bodyContact))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	// Create contact person
	bodyPerson, _ := json.Marshal(dto.CreateContactPersonRequest{
		FirstName: "John",
	})
	reqPerson := httptest.NewRequest("POST", "/api/v1/parties/550e8400-e29b-41d4-a716-446655440000/contact-persons", bytes.NewReader(bodyPerson))
	reqPerson.Header.Set("Content-Type", "application/json")
	respPerson, _ := app.Test(reqPerson)
	if respPerson.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", respPerson.StatusCode)
	}
}
