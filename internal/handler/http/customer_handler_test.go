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

type mockCustomerService struct {
	createCustomerFunc                  func(ctx context.Context, input domain.CreateCustomerInput) (db.Customer, error)
	getCustomerByIDFunc                 func(ctx context.Context, id pgtype.UUID) (db.Customer, error)
	getCustomerByPartyIDFunc            func(ctx context.Context, partyID pgtype.UUID) (db.Customer, error)
	listCustomersFunc                   func(ctx context.Context, params db.ListCustomersParams) ([]db.Customer, error)
	updateCustomerFunc                  func(ctx context.Context, input domain.UpdateCustomerInput) (db.Customer, error)
	createCustomerUserAccountFunc       func(ctx context.Context, input domain.CreateCustomerUserAccountInput) (db.CustomerUserAccount, error)
	listCustomerUserAccountsFunc        func(ctx context.Context, customerID pgtype.UUID) ([]db.CustomerUserAccount, error)
	updateCustomerUserAccountStatusFunc func(ctx context.Context, id pgtype.UUID, isActive bool) (db.CustomerUserAccount, error)
}

func (m *mockCustomerService) CreateCustomer(ctx context.Context, input domain.CreateCustomerInput) (db.Customer, error) {
	if m.createCustomerFunc != nil {
		return m.createCustomerFunc(ctx, input)
	}
	return db.Customer{}, nil
}

func (m *mockCustomerService) GetCustomerByID(ctx context.Context, id pgtype.UUID) (db.Customer, error) {
	if m.getCustomerByIDFunc != nil {
		return m.getCustomerByIDFunc(ctx, id)
	}
	return db.Customer{}, nil
}

func (m *mockCustomerService) GetCustomerByPartyID(ctx context.Context, partyID pgtype.UUID) (db.Customer, error) {
	if m.getCustomerByPartyIDFunc != nil {
		return m.getCustomerByPartyIDFunc(ctx, partyID)
	}
	return db.Customer{}, nil
}

func (m *mockCustomerService) ListCustomers(ctx context.Context, params db.ListCustomersParams) ([]db.Customer, error) {
	if m.listCustomersFunc != nil {
		return m.listCustomersFunc(ctx, params)
	}
	return []db.Customer{}, nil
}

func (m *mockCustomerService) UpdateCustomer(ctx context.Context, input domain.UpdateCustomerInput) (db.Customer, error) {
	if m.updateCustomerFunc != nil {
		return m.updateCustomerFunc(ctx, input)
	}
	return db.Customer{}, nil
}

func (m *mockCustomerService) CreateCustomerUserAccount(ctx context.Context, input domain.CreateCustomerUserAccountInput) (db.CustomerUserAccount, error) {
	if m.createCustomerUserAccountFunc != nil {
		return m.createCustomerUserAccountFunc(ctx, input)
	}
	return db.CustomerUserAccount{}, nil
}

func (m *mockCustomerService) ListCustomerUserAccounts(ctx context.Context, customerID pgtype.UUID) ([]db.CustomerUserAccount, error) {
	if m.listCustomerUserAccountsFunc != nil {
		return m.listCustomerUserAccountsFunc(ctx, customerID)
	}
	return []db.CustomerUserAccount{}, nil
}

func (m *mockCustomerService) UpdateCustomerUserAccountStatus(ctx context.Context, id pgtype.UUID, isActive bool) (db.CustomerUserAccount, error) {
	if m.updateCustomerUserAccountStatusFunc != nil {
		return m.updateCustomerUserAccountStatusFunc(ctx, id, isActive)
	}
	return db.CustomerUserAccount{}, nil
}

func TestCustomerHandler_CreateCustomer(t *testing.T) {
	app := fiber.New()
	var partyID pgtype.UUID
	partyID.Scan("550e8400-e29b-41d4-a716-446655440000")

	mockSvc := &mockCustomerService{
		createCustomerFunc: func(ctx context.Context, input domain.CreateCustomerInput) (db.Customer, error) {
			if input.CustomerCode == "EXISTS" {
				return db.Customer{}, service.ErrDuplicateCustomerCode
			}
			return db.Customer{
				ID:           partyID,
				PartyID:      partyID,
				CustomerCode: input.CustomerCode,
				CreditStatus: "ACTIVE",
			}, nil
		},
	}
	handler := NewCustomerHandler(mockSvc)
	app.Post("/api/v1/customers", handler.CreateCustomer)

	// Success
	body, _ := json.Marshal(dto.CreateCustomerRequest{
		PartyID:      "550e8400-e29b-41d4-a716-446655440000",
		CustomerCode: "CUST-1",
	})
	req := httptest.NewRequest("POST", "/api/v1/customers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	// Conflict
	body2, _ := json.Marshal(dto.CreateCustomerRequest{
		PartyID:      "550e8400-e29b-41d4-a716-446655440000",
		CustomerCode: "EXISTS",
	})
	req2 := httptest.NewRequest("POST", "/api/v1/customers", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409, got %d", resp2.StatusCode)
	}
}

func TestCustomerHandler_GetCustomerByID(t *testing.T) {
	app := fiber.New()
	var custID pgtype.UUID
	custID.Scan("550e8400-e29b-41d4-a716-446655440000")

	mockSvc := &mockCustomerService{
		getCustomerByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Customer, error) {
			return db.Customer{ID: id, CustomerCode: "CUST-10"}, nil
		},
	}
	handler := NewCustomerHandler(mockSvc)
	app.Get("/api/v1/customers/:id", handler.GetCustomerByID)

	req := httptest.NewRequest("GET", "/api/v1/customers/550e8400-e29b-41d4-a716-446655440000", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestCustomerHandler_UserAccounts(t *testing.T) {
	app := fiber.New()
	var custID pgtype.UUID
	custID.Scan("550e8400-e29b-41d4-a716-446655440000")

	mockSvc := &mockCustomerService{
		createCustomerUserAccountFunc: func(ctx context.Context, input domain.CreateCustomerUserAccountInput) (db.CustomerUserAccount, error) {
			return db.CustomerUserAccount{
				ID:         custID,
				CustomerID: custID,
				IsActive:   true,
			}, nil
		},
		listCustomerUserAccountsFunc: func(ctx context.Context, customerID pgtype.UUID) ([]db.CustomerUserAccount, error) {
			return []db.CustomerUserAccount{
				{ID: custID, CustomerID: customerID, IsActive: true},
			}, nil
		},
		updateCustomerUserAccountStatusFunc: func(ctx context.Context, id pgtype.UUID, isActive bool) (db.CustomerUserAccount, error) {
			return db.CustomerUserAccount{ID: id, IsActive: isActive}, nil
		},
	}
	handler := NewCustomerHandler(mockSvc)
	app.Post("/api/v1/customers/:id/user-accounts", handler.CreateCustomerUserAccount)
	app.Get("/api/v1/customers/:id/user-accounts", handler.ListCustomerUserAccounts)
	app.Patch("/api/v1/customers/user-accounts/:account_id/status", handler.UpdateCustomerUserAccountStatus)

	// Create user account
	body, _ := json.Marshal(dto.CreateCustomerUserAccountRequest{
		UserID: "660e8400-e29b-41d4-a716-446655440000",
	})
	req := httptest.NewRequest("POST", "/api/v1/customers/550e8400-e29b-41d4-a716-446655440000/user-accounts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	// List
	reqList := httptest.NewRequest("GET", "/api/v1/customers/550e8400-e29b-41d4-a716-446655440000/user-accounts", nil)
	respList, _ := app.Test(reqList)
	if respList.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", respList.StatusCode)
	}

	// Update status
	bodyStatus, _ := json.Marshal(dto.UpdateUserAccountStatusRequest{IsActive: false})
	reqPatch := httptest.NewRequest("PATCH", "/api/v1/customers/user-accounts/550e8400-e29b-41d4-a716-446655440000/status", bytes.NewReader(bodyStatus))
	reqPatch.Header.Set("Content-Type", "application/json")
	respPatch, _ := app.Test(reqPatch)
	if respPatch.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", respPatch.StatusCode)
	}
}
