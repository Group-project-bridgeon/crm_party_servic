package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/group-project/crm_party_service/internal/handler/http/dto"
)

func TestCompositeEndpoints(t *testing.T) {
	app, _ := setupTestApp()

	// 1. Party 360 - Party Not Found
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/parties/00000000-0000-0000-0000-000000000000/360", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", resp.StatusCode)
		}
	}

	// 2. Setup a full Party with Org, Customer, Address, Contact, Note
	var party dto.PartyResponse
	{
		partyBody := []byte(`{"party_type":"ORGANIZATION","display_name":"Acme Global Enterprise"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/parties", bytes.NewReader(partyBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		_ = json.NewDecoder(resp.Body).Decode(&party)

		// Create Org profile
		orgBody := []byte(`{"party_id":"` + party.ID + `","industry":"Aerospace"}`)
		reqOrg := httptest.NewRequest(http.MethodPost, "/api/v1/organizations", bytes.NewReader(orgBody))
		reqOrg.Header.Set("Content-Type", "application/json")
		_, _ = app.Test(reqOrg)

		// Create Customer profile
		custBody := []byte(`{"party_id":"` + party.ID + `","customer_code":"CUST-360","customer_group":"GLOBAL"}`)
		reqCust := httptest.NewRequest(http.MethodPost, "/api/v1/customers", bytes.NewReader(custBody))
		reqCust.Header.Set("Content-Type", "application/json")
		_, _ = app.Test(reqCust)

		// Create Address
		addrBody := []byte(`{"address_type":"BILLING","address_line1":"1 Launchpad Way","city":"Houston","country_code":"US"}`)
		reqAddr := httptest.NewRequest(http.MethodPost, "/api/v1/parties/"+party.ID+"/addresses", bytes.NewReader(addrBody))
		reqAddr.Header.Set("Content-Type", "application/json")
		_, _ = app.Test(reqAddr)

		// Create Contact
		contactBody := []byte(`{"contact_type":"EMAIL","contact_value":"ops@acmeglobal.com"}`)
		reqContact := httptest.NewRequest(http.MethodPost, "/api/v1/parties/"+party.ID+"/contacts", bytes.NewReader(contactBody))
		reqContact.Header.Set("Content-Type", "application/json")
		_, _ = app.Test(reqContact)

		// Create Note
		noteBody := []byte(`{"content":"Strategic account tier 1"}`)
		reqNote := httptest.NewRequest(http.MethodPost, "/api/v1/parties/"+party.ID+"/notes", bytes.NewReader(noteBody))
		reqNote.Header.Set("Content-Type", "application/json")
		_, _ = app.Test(reqNote)
	}

	// 3. Party 360 - Success with all relations
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/parties/"+party.ID+"/360", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var view dto.Party360Response
		if err := json.NewDecoder(resp.Body).Decode(&view); err != nil {
			t.Fatalf("failed to decode 360 response: %v", err)
		}
		if view.Party.ID != party.ID {
			t.Fatalf("expected party ID %s, got %s", party.ID, view.Party.ID)
		}
		if view.Organization == nil {
			t.Fatal("expected organization profile in 360 view")
		}
		if view.Customer == nil || view.Customer.CustomerCode != "CUST-360" {
			t.Fatal("expected customer in 360 view with code CUST-360")
		}
		if len(view.Addresses) != 1 {
			t.Fatalf("expected 1 address, got %d", len(view.Addresses))
		}
		if len(view.Contacts) != 1 {
			t.Fatalf("expected 1 contact, got %d", len(view.Contacts))
		}
		if len(view.Notes) != 1 {
			t.Fatalf("expected 1 note, got %d", len(view.Notes))
		}
	}

	// 4. Composite Customer - Validation failures
	{
		// Invalid party type
		body := []byte(`{"party_type":"UNKNOWN","display_name":"Acme","customer_code":"C-1"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/customers/composite", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status 400 for invalid party type, got %d", resp.StatusCode)
		}

		// Empty customer code
		bodyEmptyCode := []byte(`{"party_type":"ORGANIZATION","display_name":"Acme","customer_code":"  "}`)
		reqEmptyCode := httptest.NewRequest(http.MethodPost, "/api/v1/customers/composite", bytes.NewReader(bodyEmptyCode))
		reqEmptyCode.Header.Set("Content-Type", "application/json")
		respEmptyCode, _ := app.Test(reqEmptyCode)
		if respEmptyCode.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status 400 for empty customer code, got %d", respEmptyCode.StatusCode)
		}

		// Individual party missing first name
		bodyIndivNoName := []byte(`{"party_type":"INDIVIDUAL","display_name":"John Doe","customer_code":"C-INDIV-1"}`)
		reqIndiv := httptest.NewRequest(http.MethodPost, "/api/v1/customers/composite", bytes.NewReader(bodyIndivNoName))
		reqIndiv.Header.Set("Content-Type", "application/json")
		respIndiv, _ := app.Test(reqIndiv)
		if respIndiv.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status 400 for missing first name, got %d", respIndiv.StatusCode)
		}
	}

	// 5. Composite Customer - Success
	{
		validBody := []byte(`{
			"party_type": "ORGANIZATION",
			"display_name": "Stark Industries",
			"customer_code": "CUST-STARK",
			"customer_group": "DEFENSE",
			"industry": "Advanced Tech",
			"credit_limit": 500000.00,
			"address": {
				"address_type": "BILLING",
				"address_line1": "10880 Malibu Point",
				"city": "Malibu",
				"country_code": "US"
			}
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/customers/composite", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status 201 for composite customer, got %d", resp.StatusCode)
		}

		var created dto.CompositeCustomerResponse
		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if created.Party.DisplayName != "Stark Industries" {
			t.Fatalf("expected display name Stark Industries, got %s", created.Party.DisplayName)
		}
		if created.Customer.CustomerCode != "CUST-STARK" {
			t.Fatalf("expected customer code CUST-STARK, got %s", created.Customer.CustomerCode)
		}
	}

	// 6. Composite Supplier - Validation failures
	{
		// Invalid currency code
		body := []byte(`{"party_type":"ORGANIZATION","display_name":"Wayne Corp","supplier_code":"S-1","currency_code":"INVALID"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/suppliers/composite", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status 400 for invalid currency, got %d", resp.StatusCode)
		}

		// Invalid rating
		bodyRating := []byte(`{"party_type":"ORGANIZATION","display_name":"Wayne Corp","supplier_code":"S-1","supplier_rating": 8.0}`)
		reqRating := httptest.NewRequest(http.MethodPost, "/api/v1/suppliers/composite", bytes.NewReader(bodyRating))
		reqRating.Header.Set("Content-Type", "application/json")
		respRating, _ := app.Test(reqRating)
		if respRating.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status 400 for invalid rating, got %d", respRating.StatusCode)
		}
	}

	// 7. Composite Supplier - Success
	{
		validBody := []byte(`{
			"party_type": "ORGANIZATION",
			"display_name": "Wayne Enterprises",
			"supplier_code": "SUPP-WAYNE",
			"supplier_group": "CONGLOMERATE",
			"currency_code": "USD",
			"supplier_rating": 4.9,
			"address": {
				"address_type": "SHIPPING",
				"address_line1": "1007 Mountain Drive",
				"city": "Gotham",
				"country_code": "US"
			}
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/suppliers/composite", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status 201 for composite supplier, got %d", resp.StatusCode)
		}

		var created dto.CompositeSupplierResponse
		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if created.Party.DisplayName != "Wayne Enterprises" {
			t.Fatalf("expected display name Wayne Enterprises, got %s", created.Party.DisplayName)
		}
		if created.Supplier.SupplierCode != "SUPP-WAYNE" {
			t.Fatalf("expected supplier code SUPP-WAYNE, got %s", created.Supplier.SupplierCode)
		}
	}
}
