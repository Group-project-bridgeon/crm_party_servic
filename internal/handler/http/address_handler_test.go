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

type mockAddressService struct {
	createAddressFunc          func(ctx context.Context, input domain.CreateAddressInput) (db.PartyAddress, error)
	getAddressByIDFunc         func(ctx context.Context, id pgtype.UUID) (db.PartyAddress, error)
	listAddressesByPartyIDFunc func(ctx context.Context, partyID pgtype.UUID, addressType *string) ([]db.PartyAddress, error)
	updateAddressFunc          func(ctx context.Context, input domain.UpdateAddressInput) (db.PartyAddress, error)
	deleteAddressFunc          func(ctx context.Context, id pgtype.UUID) error
}

func (m *mockAddressService) CreateAddress(ctx context.Context, input domain.CreateAddressInput) (db.PartyAddress, error) {
	if m.createAddressFunc != nil {
		return m.createAddressFunc(ctx, input)
	}
	return db.PartyAddress{}, nil
}

func (m *mockAddressService) GetAddressByID(ctx context.Context, id pgtype.UUID) (db.PartyAddress, error) {
	if m.getAddressByIDFunc != nil {
		return m.getAddressByIDFunc(ctx, id)
	}
	return db.PartyAddress{}, nil
}

func (m *mockAddressService) ListAddressesByPartyID(ctx context.Context, partyID pgtype.UUID, addressType *string) ([]db.PartyAddress, error) {
	if m.listAddressesByPartyIDFunc != nil {
		return m.listAddressesByPartyIDFunc(ctx, partyID, addressType)
	}
	return []db.PartyAddress{}, nil
}

func (m *mockAddressService) UpdateAddress(ctx context.Context, input domain.UpdateAddressInput) (db.PartyAddress, error) {
	if m.updateAddressFunc != nil {
		return m.updateAddressFunc(ctx, input)
	}
	return db.PartyAddress{}, nil
}

func (m *mockAddressService) DeleteAddress(ctx context.Context, id pgtype.UUID) error {
	if m.deleteAddressFunc != nil {
		return m.deleteAddressFunc(ctx, id)
	}
	return nil
}

func TestAddressHandler_Flow(t *testing.T) {
	app := fiber.New()
	var partyID pgtype.UUID
	partyID.Scan("550e8400-e29b-41d4-a716-446655440000")

	mockSvc := &mockAddressService{
		createAddressFunc: func(ctx context.Context, input domain.CreateAddressInput) (db.PartyAddress, error) {
			return db.PartyAddress{
				ID:           partyID,
				PartyID:      partyID,
				AddressType:  input.AddressType,
				AddressLine1: input.AddressLine1,
				City:         input.City,
				CountryCode:  input.CountryCode,
			}, nil
		},
		listAddressesByPartyIDFunc: func(ctx context.Context, pID pgtype.UUID, addressType *string) ([]db.PartyAddress, error) {
			return []db.PartyAddress{
				{ID: partyID, PartyID: pID, City: "Chicago"},
			}, nil
		},
		deleteAddressFunc: func(ctx context.Context, id pgtype.UUID) error {
			return nil
		},
	}
	handler := NewAddressHandler(mockSvc)
	app.Post("/api/v1/parties/:party_id/addresses", handler.CreateAddress)
	app.Get("/api/v1/parties/:party_id/addresses", handler.ListAddressesByPartyID)
	app.Delete("/api/v1/addresses/:address_id", handler.DeleteAddress)

	// Create
	body, _ := json.Marshal(dto.CreateAddressRequest{
		AddressType:  "BILLING",
		AddressLine1: "123 Main St",
		City:         "Chicago",
		CountryCode:  "US",
	})
	req := httptest.NewRequest("POST", "/api/v1/parties/550e8400-e29b-41d4-a716-446655440000/addresses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	// List
	reqList := httptest.NewRequest("GET", "/api/v1/parties/550e8400-e29b-41d4-a716-446655440000/addresses", nil)
	respList, _ := app.Test(reqList)
	if respList.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", respList.StatusCode)
	}

	// Delete
	reqDel := httptest.NewRequest("DELETE", "/api/v1/addresses/550e8400-e29b-41d4-a716-446655440000", nil)
	respDel, _ := app.Test(reqDel)
	if respDel.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", respDel.StatusCode)
	}
}
