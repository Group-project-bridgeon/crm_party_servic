package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/group-project/crm_party_service/internal/handler/http/dto"
)

func TestCustomerEndpoints(t *testing.T) {
	app, _ := setupTestApp()

	// 1. Create Base Party
	var party dto.PartyResponse
	{
		partyBody := []byte(`{"party_type":"ORGANIZATION","display_name":"Acme Retailers"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/parties", bytes.NewReader(partyBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusCreated {
			t.Fatalf("failed to create party for customer: %v", err)
		}
		_ = json.NewDecoder(resp.Body).Decode(&party)
	}

	// 2. Create Customer - Validation Failure (empty code)
	{
		invalidBody := []byte(`{
			"party_id": "` + party.ID + `",
			"customer_code": "   "
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/customers", bytes.NewReader(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status 400 for empty customer code, got %d", resp.StatusCode)
		}
	}

	// 3. Create Customer - Validation Failure (invalid credit status)
	{
		invalidBody := []byte(`{
			"party_id": "` + party.ID + `",
			"customer_code": "CUST-001",
			"credit_status": "SUPER_ACTIVE"
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/customers", bytes.NewReader(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status 400 for invalid credit status, got %d", resp.StatusCode)
		}
	}

	// 4. Create Customer - Success
	var customer dto.CustomerResponse
	{
		validBody := []byte(`{
			"party_id": "` + party.ID + `",
			"customer_code": "CUST-001",
			"customer_group": "RETAIL",
			"credit_limit": 50000.00,
			"payment_terms_days": 30,
			"credit_status": "ACTIVE"
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/customers", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", resp.StatusCode)
		}

		if err := json.NewDecoder(resp.Body).Decode(&customer); err != nil {
			t.Fatalf("failed to decode customer: %v", err)
		}
		if customer.CustomerCode != "CUST-001" {
			t.Fatalf("expected customer code CUST-001, got %s", customer.CustomerCode)
		}
		if customer.CreditLimit != 50000.00 {
			t.Fatalf("expected credit limit 50000, got %f", customer.CreditLimit)
		}
	}

	// 5. Create Customer - Duplicate Code Conflict
	{
		dupBody := []byte(`{
			"party_id": "` + party.ID + `",
			"customer_code": "CUST-001"
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/customers", bytes.NewReader(dupBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusConflict {
			t.Fatalf("expected status 409 for duplicate customer code, got %d", resp.StatusCode)
		}
	}

	// 6. Get Customer by ID - Success
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/customers/"+customer.ID, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var fetched dto.CustomerResponse
		_ = json.NewDecoder(resp.Body).Decode(&fetched)
		if fetched.ID != customer.ID {
			t.Fatalf("expected ID %s, got %s", customer.ID, fetched.ID)
		}
	}

	// 7. Get Customer by Party ID - Success
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/customers/party/"+party.ID, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var fetched dto.CustomerResponse
		_ = json.NewDecoder(resp.Body).Decode(&fetched)
		if fetched.PartyID != party.ID {
			t.Fatalf("expected party ID %s, got %s", party.ID, fetched.PartyID)
		}
	}

	// 8. List Customers - Success
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/customers?limit=10&offset=0", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var listResp dto.ListCustomersResponse
		_ = json.NewDecoder(resp.Body).Decode(&listResp)
		if len(listResp.Data) == 0 {
			t.Fatal("expected at least 1 customer")
		}
	}

	// 9. Update Customer - Success
	{
		updateBody := []byte(`{
			"customer_group": "ENTERPRISE",
			"credit_limit": 100000.00,
			"payment_terms_days": 45,
			"credit_status": "ACTIVE"
		}`)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/customers/"+customer.ID, bytes.NewReader(updateBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var updated dto.CustomerResponse
		_ = json.NewDecoder(resp.Body).Decode(&updated)
		if updated.CustomerGroup == nil || *updated.CustomerGroup != "ENTERPRISE" {
			t.Fatalf("expected customer group ENTERPRISE")
		}
		if updated.CreditLimit != 100000.00 {
			t.Fatalf("expected credit limit 100000, got %f", updated.CreditLimit)
		}
	}

	// 10. Customer User Accounts - Link Portal User
	var account dto.CustomerUserAccountResponse
	{
		accBody := []byte(`{"user_id": "99999999-9999-9999-9999-999999999999"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/customers/"+customer.ID+"/user-accounts", bytes.NewReader(accBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status 201 for user account creation, got %d", resp.StatusCode)
		}
		_ = json.NewDecoder(resp.Body).Decode(&account)
		if account.CustomerID != customer.ID {
			t.Fatalf("expected customer ID %s, got %s", customer.ID, account.CustomerID)
		}
	}

	// 11. List Customer User Accounts
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/customers/"+customer.ID+"/user-accounts", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var accounts []dto.CustomerUserAccountResponse
		_ = json.NewDecoder(resp.Body).Decode(&accounts)
		if len(accounts) == 0 {
			t.Fatal("expected at least 1 linked user account")
		}
	}

	// 12. Update User Account Status
	{
		statusBody := []byte(`{"is_active": false}`)
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/customers/user-accounts/"+account.ID+"/status", bytes.NewReader(statusBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var updatedAcc dto.CustomerUserAccountResponse
		_ = json.NewDecoder(resp.Body).Decode(&updatedAcc)
		if updatedAcc.IsActive != false {
			t.Fatal("expected account is_active to be false")
		}
	}
}
