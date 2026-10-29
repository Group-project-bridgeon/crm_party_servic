package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/group-project/crm_party_service/internal/handler/http/dto"
)

func TestAddressEndpoints(t *testing.T) {
	app, _ := setupTestApp()

	// Create Party
	var party dto.PartyResponse
	{
		partyBody := []byte(`{"party_type":"ORGANIZATION","display_name":"Acme HQ"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/parties", bytes.NewReader(partyBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		_ = json.NewDecoder(resp.Body).Decode(&party)
	}

	// 1. Create Address - Validation failure (invalid address type)
	{
		invalidBody := []byte(`{
			"address_type": "UNKNOWN_TYPE",
			"address_line1": "123 Main St",
			"city": "Springfield",
			"country_code": "US"
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/parties/"+party.ID+"/addresses", bytes.NewReader(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status 400 for invalid address type, got %d", resp.StatusCode)
		}
	}

	// 2. Create Address - Success
	var address dto.AddressResponse
	{
		validBody := []byte(`{
			"address_type": "BILLING",
			"address_line1": "500 Market Street",
			"city": "San Francisco",
			"state": "CA",
			"postal_code": "94105",
			"country_code": "US",
			"is_primary": true
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/parties/"+party.ID+"/addresses", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", resp.StatusCode)
		}
		_ = json.NewDecoder(resp.Body).Decode(&address)
		if address.City != "San Francisco" {
			t.Fatalf("expected city San Francisco, got %s", address.City)
		}
		if !address.IsPrimary {
			t.Fatal("expected is_primary to be true")
		}
	}

	// 3. Get Address by ID - Success
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/addresses/"+address.ID, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}
	}

	// 4. List Addresses by Party ID
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/parties/"+party.ID+"/addresses", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}
		var list []dto.AddressResponse
		_ = json.NewDecoder(resp.Body).Decode(&list)
		if len(list) == 0 {
			t.Fatal("expected at least 1 address")
		}
	}

	// 5. Update Address
	{
		updateBody := []byte(`{
			"address_type": "BILLING",
			"address_line1": "550 Market Street",
			"city": "San Francisco",
			"country_code": "US",
			"is_primary": true,
			"is_active": true
		}`)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/addresses/"+address.ID, bytes.NewReader(updateBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}
		var updated dto.AddressResponse
		_ = json.NewDecoder(resp.Body).Decode(&updated)
		if updated.AddressLine1 != "550 Market Street" {
			t.Fatalf("expected updated address line 1, got %s", updated.AddressLine1)
		}
	}

	// 6. Delete Address
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/addresses/"+address.ID, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}
	}
}

func TestContactEndpoints(t *testing.T) {
	app, _ := setupTestApp()

	// Create Party
	var party dto.PartyResponse
	{
		partyBody := []byte(`{"party_type":"INDIVIDUAL","display_name":"Alice Smith"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/parties", bytes.NewReader(partyBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		_ = json.NewDecoder(resp.Body).Decode(&party)
	}

	// 1. Create Contact - Validation failure (invalid contact type)
	{
		invalidBody := []byte(`{
			"contact_type": "TELEPATHY",
			"contact_value": "alice"
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/parties/"+party.ID+"/contacts", bytes.NewReader(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status 400 for invalid contact type, got %d", resp.StatusCode)
		}
	}

	// 2. Create Contact - Success
	var contact dto.ContactResponse
	{
		validBody := []byte(`{
			"contact_type": "EMAIL",
			"contact_value": "alice@example.com",
			"label": "Work",
			"is_primary": true
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/parties/"+party.ID+"/contacts", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", resp.StatusCode)
		}
		_ = json.NewDecoder(resp.Body).Decode(&contact)
		if contact.ContactValue != "alice@example.com" {
			t.Fatalf("expected contact value alice@example.com, got %s", contact.ContactValue)
		}
	}

	// 3. Get Contact by ID
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/contacts/"+contact.ID, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}
	}

	// 4. List Contacts by Party ID
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/parties/"+party.ID+"/contacts", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}
		var list []dto.ContactResponse
		_ = json.NewDecoder(resp.Body).Decode(&list)
		if len(list) == 0 {
			t.Fatal("expected at least 1 contact")
		}
	}

	// 5. Update Contact
	{
		updateBody := []byte(`{
			"contact_type": "EMAIL",
			"contact_value": "alice.smith@example.com",
			"is_primary": true,
			"is_active": true
		}`)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/contacts/"+contact.ID, bytes.NewReader(updateBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}
		var updated dto.ContactResponse
		_ = json.NewDecoder(resp.Body).Decode(&updated)
		if updated.ContactValue != "alice.smith@example.com" {
			t.Fatalf("expected updated contact value, got %s", updated.ContactValue)
		}
	}

	// 6. Delete Contact
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/contacts/"+contact.ID, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}
	}
}

func TestContactPersonEndpoints(t *testing.T) {
	app, _ := setupTestApp()

	// Create Party
	var party dto.PartyResponse
	{
		partyBody := []byte(`{"party_type":"ORGANIZATION","display_name":"Omni Consumer Products"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/parties", bytes.NewReader(partyBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		_ = json.NewDecoder(resp.Body).Decode(&party)
	}

	// 1. Create Contact Person - Validation failure (missing first_name)
	{
		invalidBody := []byte(`{"first_name": "  "}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/parties/"+party.ID+"/contact-persons", bytes.NewReader(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status 400 for empty first name, got %d", resp.StatusCode)
		}
	}

	// 2. Create Contact Person - Success
	var person dto.ContactPersonResponse
	{
		validBody := []byte(`{
			"first_name": "Bob",
			"last_name": "Morton",
			"designation": "VP Executive",
			"department": "Security Concepts",
			"email": "bmorton@ocp.com",
			"is_primary": true
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/parties/"+party.ID+"/contact-persons", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", resp.StatusCode)
		}
		_ = json.NewDecoder(resp.Body).Decode(&person)
		if person.FirstName != "Bob" {
			t.Fatalf("expected first name Bob, got %s", person.FirstName)
		}
	}

	// 3. Get Contact Person by ID
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/contact-persons/"+person.ID, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}
	}

	// 4. List Contact Persons by Party ID
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/parties/"+party.ID+"/contact-persons", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}
		var list []dto.ContactPersonResponse
		_ = json.NewDecoder(resp.Body).Decode(&list)
		if len(list) == 0 {
			t.Fatal("expected at least 1 contact person")
		}
	}

	// 5. Update Contact Person
	{
		updateBody := []byte(`{
			"first_name": "Robert",
			"last_name": "Morton",
			"designation": "Senior VP",
			"is_primary": true,
			"is_active": true
		}`)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/contact-persons/"+person.ID, bytes.NewReader(updateBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}
		var updated dto.ContactPersonResponse
		_ = json.NewDecoder(resp.Body).Decode(&updated)
		if updated.FirstName != "Robert" {
			t.Fatalf("expected updated first name Robert, got %s", updated.FirstName)
		}
	}

	// 6. Delete Contact Person
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/contact-persons/"+person.ID, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}
	}
}

func TestNoteEndpoints(t *testing.T) {
	app, _ := setupTestApp()

	// Create Party
	var party dto.PartyResponse
	{
		partyBody := []byte(`{"party_type":"INDIVIDUAL","display_name":"Peter Parker"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/parties", bytes.NewReader(partyBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		_ = json.NewDecoder(resp.Body).Decode(&party)
	}

	// 1. Create Note - Validation failure (empty content)
	{
		invalidBody := []byte(`{"content": "   "}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/parties/"+party.ID+"/notes", bytes.NewReader(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status 400 for empty note content, got %d", resp.StatusCode)
		}
	}

	// 2. Create Note - Success
	var note dto.NoteResponse
	{
		validBody := []byte(`{
			"note_type": "GENERAL",
			"content": "Requested portfolio brochure by email.",
			"created_by": "00000000-0000-0000-0000-000000000001"
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/parties/"+party.ID+"/notes", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", resp.StatusCode)
		}
		_ = json.NewDecoder(resp.Body).Decode(&note)
		if note.Content != "Requested portfolio brochure by email." {
			t.Fatalf("expected correct content, got %s", note.Content)
		}
	}

	// 3. Get Note by ID
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/notes/"+note.ID, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}
	}

	// 4. List Notes by Party ID
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/parties/"+party.ID+"/notes", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}
		var list []dto.NoteResponse
		_ = json.NewDecoder(resp.Body).Decode(&list)
		if len(list) == 0 {
			t.Fatal("expected at least 1 note")
		}
	}

	// 5. Delete Note
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/notes/"+note.ID, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}
	}
}
