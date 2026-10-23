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

type mockSupplierService struct {
	createSupplierFunc                  func(ctx context.Context, input domain.CreateSupplierInput) (db.Supplier, error)
	getSupplierByIDFunc                 func(ctx context.Context, id pgtype.UUID) (db.Supplier, error)
	getSupplierByPartyIDFunc            func(ctx context.Context, partyID pgtype.UUID) (db.Supplier, error)
	listSuppliersFunc                   func(ctx context.Context, params db.ListSuppliersParams) ([]db.Supplier, error)
	updateSupplierFunc                  func(ctx context.Context, input domain.UpdateSupplierInput) (db.Supplier, error)
	createSupplierUserAccountFunc       func(ctx context.Context, input domain.CreateSupplierUserAccountInput) (db.SupplierUserAccount, error)
	listSupplierUserAccountsFunc       func(ctx context.Context, supplierID pgtype.UUID) ([]db.SupplierUserAccount, error)
	updateSupplierUserAccountStatusFunc func(ctx context.Context, id pgtype.UUID, isActive bool) (db.SupplierUserAccount, error)
}

func (m *mockSupplierService) CreateSupplier(ctx context.Context, input domain.CreateSupplierInput) (db.Supplier, error) {
	if m.createSupplierFunc != nil {
		return m.createSupplierFunc(ctx, input)
	}
	return db.Supplier{}, nil
}

func (m *mockSupplierService) GetSupplierByID(ctx context.Context, id pgtype.UUID) (db.Supplier, error) {
	if m.getSupplierByIDFunc != nil {
		return m.getSupplierByIDFunc(ctx, id)
	}
	return db.Supplier{}, nil
}

func (m *mockSupplierService) GetSupplierByPartyID(ctx context.Context, partyID pgtype.UUID) (db.Supplier, error) {
	if m.getSupplierByPartyIDFunc != nil {
		return m.getSupplierByPartyIDFunc(ctx, partyID)
	}
	return db.Supplier{}, nil
}

func (m *mockSupplierService) ListSuppliers(ctx context.Context, params db.ListSuppliersParams) ([]db.Supplier, error) {
	if m.listSuppliersFunc != nil {
		return m.listSuppliersFunc(ctx, params)
	}
	return []db.Supplier{}, nil
}

func (m *mockSupplierService) UpdateSupplier(ctx context.Context, input domain.UpdateSupplierInput) (db.Supplier, error) {
	if m.updateSupplierFunc != nil {
		return m.updateSupplierFunc(ctx, input)
	}
	return db.Supplier{}, nil
}

func (m *mockSupplierService) CreateSupplierUserAccount(ctx context.Context, input domain.CreateSupplierUserAccountInput) (db.SupplierUserAccount, error) {
	if m.createSupplierUserAccountFunc != nil {
		return m.createSupplierUserAccountFunc(ctx, input)
	}
	return db.SupplierUserAccount{}, nil
}

func (m *mockSupplierService) ListSupplierUserAccounts(ctx context.Context, supplierID pgtype.UUID) ([]db.SupplierUserAccount, error) {
	if m.listSupplierUserAccountsFunc != nil {
		return m.listSupplierUserAccountsFunc(ctx, supplierID)
	}
	return []db.SupplierUserAccount{}, nil
}

func (m *mockSupplierService) UpdateSupplierUserAccountStatus(ctx context.Context, id pgtype.UUID, isActive bool) (db.SupplierUserAccount, error) {
	if m.updateSupplierUserAccountStatusFunc != nil {
		return m.updateSupplierUserAccountStatusFunc(ctx, id, isActive)
	}
	return db.SupplierUserAccount{}, nil
}

func TestSupplierHandler_CreateSupplier(t *testing.T) {
	app := fiber.New()
	var partyID pgtype.UUID
	partyID.Scan("550e8400-e29b-41d4-a716-446655440000")

	mockSvc := &mockSupplierService{
		createSupplierFunc: func(ctx context.Context, input domain.CreateSupplierInput) (db.Supplier, error) {
			if input.SupplierCode == "EXISTS" {
				return db.Supplier{}, service.ErrDuplicateSupplierCode
			}
			return db.Supplier{
				ID:           partyID,
				PartyID:      partyID,
				SupplierCode: input.SupplierCode,
			}, nil
		},
	}
	handler := NewSupplierHandler(mockSvc)
	app.Post("/api/v1/suppliers", handler.CreateSupplier)

	// Success
	body, _ := json.Marshal(dto.CreateSupplierRequest{
		PartyID:      "550e8400-e29b-41d4-a716-446655440000",
		SupplierCode: "SUPP-1",
	})
	req := httptest.NewRequest("POST", "/api/v1/suppliers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	// Conflict
	body2, _ := json.Marshal(dto.CreateSupplierRequest{
		PartyID:      "550e8400-e29b-41d4-a716-446655440000",
		SupplierCode: "EXISTS",
	})
	req2 := httptest.NewRequest("POST", "/api/v1/suppliers", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409, got %d", resp2.StatusCode)
	}
}

func TestSupplierHandler_GetSupplierByID(t *testing.T) {
	app := fiber.New()
	var suppID pgtype.UUID
	suppID.Scan("550e8400-e29b-41d4-a716-446655440000")

	mockSvc := &mockSupplierService{
		getSupplierByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Supplier, error) {
			return db.Supplier{ID: id, SupplierCode: "SUPP-10"}, nil
		},
	}
	handler := NewSupplierHandler(mockSvc)
	app.Get("/api/v1/suppliers/:id", handler.GetSupplierByID)

	req := httptest.NewRequest("GET", "/api/v1/suppliers/550e8400-e29b-41d4-a716-446655440000", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestSupplierHandler_UserAccounts(t *testing.T) {
	app := fiber.New()
	var suppID pgtype.UUID
	suppID.Scan("550e8400-e29b-41d4-a716-446655440000")

	mockSvc := &mockSupplierService{
		createSupplierUserAccountFunc: func(ctx context.Context, input domain.CreateSupplierUserAccountInput) (db.SupplierUserAccount, error) {
			return db.SupplierUserAccount{
				ID:         suppID,
				SupplierID: suppID,
				IsActive:   true,
			}, nil
		},
		listSupplierUserAccountsFunc: func(ctx context.Context, supplierID pgtype.UUID) ([]db.SupplierUserAccount, error) {
			return []db.SupplierUserAccount{
				{ID: suppID, SupplierID: supplierID, IsActive: true},
			}, nil
		},
		updateSupplierUserAccountStatusFunc: func(ctx context.Context, id pgtype.UUID, isActive bool) (db.SupplierUserAccount, error) {
			return db.SupplierUserAccount{ID: id, IsActive: isActive}, nil
		},
	}
	handler := NewSupplierHandler(mockSvc)
	app.Post("/api/v1/suppliers/:id/user-accounts", handler.CreateSupplierUserAccount)
	app.Get("/api/v1/suppliers/:id/user-accounts", handler.ListSupplierUserAccounts)
	app.Patch("/api/v1/suppliers/user-accounts/:account_id/status", handler.UpdateSupplierUserAccountStatus)

	// Create user account
	body, _ := json.Marshal(dto.CreateSupplierUserAccountRequest{
		UserID: "660e8400-e29b-41d4-a716-446655440000",
	})
	req := httptest.NewRequest("POST", "/api/v1/suppliers/550e8400-e29b-41d4-a716-446655440000/user-accounts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	// List
	reqList := httptest.NewRequest("GET", "/api/v1/suppliers/550e8400-e29b-41d4-a716-446655440000/user-accounts", nil)
	respList, _ := app.Test(reqList)
	if respList.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", respList.StatusCode)
	}

	// Update status
	bodyStatus, _ := json.Marshal(dto.UpdateUserAccountStatusRequest{IsActive: false})
	reqPatch := httptest.NewRequest("PATCH", "/api/v1/suppliers/user-accounts/550e8400-e29b-41d4-a716-446655440000/status", bytes.NewReader(bodyStatus))
	reqPatch.Header.Set("Content-Type", "application/json")
	respPatch, _ := app.Test(reqPatch)
	if respPatch.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", respPatch.StatusCode)
	}
}
