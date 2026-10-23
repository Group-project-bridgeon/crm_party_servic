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

type mockIndividualService struct {
	createIndividualFunc func(
		ctx context.Context,
		input domain.CreateIndividualInput,
	) (db.Individual, error)

	getIndividualByPartyIDFunc func(
		ctx context.Context,
		partyID pgtype.UUID,
	) (db.Individual, error)

	updateIndividualFunc func(
		ctx context.Context,
		input domain.UpdateIndividualInput,
	) (db.Individual, error)
}

func (m *mockIndividualService) CreateIndividual(
	ctx context.Context,
	input domain.CreateIndividualInput,
) (db.Individual, error) {
	if m.createIndividualFunc != nil {
		return m.createIndividualFunc(ctx, input)
	}
	return db.Individual{}, nil
}

func (m *mockIndividualService) GetIndividualByPartyID(
	ctx context.Context,
	partyID pgtype.UUID,
) (db.Individual, error) {
	if m.getIndividualByPartyIDFunc != nil {
		return m.getIndividualByPartyIDFunc(ctx, partyID)
	}
	return db.Individual{}, nil
}

func (m *mockIndividualService) UpdateIndividual(
	ctx context.Context,
	input domain.UpdateIndividualInput,
) (db.Individual, error) {
	if m.updateIndividualFunc != nil {
		return m.updateIndividualFunc(ctx, input)
	}
	return db.Individual{}, nil
}

func TestCreateIndividual_HandlerSuccess(t *testing.T) {
	app := fiber.New()
	var partyID pgtype.UUID
	partyID.Scan("550e8400-e29b-41d4-a716-446655440000")

	mockSvc := &mockIndividualService{
		createIndividualFunc: func(ctx context.Context, input domain.CreateIndividualInput) (db.Individual, error) {
			return db.Individual{
				ID:        partyID,
				PartyID:   partyID,
				FirstName: input.FirstName,
			}, nil
		},
	}
	handler := NewIndividualHandler(mockSvc)
	app.Post("/api/v1/individuals", handler.CreateIndividual)

	body, _ := json.Marshal(dto.CreateIndividualRequest{
		PartyID:   "550e8400-e29b-41d4-a716-446655440000",
		FirstName: "John",
	})
	req := httptest.NewRequest("POST", "/api/v1/individuals", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}

func TestCreateIndividual_HandlerConflict(t *testing.T) {
	app := fiber.New()
	mockSvc := &mockIndividualService{
		createIndividualFunc: func(ctx context.Context, input domain.CreateIndividualInput) (db.Individual, error) {
			return db.Individual{}, service.ErrIndividualAlreadyExists
		},
	}
	handler := NewIndividualHandler(mockSvc)
	app.Post("/api/v1/individuals", handler.CreateIndividual)

	body, _ := json.Marshal(dto.CreateIndividualRequest{
		PartyID:   "550e8400-e29b-41d4-a716-446655440000",
		FirstName: "John",
	})
	req := httptest.NewRequest("POST", "/api/v1/individuals", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409, got %d", resp.StatusCode)
	}
}

func TestGetIndividualByPartyID_HandlerSuccess(t *testing.T) {
	app := fiber.New()
	var partyID pgtype.UUID
	partyID.Scan("550e8400-e29b-41d4-a716-446655440000")

	mockSvc := &mockIndividualService{
		getIndividualByPartyIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Individual, error) {
			return db.Individual{
				ID:        id,
				PartyID:   id,
				FirstName: "Alice",
			}, nil
		},
	}
	handler := NewIndividualHandler(mockSvc)
	app.Get("/api/v1/individuals/:party_id", handler.GetIndividualByPartyID)

	req := httptest.NewRequest("GET", "/api/v1/individuals/550e8400-e29b-41d4-a716-446655440000", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestGetIndividualByPartyID_HandlerNotFound(t *testing.T) {
	app := fiber.New()
	mockSvc := &mockIndividualService{
		getIndividualByPartyIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Individual, error) {
			return db.Individual{}, service.ErrIndividualNotFound
		},
	}
	handler := NewIndividualHandler(mockSvc)
	app.Get("/api/v1/individuals/:party_id", handler.GetIndividualByPartyID)

	req := httptest.NewRequest("GET", "/api/v1/individuals/550e8400-e29b-41d4-a716-446655440000", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestUpdateIndividual_HandlerSuccess(t *testing.T) {
	app := fiber.New()
	var partyID pgtype.UUID
	partyID.Scan("550e8400-e29b-41d4-a716-446655440000")

	mockSvc := &mockIndividualService{
		updateIndividualFunc: func(ctx context.Context, input domain.UpdateIndividualInput) (db.Individual, error) {
			return db.Individual{
				ID:        partyID,
				PartyID:   partyID,
				FirstName: input.FirstName,
			}, nil
		},
	}
	handler := NewIndividualHandler(mockSvc)
	app.Put("/api/v1/individuals/:party_id", handler.UpdateIndividual)

	body, _ := json.Marshal(dto.UpdateIndividualRequest{
		FirstName: "UpdatedAlice",
	})
	req := httptest.NewRequest("PUT", "/api/v1/individuals/550e8400-e29b-41d4-a716-446655440000", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
