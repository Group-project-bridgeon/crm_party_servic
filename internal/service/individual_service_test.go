package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
	"github.com/group-project/crm_party_service/internal/domain"
)

type mockIndividualRepository struct {
	createIndividualCalled     bool
	createIndividualParams     db.CreateIndividualParams
	getIndividualByPartyIDFunc func(ctx context.Context, partyID pgtype.UUID) (db.Individual, error)
	updateIndividualCalled     bool
	updateIndividualParams     db.UpdateIndividualParams
}

func (m *mockIndividualRepository) CreateIndividual(
	ctx context.Context,
	params db.CreateIndividualParams,
) (db.Individual, error) {
	m.createIndividualCalled = true
	m.createIndividualParams = params
	return db.Individual{
		ID:          pgtype.UUID{Valid: true},
		PartyID:     params.PartyID,
		FirstName:   params.FirstName,
		MiddleName:  params.MiddleName,
		LastName:    params.LastName,
		DateOfBirth: params.DateOfBirth,
		Gender:      params.Gender,
		JobTitle:    params.JobTitle,
		NationalID:  params.NationalID,
	}, nil
}

func (m *mockIndividualRepository) GetIndividualByPartyID(
	ctx context.Context,
	partyID pgtype.UUID,
) (db.Individual, error) {
	if m.getIndividualByPartyIDFunc != nil {
		return m.getIndividualByPartyIDFunc(ctx, partyID)
	}
	return db.Individual{}, pgx.ErrNoRows
}

func (m *mockIndividualRepository) UpdateIndividual(
	ctx context.Context,
	params db.UpdateIndividualParams,
) (db.Individual, error) {
	m.updateIndividualCalled = true
	m.updateIndividualParams = params
	return db.Individual{
		PartyID:     params.PartyID,
		FirstName:   params.FirstName,
		MiddleName:  params.MiddleName,
		LastName:    params.LastName,
		DateOfBirth: params.DateOfBirth,
		Gender:      params.Gender,
		JobTitle:    params.JobTitle,
		NationalID:  params.NationalID,
	}, nil
}

func TestCreateIndividual_Success(t *testing.T) {
	ctx := context.Background()
	partyUUID := "550e8400-e29b-41d4-a716-446655440000"

	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{
				ID:        id,
				PartyType: string(domain.PartyTypeIndividual),
			}, nil
		},
	}
	indivRepo := &mockIndividualRepository{}

	svc := NewIndividualService(partyRepo, indivRepo)

	dob := "1990-01-15"
	input := domain.CreateIndividualInput{
		PartyID:     partyUUID,
		FirstName:   "Jane",
		DateOfBirth: &dob,
	}

	individual, err := svc.CreateIndividual(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if individual.FirstName != "Jane" {
		t.Fatalf("expected first name Jane, got %s", individual.FirstName)
	}
	if !indivRepo.createIndividualCalled {
		t.Fatal("expected CreateIndividual on repository to be called")
	}
}

func TestCreateIndividual_InvalidPartyID(t *testing.T) {
	svc := NewIndividualService(&mockPartyRepository{}, &mockIndividualRepository{})
	_, err := svc.CreateIndividual(context.Background(), domain.CreateIndividualInput{
		PartyID:   "invalid-uuid",
		FirstName: "John",
	})
	if !errors.Is(err, ErrInvalidPartyID) {
		t.Fatalf("expected ErrInvalidPartyID, got %v", err)
	}
}

func TestCreateIndividual_EmptyFirstName(t *testing.T) {
	svc := NewIndividualService(&mockPartyRepository{}, &mockIndividualRepository{})
	_, err := svc.CreateIndividual(context.Background(), domain.CreateIndividualInput{
		PartyID:   "550e8400-e29b-41d4-a716-446655440000",
		FirstName: "   ",
	})
	if !errors.Is(err, ErrEmptyFirstName) {
		t.Fatalf("expected ErrEmptyFirstName, got %v", err)
	}
}

func TestCreateIndividual_PartyNotFound(t *testing.T) {
	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{}, pgx.ErrNoRows
		},
	}
	svc := NewIndividualService(partyRepo, &mockIndividualRepository{})
	_, err := svc.CreateIndividual(context.Background(), domain.CreateIndividualInput{
		PartyID:   "550e8400-e29b-41d4-a716-446655440000",
		FirstName: "John",
	})
	if !errors.Is(err, ErrPartyNotFound) {
		t.Fatalf("expected ErrPartyNotFound, got %v", err)
	}
}

func TestCreateIndividual_PartyNotIndividual(t *testing.T) {
	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{
				ID:        id,
				PartyType: string(domain.PartyTypeOrganization),
			}, nil
		},
	}
	svc := NewIndividualService(partyRepo, &mockIndividualRepository{})
	_, err := svc.CreateIndividual(context.Background(), domain.CreateIndividualInput{
		PartyID:   "550e8400-e29b-41d4-a716-446655440000",
		FirstName: "John",
	})
	if !errors.Is(err, ErrPartyNotIndividual) {
		t.Fatalf("expected ErrPartyNotIndividual, got %v", err)
	}
}

func TestCreateIndividual_AlreadyExists(t *testing.T) {
	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{
				ID:        id,
				PartyType: string(domain.PartyTypeIndividual),
			}, nil
		},
	}
	indivRepo := &mockIndividualRepository{
		getIndividualByPartyIDFunc: func(ctx context.Context, partyID pgtype.UUID) (db.Individual, error) {
			return db.Individual{PartyID: partyID}, nil
		},
	}
	svc := NewIndividualService(partyRepo, indivRepo)
	_, err := svc.CreateIndividual(context.Background(), domain.CreateIndividualInput{
		PartyID:   "550e8400-e29b-41d4-a716-446655440000",
		FirstName: "John",
	})
	if !errors.Is(err, ErrIndividualAlreadyExists) {
		t.Fatalf("expected ErrIndividualAlreadyExists, got %v", err)
	}
}

func TestCreateIndividual_InvalidDateOfBirth(t *testing.T) {
	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{
				ID:        id,
				PartyType: string(domain.PartyTypeIndividual),
			}, nil
		},
	}
	svc := NewIndividualService(partyRepo, &mockIndividualRepository{})
	badDob := "not-a-date"
	_, err := svc.CreateIndividual(context.Background(), domain.CreateIndividualInput{
		PartyID:     "550e8400-e29b-41d4-a716-446655440000",
		FirstName:   "John",
		DateOfBirth: &badDob,
	})
	if !errors.Is(err, ErrInvalidDateOfBirth) {
		t.Fatalf("expected ErrInvalidDateOfBirth, got %v", err)
	}
}

func TestGetIndividualByPartyID_Success(t *testing.T) {
	var partyID pgtype.UUID
	partyID.Scan("550e8400-e29b-41d4-a716-446655440000")

	indivRepo := &mockIndividualRepository{
		getIndividualByPartyIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Individual, error) {
			return db.Individual{PartyID: id, FirstName: "Alice"}, nil
		},
	}
	svc := NewIndividualService(&mockPartyRepository{}, indivRepo)
	indiv, err := svc.GetIndividualByPartyID(context.Background(), partyID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if indiv.FirstName != "Alice" {
		t.Fatalf("expected first name Alice, got %s", indiv.FirstName)
	}
}

func TestGetIndividualByPartyID_NotFound(t *testing.T) {
	var partyID pgtype.UUID
	partyID.Scan("550e8400-e29b-41d4-a716-446655440000")

	svc := NewIndividualService(&mockPartyRepository{}, &mockIndividualRepository{})
	_, err := svc.GetIndividualByPartyID(context.Background(), partyID)
	if !errors.Is(err, ErrIndividualNotFound) {
		t.Fatalf("expected ErrIndividualNotFound, got %v", err)
	}
}

func TestUpdateIndividual_Success(t *testing.T) {
	partyUUID := "550e8400-e29b-41d4-a716-446655440000"
	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{
				ID:        id,
				PartyType: string(domain.PartyTypeIndividual),
			}, nil
		},
	}
	indivRepo := &mockIndividualRepository{
		getIndividualByPartyIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Individual, error) {
			return db.Individual{PartyID: id, FirstName: "OldName"}, nil
		},
	}

	svc := NewIndividualService(partyRepo, indivRepo)
	input := domain.UpdateIndividualInput{
		PartyID:   partyUUID,
		FirstName: "NewName",
	}

	indiv, err := svc.UpdateIndividual(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if indiv.FirstName != "NewName" {
		t.Fatalf("expected NewName, got %s", indiv.FirstName)
	}
	if !indivRepo.updateIndividualCalled {
		t.Fatal("expected update to be called")
	}
}
