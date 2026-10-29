package service

import (
	"context"
	"errors"
	"testing"

	"github.com/group-project/crm_party_service/internal/db"
	"github.com/group-project/crm_party_service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type mockPartyRepository struct {
	createPartyCalled bool
	createPartyParams db.CreatePartyParams

	listPartiesCalled bool
	listPartiesParams db.ListPartiesParams

	getPartyByIDFunc func(
		ctx context.Context,
		id pgtype.UUID,
	) (db.Party, error)

	updatePartyFunc func(
		ctx context.Context,
		params db.UpdatePartyParams,
	) (db.Party, error)

	updatePartyStatusFunc func(
		ctx context.Context,
		params db.UpdatePartyStatusParams,
	) (db.Party, error)
}

func (m *mockPartyRepository) ListParties(
	ctx context.Context,
	params db.ListPartiesParams,
) ([]db.Party, error) {
	m.listPartiesCalled = true
	m.listPartiesParams = params

	return []db.Party{}, nil
}

func (m *mockPartyRepository) CreateParty(
	ctx context.Context,
	params db.CreatePartyParams,
) (db.Party, error) {
	m.createPartyCalled = true
	m.createPartyParams = params

	return db.Party{
		ID:          pgtype.UUID{Valid: true},
		PartyType:   params.PartyType,
		DisplayName: params.DisplayName,
		IsActive:    true,
	}, nil
}

func (m *mockPartyRepository) GetPartyByID(
	ctx context.Context,
	id pgtype.UUID,
) (db.Party, error) {
	if m.getPartyByIDFunc != nil {
		return m.getPartyByIDFunc(ctx, id)
	}

	return db.Party{}, nil
}

func (m *mockPartyRepository) UpdateParty(
	ctx context.Context,
	params db.UpdatePartyParams,
) (db.Party, error) {
	if m.updatePartyFunc != nil {
		return m.updatePartyFunc(ctx, params)
	}

	return db.Party{
		ID:          params.ID,
		DisplayName: params.DisplayName,
		LegalName:   params.LegalName,
		Email:       params.Email,
		Phone:       params.Phone,
		Website:     params.Website,
		TaxNumber:   params.TaxNumber,
		Notes:       params.Notes,
		IsActive:    true,
	}, nil
}

func (m *mockPartyRepository) UpdatePartyStatus(
	ctx context.Context,
	params db.UpdatePartyStatusParams,
) (db.Party, error) {
	if m.updatePartyStatusFunc != nil {
		return m.updatePartyStatusFunc(ctx, params)
	}

	return db.Party{
		ID:       params.ID,
		IsActive: params.IsActive,
	}, nil
}

func TestCreateParty_ValidParty(t *testing.T) {
	repo := &mockPartyRepository{}
	service := NewPartyService(repo)

	input := domain.CreatePartyInput{
		PartyType:   domain.PartyTypeOrganization,
		DisplayName: "  ABC Technologies  ",
	}

	party, err := service.CreateParty(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !repo.createPartyCalled {
		t.Fatal("expected repository CreateParty to be called")
	}

	if party.DisplayName != "ABC Technologies" {
		t.Fatalf(
			"expected trimmed display name %q, got %q",
			"ABC Technologies",
			party.DisplayName,
		)
	}
}

func TestListParties_DefaultPagination(t *testing.T) {
	repo := &mockPartyRepository{}
	service := NewPartyService(repo)

	params := db.ListPartiesParams{}

	_, err := service.ListParties(context.Background(), params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.listPartiesParams.Limit != 20 {
		t.Fatalf(
			"expected default limit 20, got %d",
			repo.listPartiesParams.Limit,
		)
	}

	if repo.listPartiesParams.Offset != 0 {
		t.Fatalf(
			"expected default offset 0, got %d",
			repo.listPartiesParams.Offset,
		)
	}
}

func TestListParties_MaximumLimit(t *testing.T) {
	repo := &mockPartyRepository{}
	service := NewPartyService(repo)

	params := db.ListPartiesParams{
		Limit: 500,
	}

	_, err := service.ListParties(context.Background(), params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.listPartiesParams.Limit != 100 {
		t.Fatalf(
			"expected maximum limit 100, got %d",
			repo.listPartiesParams.Limit,
		)
	}
}

func TestListParties_NegativeOffset(t *testing.T) {
	repo := &mockPartyRepository{}
	service := NewPartyService(repo)

	params := db.ListPartiesParams{
		Offset: -1,
		Limit:  20,
	}

	_, err := service.ListParties(context.Background(), params)

	if err != ErrInvalidPagination {
		t.Fatalf(
			"expected ErrInvalidPagination, got %v",
			err,
		)
	}
}

func TestListParties_EmptySearch(t *testing.T) {
	repo := &mockPartyRepository{}
	service := NewPartyService(repo)

	search := "   "

	params := db.ListPartiesParams{
		Search: &search,
		Limit:  20,
	}

	_, err := service.ListParties(context.Background(), params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.listPartiesParams.Search != nil {
		t.Fatal("expected empty search to become nil")
	}
}

func TestCreateParty_InvalidPartyType(t *testing.T) {
	repo := &mockPartyRepository{}
	service := NewPartyService(repo)

	input := domain.CreatePartyInput{
		PartyType:   domain.PartyType("CUSTOMER"),
		DisplayName: "ABC Technologies",
	}

	_, err := service.CreateParty(context.Background(), input)

	if err != ErrInvalidPartyType {
		t.Fatalf(
			"expected ErrInvalidPartyType, got %v",
			err,
		)
	}

	if repo.createPartyCalled {
		t.Fatal("repository should not be called for invalid party type")
	}
}

func TestCreateParty_EmptyDisplayName(t *testing.T) {
	repo := &mockPartyRepository{}
	service := NewPartyService(repo)

	input := domain.CreatePartyInput{
		PartyType:   domain.PartyTypeIndividual,
		DisplayName: "   ",
	}

	_, err := service.CreateParty(context.Background(), input)

	if err != ErrEmptyDisplayName {
		t.Fatalf(
			"expected ErrEmptyDisplayName, got %v",
			err,
		)
	}

	if repo.createPartyCalled {
		t.Fatal("repository should not be called for empty display name")
	}
}

func TestUpdateParty_Success(t *testing.T) {
	partyUUID := "550e8400-e29b-41d4-a716-446655440000"
	repo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{ID: id, DisplayName: "Old Name"}, nil
		},
	}
	svc := NewPartyService(repo)

	input := domain.UpdatePartyInput{
		ID:          partyUUID,
		DisplayName: "New Name",
	}

	party, err := svc.UpdateParty(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if party.DisplayName != "New Name" {
		t.Fatalf("expected New Name, got %s", party.DisplayName)
	}
}

func TestUpdateParty_InvalidPartyID(t *testing.T) {
	svc := NewPartyService(&mockPartyRepository{})
	_, err := svc.UpdateParty(context.Background(), domain.UpdatePartyInput{
		ID:          "bad-uuid",
		DisplayName: "Name",
	})
	if !errors.Is(err, ErrInvalidPartyID) {
		t.Fatalf("expected ErrInvalidPartyID, got %v", err)
	}
}

func TestUpdateParty_EmptyDisplayName(t *testing.T) {
	svc := NewPartyService(&mockPartyRepository{})
	_, err := svc.UpdateParty(context.Background(), domain.UpdatePartyInput{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		DisplayName: "   ",
	})
	if !errors.Is(err, ErrEmptyDisplayName) {
		t.Fatalf("expected ErrEmptyDisplayName, got %v", err)
	}
}

func TestUpdateParty_NotFound(t *testing.T) {
	repo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{}, pgx.ErrNoRows
		},
	}
	svc := NewPartyService(repo)
	_, err := svc.UpdateParty(context.Background(), domain.UpdatePartyInput{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		DisplayName: "Name",
	})
	if !errors.Is(err, ErrPartyNotFound) {
		t.Fatalf("expected ErrPartyNotFound, got %v", err)
	}
}

func TestUpdatePartyStatus_Success(t *testing.T) {
	var partyUUID pgtype.UUID
	partyUUID.Scan("550e8400-e29b-41d4-a716-446655440000")
	repo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{ID: id, IsActive: true}, nil
		},
	}
	svc := NewPartyService(repo)

	party, err := svc.UpdatePartyStatus(context.Background(), partyUUID, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if party.IsActive != false {
		t.Fatalf("expected false, got %v", party.IsActive)
	}
}

func TestUpdatePartyStatus_NotFound(t *testing.T) {
	var partyUUID pgtype.UUID
	partyUUID.Scan("550e8400-e29b-41d4-a716-446655440000")
	repo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{}, pgx.ErrNoRows
		},
	}
	svc := NewPartyService(repo)

	_, err := svc.UpdatePartyStatus(context.Background(), partyUUID, false)
	if !errors.Is(err, ErrPartyNotFound) {
		t.Fatalf("expected ErrPartyNotFound, got %v", err)
	}
}
