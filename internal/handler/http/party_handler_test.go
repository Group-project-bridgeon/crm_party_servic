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
	"github.com/group-project/crm_party_service/internal/service"
)

type mockPartyService struct {
	createPartyFunc       func(ctx context.Context, input domain.CreatePartyInput) (db.Party, error)
	getPartyByIDFunc      func(ctx context.Context, id pgtype.UUID) (db.Party, error)
	listPartiesFunc       func(ctx context.Context, params db.ListPartiesParams) ([]db.Party, error)
	updatePartyFunc       func(ctx context.Context, input domain.UpdatePartyInput) (db.Party, error)
	updatePartyStatusFunc func(ctx context.Context, id pgtype.UUID, isActive bool) (db.Party, error)
}

func (m *mockPartyService) CreateParty(ctx context.Context, input domain.CreatePartyInput) (db.Party, error) {
	if m.createPartyFunc != nil {
		return m.createPartyFunc(ctx, input)
	}
	return db.Party{}, nil
}

func (m *mockPartyService) GetPartyByID(ctx context.Context, id pgtype.UUID) (db.Party, error) {
	if m.getPartyByIDFunc != nil {
		return m.getPartyByIDFunc(ctx, id)
	}
	return db.Party{}, nil
}

func (m *mockPartyService) ListParties(ctx context.Context, params db.ListPartiesParams) ([]db.Party, error) {
	if m.listPartiesFunc != nil {
		return m.listPartiesFunc(ctx, params)
	}
	return []db.Party{}, nil
}

func (m *mockPartyService) UpdateParty(ctx context.Context, input domain.UpdatePartyInput) (db.Party, error) {
	if m.updatePartyFunc != nil {
		return m.updatePartyFunc(ctx, input)
	}
	return db.Party{}, nil
}

func (m *mockPartyService) UpdatePartyStatus(ctx context.Context, id pgtype.UUID, isActive bool) (db.Party, error) {
	if m.updatePartyStatusFunc != nil {
		return m.updatePartyStatusFunc(ctx, id, isActive)
	}
	return db.Party{}, nil
}

func TestPartyHandler_Flow(t *testing.T) {
	app := fiber.New()
	var partyID pgtype.UUID
	partyID.Scan("550e8400-e29b-41d4-a716-446655440000")

	mockSvc := &mockPartyService{
		createPartyFunc: func(ctx context.Context, input domain.CreatePartyInput) (db.Party, error) {
			if input.DisplayName == "" {
				return db.Party{}, service.ErrEmptyDisplayName
			}
			return db.Party{
				ID:          partyID,
				PartyType:   string(input.PartyType),
				DisplayName: input.DisplayName,
				IsActive:    true,
			}, nil
		},
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{
				ID:          id,
				PartyType:   "ORGANIZATION",
				DisplayName: "Acme Corp",
				IsActive:    true,
			}, nil
		},
		listPartiesFunc: func(ctx context.Context, params db.ListPartiesParams) ([]db.Party, error) {
			return []db.Party{
				{ID: partyID, DisplayName: "Acme Corp", IsActive: true},
			}, nil
		},
		updatePartyFunc: func(ctx context.Context, input domain.UpdatePartyInput) (db.Party, error) {
			return db.Party{
				ID:          partyID,
				DisplayName: input.DisplayName,
				IsActive:    true,
			}, nil
		},
		updatePartyStatusFunc: func(ctx context.Context, id pgtype.UUID, isActive bool) (db.Party, error) {
			return db.Party{
				ID:       id,
				IsActive: isActive,
			}, nil
		},
	}
	handler := NewPartyHandler(mockSvc)
	app.Post("/api/v1/parties", handler.CreateParty)
	app.Get("/api/v1/parties", handler.ListParties)
	app.Get("/api/v1/parties/:id", handler.GetPartyByID)
	app.Put("/api/v1/parties/:id", handler.UpdateParty)
	app.Patch("/api/v1/parties/:id/status", handler.UpdatePartyStatus)

	// 1. Create Party
	body, _ := json.Marshal(dto.CreatePartyRequest{
		PartyType:   "ORGANIZATION",
		DisplayName: "Acme Corp",
	})
	req := httptest.NewRequest("POST", "/api/v1/parties", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	// 2. Get Party
	reqGet := httptest.NewRequest("GET", "/api/v1/parties/550e8400-e29b-41d4-a716-446655440000", nil)
	respGet, _ := app.Test(reqGet)
	if respGet.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", respGet.StatusCode)
	}

	// 3. List Parties
	reqList := httptest.NewRequest("GET", "/api/v1/parties?limit=10&offset=0", nil)
	respList, _ := app.Test(reqList)
	if respList.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", respList.StatusCode)
	}

	// 4. Update Party
	updateBody, _ := json.Marshal(dto.UpdatePartyRequest{
		DisplayName: "Acme Global Corp",
	})
	reqPut := httptest.NewRequest("PUT", "/api/v1/parties/550e8400-e29b-41d4-a716-446655440000", bytes.NewReader(updateBody))
	reqPut.Header.Set("Content-Type", "application/json")
	respPut, _ := app.Test(reqPut)
	if respPut.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", respPut.StatusCode)
	}

	// 5. Update Status
	statusBody, _ := json.Marshal(dto.UpdatePartyStatusRequest{
		IsActive: false,
	})
	reqPatch := httptest.NewRequest("PATCH", "/api/v1/parties/550e8400-e29b-41d4-a716-446655440000/status", bytes.NewReader(statusBody))
	reqPatch.Header.Set("Content-Type", "application/json")
	respPatch, _ := app.Test(reqPatch)
	if respPatch.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", respPatch.StatusCode)
	}
}
