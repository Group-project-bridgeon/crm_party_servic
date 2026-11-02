package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
	"github.com/group-project/crm_party_service/internal/domain"
	"github.com/group-project/crm_party_service/internal/repository"
)

type CompositeService struct {
	txManager   repository.TransactionManager
	partyRepo   repository.PartyRepository
	orgRepo     repository.OrganizationRepository
	indivRepo   repository.IndividualRepository
	custRepo    repository.CustomerRepository
	suppRepo    repository.SupplierRepository
	addrRepo    repository.AddressRepository
	contactRepo repository.ContactRepository
	noteRepo    repository.NoteRepository
}

func NewCompositeService(
	txManager repository.TransactionManager,
	partyRepo repository.PartyRepository,
	orgRepo repository.OrganizationRepository,
	indivRepo repository.IndividualRepository,
	custRepo repository.CustomerRepository,
	suppRepo repository.SupplierRepository,
	addrRepo repository.AddressRepository,
	contactRepo repository.ContactRepository,
	noteRepo repository.NoteRepository,
) *CompositeService {
	return &CompositeService{
		txManager:   txManager,
		partyRepo:   partyRepo,
		orgRepo:     orgRepo,
		indivRepo:   indivRepo,
		custRepo:    custRepo,
		suppRepo:    suppRepo,
		addrRepo:    addrRepo,
		contactRepo: contactRepo,
		noteRepo:    noteRepo,
	}
}

func (s *CompositeService) GetParty360(
	ctx context.Context,
	partyID pgtype.UUID,
) (*domain.Party360, error) {
	party, err := s.partyRepo.GetPartyByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPartyNotFound
		}
		return nil, err
	}

	result := &domain.Party360{
		Party: party,
	}

	// 1. Fetch Profile (Organization or Individual)
	if party.PartyType == string(domain.PartyTypeOrganization) {
		org, err := s.orgRepo.GetOrganizationByPartyID(ctx, partyID)
		if err == nil {
			result.Organization = &org
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	} else if party.PartyType == string(domain.PartyTypeIndividual) {
		indiv, err := s.indivRepo.GetIndividualByPartyID(ctx, partyID)
		if err == nil {
			result.Individual = &indiv
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	// 2. Fetch Customer role (if exists)
	cust, err := s.custRepo.GetCustomerByPartyID(ctx, partyID)
	if err == nil {
		result.Customer = &cust
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// 3. Fetch Supplier role (if exists)
	supp, err := s.suppRepo.GetSupplierByPartyID(ctx, partyID)
	if err == nil {
		result.Supplier = &supp
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// 4. Fetch Addresses
	addresses, err := s.addrRepo.ListPartyAddressesByPartyID(ctx, db.ListPartyAddressesByPartyIDParams{
		PartyID: partyID,
	})
	if err == nil {
		result.Addresses = addresses
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// 5. Fetch Contacts
	contacts, err := s.contactRepo.ListPartyContactsByPartyID(ctx, db.ListPartyContactsByPartyIDParams{
		PartyID: partyID,
	})
	if err == nil {
		result.Contacts = contacts
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// 6. Fetch Contact Persons
	persons, err := s.contactRepo.ListPartyContactPersonsByPartyID(ctx, partyID)
	if err == nil {
		result.ContactPersons = persons
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// 7. Fetch Notes
	notes, err := s.noteRepo.ListPartyNotesByPartyID(ctx, partyID)
	if err == nil {
		result.Notes = notes
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return result, nil
}

func (s *CompositeService) CreateCompositeCustomer(
	ctx context.Context,
	input domain.CompositeCustomerInput,
) (*domain.CompositeCustomerResult, error) {
	// 1. Validate party type and display name
	if !input.PartyType.IsValid() {
		return nil, ErrInvalidPartyType
	}
	displayName := strings.TrimSpace(input.DisplayName)
	if displayName == "" {
		return nil, ErrEmptyDisplayName
	}

	// 2. Validate customer code
	customerCode := strings.TrimSpace(input.CustomerCode)
	if customerCode == "" {
		return nil, ErrEmptyCustomerCode
	}

	// 3. Check duplicate customer code
	_, err := s.custRepo.GetCustomerByCode(ctx, customerCode)
	if err == nil {
		return nil, ErrDuplicateCustomerCode
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// 4. Validate profile-specific inputs
	var foundedDate pgtype.Date
	var employeeCount *int32
	var dob pgtype.Date
	var firstName string

	if input.PartyType == domain.PartyTypeOrganization {
		if input.EmployeeCount != nil {
			if *input.EmployeeCount < 0 || int64(*input.EmployeeCount) > math.MaxInt32 {
				return nil, ErrInvalidEmployeeCount
			}
			val := int32(*input.EmployeeCount)
			employeeCount = &val
		}
		if input.FoundedDate != nil && *input.FoundedDate != "" {
			t, err := time.Parse("2006-01-02", *input.FoundedDate)
			if err != nil {
				return nil, ErrInvalidFoundedDate
			}
			foundedDate = pgtype.Date{Time: t, Valid: true}
		}
	} else if input.PartyType == domain.PartyTypeIndividual {
		if input.FirstName == nil || strings.TrimSpace(*input.FirstName) == "" {
			return nil, ErrEmptyFirstName
		}
		firstName = strings.TrimSpace(*input.FirstName)
		if input.DateOfBirth != nil && *input.DateOfBirth != "" {
			t, err := time.Parse("2006-01-02", *input.DateOfBirth)
			if err != nil {
				return nil, ErrInvalidDateOfBirth
			}
			dob = pgtype.Date{Time: t, Valid: true}
		}
	}

	// 5. Validate customer finance & status
	var creditLimit pgtype.Numeric
	if input.CreditLimit != nil {
		if *input.CreditLimit < 0 {
			return nil, ErrInvalidCreditLimit
		}
		_ = creditLimit.Scan(fmt.Sprintf("%.2f", *input.CreditLimit))
	} else {
		_ = creditLimit.Scan("0")
	}

	var paymentTerms int32 = 0
	if input.PaymentTermsDays != nil {
		if *input.PaymentTermsDays < 0 {
			return nil, ErrInvalidPaymentTerms
		}
		paymentTerms = int32(*input.PaymentTermsDays)
	}

	creditStatus := string(domain.CreditStatusActive)
	if input.CreditStatus != nil && strings.TrimSpace(*input.CreditStatus) != "" {
		st := strings.ToUpper(strings.TrimSpace(*input.CreditStatus))
		if !domain.CreditStatus(st).IsValid() {
			return nil, ErrInvalidCreditStatus
		}
		creditStatus = st
	}

	var priceListID pgtype.UUID
	if input.PriceListID != nil && *input.PriceListID != "" {
		if err := priceListID.Scan(*input.PriceListID); err != nil {
			return nil, err
		}
	}

	var salesRepID pgtype.UUID
	if input.SalesRepID != nil && *input.SalesRepID != "" {
		if err := salesRepID.Scan(*input.SalesRepID); err != nil {
			return nil, err
		}
	}

	taxExempt := false
	if input.TaxExempt != nil {
		taxExempt = *input.TaxExempt
	}

	// 6. Validate initial address if provided
	if input.Address != nil {
		addrType := strings.ToUpper(strings.TrimSpace(input.Address.AddressType))
		if addrType == "" {
			addrType = string(domain.AddressTypeBilling)
		}
		if !domain.AddressType(addrType).IsValid() {
			return nil, ErrInvalidAddressType
		}
		input.Address.AddressType = addrType

		if strings.TrimSpace(input.Address.AddressLine1) == "" {
			return nil, ErrEmptyAddressLine1
		}
		if strings.TrimSpace(input.Address.City) == "" {
			return nil, ErrEmptyCity
		}
		cc := strings.ToUpper(strings.TrimSpace(input.Address.CountryCode))
		if len(cc) != 2 {
			return nil, ErrInvalidCountryCode
		}
		input.Address.CountryCode = cc
	}

	// 7. Atomic transaction execution
	result := &domain.CompositeCustomerResult{}

	err = s.txManager.ExecTx(ctx, func(q *db.Queries) error {
		// A. Create Party
		p, err := q.CreateParty(ctx, db.CreatePartyParams{
			PartyType:   string(input.PartyType),
			DisplayName: displayName,
			LegalName:   input.LegalName,
			Email:       input.Email,
			Phone:       input.Phone,
			Website:     input.Website,
			TaxNumber:   input.TaxNumber,
			Notes:       input.Notes,
		})
		if err != nil {
			return err
		}
		result.Party = p

		// B. Create Organization or Individual
		if input.PartyType == domain.PartyTypeOrganization {
			org, err := q.CreateOrganization(ctx, db.CreateOrganizationParams{
				PartyID:            p.ID,
				RegistrationNumber: input.RegistrationNumber,
				Industry:           input.Industry,
				OrganizationType:   input.OrganizationType,
				FoundedDate:        foundedDate,
				EmployeeCount:      employeeCount,
				Description:        input.Description,
			})
			if err != nil {
				return err
			}
			result.Organization = &org
		} else if input.PartyType == domain.PartyTypeIndividual {
			indiv, err := q.CreateIndividual(ctx, db.CreateIndividualParams{
				PartyID:     p.ID,
				FirstName:   firstName,
				MiddleName:  input.MiddleName,
				LastName:    input.LastName,
				DateOfBirth: dob,
				Gender:      input.Gender,
				JobTitle:    input.JobTitle,
				NationalID:  input.NationalID,
			})
			if err != nil {
				return err
			}
			result.Individual = &indiv
		}

		// C. Create Customer
		cust, err := q.CreateCustomer(ctx, db.CreateCustomerParams{
			PartyID:            p.ID,
			CustomerCode:       customerCode,
			CustomerGroup:      input.CustomerGroup,
			PriceListID:        priceListID,
			CreditLimit:        creditLimit,
			PaymentTermsDays:   paymentTerms,
			CreditStatus:       creditStatus,
			SalesRepID:         salesRepID,
			TaxExempt:          taxExempt,
			TaxExemptionNumber: input.TaxExemptionNumber,
		})
		if err != nil {
			return err
		}
		result.Customer = cust

		// D. Optional Primary Address
		if input.Address != nil {
			var lat pgtype.Numeric
			if input.Address.Latitude != nil {
				_ = lat.Scan(fmt.Sprintf("%.7f", *input.Address.Latitude))
			}
			var lon pgtype.Numeric
			if input.Address.Longitude != nil {
				_ = lon.Scan(fmt.Sprintf("%.7f", *input.Address.Longitude))
			}

			addr, err := q.CreatePartyAddress(ctx, db.CreatePartyAddressParams{
				PartyID:      p.ID,
				AddressType:  input.Address.AddressType,
				AddressName:  input.Address.AddressName,
				AddressLine1: input.Address.AddressLine1,
				AddressLine2: input.Address.AddressLine2,
				Landmark:     input.Address.Landmark,
				City:         input.Address.City,
				State:        input.Address.State,
				PostalCode:   input.Address.PostalCode,
				CountryCode:  input.Address.CountryCode,
				Latitude:     lat,
				Longitude:    lon,
				IsPrimary:    true,
				IsActive:     true,
			})
			if err != nil {
				return err
			}
			result.Address = &addr
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *CompositeService) CreateCompositeSupplier(
	ctx context.Context,
	input domain.CompositeSupplierInput,
) (*domain.CompositeSupplierResult, error) {
	// 1. Validate party type and display name
	if !input.PartyType.IsValid() {
		return nil, ErrInvalidPartyType
	}
	displayName := strings.TrimSpace(input.DisplayName)
	if displayName == "" {
		return nil, ErrEmptyDisplayName
	}

	// 2. Validate supplier code
	supplierCode := strings.TrimSpace(input.SupplierCode)
	if supplierCode == "" {
		return nil, ErrEmptySupplierCode
	}

	// 3. Check duplicate supplier code
	_, err := s.suppRepo.GetSupplierByCode(ctx, supplierCode)
	if err == nil {
		return nil, ErrDuplicateSupplierCode
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// 4. Validate profile-specific inputs
	var foundedDate pgtype.Date
	var employeeCount *int32
	var dob pgtype.Date
	var firstName string

	if input.PartyType == domain.PartyTypeOrganization {
		if input.EmployeeCount != nil {
			if *input.EmployeeCount < 0 || int64(*input.EmployeeCount) > math.MaxInt32 {
				return nil, ErrInvalidEmployeeCount
			}
			val := int32(*input.EmployeeCount)
			employeeCount = &val
		}
		if input.FoundedDate != nil && *input.FoundedDate != "" {
			t, err := time.Parse("2006-01-02", *input.FoundedDate)
			if err != nil {
				return nil, ErrInvalidFoundedDate
			}
			foundedDate = pgtype.Date{Time: t, Valid: true}
		}
	} else if input.PartyType == domain.PartyTypeIndividual {
		if input.FirstName == nil || strings.TrimSpace(*input.FirstName) == "" {
			return nil, ErrEmptyFirstName
		}
		firstName = strings.TrimSpace(*input.FirstName)
		if input.DateOfBirth != nil && *input.DateOfBirth != "" {
			t, err := time.Parse("2006-01-02", *input.DateOfBirth)
			if err != nil {
				return nil, ErrInvalidDateOfBirth
			}
			dob = pgtype.Date{Time: t, Valid: true}
		}
	}

	// 5. Validate supplier terms, currency, rating
	var paymentTerms int32 = 0
	if input.PaymentTermsDays != nil {
		if *input.PaymentTermsDays < 0 {
			return nil, ErrInvalidPaymentTerms
		}
		paymentTerms = int32(*input.PaymentTermsDays)
	}

	var currencyCode *string
	if input.CurrencyCode != nil && strings.TrimSpace(*input.CurrencyCode) != "" {
		c := strings.ToUpper(strings.TrimSpace(*input.CurrencyCode))
		if len(c) != 3 {
			return nil, ErrInvalidCurrencyCode
		}
		currencyCode = &c
	}

	var supplierRating pgtype.Numeric
	if input.SupplierRating != nil {
		if *input.SupplierRating < 0 || *input.SupplierRating > 5 {
			return nil, ErrInvalidSupplierRating
		}
		_ = supplierRating.Scan(fmt.Sprintf("%.2f", *input.SupplierRating))
	}

	var purchaseRepID pgtype.UUID
	if input.PurchaseRepID != nil && *input.PurchaseRepID != "" {
		if err := purchaseRepID.Scan(*input.PurchaseRepID); err != nil {
			return nil, err
		}
	}

	taxExempt := false
	if input.TaxExempt != nil {
		taxExempt = *input.TaxExempt
	}

	// 6. Validate initial address if provided
	if input.Address != nil {
		addrType := strings.ToUpper(strings.TrimSpace(input.Address.AddressType))
		if addrType == "" {
			addrType = string(domain.AddressTypeOffice)
		}
		if !domain.AddressType(addrType).IsValid() {
			return nil, ErrInvalidAddressType
		}
		input.Address.AddressType = addrType

		if strings.TrimSpace(input.Address.AddressLine1) == "" {
			return nil, ErrEmptyAddressLine1
		}
		if strings.TrimSpace(input.Address.City) == "" {
			return nil, ErrEmptyCity
		}
		cc := strings.ToUpper(strings.TrimSpace(input.Address.CountryCode))
		if len(cc) != 2 {
			return nil, ErrInvalidCountryCode
		}
		input.Address.CountryCode = cc
	}

	// 7. Atomic transaction execution
	result := &domain.CompositeSupplierResult{}

	err = s.txManager.ExecTx(ctx, func(q *db.Queries) error {
		// A. Create Party
		p, err := q.CreateParty(ctx, db.CreatePartyParams{
			PartyType:   string(input.PartyType),
			DisplayName: displayName,
			LegalName:   input.LegalName,
			Email:       input.Email,
			Phone:       input.Phone,
			Website:     input.Website,
			TaxNumber:   input.TaxNumber,
			Notes:       input.Notes,
		})
		if err != nil {
			return err
		}
		result.Party = p

		// B. Create Organization or Individual
		if input.PartyType == domain.PartyTypeOrganization {
			org, err := q.CreateOrganization(ctx, db.CreateOrganizationParams{
				PartyID:            p.ID,
				RegistrationNumber: input.RegistrationNumber,
				Industry:           input.Industry,
				OrganizationType:   input.OrganizationType,
				FoundedDate:        foundedDate,
				EmployeeCount:      employeeCount,
				Description:        input.Description,
			})
			if err != nil {
				return err
			}
			result.Organization = &org
		} else if input.PartyType == domain.PartyTypeIndividual {
			indiv, err := q.CreateIndividual(ctx, db.CreateIndividualParams{
				PartyID:     p.ID,
				FirstName:   firstName,
				MiddleName:  input.MiddleName,
				LastName:    input.LastName,
				DateOfBirth: dob,
				Gender:      input.Gender,
				JobTitle:    input.JobTitle,
				NationalID:  input.NationalID,
			})
			if err != nil {
				return err
			}
			result.Individual = &indiv
		}

		// C. Create Supplier
		supp, err := q.CreateSupplier(ctx, db.CreateSupplierParams{
			PartyID:            p.ID,
			SupplierCode:       supplierCode,
			SupplierGroup:      input.SupplierGroup,
			PaymentTermsDays:   paymentTerms,
			PaymentMethod:      input.PaymentMethod,
			CurrencyCode:       currencyCode,
			PurchaseRepID:      purchaseRepID,
			TaxExempt:          taxExempt,
			TaxExemptionNumber: input.TaxExemptionNumber,
			SupplierRating:     supplierRating,
		})
		if err != nil {
			return err
		}
		result.Supplier = supp

		// D. Optional Primary Address
		if input.Address != nil {
			var lat pgtype.Numeric
			if input.Address.Latitude != nil {
				_ = lat.Scan(fmt.Sprintf("%.7f", *input.Address.Latitude))
			}
			var lon pgtype.Numeric
			if input.Address.Longitude != nil {
				_ = lon.Scan(fmt.Sprintf("%.7f", *input.Address.Longitude))
			}

			addr, err := q.CreatePartyAddress(ctx, db.CreatePartyAddressParams{
				PartyID:      p.ID,
				AddressType:  input.Address.AddressType,
				AddressName:  input.Address.AddressName,
				AddressLine1: input.Address.AddressLine1,
				AddressLine2: input.Address.AddressLine2,
				Landmark:     input.Address.Landmark,
				City:         input.Address.City,
				State:        input.Address.State,
				PostalCode:   input.Address.PostalCode,
				CountryCode:  input.Address.CountryCode,
				Latitude:     lat,
				Longitude:    lon,
				IsPrimary:    true,
				IsActive:     true,
			})
			if err != nil {
				return err
			}
			result.Address = &addr
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
