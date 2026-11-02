package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
	"github.com/group-project/crm_party_service/internal/domain"
)

type mockContactRepository struct {
	createPartyContactCalled       bool
	createPartyContactParams       db.CreatePartyContactParams
	getPartyContactByIDFunc        func(ctx context.Context, id pgtype.UUID) (db.PartyContact, error)
	listPartyContactsByPartyIDFunc func(ctx context.Context, params db.ListPartyContactsByPartyIDParams) ([]db.PartyContact, error)
	updatePartyContactFunc         func(ctx context.Context, params db.UpdatePartyContactParams) (db.PartyContact, error)
	deletePartyContactFunc         func(ctx context.Context, id pgtype.UUID) error

	createPartyContactPersonCalled       bool
	createPartyContactPersonParams       db.CreatePartyContactPersonParams
	getPartyContactPersonByIDFunc        func(ctx context.Context, id pgtype.UUID) (db.PartyContactPerson, error)
	listPartyContactPersonsByPartyIDFunc func(ctx context.Context, partyID pgtype.UUID) ([]db.PartyContactPerson, error)
	updatePartyContactPersonFunc         func(ctx context.Context, params db.UpdatePartyContactPersonParams) (db.PartyContactPerson, error)
	deletePartyContactPersonFunc         func(ctx context.Context, id pgtype.UUID) error
}

func (m *mockContactRepository) CreatePartyContact(ctx context.Context, params db.CreatePartyContactParams) (db.PartyContact, error) {
	m.createPartyContactCalled = true
	m.createPartyContactParams = params
	return db.PartyContact{
		ID:           pgtype.UUID{Valid: true},
		PartyID:      params.PartyID,
		ContactType:  params.ContactType,
		ContactValue: params.ContactValue,
	}, nil
}

func (m *mockContactRepository) GetPartyContactByID(ctx context.Context, id pgtype.UUID) (db.PartyContact, error) {
	if m.getPartyContactByIDFunc != nil {
		return m.getPartyContactByIDFunc(ctx, id)
	}
	return db.PartyContact{}, pgx.ErrNoRows
}

func (m *mockContactRepository) ListPartyContactsByPartyID(ctx context.Context, params db.ListPartyContactsByPartyIDParams) ([]db.PartyContact, error) {
	if m.listPartyContactsByPartyIDFunc != nil {
		return m.listPartyContactsByPartyIDFunc(ctx, params)
	}
	return []db.PartyContact{}, nil
}

func (m *mockContactRepository) UpdatePartyContact(ctx context.Context, params db.UpdatePartyContactParams) (db.PartyContact, error) {
	if m.updatePartyContactFunc != nil {
		return m.updatePartyContactFunc(ctx, params)
	}
	return db.PartyContact{
		ID:           params.ID,
		ContactType:  params.ContactType,
		ContactValue: params.ContactValue,
	}, nil
}

func (m *mockContactRepository) ResetPrimaryPartyContact(ctx context.Context, params db.ResetPrimaryPartyContactParams) error {
	return nil
}

func (m *mockContactRepository) DeletePartyContact(ctx context.Context, id pgtype.UUID) error {
	if m.deletePartyContactFunc != nil {
		return m.deletePartyContactFunc(ctx, id)
	}
	return nil
}

func (m *mockContactRepository) CreatePartyContactPerson(ctx context.Context, params db.CreatePartyContactPersonParams) (db.PartyContactPerson, error) {
	m.createPartyContactPersonCalled = true
	m.createPartyContactPersonParams = params
	return db.PartyContactPerson{
		ID:        pgtype.UUID{Valid: true},
		PartyID:   params.PartyID,
		FirstName: params.FirstName,
		LastName:  params.LastName,
	}, nil
}

func (m *mockContactRepository) GetPartyContactPersonByID(ctx context.Context, id pgtype.UUID) (db.PartyContactPerson, error) {
	if m.getPartyContactPersonByIDFunc != nil {
		return m.getPartyContactPersonByIDFunc(ctx, id)
	}
	return db.PartyContactPerson{}, pgx.ErrNoRows
}

func (m *mockContactRepository) ListPartyContactPersonsByPartyID(ctx context.Context, partyID pgtype.UUID) ([]db.PartyContactPerson, error) {
	if m.listPartyContactPersonsByPartyIDFunc != nil {
		return m.listPartyContactPersonsByPartyIDFunc(ctx, partyID)
	}
	return []db.PartyContactPerson{}, nil
}

func (m *mockContactRepository) UpdatePartyContactPerson(ctx context.Context, params db.UpdatePartyContactPersonParams) (db.PartyContactPerson, error) {
	if m.updatePartyContactPersonFunc != nil {
		return m.updatePartyContactPersonFunc(ctx, params)
	}
	return db.PartyContactPerson{
		ID:        params.ID,
		FirstName: params.FirstName,
	}, nil
}

func (m *mockContactRepository) ResetPrimaryPartyContactPerson(ctx context.Context, partyID pgtype.UUID) error {
	return nil
}

func (m *mockContactRepository) DeletePartyContactPerson(ctx context.Context, id pgtype.UUID) error {
	if m.deletePartyContactPersonFunc != nil {
		return m.deletePartyContactPersonFunc(ctx, id)
	}
	return nil
}

func TestDirectContact_Flow(t *testing.T) {
	partyUUID := "550e8400-e29b-41d4-a716-446655440000"
	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{ID: id}, nil
		},
	}
	contactRepo := &mockContactRepository{}
	svc := NewContactService(partyRepo, contactRepo)

	// Create
	contact, err := svc.CreateContact(context.Background(), domain.CreateContactInput{
		PartyID:      partyUUID,
		ContactType:  "EMAIL",
		ContactValue: "sales@example.com",
		IsPrimary:    true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if contact.ContactValue != "sales@example.com" {
		t.Fatalf("expected sales@example.com, got %s", contact.ContactValue)
	}

	// Validation: Invalid type
	_, err = svc.CreateContact(context.Background(), domain.CreateContactInput{
		PartyID:      partyUUID,
		ContactType:  "PAGER",
		ContactValue: "1234",
	})
	if !errors.Is(err, ErrInvalidContactType) {
		t.Fatalf("expected ErrInvalidContactType, got %v", err)
	}
}

func TestContactPerson_Flow(t *testing.T) {
	partyUUID := "550e8400-e29b-41d4-a716-446655440000"
	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{ID: id}, nil
		},
	}
	contactRepo := &mockContactRepository{}
	svc := NewContactService(partyRepo, contactRepo)

	// Create
	person, err := svc.CreateContactPerson(context.Background(), domain.CreateContactPersonInput{
		PartyID:   partyUUID,
		FirstName: "Robert",
		IsPrimary: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if person.FirstName != "Robert" {
		t.Fatalf("expected Robert, got %s", person.FirstName)
	}

	// Validation: Empty first name
	_, err = svc.CreateContactPerson(context.Background(), domain.CreateContactPersonInput{
		PartyID:   partyUUID,
		FirstName: "   ",
	})
	if !errors.Is(err, ErrEmptyFirstName) {
		t.Fatalf("expected ErrEmptyFirstName, got %v", err)
	}
}
