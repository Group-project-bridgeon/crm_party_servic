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

type mockCompositeService struct {
	getParty360Func             func(ctx context.Context, partyID pgtype.UUID) (*domain.Party360, error)
	createCompositeCustomerFunc func(ctx context.Context, input domain.CompositeCustomerInput) (*domain.CompositeCustomerResult, error)
	createCompositeSupplierFunc func(ctx context.Context, input domain.CompositeSupplierInput) (*domain.CompositeSupplierResult, error)
}

func (m *mockCompositeService) GetParty360(ctx context.Context, partyID pgtype.UUID) (*domain.Party360, error) {
	if m.getParty360Func != nil {
		return m.getParty360Func(ctx, partyID)
	}
	return &domain.Party360{}, nil
}

func (m *mockCompositeService) CreateCompositeCustomer(ctx context.Context, input domain.CompositeCustomerInput) (*domain.CompositeCustomerResult, error) {
	if m.createCompositeCustomerFunc != nil {
		return m.createCompositeCustomerFunc(ctx, input)
	}
	return &domain.CompositeCustomerResult{}, nil
}

func (m *mockCompositeService) CreateCompositeSupplier(ctx context.Context, input domain.CompositeSupplierInput) (*domain.CompositeSupplierResult, error) {
	if m.createCompositeSupplierFunc != nil {
		return m.createCompositeSupplierFunc(ctx, input)
	}
	return &domain.CompositeSupplierResult{}, nil
}

func TestCompositeHandler_GetParty360(t *testing.T) {
	app := fiber.New()
	var partyID pgtype.UUID
	partyID.Scan("550e8400-e29b-41d4-a716-446655440000")

	mockSvc := &mockCompositeService{
		getParty360Func: func(ctx context.Context, id pgtype.UUID) (*domain.Party360, error) {
			return &domain.Party360{
				Party: db.Party{
					ID:          id,
					PartyType:   "ORGANIZATION",
					DisplayName: "Star Logistics",
				},
			}, nil
		},
	}
	handler := NewCompositeHandler(mockSvc)
	app.Get("/api/v1/parties/:id/360", handler.GetParty360)

	// Success
	req := httptest.NewRequest("GET", "/api/v1/parties/550e8400-e29b-41d4-a716-446655440000/360", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// Not Found
	mockSvc.getParty360Func = func(ctx context.Context, id pgtype.UUID) (*domain.Party360, error) {
		return nil, service.ErrPartyNotFound
	}
	respNotFound, _ := app.Test(req)
	if respNotFound.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", respNotFound.StatusCode)
	}
}

func TestCompositeHandler_Onboarding(t *testing.T) {
	app := fiber.New()
	var partyID pgtype.UUID
	partyID.Scan("550e8400-e29b-41d4-a716-446655440000")

	mockSvc := &mockCompositeService{
		createCompositeCustomerFunc: func(ctx context.Context, input domain.CompositeCustomerInput) (*domain.CompositeCustomerResult, error) {
			if input.CustomerCode == "EXISTS" {
				return nil, service.ErrDuplicateCustomerCode
			}
			return &domain.CompositeCustomerResult{
				Party:    db.Party{ID: partyID, DisplayName: input.DisplayName},
				Customer: db.Customer{ID: partyID, CustomerCode: input.CustomerCode},
			}, nil
		},
		createCompositeSupplierFunc: func(ctx context.Context, input domain.CompositeSupplierInput) (*domain.CompositeSupplierResult, error) {
			return &domain.CompositeSupplierResult{
				Party:    db.Party{ID: partyID, DisplayName: input.DisplayName},
				Supplier: db.Supplier{ID: partyID, SupplierCode: input.SupplierCode},
			}, nil
		},
	}
	handler := NewCompositeHandler(mockSvc)
	app.Post("/api/v1/customers/composite", handler.CreateCompositeCustomer)
	app.Post("/api/v1/suppliers/composite", handler.CreateCompositeSupplier)

	// Customer onboarding success
	custBody, _ := json.Marshal(dto.CreateCompositeCustomerRequest{
		PartyType:    "ORGANIZATION",
		DisplayName:  "Prime Retail",
		CustomerCode: "CUST-PRIME",
	})
	reqCust := httptest.NewRequest("POST", "/api/v1/customers/composite", bytes.NewReader(custBody))
	reqCust.Header.Set("Content-Type", "application/json")
	respCust, _ := app.Test(reqCust)
	if respCust.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", respCust.StatusCode)
	}

	// Supplier onboarding success
	suppBody, _ := json.Marshal(dto.CreateCompositeSupplierRequest{
		PartyType:    "ORGANIZATION",
		DisplayName:  "Global Parts",
		SupplierCode: "SUPP-PARTS",
	})
	reqSupp := httptest.NewRequest("POST", "/api/v1/suppliers/composite", bytes.NewReader(suppBody))
	reqSupp.Header.Set("Content-Type", "application/json")
	respSupp, _ := app.Test(reqSupp)
	if respSupp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", respSupp.StatusCode)
	}
}
