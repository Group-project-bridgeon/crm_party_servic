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

type mockSupplierRepository struct {
	createSupplierCalled                bool
	createSupplierParams                db.CreateSupplierParams
	getSupplierByIDFunc                 func(ctx context.Context, id pgtype.UUID) (db.Supplier, error)
	getSupplierByPartyIDFunc            func(ctx context.Context, partyID pgtype.UUID) (db.Supplier, error)
	getSupplierByCodeFunc               func(ctx context.Context, code string) (db.Supplier, error)
	listSuppliersFunc                   func(ctx context.Context, params db.ListSuppliersParams) ([]db.Supplier, error)
	updateSupplierFunc                  func(ctx context.Context, params db.UpdateSupplierParams) (db.Supplier, error)
	createSupplierUserAccountFunc       func(ctx context.Context, params db.CreateSupplierUserAccountParams) (db.SupplierUserAccount, error)
	listSupplierUserAccountsFunc        func(ctx context.Context, supplierID pgtype.UUID) ([]db.SupplierUserAccount, error)
	updateSupplierUserAccountStatusFunc func(ctx context.Context, params db.UpdateSupplierUserAccountStatusParams) (db.SupplierUserAccount, error)
}

func (m *mockSupplierRepository) CreateSupplier(ctx context.Context, params db.CreateSupplierParams) (db.Supplier, error) {
	m.createSupplierCalled = true
	m.createSupplierParams = params
	return db.Supplier{
		ID:               pgtype.UUID{Valid: true},
		PartyID:          params.PartyID,
		SupplierCode:     params.SupplierCode,
		SupplierGroup:    params.SupplierGroup,
		PaymentTermsDays: params.PaymentTermsDays,
		PaymentMethod:    params.PaymentMethod,
		CurrencyCode:     params.CurrencyCode,
		SupplierRating:   params.SupplierRating,
	}, nil
}

func (m *mockSupplierRepository) GetSupplierByID(ctx context.Context, id pgtype.UUID) (db.Supplier, error) {
	if m.getSupplierByIDFunc != nil {
		return m.getSupplierByIDFunc(ctx, id)
	}
	return db.Supplier{}, pgx.ErrNoRows
}

func (m *mockSupplierRepository) GetSupplierByPartyID(ctx context.Context, partyID pgtype.UUID) (db.Supplier, error) {
	if m.getSupplierByPartyIDFunc != nil {
		return m.getSupplierByPartyIDFunc(ctx, partyID)
	}
	return db.Supplier{}, pgx.ErrNoRows
}

func (m *mockSupplierRepository) GetSupplierByCode(ctx context.Context, code string) (db.Supplier, error) {
	if m.getSupplierByCodeFunc != nil {
		return m.getSupplierByCodeFunc(ctx, code)
	}
	return db.Supplier{}, pgx.ErrNoRows
}

func (m *mockSupplierRepository) ListSuppliers(ctx context.Context, params db.ListSuppliersParams) ([]db.Supplier, error) {
	if m.listSuppliersFunc != nil {
		return m.listSuppliersFunc(ctx, params)
	}
	return []db.Supplier{}, nil
}

func (m *mockSupplierRepository) UpdateSupplier(ctx context.Context, params db.UpdateSupplierParams) (db.Supplier, error) {
	if m.updateSupplierFunc != nil {
		return m.updateSupplierFunc(ctx, params)
	}
	return db.Supplier{
		ID:               params.ID,
		SupplierGroup:    params.SupplierGroup,
		PaymentTermsDays: params.PaymentTermsDays,
		PaymentMethod:    params.PaymentMethod,
		CurrencyCode:     params.CurrencyCode,
		SupplierRating:   params.SupplierRating,
	}, nil
}

func (m *mockSupplierRepository) CreateSupplierUserAccount(ctx context.Context, params db.CreateSupplierUserAccountParams) (db.SupplierUserAccount, error) {
	if m.createSupplierUserAccountFunc != nil {
		return m.createSupplierUserAccountFunc(ctx, params)
	}
	return db.SupplierUserAccount{
		ID:         pgtype.UUID{Valid: true},
		SupplierID: params.SupplierID,
		UserID:     params.UserID,
		IsActive:   true,
	}, nil
}

func (m *mockSupplierRepository) ListSupplierUserAccountsBySupplierID(ctx context.Context, supplierID pgtype.UUID) ([]db.SupplierUserAccount, error) {
	if m.listSupplierUserAccountsFunc != nil {
		return m.listSupplierUserAccountsFunc(ctx, supplierID)
	}
	return []db.SupplierUserAccount{}, nil
}

func (m *mockSupplierRepository) UpdateSupplierUserAccountStatus(ctx context.Context, params db.UpdateSupplierUserAccountStatusParams) (db.SupplierUserAccount, error) {
	if m.updateSupplierUserAccountStatusFunc != nil {
		return m.updateSupplierUserAccountStatusFunc(ctx, params)
	}
	return db.SupplierUserAccount{
		ID:       params.ID,
		IsActive: params.IsActive,
	}, nil
}

func TestCreateSupplier_Success(t *testing.T) {
	partyUUID := "550e8400-e29b-41d4-a716-446655440000"
	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{ID: id, DisplayName: "Supplier Inc"}, nil
		},
	}
	suppRepo := &mockSupplierRepository{}
	svc := NewSupplierService(partyRepo, suppRepo)

	rating := 4.5
	currency := "USD"
	terms := 45
	input := domain.CreateSupplierInput{
		PartyID:          partyUUID,
		SupplierCode:     "SUPP-001",
		SupplierRating:   &rating,
		CurrencyCode:     &currency,
		PaymentTermsDays: &terms,
	}

	supplier, err := svc.CreateSupplier(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if supplier.SupplierCode != "SUPP-001" {
		t.Fatalf("expected SUPP-001, got %s", supplier.SupplierCode)
	}
	if !suppRepo.createSupplierCalled {
		t.Fatal("expected CreateSupplier to be called")
	}
}

func TestCreateSupplier_ValidationErrors(t *testing.T) {
	partyUUID := "550e8400-e29b-41d4-a716-446655440000"
	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{ID: id}, nil
		},
	}
	suppRepo := &mockSupplierRepository{}
	svc := NewSupplierService(partyRepo, suppRepo)

	// Empty supplier code
	_, err := svc.CreateSupplier(context.Background(), domain.CreateSupplierInput{
		PartyID:      partyUUID,
		SupplierCode: "  ",
	})
	if !errors.Is(err, ErrEmptySupplierCode) {
		t.Fatalf("expected ErrEmptySupplierCode, got %v", err)
	}

	// Negative payment terms
	negTerms := -5
	_, err = svc.CreateSupplier(context.Background(), domain.CreateSupplierInput{
		PartyID:          partyUUID,
		SupplierCode:     "SUPP-002",
		PaymentTermsDays: &negTerms,
	})
	if !errors.Is(err, ErrInvalidPaymentTerms) {
		t.Fatalf("expected ErrInvalidPaymentTerms, got %v", err)
	}

	// Invalid currency
	badCurrency := "USDD"
	_, err = svc.CreateSupplier(context.Background(), domain.CreateSupplierInput{
		PartyID:      partyUUID,
		SupplierCode: "SUPP-002",
		CurrencyCode: &badCurrency,
	})
	if !errors.Is(err, ErrInvalidCurrencyCode) {
		t.Fatalf("expected ErrInvalidCurrencyCode, got %v", err)
	}

	// Invalid rating (<0 or >5)
	badRating := 5.5
	_, err = svc.CreateSupplier(context.Background(), domain.CreateSupplierInput{
		PartyID:        partyUUID,
		SupplierCode:   "SUPP-002",
		SupplierRating: &badRating,
	})
	if !errors.Is(err, ErrInvalidSupplierRating) {
		t.Fatalf("expected ErrInvalidSupplierRating, got %v", err)
	}
}

func TestCreateSupplier_DuplicateChecks(t *testing.T) {
	partyUUID := "550e8400-e29b-41d4-a716-446655440000"
	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{ID: id}, nil
		},
	}

	// Supplier already exists for party
	suppRepo1 := &mockSupplierRepository{
		getSupplierByPartyIDFunc: func(ctx context.Context, partyID pgtype.UUID) (db.Supplier, error) {
			return db.Supplier{PartyID: partyID}, nil
		},
	}
	svc1 := NewSupplierService(partyRepo, suppRepo1)
	_, err := svc1.CreateSupplier(context.Background(), domain.CreateSupplierInput{
		PartyID:      partyUUID,
		SupplierCode: "SUPP-003",
	})
	if !errors.Is(err, ErrSupplierAlreadyExists) {
		t.Fatalf("expected ErrSupplierAlreadyExists, got %v", err)
	}

	// Supplier code already in use
	suppRepo2 := &mockSupplierRepository{
		getSupplierByCodeFunc: func(ctx context.Context, code string) (db.Supplier, error) {
			return db.Supplier{SupplierCode: code}, nil
		},
	}
	svc2 := NewSupplierService(partyRepo, suppRepo2)
	_, err = svc2.CreateSupplier(context.Background(), domain.CreateSupplierInput{
		PartyID:      partyUUID,
		SupplierCode: "SUPP-003",
	})
	if !errors.Is(err, ErrDuplicateSupplierCode) {
		t.Fatalf("expected ErrDuplicateSupplierCode, got %v", err)
	}
}

func TestSupplierUserAccount_Flow(t *testing.T) {
	suppUUID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID := "660e8400-e29b-41d4-a716-446655440000"

	suppRepo := &mockSupplierRepository{
		getSupplierByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Supplier, error) {
			return db.Supplier{ID: id}, nil
		},
	}
	svc := NewSupplierService(&mockPartyRepository{}, suppRepo)

	// Create
	acc, err := svc.CreateSupplierUserAccount(context.Background(), domain.CreateSupplierUserAccountInput{
		SupplierID: suppUUID,
		UserID:     userUUID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !acc.IsActive {
		t.Fatal("expected active account")
	}

	// List
	var sID pgtype.UUID
	sID.Scan(suppUUID)
	accounts, err := svc.ListSupplierUserAccounts(context.Background(), sID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if accounts == nil {
		t.Fatal("expected accounts slice")
	}

	// Update status
	updated, err := svc.UpdateSupplierUserAccountStatus(context.Background(), acc.ID, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.IsActive != false {
		t.Fatal("expected inactive account")
	}
}
