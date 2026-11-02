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

type mockCustomerRepository struct {
	createCustomerCalled                bool
	createCustomerParams                db.CreateCustomerParams
	getCustomerByIDFunc                 func(ctx context.Context, id pgtype.UUID) (db.Customer, error)
	getCustomerByPartyIDFunc            func(ctx context.Context, partyID pgtype.UUID) (db.Customer, error)
	getCustomerByCodeFunc               func(ctx context.Context, code string) (db.Customer, error)
	listCustomersFunc                   func(ctx context.Context, params db.ListCustomersParams) ([]db.Customer, error)
	updateCustomerFunc                  func(ctx context.Context, params db.UpdateCustomerParams) (db.Customer, error)
	createCustomerUserAccountFunc       func(ctx context.Context, params db.CreateCustomerUserAccountParams) (db.CustomerUserAccount, error)
	listCustomerUserAccountsFunc        func(ctx context.Context, customerID pgtype.UUID) ([]db.CustomerUserAccount, error)
	updateCustomerUserAccountStatusFunc func(ctx context.Context, params db.UpdateCustomerUserAccountStatusParams) (db.CustomerUserAccount, error)
}

func (m *mockCustomerRepository) CreateCustomer(ctx context.Context, params db.CreateCustomerParams) (db.Customer, error) {
	m.createCustomerCalled = true
	m.createCustomerParams = params
	return db.Customer{
		ID:               pgtype.UUID{Valid: true},
		PartyID:          params.PartyID,
		CustomerCode:     params.CustomerCode,
		CustomerGroup:    params.CustomerGroup,
		CreditLimit:      params.CreditLimit,
		PaymentTermsDays: params.PaymentTermsDays,
		CreditStatus:     params.CreditStatus,
	}, nil
}

func (m *mockCustomerRepository) GetCustomerByID(ctx context.Context, id pgtype.UUID) (db.Customer, error) {
	if m.getCustomerByIDFunc != nil {
		return m.getCustomerByIDFunc(ctx, id)
	}
	return db.Customer{}, pgx.ErrNoRows
}

func (m *mockCustomerRepository) GetCustomerByPartyID(ctx context.Context, partyID pgtype.UUID) (db.Customer, error) {
	if m.getCustomerByPartyIDFunc != nil {
		return m.getCustomerByPartyIDFunc(ctx, partyID)
	}
	return db.Customer{}, pgx.ErrNoRows
}

func (m *mockCustomerRepository) GetCustomerByCode(ctx context.Context, code string) (db.Customer, error) {
	if m.getCustomerByCodeFunc != nil {
		return m.getCustomerByCodeFunc(ctx, code)
	}
	return db.Customer{}, pgx.ErrNoRows
}

func (m *mockCustomerRepository) ListCustomers(ctx context.Context, params db.ListCustomersParams) ([]db.Customer, error) {
	if m.listCustomersFunc != nil {
		return m.listCustomersFunc(ctx, params)
	}
	return []db.Customer{}, nil
}

func (m *mockCustomerRepository) UpdateCustomer(ctx context.Context, params db.UpdateCustomerParams) (db.Customer, error) {
	if m.updateCustomerFunc != nil {
		return m.updateCustomerFunc(ctx, params)
	}
	return db.Customer{
		ID:               params.ID,
		CustomerGroup:    params.CustomerGroup,
		CreditLimit:      params.CreditLimit,
		PaymentTermsDays: params.PaymentTermsDays,
		CreditStatus:     params.CreditStatus,
	}, nil
}

func (m *mockCustomerRepository) CreateCustomerUserAccount(ctx context.Context, params db.CreateCustomerUserAccountParams) (db.CustomerUserAccount, error) {
	if m.createCustomerUserAccountFunc != nil {
		return m.createCustomerUserAccountFunc(ctx, params)
	}
	return db.CustomerUserAccount{
		ID:         pgtype.UUID{Valid: true},
		CustomerID: params.CustomerID,
		UserID:     params.UserID,
		IsActive:   true,
	}, nil
}

func (m *mockCustomerRepository) ListCustomerUserAccountsByCustomerID(ctx context.Context, customerID pgtype.UUID) ([]db.CustomerUserAccount, error) {
	if m.listCustomerUserAccountsFunc != nil {
		return m.listCustomerUserAccountsFunc(ctx, customerID)
	}
	return []db.CustomerUserAccount{}, nil
}

func (m *mockCustomerRepository) UpdateCustomerUserAccountStatus(ctx context.Context, params db.UpdateCustomerUserAccountStatusParams) (db.CustomerUserAccount, error) {
	if m.updateCustomerUserAccountStatusFunc != nil {
		return m.updateCustomerUserAccountStatusFunc(ctx, params)
	}
	return db.CustomerUserAccount{
		ID:       params.ID,
		IsActive: params.IsActive,
	}, nil
}

func TestCreateCustomer_Success(t *testing.T) {
	partyUUID := "550e8400-e29b-41d4-a716-446655440000"
	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{ID: id, DisplayName: "Acme Corp"}, nil
		},
	}
	custRepo := &mockCustomerRepository{}
	svc := NewCustomerService(partyRepo, custRepo)

	limit := 10000.50
	terms := 30
	status := "ACTIVE"
	input := domain.CreateCustomerInput{
		PartyID:          partyUUID,
		CustomerCode:     "CUST-001",
		CreditLimit:      &limit,
		PaymentTermsDays: &terms,
		CreditStatus:     &status,
	}

	cust, err := svc.CreateCustomer(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cust.CustomerCode != "CUST-001" {
		t.Fatalf("expected CUST-001, got %s", cust.CustomerCode)
	}
	if !custRepo.createCustomerCalled {
		t.Fatal("expected CreateCustomer to be called")
	}
}

func TestCreateCustomer_ValidationErrors(t *testing.T) {
	partyUUID := "550e8400-e29b-41d4-a716-446655440000"
	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{ID: id}, nil
		},
	}
	custRepo := &mockCustomerRepository{}
	svc := NewCustomerService(partyRepo, custRepo)

	// Empty customer code
	_, err := svc.CreateCustomer(context.Background(), domain.CreateCustomerInput{
		PartyID:      partyUUID,
		CustomerCode: "  ",
	})
	if !errors.Is(err, ErrEmptyCustomerCode) {
		t.Fatalf("expected ErrEmptyCustomerCode, got %v", err)
	}

	// Negative credit limit
	negLimit := -10.0
	_, err = svc.CreateCustomer(context.Background(), domain.CreateCustomerInput{
		PartyID:      partyUUID,
		CustomerCode: "CUST-002",
		CreditLimit:  &negLimit,
	})
	if !errors.Is(err, ErrInvalidCreditLimit) {
		t.Fatalf("expected ErrInvalidCreditLimit, got %v", err)
	}

	// Negative payment terms
	negTerms := -5
	_, err = svc.CreateCustomer(context.Background(), domain.CreateCustomerInput{
		PartyID:          partyUUID,
		CustomerCode:     "CUST-002",
		PaymentTermsDays: &negTerms,
	})
	if !errors.Is(err, ErrInvalidPaymentTerms) {
		t.Fatalf("expected ErrInvalidPaymentTerms, got %v", err)
	}

	// Invalid credit status
	badStatus := "INVALID"
	_, err = svc.CreateCustomer(context.Background(), domain.CreateCustomerInput{
		PartyID:      partyUUID,
		CustomerCode: "CUST-002",
		CreditStatus: &badStatus,
	})
	if !errors.Is(err, ErrInvalidCreditStatus) {
		t.Fatalf("expected ErrInvalidCreditStatus, got %v", err)
	}
}

func TestCreateCustomer_DuplicateChecks(t *testing.T) {
	partyUUID := "550e8400-e29b-41d4-a716-446655440000"
	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{ID: id}, nil
		},
	}

	// Customer already exists for party
	custRepo1 := &mockCustomerRepository{
		getCustomerByPartyIDFunc: func(ctx context.Context, partyID pgtype.UUID) (db.Customer, error) {
			return db.Customer{PartyID: partyID}, nil
		},
	}
	svc1 := NewCustomerService(partyRepo, custRepo1)
	_, err := svc1.CreateCustomer(context.Background(), domain.CreateCustomerInput{
		PartyID:      partyUUID,
		CustomerCode: "CUST-003",
	})
	if !errors.Is(err, ErrCustomerAlreadyExists) {
		t.Fatalf("expected ErrCustomerAlreadyExists, got %v", err)
	}

	// Customer code already in use
	custRepo2 := &mockCustomerRepository{
		getCustomerByCodeFunc: func(ctx context.Context, code string) (db.Customer, error) {
			return db.Customer{CustomerCode: code}, nil
		},
	}
	svc2 := NewCustomerService(partyRepo, custRepo2)
	_, err = svc2.CreateCustomer(context.Background(), domain.CreateCustomerInput{
		PartyID:      partyUUID,
		CustomerCode: "CUST-003",
	})
	if !errors.Is(err, ErrDuplicateCustomerCode) {
		t.Fatalf("expected ErrDuplicateCustomerCode, got %v", err)
	}
}

func TestGetCustomerByID_SuccessAndNotFound(t *testing.T) {
	var custID pgtype.UUID
	custID.Scan("550e8400-e29b-41d4-a716-446655440000")

	custRepo := &mockCustomerRepository{
		getCustomerByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Customer, error) {
			return db.Customer{ID: id, CustomerCode: "CUST-100"}, nil
		},
	}
	svc := NewCustomerService(&mockPartyRepository{}, custRepo)

	cust, err := svc.GetCustomerByID(context.Background(), custID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cust.CustomerCode != "CUST-100" {
		t.Fatalf("expected CUST-100, got %s", cust.CustomerCode)
	}

	// Not found
	custRepo.getCustomerByIDFunc = nil
	_, err = svc.GetCustomerByID(context.Background(), custID)
	if !errors.Is(err, ErrCustomerNotFound) {
		t.Fatalf("expected ErrCustomerNotFound, got %v", err)
	}
}

func TestUpdateCustomer_Success(t *testing.T) {
	custUUID := "550e8400-e29b-41d4-a716-446655440000"
	custRepo := &mockCustomerRepository{
		getCustomerByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Customer, error) {
			return db.Customer{ID: id, CustomerCode: "CUST-100"}, nil
		},
	}
	svc := NewCustomerService(&mockPartyRepository{}, custRepo)

	group := "VIP"
	cust, err := svc.UpdateCustomer(context.Background(), domain.UpdateCustomerInput{
		ID:            custUUID,
		CustomerGroup: &group,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cust.CustomerGroup == nil || *cust.CustomerGroup != "VIP" {
		t.Fatalf("expected VIP, got %v", cust.CustomerGroup)
	}
}

func TestCustomerUserAccount_Flow(t *testing.T) {
	custUUID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID := "660e8400-e29b-41d4-a716-446655440000"

	custRepo := &mockCustomerRepository{
		getCustomerByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Customer, error) {
			return db.Customer{ID: id}, nil
		},
	}
	svc := NewCustomerService(&mockPartyRepository{}, custRepo)

	// Create
	acc, err := svc.CreateCustomerUserAccount(context.Background(), domain.CreateCustomerUserAccountInput{
		CustomerID: custUUID,
		UserID:     userUUID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !acc.IsActive {
		t.Fatal("expected active account")
	}

	// List
	var cID pgtype.UUID
	cID.Scan(custUUID)
	accounts, err := svc.ListCustomerUserAccounts(context.Background(), cID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if accounts == nil {
		t.Fatal("expected accounts slice")
	}

	// Update status
	updated, err := svc.UpdateCustomerUserAccountStatus(context.Background(), acc.ID, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.IsActive != false {
		t.Fatal("expected inactive account")
	}
}
