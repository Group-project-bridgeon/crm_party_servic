package api_test

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
	httpHandler "github.com/group-project/crm_party_service/internal/handler/http"
	"github.com/group-project/crm_party_service/internal/service"
)

type inMemoryStore struct {
	mu             sync.RWMutex
	parties        map[string]db.Party
	organizations  map[string]db.Organization
	individuals    map[string]db.Individual
	customers      map[string]db.Customer
	suppliers      map[string]db.Supplier
	addresses      map[string]db.PartyAddress
	contacts       map[string]db.PartyContact
	contactPersons map[string]db.PartyContactPerson
	notes          map[string]db.PartyNote
	custAccounts   map[string]db.CustomerUserAccount
	suppAccounts   map[string]db.SupplierUserAccount
}

func newInMemoryStore() *inMemoryStore {
	return &inMemoryStore{
		parties:        make(map[string]db.Party),
		organizations:  make(map[string]db.Organization),
		individuals:    make(map[string]db.Individual),
		customers:      make(map[string]db.Customer),
		suppliers:      make(map[string]db.Supplier),
		addresses:      make(map[string]db.PartyAddress),
		contacts:       make(map[string]db.PartyContact),
		contactPersons: make(map[string]db.PartyContactPerson),
		notes:          make(map[string]db.PartyNote),
		custAccounts:   make(map[string]db.CustomerUserAccount),
		suppAccounts:   make(map[string]db.SupplierUserAccount),
	}
}

// Memory Party Repo
type memPartyRepo struct{ store *inMemoryStore }

func (r *memPartyRepo) CreateParty(ctx context.Context, params db.CreatePartyParams) (db.Party, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	var id pgtype.UUID
	_ = id.Scan(uuid.New().String())
	p := db.Party{
		ID:          id,
		PartyType:   params.PartyType,
		DisplayName: params.DisplayName,
		LegalName:   params.LegalName,
		Email:       params.Email,
		Phone:       params.Phone,
		Website:     params.Website,
		TaxNumber:   params.TaxNumber,
		Notes:       params.Notes,
		IsActive:    true,
		CreatedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	r.store.parties[id.String()] = p
	return p, nil
}

func (r *memPartyRepo) GetPartyByID(ctx context.Context, id pgtype.UUID) (db.Party, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	p, ok := r.store.parties[id.String()]
	if !ok {
		return db.Party{}, pgx.ErrNoRows
	}
	return p, nil
}

func (r *memPartyRepo) ListParties(ctx context.Context, params db.ListPartiesParams) ([]db.Party, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	res := make([]db.Party, 0, len(r.store.parties))
	for _, p := range r.store.parties {
		res = append(res, p)
	}
	return res, nil
}

func (r *memPartyRepo) UpdateParty(ctx context.Context, params db.UpdatePartyParams) (db.Party, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	p, ok := r.store.parties[params.ID.String()]
	if !ok {
		return db.Party{}, pgx.ErrNoRows
	}
	p.DisplayName = params.DisplayName
	p.LegalName = params.LegalName
	p.Email = params.Email
	p.Phone = params.Phone
	p.Website = params.Website
	p.TaxNumber = params.TaxNumber
	p.Notes = params.Notes
	p.UpdatedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	r.store.parties[params.ID.String()] = p
	return p, nil
}

func (r *memPartyRepo) UpdatePartyStatus(ctx context.Context, params db.UpdatePartyStatusParams) (db.Party, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	p, ok := r.store.parties[params.ID.String()]
	if !ok {
		return db.Party{}, pgx.ErrNoRows
	}
	p.IsActive = params.IsActive
	p.UpdatedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	r.store.parties[params.ID.String()] = p
	return p, nil
}

// Memory Organization Repo
type memOrgRepo struct{ store *inMemoryStore }

func (r *memOrgRepo) CreateOrganization(ctx context.Context, params db.CreateOrganizationParams) (db.Organization, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	var id pgtype.UUID
	_ = id.Scan(uuid.New().String())
	org := db.Organization{
		ID:                 id,
		PartyID:            params.PartyID,
		RegistrationNumber: params.RegistrationNumber,
		Industry:           params.Industry,
		OrganizationType:   params.OrganizationType,
		FoundedDate:        params.FoundedDate,
		EmployeeCount:      params.EmployeeCount,
		Description:        params.Description,
		CreatedAt:          pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:          pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	r.store.organizations[params.PartyID.String()] = org
	return org, nil
}

func (r *memOrgRepo) GetOrganizationByPartyID(ctx context.Context, partyID pgtype.UUID) (db.Organization, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	org, ok := r.store.organizations[partyID.String()]
	if !ok {
		return db.Organization{}, pgx.ErrNoRows
	}
	return org, nil
}

func (r *memOrgRepo) UpdateOrganization(ctx context.Context, params db.UpdateOrganizationParams) (db.Organization, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	org, ok := r.store.organizations[params.PartyID.String()]
	if !ok {
		return db.Organization{}, pgx.ErrNoRows
	}
	org.RegistrationNumber = params.RegistrationNumber
	org.Industry = params.Industry
	org.OrganizationType = params.OrganizationType
	org.FoundedDate = params.FoundedDate
	org.EmployeeCount = params.EmployeeCount
	org.Description = params.Description
	org.UpdatedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	r.store.organizations[params.PartyID.String()] = org
	return org, nil
}

// Memory Individual Repo
type memIndivRepo struct{ store *inMemoryStore }

func (r *memIndivRepo) CreateIndividual(ctx context.Context, params db.CreateIndividualParams) (db.Individual, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	var id pgtype.UUID
	_ = id.Scan(uuid.New().String())
	indiv := db.Individual{
		ID:          id,
		PartyID:     params.PartyID,
		FirstName:   params.FirstName,
		MiddleName:  params.MiddleName,
		LastName:    params.LastName,
		DateOfBirth: params.DateOfBirth,
		Gender:      params.Gender,
		JobTitle:    params.JobTitle,
		NationalID:  params.NationalID,
		CreatedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	r.store.individuals[params.PartyID.String()] = indiv
	return indiv, nil
}

func (r *memIndivRepo) GetIndividualByPartyID(ctx context.Context, partyID pgtype.UUID) (db.Individual, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	indiv, ok := r.store.individuals[partyID.String()]
	if !ok {
		return db.Individual{}, pgx.ErrNoRows
	}
	return indiv, nil
}

func (r *memIndivRepo) UpdateIndividual(ctx context.Context, params db.UpdateIndividualParams) (db.Individual, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	indiv, ok := r.store.individuals[params.PartyID.String()]
	if !ok {
		return db.Individual{}, pgx.ErrNoRows
	}
	indiv.FirstName = params.FirstName
	indiv.MiddleName = params.MiddleName
	indiv.LastName = params.LastName
	indiv.DateOfBirth = params.DateOfBirth
	indiv.Gender = params.Gender
	indiv.JobTitle = params.JobTitle
	indiv.NationalID = params.NationalID
	indiv.UpdatedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	r.store.individuals[params.PartyID.String()] = indiv
	return indiv, nil
}

// Memory Customer Repo
type memCustRepo struct{ store *inMemoryStore }

func (r *memCustRepo) CreateCustomer(ctx context.Context, params db.CreateCustomerParams) (db.Customer, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	var id pgtype.UUID
	_ = id.Scan(uuid.New().String())
	c := db.Customer{
		ID:                 id,
		PartyID:            params.PartyID,
		CustomerCode:       params.CustomerCode,
		CustomerGroup:      params.CustomerGroup,
		PriceListID:        params.PriceListID,
		CreditLimit:        params.CreditLimit,
		PaymentTermsDays:   params.PaymentTermsDays,
		CreditStatus:       params.CreditStatus,
		SalesRepID:         params.SalesRepID,
		TaxExempt:          params.TaxExempt,
		TaxExemptionNumber: params.TaxExemptionNumber,
		CreatedAt:          pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:          pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	r.store.customers[id.String()] = c
	return c, nil
}

func (r *memCustRepo) GetCustomerByID(ctx context.Context, id pgtype.UUID) (db.Customer, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	c, ok := r.store.customers[id.String()]
	if !ok {
		return db.Customer{}, pgx.ErrNoRows
	}
	return c, nil
}

func (r *memCustRepo) GetCustomerByPartyID(ctx context.Context, partyID pgtype.UUID) (db.Customer, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	for _, c := range r.store.customers {
		if c.PartyID == partyID {
			return c, nil
		}
	}
	return db.Customer{}, pgx.ErrNoRows
}

func (r *memCustRepo) GetCustomerByCode(ctx context.Context, code string) (db.Customer, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	for _, c := range r.store.customers {
		if strings.EqualFold(c.CustomerCode, code) {
			return c, nil
		}
	}
	return db.Customer{}, pgx.ErrNoRows
}

func (r *memCustRepo) ListCustomers(ctx context.Context, params db.ListCustomersParams) ([]db.Customer, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	res := make([]db.Customer, 0, len(r.store.customers))
	for _, c := range r.store.customers {
		res = append(res, c)
	}
	return res, nil
}

func (r *memCustRepo) UpdateCustomer(ctx context.Context, params db.UpdateCustomerParams) (db.Customer, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	c, ok := r.store.customers[params.ID.String()]
	if !ok {
		return db.Customer{}, pgx.ErrNoRows
	}
	c.CustomerGroup = params.CustomerGroup
	c.PriceListID = params.PriceListID
	c.CreditLimit = params.CreditLimit
	c.PaymentTermsDays = params.PaymentTermsDays
	c.CreditStatus = params.CreditStatus
	c.SalesRepID = params.SalesRepID
	c.TaxExempt = params.TaxExempt
	c.TaxExemptionNumber = params.TaxExemptionNumber
	c.UpdatedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	r.store.customers[params.ID.String()] = c
	return c, nil
}

func (r *memCustRepo) CreateCustomerUserAccount(ctx context.Context, params db.CreateCustomerUserAccountParams) (db.CustomerUserAccount, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	var id pgtype.UUID
	_ = id.Scan(uuid.New().String())
	acc := db.CustomerUserAccount{
		ID:         id,
		CustomerID: params.CustomerID,
		UserID:     params.UserID,
		IsActive:   true,
		CreatedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	r.store.custAccounts[id.String()] = acc
	return acc, nil
}

func (r *memCustRepo) ListCustomerUserAccountsByCustomerID(ctx context.Context, customerID pgtype.UUID) ([]db.CustomerUserAccount, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	res := make([]db.CustomerUserAccount, 0)
	for _, a := range r.store.custAccounts {
		if a.CustomerID == customerID {
			res = append(res, a)
		}
	}
	return res, nil
}

func (r *memCustRepo) UpdateCustomerUserAccountStatus(ctx context.Context, params db.UpdateCustomerUserAccountStatusParams) (db.CustomerUserAccount, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	acc, ok := r.store.custAccounts[params.ID.String()]
	if !ok {
		return db.CustomerUserAccount{}, pgx.ErrNoRows
	}
	acc.IsActive = params.IsActive
	acc.UpdatedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	r.store.custAccounts[params.ID.String()] = acc
	return acc, nil
}

// Memory Supplier Repo
type memSuppRepo struct{ store *inMemoryStore }

func (r *memSuppRepo) CreateSupplier(ctx context.Context, params db.CreateSupplierParams) (db.Supplier, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	var id pgtype.UUID
	_ = id.Scan(uuid.New().String())
	s := db.Supplier{
		ID:                 id,
		PartyID:            params.PartyID,
		SupplierCode:       params.SupplierCode,
		SupplierGroup:      params.SupplierGroup,
		PaymentTermsDays:   params.PaymentTermsDays,
		PaymentMethod:      params.PaymentMethod,
		CurrencyCode:       params.CurrencyCode,
		PurchaseRepID:      params.PurchaseRepID,
		TaxExempt:          params.TaxExempt,
		TaxExemptionNumber: params.TaxExemptionNumber,
		SupplierRating:     params.SupplierRating,
		CreatedAt:          pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:          pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	r.store.suppliers[id.String()] = s
	return s, nil
}

func (r *memSuppRepo) GetSupplierByID(ctx context.Context, id pgtype.UUID) (db.Supplier, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	s, ok := r.store.suppliers[id.String()]
	if !ok {
		return db.Supplier{}, pgx.ErrNoRows
	}
	return s, nil
}

func (r *memSuppRepo) GetSupplierByPartyID(ctx context.Context, partyID pgtype.UUID) (db.Supplier, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	for _, s := range r.store.suppliers {
		if s.PartyID == partyID {
			return s, nil
		}
	}
	return db.Supplier{}, pgx.ErrNoRows
}

func (r *memSuppRepo) GetSupplierByCode(ctx context.Context, code string) (db.Supplier, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	for _, s := range r.store.suppliers {
		if strings.EqualFold(s.SupplierCode, code) {
			return s, nil
		}
	}
	return db.Supplier{}, pgx.ErrNoRows
}

func (r *memSuppRepo) ListSuppliers(ctx context.Context, params db.ListSuppliersParams) ([]db.Supplier, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	res := make([]db.Supplier, 0, len(r.store.suppliers))
	for _, s := range r.store.suppliers {
		res = append(res, s)
	}
	return res, nil
}

func (r *memSuppRepo) UpdateSupplier(ctx context.Context, params db.UpdateSupplierParams) (db.Supplier, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	s, ok := r.store.suppliers[params.ID.String()]
	if !ok {
		return db.Supplier{}, pgx.ErrNoRows
	}
	s.SupplierGroup = params.SupplierGroup
	s.PaymentTermsDays = params.PaymentTermsDays
	s.PaymentMethod = params.PaymentMethod
	s.CurrencyCode = params.CurrencyCode
	s.PurchaseRepID = params.PurchaseRepID
	s.TaxExempt = params.TaxExempt
	s.TaxExemptionNumber = params.TaxExemptionNumber
	s.SupplierRating = params.SupplierRating
	s.UpdatedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	r.store.suppliers[params.ID.String()] = s
	return s, nil
}

func (r *memSuppRepo) CreateSupplierUserAccount(ctx context.Context, params db.CreateSupplierUserAccountParams) (db.SupplierUserAccount, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	var id pgtype.UUID
	_ = id.Scan(uuid.New().String())
	acc := db.SupplierUserAccount{
		ID:         id,
		SupplierID: params.SupplierID,
		UserID:     params.UserID,
		IsActive:   true,
		CreatedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	r.store.suppAccounts[id.String()] = acc
	return acc, nil
}

func (r *memSuppRepo) ListSupplierUserAccountsBySupplierID(ctx context.Context, supplierID pgtype.UUID) ([]db.SupplierUserAccount, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	res := make([]db.SupplierUserAccount, 0)
	for _, a := range r.store.suppAccounts {
		if a.SupplierID == supplierID {
			res = append(res, a)
		}
	}
	return res, nil
}

func (r *memSuppRepo) UpdateSupplierUserAccountStatus(ctx context.Context, params db.UpdateSupplierUserAccountStatusParams) (db.SupplierUserAccount, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	acc, ok := r.store.suppAccounts[params.ID.String()]
	if !ok {
		return db.SupplierUserAccount{}, pgx.ErrNoRows
	}
	acc.IsActive = params.IsActive
	acc.UpdatedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	r.store.suppAccounts[params.ID.String()] = acc
	return acc, nil
}

// Memory Address Repo
type memAddrRepo struct{ store *inMemoryStore }

func (r *memAddrRepo) CreatePartyAddress(ctx context.Context, params db.CreatePartyAddressParams) (db.PartyAddress, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	var id pgtype.UUID
	_ = id.Scan(uuid.New().String())
	addr := db.PartyAddress{
		ID:           id,
		PartyID:      params.PartyID,
		AddressType:  params.AddressType,
		AddressName:  params.AddressName,
		AddressLine1: params.AddressLine1,
		AddressLine2: params.AddressLine2,
		Landmark:     params.Landmark,
		City:         params.City,
		State:        params.State,
		PostalCode:   params.PostalCode,
		CountryCode:  params.CountryCode,
		Latitude:     params.Latitude,
		Longitude:    params.Longitude,
		IsPrimary:    params.IsPrimary,
		IsActive:     params.IsActive,
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	r.store.addresses[id.String()] = addr
	return addr, nil
}

func (r *memAddrRepo) GetPartyAddressByID(ctx context.Context, id pgtype.UUID) (db.PartyAddress, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	a, ok := r.store.addresses[id.String()]
	if !ok {
		return db.PartyAddress{}, pgx.ErrNoRows
	}
	return a, nil
}

func (r *memAddrRepo) ListPartyAddressesByPartyID(ctx context.Context, params db.ListPartyAddressesByPartyIDParams) ([]db.PartyAddress, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	res := make([]db.PartyAddress, 0)
	for _, a := range r.store.addresses {
		if a.PartyID == params.PartyID {
			if params.AddressType == nil || a.AddressType == *params.AddressType {
				res = append(res, a)
			}
		}
	}
	return res, nil
}

func (r *memAddrRepo) UpdatePartyAddress(ctx context.Context, params db.UpdatePartyAddressParams) (db.PartyAddress, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	a, ok := r.store.addresses[params.ID.String()]
	if !ok {
		return db.PartyAddress{}, pgx.ErrNoRows
	}
	a.AddressType = params.AddressType
	a.AddressName = params.AddressName
	a.AddressLine1 = params.AddressLine1
	a.AddressLine2 = params.AddressLine2
	a.Landmark = params.Landmark
	a.City = params.City
	a.State = params.State
	a.PostalCode = params.PostalCode
	a.CountryCode = params.CountryCode
	a.Latitude = params.Latitude
	a.Longitude = params.Longitude
	a.IsPrimary = params.IsPrimary
	a.IsActive = params.IsActive
	a.UpdatedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	r.store.addresses[params.ID.String()] = a
	return a, nil
}

func (r *memAddrRepo) ResetPrimaryPartyAddress(ctx context.Context, params db.ResetPrimaryPartyAddressParams) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	for k, a := range r.store.addresses {
		if a.PartyID == params.PartyID && a.AddressType == params.AddressType {
			a.IsPrimary = false
			r.store.addresses[k] = a
		}
	}
	return nil
}

func (r *memAddrRepo) DeletePartyAddress(ctx context.Context, id pgtype.UUID) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	delete(r.store.addresses, id.String())
	return nil
}

// Memory Contact Repo
type memContactRepo struct{ store *inMemoryStore }

func (r *memContactRepo) CreatePartyContact(ctx context.Context, params db.CreatePartyContactParams) (db.PartyContact, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	var id pgtype.UUID
	_ = id.Scan(uuid.New().String())
	c := db.PartyContact{
		ID:           id,
		PartyID:      params.PartyID,
		ContactType:  params.ContactType,
		ContactValue: params.ContactValue,
		Label:        params.Label,
		IsPrimary:    params.IsPrimary,
		IsActive:     params.IsActive,
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	r.store.contacts[id.String()] = c
	return c, nil
}

func (r *memContactRepo) GetPartyContactByID(ctx context.Context, id pgtype.UUID) (db.PartyContact, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	c, ok := r.store.contacts[id.String()]
	if !ok {
		return db.PartyContact{}, pgx.ErrNoRows
	}
	return c, nil
}

func (r *memContactRepo) ListPartyContactsByPartyID(ctx context.Context, params db.ListPartyContactsByPartyIDParams) ([]db.PartyContact, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	res := make([]db.PartyContact, 0)
	for _, c := range r.store.contacts {
		if c.PartyID == params.PartyID {
			if params.ContactType == nil || c.ContactType == *params.ContactType {
				res = append(res, c)
			}
		}
	}
	return res, nil
}

func (r *memContactRepo) UpdatePartyContact(ctx context.Context, params db.UpdatePartyContactParams) (db.PartyContact, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	c, ok := r.store.contacts[params.ID.String()]
	if !ok {
		return db.PartyContact{}, pgx.ErrNoRows
	}
	c.ContactType = params.ContactType
	c.ContactValue = params.ContactValue
	c.Label = params.Label
	c.IsPrimary = params.IsPrimary
	c.IsActive = params.IsActive
	c.UpdatedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	r.store.contacts[params.ID.String()] = c
	return c, nil
}

func (r *memContactRepo) ResetPrimaryPartyContact(ctx context.Context, params db.ResetPrimaryPartyContactParams) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	for k, c := range r.store.contacts {
		if c.PartyID == params.PartyID && c.ContactType == params.ContactType {
			c.IsPrimary = false
			r.store.contacts[k] = c
		}
	}
	return nil
}

func (r *memContactRepo) DeletePartyContact(ctx context.Context, id pgtype.UUID) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	delete(r.store.contacts, id.String())
	return nil
}

func (r *memContactRepo) CreatePartyContactPerson(ctx context.Context, params db.CreatePartyContactPersonParams) (db.PartyContactPerson, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	var id pgtype.UUID
	_ = id.Scan(uuid.New().String())
	p := db.PartyContactPerson{
		ID:          id,
		PartyID:     params.PartyID,
		FirstName:   params.FirstName,
		MiddleName:  params.MiddleName,
		LastName:    params.LastName,
		Designation: params.Designation,
		Department:  params.Department,
		Email:       params.Email,
		Phone:       params.Phone,
		Mobile:      params.Mobile,
		IsPrimary:   params.IsPrimary,
		IsActive:    params.IsActive,
		Notes:       params.Notes,
		CreatedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	r.store.contactPersons[id.String()] = p
	return p, nil
}

func (r *memContactRepo) GetPartyContactPersonByID(ctx context.Context, id pgtype.UUID) (db.PartyContactPerson, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	p, ok := r.store.contactPersons[id.String()]
	if !ok {
		return db.PartyContactPerson{}, pgx.ErrNoRows
	}
	return p, nil
}

func (r *memContactRepo) ListPartyContactPersonsByPartyID(ctx context.Context, partyID pgtype.UUID) ([]db.PartyContactPerson, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	res := make([]db.PartyContactPerson, 0)
	for _, p := range r.store.contactPersons {
		if p.PartyID == partyID {
			res = append(res, p)
		}
	}
	return res, nil
}

func (r *memContactRepo) UpdatePartyContactPerson(ctx context.Context, params db.UpdatePartyContactPersonParams) (db.PartyContactPerson, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	p, ok := r.store.contactPersons[params.ID.String()]
	if !ok {
		return db.PartyContactPerson{}, pgx.ErrNoRows
	}
	p.FirstName = params.FirstName
	p.MiddleName = params.MiddleName
	p.LastName = params.LastName
	p.Designation = params.Designation
	p.Department = params.Department
	p.Email = params.Email
	p.Phone = params.Phone
	p.Mobile = params.Mobile
	p.IsPrimary = params.IsPrimary
	p.IsActive = params.IsActive
	p.Notes = params.Notes
	p.UpdatedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	r.store.contactPersons[params.ID.String()] = p
	return p, nil
}

func (r *memContactRepo) ResetPrimaryPartyContactPerson(ctx context.Context, partyID pgtype.UUID) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	for k, p := range r.store.contactPersons {
		if p.PartyID == partyID {
			p.IsPrimary = false
			r.store.contactPersons[k] = p
		}
	}
	return nil
}

func (r *memContactRepo) DeletePartyContactPerson(ctx context.Context, id pgtype.UUID) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	delete(r.store.contactPersons, id.String())
	return nil
}

// Memory Note Repo
type memNoteRepo struct{ store *inMemoryStore }

func (r *memNoteRepo) CreatePartyNote(ctx context.Context, params db.CreatePartyNoteParams) (db.PartyNote, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	var id pgtype.UUID
	_ = id.Scan(uuid.New().String())
	n := db.PartyNote{
		ID:        id,
		PartyID:   params.PartyID,
		NoteType:  params.NoteType,
		Content:   params.Content,
		CreatedBy: params.CreatedBy,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	r.store.notes[id.String()] = n
	return n, nil
}

func (r *memNoteRepo) GetPartyNoteByID(ctx context.Context, id pgtype.UUID) (db.PartyNote, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	n, ok := r.store.notes[id.String()]
	if !ok {
		return db.PartyNote{}, pgx.ErrNoRows
	}
	return n, nil
}

func (r *memNoteRepo) ListPartyNotesByPartyID(ctx context.Context, partyID pgtype.UUID) ([]db.PartyNote, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	res := make([]db.PartyNote, 0)
	for _, n := range r.store.notes {
		if n.PartyID == partyID {
			res = append(res, n)
		}
	}
	return res, nil
}

func (r *memNoteRepo) DeletePartyNote(ctx context.Context, id pgtype.UUID) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	delete(r.store.notes, id.String())
	return nil
}

type mockRow struct {
	sql  string
	args []any
}

func (m *mockRow) Scan(dest ...any) error {
	for _, d := range dest {
		switch v := d.(type) {
		case *pgtype.UUID:
			*v = pgtype.UUID{Bytes: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, Valid: true}
		case *pgtype.Timestamptz:
			*v = pgtype.Timestamptz{Time: time.Now(), Valid: true}
		case *pgtype.Date:
			*v = pgtype.Date{Time: time.Now(), Valid: true}
		case *pgtype.Numeric:
			_ = v.Scan("0")
		case *bool:
			*v = true
		}
	}

	if strings.Contains(m.sql, "INSERT INTO parties") && len(m.args) >= 2 {
		if pt, ok := dest[1].(*string); ok {
			if s, ok := m.args[0].(string); ok {
				*pt = s
			}
		}
		if dn, ok := dest[2].(*string); ok {
			if s, ok := m.args[1].(string); ok {
				*dn = s
			}
		}
	} else if strings.Contains(m.sql, "INSERT INTO customers") && len(m.args) >= 2 {
		if cc, ok := dest[2].(*string); ok {
			if s, ok := m.args[1].(string); ok {
				*cc = s
			}
		}
	} else if strings.Contains(m.sql, "INSERT INTO suppliers") && len(m.args) >= 2 {
		if sc, ok := dest[2].(*string); ok {
			if s, ok := m.args[1].(string); ok {
				*sc = s
			}
		}
	}
	return nil
}

type mockDBTX struct{}

func (m *mockDBTX) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (m *mockDBTX) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}

func (m *mockDBTX) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return &mockRow{sql: sql, args: args}
}

// Memory Transaction Manager
type memTxManager struct{}

func (m *memTxManager) ExecTx(ctx context.Context, fn func(queries *db.Queries) error) error {
	q := db.New(&mockDBTX{})
	return fn(q)
}

func setupTestApp() (*fiber.App, *inMemoryStore) {
	store := newInMemoryStore()

	partyRepo := &memPartyRepo{store: store}
	partyService := service.NewPartyService(partyRepo)
	partyHandler := httpHandler.NewPartyHandler(partyService)

	orgRepo := &memOrgRepo{store: store}
	orgService := service.NewOrganizationService(partyRepo, orgRepo)
	organizationHandler := httpHandler.NewOrganizationHandler(orgService)

	indivRepo := &memIndivRepo{store: store}
	indivService := service.NewIndividualService(partyRepo, indivRepo)
	individualHandler := httpHandler.NewIndividualHandler(indivService)

	custRepo := &memCustRepo{store: store}
	custService := service.NewCustomerService(partyRepo, custRepo)
	customerHandler := httpHandler.NewCustomerHandler(custService)

	suppRepo := &memSuppRepo{store: store}
	suppService := service.NewSupplierService(partyRepo, suppRepo)
	supplierHandler := httpHandler.NewSupplierHandler(suppService)

	addrRepo := &memAddrRepo{store: store}
	addrService := service.NewAddressService(partyRepo, addrRepo)
	addressHandler := httpHandler.NewAddressHandler(addrService)

	contactRepo := &memContactRepo{store: store}
	contactService := service.NewContactService(partyRepo, contactRepo)
	contactHandler := httpHandler.NewContactHandler(contactService)

	noteRepo := &memNoteRepo{store: store}
	noteService := service.NewNoteService(partyRepo, noteRepo)
	noteHandler := httpHandler.NewNoteHandler(noteService)

	txMgr := &memTxManager{}
	compositeService := service.NewCompositeService(
		txMgr,
		partyRepo,
		orgRepo,
		indivRepo,
		custRepo,
		suppRepo,
		addrRepo,
		contactRepo,
		noteRepo,
	)
	compositeHandler := httpHandler.NewCompositeHandler(compositeService)

	app := fiber.New()

	httpHandler.RegisterRoutes(app, httpHandler.Handlers{
		Party:        partyHandler,
		Organization: organizationHandler,
		Individual:   individualHandler,
		Customer:     customerHandler,
		Supplier:     supplierHandler,
		Address:      addressHandler,
		Contact:      contactHandler,
		Note:         noteHandler,
		Composite:    compositeHandler,
	})

	return app, store
}
