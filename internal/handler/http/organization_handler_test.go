package http

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/group-project/crm_party_service/internal/db"
	"github.com/group-project/crm_party_service/internal/domain"
	"github.com/group-project/crm_party_service/internal/service"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockOrganizationService struct {
	createOrganizationFunc func(
		ctx context.Context,
		input domain.CreateOrganizationInput,
	) (db.Organization, error)

	getOrganizationByPartyIDFunc func(
		ctx context.Context,
		partyID pgtype.UUID,
	) (db.Organization, error)

	updateOrganizationFunc func(
		ctx context.Context,
		input domain.UpdateOrganizationInput,
	) (db.Organization, error)
}

func (m *mockOrganizationService) CreateOrganization(
	ctx context.Context,
	input domain.CreateOrganizationInput,
) (db.Organization, error) {
	if m.createOrganizationFunc != nil {
		return m.createOrganizationFunc(ctx, input)
	}

	return db.Organization{}, nil
}

func (m *mockOrganizationService) UpdateOrganization(
	ctx context.Context,
	input domain.UpdateOrganizationInput,
) (db.Organization, error) {
	if m.updateOrganizationFunc != nil {
		return m.updateOrganizationFunc(ctx, input)
	}

	return db.Organization{}, nil
}



func TestCreateOrganization_Success(t *testing.T) {
	app := fiber.New()

	expectedPartyID := pgtype.UUID{
		Bytes: [16]byte{1},
		Valid: true,
	}

	expectedOrganizationID := pgtype.UUID{
		Bytes: [16]byte{2},
		Valid: true,
	}

	employeeCount := int32(50)

	mockService := &mockOrganizationService{
		createOrganizationFunc: func(
			ctx context.Context,
			input domain.CreateOrganizationInput,
		) (db.Organization, error) {
			if input.PartyID != expectedPartyID.String() {
				t.Fatalf("expected party ID %s, got %s",
					expectedPartyID.String(),
					input.PartyID,
				)
			}

			if input.EmployeeCount == nil {
				t.Fatal("expected employee count")
			}

			if *input.EmployeeCount != 50 {
				t.Fatalf("expected employee count 50, got %d",
					*input.EmployeeCount,
				)
			}

			return db.Organization{
				ID:                 expectedOrganizationID,
				PartyID:            expectedPartyID,
				RegistrationNumber: stringPtr("REG-001"),
				Industry:           stringPtr("Technology"),
				EmployeeCount:      &employeeCount,
			}, nil
		},
	}

	handler := NewOrganizationHandler(mockService)

	app.Post("/api/v1/organizations", handler.CreateOrganization)

	body := `{
		"party_id": "` + expectedPartyID.String() + `",
		"registration_number": "REG-001",
		"industry": "Technology",
		"employee_count": 50
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/organizations",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}
}

func stringPtr(value string) *string {
	return &value
}

func TestCreateOrganization_InvalidPartyID(t *testing.T) {
	app := fiber.New()

	mockService := &mockOrganizationService{
		createOrganizationFunc: func(
			ctx context.Context,
			input domain.CreateOrganizationInput,
		) (db.Organization, error) {
			return db.Organization{}, service.ErrInvalidPartyID
		},
	}

	handler := NewOrganizationHandler(mockService)

	app.Post("/api/v1/organizations", handler.CreateOrganization)

	body := `{
		"party_id": "invalid-uuid"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/organizations",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreateOrganization_PartyNotFound(t *testing.T) {
	app := fiber.New()

	mockService := &mockOrganizationService{
		createOrganizationFunc: func(
			ctx context.Context,
			input domain.CreateOrganizationInput,
		) (db.Organization, error) {
			return db.Organization{}, service.ErrPartyNotFound
		},
	}

	handler := NewOrganizationHandler(mockService)

	app.Post("/api/v1/organizations", handler.CreateOrganization)

	body := `{
		"party_id": "550e8400-e29b-41d4-a716-446655440000"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/organizations",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestCreateOrganization_PartyNotOrganization(t *testing.T) {
	app := fiber.New()

	mockService := &mockOrganizationService{
		createOrganizationFunc: func(
			ctx context.Context,
			input domain.CreateOrganizationInput,
		) (db.Organization, error) {
			return db.Organization{}, service.ErrPartyNotOrganization
		},
	}

	handler := NewOrganizationHandler(mockService)

	app.Post("/api/v1/organizations", handler.CreateOrganization)

	body := `{
		"party_id": "550e8400-e29b-41d4-a716-446655440000"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/organizations",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreateOrganization_AlreadyExists(t *testing.T) {
	app := fiber.New()

	mockService := &mockOrganizationService{
		createOrganizationFunc: func(
			ctx context.Context,
			input domain.CreateOrganizationInput,
		) (db.Organization, error) {
			return db.Organization{}, service.ErrOrganizationAlreadyExists
		},
	}

	handler := NewOrganizationHandler(mockService)

	app.Post("/api/v1/organizations", handler.CreateOrganization)

	body := `{
		"party_id": "550e8400-e29b-41d4-a716-446655440000"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/organizations",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", resp.StatusCode)
	}
}

func TestCreateOrganization_InvalidEmployeeCount(t *testing.T) {
	app := fiber.New()

	mockService := &mockOrganizationService{
		createOrganizationFunc: func(
			ctx context.Context,
			input domain.CreateOrganizationInput,
		) (db.Organization, error) {
			return db.Organization{}, service.ErrInvalidEmployeeCount
		},
	}

	handler := NewOrganizationHandler(mockService)

	app.Post("/api/v1/organizations", handler.CreateOrganization)

	body := `{
		"party_id": "550e8400-e29b-41d4-a716-446655440000",
		"employee_count": -10
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/organizations",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreateOrganization_InvalidFoundedDate(t *testing.T) {
	app := fiber.New()

	mockService := &mockOrganizationService{
		createOrganizationFunc: func(
			ctx context.Context,
			input domain.CreateOrganizationInput,
		) (db.Organization, error) {
			return db.Organization{}, service.ErrInvalidFoundedDate
		},
	}

	handler := NewOrganizationHandler(mockService)

	app.Post("/api/v1/organizations", handler.CreateOrganization)

	body := `{
		"party_id": "550e8400-e29b-41d4-a716-446655440000",
		"founded_date": "2024-99-99"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/organizations",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreateOrganization_InternalServerError(t *testing.T) {
	app := fiber.New()

	mockService := &mockOrganizationService{
		createOrganizationFunc: func(
			ctx context.Context,
			input domain.CreateOrganizationInput,
		) (db.Organization, error) {
			return db.Organization{}, errors.New("database connection failed")
		},
	}

	handler := NewOrganizationHandler(mockService)

	app.Post("/api/v1/organizations", handler.CreateOrganization)

	body := `{
		"party_id": "550e8400-e29b-41d4-a716-446655440000"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/organizations",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", resp.StatusCode)
	}
}

func TestCreateOrganization_InvalidJSON(t *testing.T) {
	app := fiber.New()

	mockService := &mockOrganizationService{
		createOrganizationFunc: func(
			ctx context.Context,
			input domain.CreateOrganizationInput,
		) (db.Organization, error) {
			t.Fatal("service should not be called for invalid JSON")
			return db.Organization{}, nil
		},
	}

	handler := NewOrganizationHandler(mockService)

	app.Post("/api/v1/organizations", handler.CreateOrganization)

	body := `{
		"party_id": "550e8400-e29b-41d4-a716-446655440000",
		"employee_count":
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/organizations",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func (m *mockOrganizationService) GetOrganizationByPartyID(
	ctx context.Context,
	partyID pgtype.UUID,
) (db.Organization, error) {
	if m.getOrganizationByPartyIDFunc != nil {
		return m.getOrganizationByPartyIDFunc(ctx, partyID)
	}

	return db.Organization{}, nil
}

func TestGetOrganizationByPartyID_Success(t *testing.T) {
	app := fiber.New()

	partyID := pgtype.UUID{
		Bytes: [16]byte{1},
		Valid: true,
	}

	organizationID := pgtype.UUID{
		Bytes: [16]byte{2},
		Valid: true,
	}

	mockService := &mockOrganizationService{
		getOrganizationByPartyIDFunc: func(
			ctx context.Context,
			id pgtype.UUID,
		) (db.Organization, error) {
			if id != partyID {
				t.Fatalf("unexpected party ID")
			}

			return db.Organization{
				ID:                 organizationID,
				PartyID:            partyID,
				RegistrationNumber: stringPtr("REG-2026-001"),
				Industry:           stringPtr("Software"),
				OrganizationType:   stringPtr("PRIVATE_LIMITED"),
				EmployeeCount:      func() *int32 { v := int32(25); return &v }(),
			}, nil
		},
	}

	handler := NewOrganizationHandler(mockService)

	app.Get(
		"/api/v1/organizations/:party_id",
		handler.GetOrganizationByPartyID,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/organizations/"+partyID.String(),
		nil,
	)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetOrganizationByPartyID_InvalidPartyID(t *testing.T) {
	app := fiber.New()

	mockService := &mockOrganizationService{
		getOrganizationByPartyIDFunc: func(
			ctx context.Context,
			id pgtype.UUID,
		) (db.Organization, error) {
			t.Fatal("service should not be called")
			return db.Organization{}, nil
		},
	}

	handler := NewOrganizationHandler(mockService)

	app.Get(
		"/api/v1/organizations/:party_id",
		handler.GetOrganizationByPartyID,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/organizations/not-a-uuid",
		nil,
	)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestGetOrganizationByPartyID_NotFound(t *testing.T) {
	app := fiber.New()

	partyID := pgtype.UUID{
		Bytes: [16]byte{1},
		Valid: true,
	}

	mockService := &mockOrganizationService{
		getOrganizationByPartyIDFunc: func(
			ctx context.Context,
			id pgtype.UUID,
		) (db.Organization, error) {
			return db.Organization{}, service.ErrOrganizationNotFound
		},
	}

	handler := NewOrganizationHandler(mockService)

	app.Get(
		"/api/v1/organizations/:party_id",
		handler.GetOrganizationByPartyID,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/organizations/"+partyID.String(),
		nil,
	)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestGetOrganizationByPartyID_InternalServerError(t *testing.T) {
	app := fiber.New()

	partyID := pgtype.UUID{
		Bytes: [16]byte{1},
		Valid: true,
	}

	mockService := &mockOrganizationService{
		getOrganizationByPartyIDFunc: func(
			ctx context.Context,
			id pgtype.UUID,
		) (db.Organization, error) {
			return db.Organization{}, errors.New("database connection failed")
		},
	}

	handler := NewOrganizationHandler(mockService)

	app.Get(
		"/api/v1/organizations/:party_id",
		handler.GetOrganizationByPartyID,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/organizations/"+partyID.String(),
		nil,
	)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", resp.StatusCode)
	}
}


