package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
	"github.com/group-project/crm_party_service/internal/domain"
)

type mockOrganizationRepository struct {
	createOrganizationCalled     bool
	createOrganizationParams     db.CreateOrganizationParams
	createOrganizationErr        error
	getOrganizationByPartyIDFunc func(ctx context.Context, partyID pgtype.UUID) (db.Organization, error)
}

func (m *mockOrganizationRepository) CreateOrganization(
	ctx context.Context,
	params db.CreateOrganizationParams,
) (db.Organization, error) {
	m.createOrganizationCalled = true
	m.createOrganizationParams = params

	if m.createOrganizationErr != nil {
		return db.Organization{}, m.createOrganizationErr
	}

	return db.Organization{
		PartyID:            params.PartyID,
		RegistrationNumber: params.RegistrationNumber,
		Industry:           params.Industry,
		OrganizationType:   params.OrganizationType,
		FoundedDate:        params.FoundedDate,
		EmployeeCount:      params.EmployeeCount,
		Description:        params.Description,
	}, nil
}
func (m *mockOrganizationRepository) GetOrganizationByPartyID(
	ctx context.Context,
	partyID pgtype.UUID,
) (db.Organization, error) {
	if m.getOrganizationByPartyIDFunc != nil {
		return m.getOrganizationByPartyIDFunc(ctx, partyID)
	}

	return db.Organization{}, pgx.ErrNoRows
}

func TestCreateOrganization_ValidOrganization(t *testing.T) {
	ctx := context.Background()

	partyID := "550e8400-e29b-41d4-a716-446655440000"
	registrationNumber := "REG-12345"
	foundedDate := "2020-05-10"
	employeeCount := 50
	industry := "Technology"

	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(
			ctx context.Context,
			id pgtype.UUID,
		) (db.Party, error) {
			return db.Party{
				ID:        id,
				PartyType: string(domain.PartyTypeOrganization),
			}, nil
		},
	}

	organizationRepo := &mockOrganizationRepository{}

	service := NewOrganizationService(
		partyRepo,
		organizationRepo,
	)

	input := domain.CreateOrganizationInput{
		PartyID:            partyID,
		RegistrationNumber: &registrationNumber,
		FoundedDate:        &foundedDate,
		EmployeeCount:      &employeeCount,
		Industry:           &industry,
	}

	result, err := service.CreateOrganization(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !organizationRepo.createOrganizationCalled {
		t.Fatal("expected CreateOrganization repository method to be called")
	}

	if result.PartyID.String() != partyID {
		t.Fatalf(
			"expected party ID %s, got %s",
			partyID,
			result.PartyID.String(),
		)
	}

	if result.RegistrationNumber == nil {
		t.Fatal("expected registration number")
	}

	if *result.RegistrationNumber != registrationNumber {
		t.Fatalf(
			"expected registration number %q, got %q",
			registrationNumber,
			*result.RegistrationNumber,
		)
	}

	if result.EmployeeCount == nil {
		t.Fatal("expected employee count")
	}

	if *result.EmployeeCount != int32(employeeCount) {
		t.Fatalf(
			"expected employee count %d, got %d",
			employeeCount,
			*result.EmployeeCount,
		)
	}

	if !result.FoundedDate.Valid {
		t.Fatal("expected founded date to be valid")
	}

	expectedDate, err := time.Parse("2006-01-02", foundedDate)
	if err != nil {
		t.Fatalf("failed to parse expected date: %v", err)
	}

	if !result.FoundedDate.Time.Equal(expectedDate) {
		t.Fatalf(
			"expected founded date %v, got %v",
			expectedDate,
			result.FoundedDate.Time,
		)
	}
}

func TestCreateOrganization_PartyNotFound(t *testing.T) {
	ctx := context.Background()

	partyID := "550e8400-e29b-41d4-a716-446655440000"

	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(
			ctx context.Context,
			id pgtype.UUID,
		) (db.Party, error) {
			return db.Party{}, pgx.ErrNoRows
		},
	}

	organizationRepo := &mockOrganizationRepository{}

	service := NewOrganizationService(
		partyRepo,
		organizationRepo,
	)

	input := domain.CreateOrganizationInput{
		PartyID: partyID,
	}

	_, err := service.CreateOrganization(ctx, input)

	if err != ErrPartyNotFound {
		t.Fatalf(
			"expected ErrPartyNotFound, got %v",
			err,
		)
	}

	if organizationRepo.createOrganizationCalled {
		t.Fatal(
			"organization should not be created when party does not exist",
		)
	}
}

func TestCreateOrganization_PartyIsIndividual(t *testing.T) {
	ctx := context.Background()

	partyID := "550e8400-e29b-41d4-a716-446655440000"

	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(
			ctx context.Context,
			id pgtype.UUID,
		) (db.Party, error) {
			return db.Party{
				ID:        id,
				PartyType: string(domain.PartyTypeIndividual),
			}, nil
		},
	}

	organizationRepo := &mockOrganizationRepository{}

	service := NewOrganizationService(
		partyRepo,
		organizationRepo,
	)

	input := domain.CreateOrganizationInput{
		PartyID: partyID,
	}

	_, err := service.CreateOrganization(ctx, input)

	if err != ErrPartyNotOrganization {
		t.Fatalf(
			"expected ErrPartyNotOrganization, got %v",
			err,
		)
	}

	if organizationRepo.createOrganizationCalled {
		t.Fatal(
			"organization should not be created for an individual party",
		)
	}
}

func TestCreateOrganization_AlreadyExists(t *testing.T) {
	ctx := context.Background()

	partyID := "550e8400-e29b-41d4-a716-446655440000"

	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(
			ctx context.Context,
			id pgtype.UUID,
		) (db.Party, error) {
			return db.Party{
				ID:        id,
				PartyType: string(domain.PartyTypeOrganization),
			}, nil
		},
	}

	organizationRepo := &mockOrganizationRepository{
		getOrganizationByPartyIDFunc: func(
			ctx context.Context,
			id pgtype.UUID,
		) (db.Organization, error) {
			return db.Organization{
				PartyID: id,
			}, nil
		},
	}

	service := NewOrganizationService(
		partyRepo,
		organizationRepo,
	)

	input := domain.CreateOrganizationInput{
		PartyID: partyID,
	}

	_, err := service.CreateOrganization(ctx, input)

	if err != ErrOrganizationAlreadyExists {
		t.Fatalf(
			"expected ErrOrganizationAlreadyExists, got %v",
			err,
		)
	}

	if organizationRepo.createOrganizationCalled {
		t.Fatal(
			"organization should not be created when one already exists",
		)
	}
}

func TestCreateOrganization_InvalidEmployeeCount(t *testing.T) {
	ctx := context.Background()

	partyID := "550e8400-e29b-41d4-a716-446655440000"
	employeeCount := -10

	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(
			ctx context.Context,
			id pgtype.UUID,
		) (db.Party, error) {
			return db.Party{
				ID:        id,
				PartyType: string(domain.PartyTypeOrganization),
			}, nil
		},
	}

	organizationRepo := &mockOrganizationRepository{}

	service := NewOrganizationService(
		partyRepo,
		organizationRepo,
	)

	input := domain.CreateOrganizationInput{
		PartyID:       partyID,
		EmployeeCount: &employeeCount,
	}

	_, err := service.CreateOrganization(ctx, input)

	if err != ErrInvalidEmployeeCount {
		t.Fatalf(
			"expected ErrInvalidEmployeeCount, got %v",
			err,
		)
	}

	if organizationRepo.createOrganizationCalled {
		t.Fatal(
			"organization should not be created with invalid employee count",
		)
	}
}

func TestCreateOrganization_InvalidFoundedDate(t *testing.T) {
	ctx := context.Background()

	partyID := "550e8400-e29b-41d4-a716-446655440000"
	foundedDate := "2024-99-99"

	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{
				ID:        id,
				PartyType: string(domain.PartyTypeOrganization),
			}, nil
		},
	}

	organizationRepo := &mockOrganizationRepository{}

	service := NewOrganizationService(
		partyRepo,
		organizationRepo,
	)

	input := domain.CreateOrganizationInput{
		PartyID:     partyID,
		FoundedDate: &foundedDate,
	}

	_, err := service.CreateOrganization(ctx, input)

	if err != ErrInvalidFoundedDate {
		t.Fatalf("expected ErrInvalidFoundedDate, got %v", err)
	}

	if organizationRepo.createOrganizationCalled {
		t.Fatal("expected organization not to be created")
	}
}

func TestCreateOrganization_InvalidPartyID(t *testing.T) {
	ctx := context.Background()

	partyRepo := &mockPartyRepository{}
	organizationRepo := &mockOrganizationRepository{}

	service := NewOrganizationService(
		partyRepo,
		organizationRepo,
	)

	input := domain.CreateOrganizationInput{
		PartyID: "not-a-valid-uuid",
	}

	_, err := service.CreateOrganization(ctx, input)

	if err != ErrInvalidPartyID {
		t.Fatalf("expected ErrInvalidPartyID, got %v", err)
	}

	if organizationRepo.createOrganizationCalled {
		t.Fatal("expected organization not to be created")
	}
}

func TestCreateOrganization_PartyRepositoryError(t *testing.T) {
	ctx := context.Background()

	partyID := "550e8400-e29b-41d4-a716-446655440000"
	expectedErr := errors.New("database connection failed")

	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{}, expectedErr
		},
	}

	organizationRepo := &mockOrganizationRepository{}

	service := NewOrganizationService(
		partyRepo,
		organizationRepo,
	)

	input := domain.CreateOrganizationInput{
		PartyID: partyID,
	}

	_, err := service.CreateOrganization(ctx, input)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected repository error, got %v", err)
	}

	if organizationRepo.createOrganizationCalled {
		t.Fatal("expected organization not to be created")
	}
}

func TestCreateOrganization_OrganizationRepositoryGetError(t *testing.T) {
	ctx := context.Background()

	partyID := "550e8400-e29b-41d4-a716-446655440000"
	expectedErr := errors.New("database query failed")

	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{
				ID:        id,
				PartyType: string(domain.PartyTypeOrganization),
			}, nil
		},
	}

	organizationRepo := &mockOrganizationRepository{
		getOrganizationByPartyIDFunc: func(ctx context.Context, partyID pgtype.UUID) (db.Organization, error) {
			return db.Organization{}, expectedErr
		},
	}

	service := NewOrganizationService(
		partyRepo,
		organizationRepo,
	)

	input := domain.CreateOrganizationInput{
		PartyID: partyID,
	}

	_, err := service.CreateOrganization(ctx, input)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected repository error, got %v", err)
	}

	if organizationRepo.createOrganizationCalled {
		t.Fatal("expected organization not to be created")
	}
}

func TestCreateOrganization_CreateRepositoryError(t *testing.T) {
	ctx := context.Background()

	partyID := "550e8400-e29b-41d4-a716-446655440000"
	expectedErr := errors.New("failed to insert organization")

	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{
				ID:        id,
				PartyType: string(domain.PartyTypeOrganization),
			}, nil
		},
	}

	organizationRepo := &mockOrganizationRepository{
		createOrganizationErr: expectedErr,
	}

	service := NewOrganizationService(
		partyRepo,
		organizationRepo,
	)

	input := domain.CreateOrganizationInput{
		PartyID: partyID,
	}

	_, err := service.CreateOrganization(ctx, input)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected repository error, got %v", err)
	}

	if !organizationRepo.createOrganizationCalled {
		t.Fatal("expected organization repository to be called")
	}
}

func TestGetOrganizationByPartyID_Success(t *testing.T) {
	partyID := pgtype.UUID{
		Bytes: [16]byte{1},
		Valid: true,
	}

	organizationID := pgtype.UUID{
		Bytes: [16]byte{2},
		Valid: true,
	}

	mockRepository := &mockOrganizationRepository{
		getOrganizationByPartyIDFunc: func(
			ctx context.Context,
			id pgtype.UUID,
		) (db.Organization, error) {
			if id != partyID {
				t.Fatalf("unexpected party ID")
			}

			return db.Organization{
				ID:      organizationID,
				PartyID: partyID,
			}, nil
		},
	}

	mockPartyRepository := &mockPartyRepository{}

	service := NewOrganizationService(
		mockPartyRepository,
		mockRepository,
	)

	organization, err := service.GetOrganizationByPartyID(
		context.Background(),
		partyID,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if organization.ID != organizationID {
		t.Fatalf("unexpected organization ID")
	}
}

func TestGetOrganizationByPartyID_NotFound(t *testing.T) {
	partyID := pgtype.UUID{
		Bytes: [16]byte{1},
		Valid: true,
	}

	mockRepository := &mockOrganizationRepository{
		getOrganizationByPartyIDFunc: func(
			ctx context.Context,
			id pgtype.UUID,
		) (db.Organization, error) {
			return db.Organization{}, pgx.ErrNoRows
		},
	}

	mockPartyRepository := &mockPartyRepository{}

	service := NewOrganizationService(
		mockPartyRepository,
		mockRepository,
	)

	_, err := service.GetOrganizationByPartyID(
		context.Background(),
		partyID,
	)

	if !errors.Is(err, ErrOrganizationNotFound) {
		t.Fatalf(
			"expected ErrOrganizationNotFound, got %v",
			err,
		)
	}
}

func TestGetOrganizationByPartyID_RepositoryError(t *testing.T) {
	partyID := pgtype.UUID{
		Bytes: [16]byte{1},
		Valid: true,
	}

	repositoryErr := errors.New("database connection failed")

	mockRepository := &mockOrganizationRepository{
		getOrganizationByPartyIDFunc: func(
			ctx context.Context,
			id pgtype.UUID,
		) (db.Organization, error) {
			return db.Organization{}, repositoryErr
		},
	}

	mockPartyRepository := &mockPartyRepository{}

	service := NewOrganizationService(
		mockPartyRepository,
		mockRepository,
	)

	_, err := service.GetOrganizationByPartyID(
		context.Background(),
		partyID,
	)

	if !errors.Is(err, repositoryErr) {
		t.Fatalf(
			"expected repository error, got %v",
			err,
		)
	}
}
