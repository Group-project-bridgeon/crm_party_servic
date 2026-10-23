package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/group-project/crm_party_service/internal/handler/http/dto"
)

func TestPartyEndpoints(t *testing.T) {
	app, _ := setupTestApp()

	// 1. Create Party - Validation Failure (invalid party type)
	{
		invalidBody := []byte(`{"party_type":"INVALID_TYPE","display_name":"Acme Corp"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/parties", bytes.NewReader(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status 400 for invalid party type, got %d", resp.StatusCode)
		}
	}

	// 2. Create Party - Validation Failure (missing display name)
	{
		invalidBody := []byte(`{"party_type":"ORGANIZATION","display_name":"   "}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/parties", bytes.NewReader(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status 400 for empty display name, got %d", resp.StatusCode)
		}
	}

	// 3. Create Party - Success
	var createdParty dto.PartyResponse
	{
		validBody := []byte(`{
			"party_type": "ORGANIZATION",
			"display_name": "Globex Corp",
			"legal_name": "Globex Corporation LLC",
			"email": "info@globex.com"
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/parties", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", resp.StatusCode)
		}

		if err := json.NewDecoder(resp.Body).Decode(&createdParty); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if createdParty.ID == "" {
			t.Fatal("expected non-empty party ID")
		}
		if createdParty.DisplayName != "Globex Corp" {
			t.Fatalf("expected display name Globex Corp, got %s", createdParty.DisplayName)
		}
	}

	// 4. Get Party by ID - Success
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/parties/"+createdParty.ID, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var party dto.PartyResponse
		if err := json.NewDecoder(resp.Body).Decode(&party); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if party.ID != createdParty.ID {
			t.Fatalf("expected ID %s, got %s", createdParty.ID, party.ID)
		}
	}

	// 5. Get Party by ID - Not Found
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/parties/00000000-0000-0000-0000-000000000000", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", resp.StatusCode)
		}
	}

	// 6. List Parties - Success
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/parties?limit=10&offset=0", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var listResp dto.ListPartiesResponse
		if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(listResp.Data) == 0 {
			t.Fatal("expected at least 1 party in list")
		}
	}

	// 7. Update Party - Success
	{
		updateBody := []byte(`{
			"display_name": "Globex Global",
			"legal_name": "Globex Global Inc"
		}`)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/parties/"+createdParty.ID, bytes.NewReader(updateBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var updated dto.PartyResponse
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if updated.DisplayName != "Globex Global" {
			t.Fatalf("expected updated display name 'Globex Global', got %s", updated.DisplayName)
		}
	}

	// 8. Update Party Status - Success
	{
		statusBody := []byte(`{"is_active": false}`)
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/parties/"+createdParty.ID+"/status", bytes.NewReader(statusBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var updated dto.PartyResponse
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if updated.IsActive != false {
			t.Fatal("expected is_active to be false")
		}
	}

	// 9. Delete Party (Deactivate) - Success
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/parties/"+createdParty.ID, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}
	}
}

func TestOrganizationEndpoints(t *testing.T) {
	app, _ := setupTestApp()

	// 1. Create Base Party first
	var party dto.PartyResponse
	{
		partyBody := []byte(`{"party_type":"ORGANIZATION","display_name":"Acme Logistics"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/parties", bytes.NewReader(partyBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusCreated {
			t.Fatalf("failed to create party for org test: %v, status: %d", err, resp.StatusCode)
		}
		_ = json.NewDecoder(resp.Body).Decode(&party)
	}

	// 2. Create Organization Profile - Success
	{
		orgBody := []byte(`{
			"party_id": "` + party.ID + `",
			"registration_number": "REG-12345",
			"industry": "Logistics",
			"organization_type": "PRIVATE_LIMITED"
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/organizations", bytes.NewReader(orgBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", resp.StatusCode)
		}
	}

	// 3. Get Organization Profile - Success
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/organizations/"+party.ID, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var org dto.OrganizationResponse
		if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if org.PartyID != party.ID {
			t.Fatalf("expected party ID %s, got %s", party.ID, org.PartyID)
		}
	}

	// 4. Update Organization Profile - Success
	{
		updateBody := []byte(`{
			"registration_number": "REG-99999",
			"industry": "Supply Chain"
		}`)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/organizations/"+party.ID, bytes.NewReader(updateBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var updated dto.OrganizationResponse
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if updated.RegistrationNumber == nil || *updated.RegistrationNumber != "REG-99999" {
			t.Fatalf("expected updated reg number REG-99999")
		}
	}
}

func TestIndividualEndpoints(t *testing.T) {
	app, _ := setupTestApp()

	// 1. Create Base Party first
	var party dto.PartyResponse
	{
		partyBody := []byte(`{"party_type":"INDIVIDUAL","display_name":"Jane Doe"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/parties", bytes.NewReader(partyBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusCreated {
			t.Fatalf("failed to create party for indiv test: %v, status: %d", err, resp.StatusCode)
		}
		_ = json.NewDecoder(resp.Body).Decode(&party)
	}

	// 2. Create Individual Profile - Validation failure (missing first_name)
	{
		invalidBody := []byte(`{
			"party_id": "` + party.ID + `",
			"first_name": "  "
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/individuals", bytes.NewReader(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status 400 for empty first name, got %d", resp.StatusCode)
		}
	}

	// 3. Create Individual Profile - Success
	{
		validBody := []byte(`{
			"party_id": "` + party.ID + `",
			"first_name": "Jane",
			"last_name": "Doe",
			"job_title": "Software Engineer"
		}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/individuals", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", resp.StatusCode)
		}
	}

	// 4. Get Individual Profile - Success
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/individuals/"+party.ID, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var indiv dto.IndividualResponse
		if err := json.NewDecoder(resp.Body).Decode(&indiv); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if indiv.PartyID != party.ID {
			t.Fatalf("expected party ID %s, got %s", party.ID, indiv.PartyID)
		}
		if indiv.FirstName != "Jane" {
			t.Fatalf("expected first name Jane, got %s", indiv.FirstName)
		}
	}

	// 5. Update Individual Profile - Success
	{
		updateBody := []byte(`{
			"first_name": "Janet",
			"last_name": "Doe",
			"job_title": "Senior Engineer"
		}`)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/individuals/"+party.ID, bytes.NewReader(updateBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var updated dto.IndividualResponse
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if updated.FirstName != "Janet" {
			t.Fatalf("expected updated first name Janet, got %s", updated.FirstName)
		}
	}
}
