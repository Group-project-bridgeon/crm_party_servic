package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
	"github.com/group-project/crm_party_service/internal/domain"
)

type mockRow struct {
	scanFunc func(dest ...any) error
}

func (m *mockRow) Scan(dest ...any) error {
	if m.scanFunc != nil {
		return m.scanFunc(dest...)
	}
	for _, d := range dest {
		if uuidPtr, ok := d.(*pgtype.UUID); ok {
			*uuidPtr = pgtype.UUID{Valid: true}
		}
	}
	return nil
}

type mockDBTX struct {
	queryRowFunc func(ctx context.Context, sql string, args ...any) pgx.Row
}

func (m *mockDBTX) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (m *mockDBTX) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}

func (m *mockDBTX) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if m.queryRowFunc != nil {
		return m.queryRowFunc(ctx, sql, args...)
	}
	return &mockRow{}
}

type mockTransactionManager struct {
	execTxFunc func(ctx context.Context, fn func(queries *db.Queries) error) error
}

func (m *mockTransactionManager) ExecTx(ctx context.Context, fn func(queries *db.Queries) error) error {
	if m.execTxFunc != nil {
		return m.execTxFunc(ctx, fn)
	}
	q := db.New(&mockDBTX{})
	return fn(q)
}

func TestGetParty360_Success(t *testing.T) {
	var partyID pgtype.UUID
	partyID.Scan("550e8400-e29b-41d4-a716-446655440000")

	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{
				ID:          id,
				PartyType:   string(domain.PartyTypeOrganization),
				DisplayName: "Acme Corp",
			}, nil
		},
	}
	orgRepo := &mockOrganizationRepository{
		getOrganizationByPartyIDFunc: func(ctx context.Context, pID pgtype.UUID) (db.Organization, error) {
			return db.Organization{PartyID: pID}, nil
		},
	}
	indivRepo := &mockIndividualRepository{}
	custRepo := &mockCustomerRepository{
		getCustomerByPartyIDFunc: func(ctx context.Context, pID pgtype.UUID) (db.Customer, error) {
			return db.Customer{PartyID: pID, CustomerCode: "CUST-100"}, nil
		},
	}
	suppRepo := &mockSupplierRepository{}
	addrRepo := &mockAddressRepository{
		listPartyAddressesByPartyIDFunc: func(ctx context.Context, params db.ListPartyAddressesByPartyIDParams) ([]db.PartyAddress, error) {
			return []db.PartyAddress{{PartyID: params.PartyID, City: "New York"}}, nil
		},
	}
	contactRepo := &mockContactRepository{
		listPartyContactsByPartyIDFunc: func(ctx context.Context, params db.ListPartyContactsByPartyIDParams) ([]db.PartyContact, error) {
			return []db.PartyContact{{PartyID: params.PartyID, ContactValue: "acme@example.com"}}, nil
		},
	}
	noteRepo := &mockNoteRepository{
		listPartyNotesByPartyIDFunc: func(ctx context.Context, pID pgtype.UUID) ([]db.PartyNote, error) {
			return []db.PartyNote{{PartyID: pID, Content: "VIP Customer"}}, nil
		},
	}

	svc := NewCompositeService(
		&mockTransactionManager{},
		partyRepo,
		orgRepo,
		indivRepo,
		custRepo,
		suppRepo,
		addrRepo,
		contactRepo,
		noteRepo,
	)

	p360, err := svc.GetParty360(context.Background(), partyID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p360.Party.DisplayName != "Acme Corp" {
		t.Fatalf("expected Acme Corp, got %s", p360.Party.DisplayName)
	}
	if p360.Organization == nil {
		t.Fatal("expected organization to be present")
	}
	if p360.Customer == nil || p360.Customer.CustomerCode != "CUST-100" {
		t.Fatal("expected customer role to be present")
	}
	if len(p360.Addresses) != 1 {
		t.Fatalf("expected 1 address, got %d", len(p360.Addresses))
	}
	if len(p360.Contacts) != 1 {
		t.Fatalf("expected 1 contact, got %d", len(p360.Contacts))
	}
	if len(p360.Notes) != 1 {
		t.Fatalf("expected 1 note, got %d", len(p360.Notes))
	}
}

func TestGetParty360_NotFound(t *testing.T) {
	var partyID pgtype.UUID
	partyID.Scan("550e8400-e29b-41d4-a716-446655440000")

	partyRepo := &mockPartyRepository{
		getPartyByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Party, error) {
			return db.Party{}, pgx.ErrNoRows
		},
	}

	svc := NewCompositeService(
		&mockTransactionManager{},
		partyRepo,
		&mockOrganizationRepository{},
		&mockIndividualRepository{},
		&mockCustomerRepository{},
		&mockSupplierRepository{},
		&mockAddressRepository{},
		&mockContactRepository{},
		&mockNoteRepository{},
	)

	_, err := svc.GetParty360(context.Background(), partyID)
	if !errors.Is(err, ErrPartyNotFound) {
		t.Fatalf("expected ErrPartyNotFound, got %v", err)
	}
}

func TestCreateCompositeCustomer_Success(t *testing.T) {
	txMgr := &mockTransactionManager{}
	custRepo := &mockCustomerRepository{
		getCustomerByCodeFunc: func(ctx context.Context, code string) (db.Customer, error) {
			return db.Customer{}, pgx.ErrNoRows
		},
	}

	svc := NewCompositeService(
		txMgr,
		&mockPartyRepository{},
		&mockOrganizationRepository{},
		&mockIndividualRepository{},
		custRepo,
		&mockSupplierRepository{},
		&mockAddressRepository{},
		&mockContactRepository{},
		&mockNoteRepository{},
	)

	input := domain.CompositeCustomerInput{
		PartyType:    domain.PartyTypeOrganization,
		DisplayName:  "MegaCorp Inc",
		CustomerCode: "CUST-MEGA",
		Address: &domain.CreateAddressInput{
			AddressType:  "BILLING",
			AddressLine1: "100 Broadway",
			City:         "New York",
			CountryCode:  "US",
		},
	}

	result, err := svc.CreateCompositeCustomer(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestCreateCompositeCustomer_ValidationErrors(t *testing.T) {
	svc := NewCompositeService(
		&mockTransactionManager{},
		&mockPartyRepository{},
		&mockOrganizationRepository{},
		&mockIndividualRepository{},
		&mockCustomerRepository{},
		&mockSupplierRepository{},
		&mockAddressRepository{},
		&mockContactRepository{},
		&mockNoteRepository{},
	)

	// Missing display name
	_, err := svc.CreateCompositeCustomer(context.Background(), domain.CompositeCustomerInput{
		PartyType:    domain.PartyTypeOrganization,
		DisplayName:  "   ",
		CustomerCode: "CUST-001",
	})
	if !errors.Is(err, ErrEmptyDisplayName) {
		t.Fatalf("expected ErrEmptyDisplayName, got %v", err)
	}

	// Missing customer code
	_, err = svc.CreateCompositeCustomer(context.Background(), domain.CompositeCustomerInput{
		PartyType:    domain.PartyTypeOrganization,
		DisplayName:  "Name",
		CustomerCode: "   ",
	})
	if !errors.Is(err, ErrEmptyCustomerCode) {
		t.Fatalf("expected ErrEmptyCustomerCode, got %v", err)
	}

	// Individual customer missing first name
	_, err = svc.CreateCompositeCustomer(context.Background(), domain.CompositeCustomerInput{
		PartyType:    domain.PartyTypeIndividual,
		DisplayName:  "John Doe",
		CustomerCode: "CUST-001",
	})
	if !errors.Is(err, ErrEmptyFirstName) {
		t.Fatalf("expected ErrEmptyFirstName, got %v", err)
	}
}

func TestCreateCompositeSupplier_Success(t *testing.T) {
	txMgr := &mockTransactionManager{}
	suppRepo := &mockSupplierRepository{
		getSupplierByCodeFunc: func(ctx context.Context, code string) (db.Supplier, error) {
			return db.Supplier{}, pgx.ErrNoRows
		},
	}

	svc := NewCompositeService(
		txMgr,
		&mockPartyRepository{},
		&mockOrganizationRepository{},
		&mockIndividualRepository{},
		&mockCustomerRepository{},
		suppRepo,
		&mockAddressRepository{},
		&mockContactRepository{},
		&mockNoteRepository{},
	)

	rating := 4.8
	curr := "USD"
	input := domain.CompositeSupplierInput{
		PartyType:      domain.PartyTypeOrganization,
		DisplayName:    "Global Supplier",
		SupplierCode:   "SUPP-GLOBAL",
		SupplierRating: &rating,
		CurrencyCode:   &curr,
	}

	result, err := svc.CreateCompositeSupplier(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}
