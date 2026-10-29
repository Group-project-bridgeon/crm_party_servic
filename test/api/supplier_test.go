package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/group-project/crm_party_service/internal/handler/http/dto"
)

func TestSupplierEndpoints(t *testing.T) {
	app, _ := setupTestApp()

	// 1. Create Base Party
	var party dto.PartyResponse
	{
		partyBody := []byte(`{"party_type":"ORGANIZATION","display_name":"Apex Components"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/parties", bytes.NewReader(partyBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusCreated {
			t.Fatalf("failed to create party for supplier: %v", err)
		}
		_ = json.NewDecoder(resp.Body).Decode(&party)
	}

	// 2. Create Supplier - Validation Failure (empty code)
	{
		invalidBody := []byte(`{
			"party_id": "` + party.ID + `",
			"supplier_code": "   "
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/suppliers", bytes.NewReader(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status 400 for empty supplier code, got %d", resp.StatusCode)
		}
	}

	// 3. Create Supplier - Validation Failure (invalid rating > 5.0)
	{
		invalidBody := []byte(`{
			"party_id": "` + party.ID + `",
			"supplier_code": "SUPP-001",
			"supplier_rating": 6.5
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/suppliers", bytes.NewReader(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status 400 for invalid rating, got %d", resp.StatusCode)
		}
	}

	// 4. Create Supplier - Validation Failure (invalid currency code len != 3)
	{
		invalidBody := []byte(`{
			"party_id": "` + party.ID + `",
			"supplier_code": "SUPP-001",
			"currency_code": "USDD"
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/suppliers", bytes.NewReader(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status 400 for invalid currency, got %d", resp.StatusCode)
		}
	}

	// 5. Create Supplier - Success
	var supplier dto.SupplierResponse
	{
		validBody := []byte(`{
			"party_id": "` + party.ID + `",
			"supplier_code": "SUPP-001",
			"supplier_group": "HARDWARE",
			"currency_code": "USD",
			"supplier_rating": 4.8,
			"payment_terms_days": 60,
			"payment_method": "WIRE"
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/suppliers", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", resp.StatusCode)
		}

		if err := json.NewDecoder(resp.Body).Decode(&supplier); err != nil {
			t.Fatalf("failed to decode supplier: %v", err)
		}
		if supplier.SupplierCode != "SUPP-001" {
			t.Fatalf("expected supplier code SUPP-001, got %s", supplier.SupplierCode)
		}
		if supplier.SupplierRating == nil || *supplier.SupplierRating != 4.8 {
			t.Fatalf("expected supplier rating 4.8")
		}
	}

	// 6. Create Supplier - Duplicate Code Conflict
	{
		dupBody := []byte(`{
			"party_id": "` + party.ID + `",
			"supplier_code": "SUPP-001"
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/suppliers", bytes.NewReader(dupBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusConflict {
			t.Fatalf("expected status 409 for duplicate supplier code, got %d", resp.StatusCode)
		}
	}

	// 7. Get Supplier by ID - Success
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/suppliers/"+supplier.ID, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var fetched dto.SupplierResponse
		_ = json.NewDecoder(resp.Body).Decode(&fetched)
		if fetched.ID != supplier.ID {
			t.Fatalf("expected ID %s, got %s", supplier.ID, fetched.ID)
		}
	}

	// 8. Get Supplier by Party ID - Success
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/suppliers/party/"+party.ID, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var fetched dto.SupplierResponse
		_ = json.NewDecoder(resp.Body).Decode(&fetched)
		if fetched.PartyID != party.ID {
			t.Fatalf("expected party ID %s, got %s", party.ID, fetched.PartyID)
		}
	}

	// 9. List Suppliers - Success
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/suppliers?limit=10&offset=0", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var listResp dto.ListSuppliersResponse
		_ = json.NewDecoder(resp.Body).Decode(&listResp)
		if len(listResp.Data) == 0 {
			t.Fatal("expected at least 1 supplier")
		}
	}

	// 10. Update Supplier - Success
	{
		updateBody := []byte(`{
			"supplier_group": "PRECISION_HARDWARE",
			"supplier_rating": 4.95,
			"payment_terms_days": 45
		}`)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/suppliers/"+supplier.ID, bytes.NewReader(updateBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var updated dto.SupplierResponse
		_ = json.NewDecoder(resp.Body).Decode(&updated)
		if updated.SupplierGroup == nil || *updated.SupplierGroup != "PRECISION_HARDWARE" {
			t.Fatalf("expected supplier group PRECISION_HARDWARE")
		}
	}

	// 11. Supplier User Accounts - Link Vendor User
	var account dto.SupplierUserAccountResponse
	{
		accBody := []byte(`{"user_id": "88888888-8888-8888-8888-888888888888"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/suppliers/"+supplier.ID+"/user-accounts", bytes.NewReader(accBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status 201 for supplier user account creation, got %d", resp.StatusCode)
		}
		_ = json.NewDecoder(resp.Body).Decode(&account)
		if account.SupplierID != supplier.ID {
			t.Fatalf("expected supplier ID %s, got %s", supplier.ID, account.SupplierID)
		}
	}

	// 12. List Supplier User Accounts
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/suppliers/"+supplier.ID+"/user-accounts", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var accounts []dto.SupplierUserAccountResponse
		_ = json.NewDecoder(resp.Body).Decode(&accounts)
		if len(accounts) == 0 {
			t.Fatal("expected at least 1 supplier linked user account")
		}
	}

	// 13. Update Supplier User Account Status
	{
		statusBody := []byte(`{"is_active": false}`)
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/suppliers/user-accounts/"+account.ID+"/status", bytes.NewReader(statusBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var updatedAcc dto.SupplierUserAccountResponse
		_ = json.NewDecoder(resp.Body).Decode(&updatedAcc)
		if updatedAcc.IsActive != false {
			t.Fatal("expected account is_active to be false")
		}
	}
}
