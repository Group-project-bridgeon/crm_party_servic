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

type mockAddressRepository struct {
	createPartyAddressCalled        bool
	createPartyAddressParams        db.CreatePartyAddressParams
	getPartyAddressByIDFunc         func(ctx context.Context, id pgtype.UUID) (db.PartyAddress, error)
	listPartyAddressesByPartyIDFunc func(ctx context.Context, params db.ListPartyAddressesByPartyIDParams) ([]db.PartyAddress, error)
	updatePartyAddressFunc          func(ctx context.Context, params db.UpdatePartyAddressParams) (db.PartyAddress, error)
	resetPrimaryPartyAddressFunc    func(ctx context.Context, params db.ResetPrimaryPartyAddressParams) error
	deletePartyAddressFunc          func(ctx context.Context, id pgtype.UUID) error
}

func (m *mockAddressRepository) CreatePartyAddress(ctx context.Context, params db.CreatePartyAddressParams) (db.PartyAddress, error) {
	m.createPartyAddressCalled = true
	m.createPartyAddressParams = params
	return db.PartyAddress{
		ID:           pgtype.UUID{Valid: true},
		PartyID:      params.PartyID,
		AddressType:  params.AddressType,
		AddressName:  params.AddressName,
		AddressLine1: params.AddressLine1,
		City:         params.City,
		CountryCode:  params.CountryCode,
		IsPrimary:    params.IsPrimary,
		IsActive:     params.IsActive,
	}, nil
}

func (m *mockAddressRepository) GetPartyAddressByID(ctx context.Context, id pgtype.UUID) (db.PartyAddress, error) {
	if m.getPartyAddressByIDFunc != nil {
		return m.getPartyAddressByIDFunc(ctx, id)
	}
	return db.PartyAddress{}, pgx.ErrNoRows
}

func (m *mockAddressRepository) ListPartyAddressesByPartyID(ctx context.Context, params db.ListPartyAddressesByPartyIDParams) ([]db.PartyAddress, error) {
	if m.listPartyAddressesByPartyIDFunc != nil {
		return m.listPartyAddressesByPartyIDFunc(ctx, params)
	}
	return []db.PartyAddress{}, nil
}

func (m *mockAddressRepository) UpdatePartyAddress(ctx context.Context, params db.UpdatePartyAddressParams) (db.PartyAddress, error) {
	if m.updatePartyAddressFunc != nil {
		return m.updatePartyAddressFunc(ctx, params)
	}
	return db.PartyAddress{
		ID:           params.ID,
		AddressType:  params.AddressType,
		AddressLine1: params.AddressLine1,
		City:         params.City,
		CountryCode:  params.CountryCode,
		IsPrimary:    params.IsPrimary,
	}, nil
}

func (m *mockAddressRepository) ResetPrimaryPartyAddress(ctx context.Context, params db.ResetPrimaryPartyAddressParams) error {
	if m.resetPrimaryPartyAddressFunc != nil {
		return m.resetPrimaryPartyAddressFunc(ctx, params)
	}
	return nil
}

func (m *mockAddressRepository) DeletePartyAddress(ctx context.Context, id pgtype.UUID) error {
	if m.deletePartyAddressFunc != nil {
		return m.deletePartyAddressFunc(ctx, id)
	}
	return nil
}

func TestCreateAddress_Success(t *testing.T) {
	partyUUID := "550e8400-e29b-41d4-a716-446655440000"
	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{ID: id}, nil
		},
	}
	addrRepo := &mockAddressRepository{}
	svc := NewAddressService(partyRepo, addrRepo)

	lat := 37.7749
	lon := -122.4194
	input := domain.CreateAddressInput{
		PartyID:      partyUUID,
		AddressType:  "BILLING",
		AddressLine1: "123 Market St",
		City:         "San Francisco",
		CountryCode:  "US",
		Latitude:     &lat,
		Longitude:    &lon,
		IsPrimary:    true,
	}

	addr, err := svc.CreateAddress(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr.City != "San Francisco" {
		t.Fatalf("expected San Francisco, got %s", addr.City)
	}
	if !addrRepo.createPartyAddressCalled {
		t.Fatal("expected CreatePartyAddress to be called")
	}
}

func TestCreateAddress_ValidationErrors(t *testing.T) {
	partyUUID := "550e8400-e29b-41d4-a716-446655440000"
	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{ID: id}, nil
		},
	}
	addrRepo := &mockAddressRepository{}
	svc := NewAddressService(partyRepo, addrRepo)

	// Invalid type
	_, err := svc.CreateAddress(context.Background(), domain.CreateAddressInput{
		PartyID:      partyUUID,
		AddressType:  "INVALID_TYPE",
		AddressLine1: "Line 1",
		City:         "City",
		CountryCode:  "US",
	})
	if !errors.Is(err, ErrInvalidAddressType) {
		t.Fatalf("expected ErrInvalidAddressType, got %v", err)
	}

	// Empty line 1
	_, err = svc.CreateAddress(context.Background(), domain.CreateAddressInput{
		PartyID:      partyUUID,
		AddressType:  "BILLING",
		AddressLine1: "   ",
		City:         "City",
		CountryCode:  "US",
	})
	if !errors.Is(err, ErrEmptyAddressLine1) {
		t.Fatalf("expected ErrEmptyAddressLine1, got %v", err)
	}

	// Empty city
	_, err = svc.CreateAddress(context.Background(), domain.CreateAddressInput{
		PartyID:      partyUUID,
		AddressType:  "BILLING",
		AddressLine1: "Line 1",
		City:         "   ",
		CountryCode:  "US",
	})
	if !errors.Is(err, ErrEmptyCity) {
		t.Fatalf("expected ErrEmptyCity, got %v", err)
	}

	// Invalid country code
	_, err = svc.CreateAddress(context.Background(), domain.CreateAddressInput{
		PartyID:      partyUUID,
		AddressType:  "BILLING",
		AddressLine1: "Line 1",
		City:         "City",
		CountryCode:  "USA",
	})
	if !errors.Is(err, ErrInvalidCountryCode) {
		t.Fatalf("expected ErrInvalidCountryCode, got %v", err)
	}

	// Invalid coordinates
	badLat := 95.0
	_, err = svc.CreateAddress(context.Background(), domain.CreateAddressInput{
		PartyID:      partyUUID,
		AddressType:  "BILLING",
		AddressLine1: "Line 1",
		City:         "City",
		CountryCode:  "US",
		Latitude:     &badLat,
	})
	if !errors.Is(err, ErrInvalidCoordinates) {
		t.Fatalf("expected ErrInvalidCoordinates, got %v", err)
	}
}

func TestAddress_GetUpdateDelete(t *testing.T) {
	addrUUID := "550e8400-e29b-41d4-a716-446655440000"
	var aID pgtype.UUID
	aID.Scan(addrUUID)

	addrRepo := &mockAddressRepository{
		getPartyAddressByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.PartyAddress, error) {
			return db.PartyAddress{ID: id, City: "Chicago"}, nil
		},
	}
	svc := NewAddressService(&mockPartyRepository{}, addrRepo)

	// Get
	addr, err := svc.GetAddressByID(context.Background(), aID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr.City != "Chicago" {
		t.Fatalf("expected Chicago, got %s", addr.City)
	}

	// Update
	updated, err := svc.UpdateAddress(context.Background(), domain.UpdateAddressInput{
		ID:           addrUUID,
		AddressType:  "SHIPPING",
		AddressLine1: "New St",
		City:         "Austin",
		CountryCode:  "US",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.City != "Austin" {
		t.Fatalf("expected Austin, got %s", updated.City)
	}

	// Delete
	err = svc.DeleteAddress(context.Background(), aID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

//helooo