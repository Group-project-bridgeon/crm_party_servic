package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
	"github.com/group-project/crm_party_service/internal/domain"
	"github.com/group-project/crm_party_service/internal/repository"
)

type CustomerService struct {
	partyRepository    repository.PartyRepository
	customerRepository repository.CustomerRepository
}

func NewCustomerService(
	partyRepository repository.PartyRepository,
	customerRepository repository.CustomerRepository,
) *CustomerService {
	return &CustomerService{
		partyRepository:    partyRepository,
		customerRepository: customerRepository,
	}
}

func (s *CustomerService) CreateCustomer(
	ctx context.Context,
	input domain.CreateCustomerInput,
) (db.Customer, error) {
	// 1. Parse party ID.
	var partyID pgtype.UUID
	if err := partyID.Scan(input.PartyID); err != nil {
		return db.Customer{}, ErrInvalidPartyID
	}

	// 2. Validate customer code.
	input.CustomerCode = strings.TrimSpace(input.CustomerCode)
	if input.CustomerCode == "" {
		return db.Customer{}, ErrEmptyCustomerCode
	}

	// 3. Ensure party exists.
	_, err := s.partyRepository.GetPartyByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Customer{}, ErrPartyNotFound
		}
		return db.Customer{}, err
	}

	// 4. Check if party is already a customer.
	_, err = s.customerRepository.GetCustomerByPartyID(ctx, partyID)
	if err == nil {
		return db.Customer{}, ErrCustomerAlreadyExists
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return db.Customer{}, err
	}

	// 5. Check if customer code is already used.
	_, err = s.customerRepository.GetCustomerByCode(ctx, input.CustomerCode)
	if err == nil {
		return db.Customer{}, ErrDuplicateCustomerCode
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return db.Customer{}, err
	}

	// 6. Credit limit validation & conversion.
	var creditLimit pgtype.Numeric
	if input.CreditLimit != nil {
		if *input.CreditLimit < 0 {
			return db.Customer{}, ErrInvalidCreditLimit
		}
		if err := creditLimit.Scan(fmt.Sprintf("%.2f", *input.CreditLimit)); err != nil {
			return db.Customer{}, ErrInvalidCreditLimit
		}
	} else {
		_ = creditLimit.Scan("0")
	}

	// 7. Payment terms validation.
	var paymentTerms int32 = 0
	if input.PaymentTermsDays != nil {
		if *input.PaymentTermsDays < 0 {
			return db.Customer{}, ErrInvalidPaymentTerms
		}
		paymentTerms = int32(*input.PaymentTermsDays)
	}

	// 8. Credit status validation.
	creditStatus := string(domain.CreditStatusActive)
	if input.CreditStatus != nil && strings.TrimSpace(*input.CreditStatus) != "" {
		status := strings.ToUpper(strings.TrimSpace(*input.CreditStatus))
		if !domain.CreditStatus(status).IsValid() {
			return db.Customer{}, ErrInvalidCreditStatus
		}
		creditStatus = status
	}

	// 9. Optional UUIDs: PriceListID and SalesRepID.
	var priceListID pgtype.UUID
	if input.PriceListID != nil && *input.PriceListID != "" {
		if err := priceListID.Scan(*input.PriceListID); err != nil {
			return db.Customer{}, err
		}
	}

	var salesRepID pgtype.UUID
	if input.SalesRepID != nil && *input.SalesRepID != "" {
		if err := salesRepID.Scan(*input.SalesRepID); err != nil {
			return db.Customer{}, err
		}
	}

	taxExempt := false
	if input.TaxExempt != nil {
		taxExempt = *input.TaxExempt
	}

	params := db.CreateCustomerParams{
		PartyID:            partyID,
		CustomerCode:       input.CustomerCode,
		CustomerGroup:      input.CustomerGroup,
		PriceListID:        priceListID,
		CreditLimit:        creditLimit,
		PaymentTermsDays:   paymentTerms,
		CreditStatus:       creditStatus,
		SalesRepID:         salesRepID,
		TaxExempt:          taxExempt,
		TaxExemptionNumber: input.TaxExemptionNumber,
	}

	return s.customerRepository.CreateCustomer(ctx, params)
}

func (s *CustomerService) GetCustomerByID(
	ctx context.Context,
	id pgtype.UUID,
) (db.Customer, error) {
	customer, err := s.customerRepository.GetCustomerByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Customer{}, ErrCustomerNotFound
		}
		return db.Customer{}, err
	}
	return customer, nil
}

func (s *CustomerService) GetCustomerByPartyID(
	ctx context.Context,
	partyID pgtype.UUID,
) (db.Customer, error) {
	customer, err := s.customerRepository.GetCustomerByPartyID(ctx, partyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Customer{}, ErrCustomerNotFound
		}
		return db.Customer{}, err
	}
	return customer, nil
}

func (s *CustomerService) ListCustomers(
	ctx context.Context,
	params db.ListCustomersParams,
) ([]db.Customer, error) {
	if params.Offset < 0 {
		return nil, ErrInvalidPagination
	}
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	if params.Search != nil {
		search := strings.TrimSpace(*params.Search)
		if search == "" {
			params.Search = nil
		} else {
			params.Search = &search
		}
	}

	return s.customerRepository.ListCustomers(ctx, params)
}

func (s *CustomerService) UpdateCustomer(
	ctx context.Context,
	input domain.UpdateCustomerInput,
) (db.Customer, error) {
	var customerID pgtype.UUID
	if err := customerID.Scan(input.ID); err != nil {
		return db.Customer{}, ErrInvalidCustomerID
	}

	existing, err := s.customerRepository.GetCustomerByID(ctx, customerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Customer{}, ErrCustomerNotFound
		}
		return db.Customer{}, err
	}

	creditLimit := existing.CreditLimit
	if input.CreditLimit != nil {
		if *input.CreditLimit < 0 {
			return db.Customer{}, ErrInvalidCreditLimit
		}
		if err := creditLimit.Scan(fmt.Sprintf("%.2f", *input.CreditLimit)); err != nil {
			return db.Customer{}, ErrInvalidCreditLimit
		}
	}

	paymentTerms := existing.PaymentTermsDays
	if input.PaymentTermsDays != nil {
		if *input.PaymentTermsDays < 0 {
			return db.Customer{}, ErrInvalidPaymentTerms
		}
		paymentTerms = int32(*input.PaymentTermsDays)
	}

	creditStatus := existing.CreditStatus
	if input.CreditStatus != nil && strings.TrimSpace(*input.CreditStatus) != "" {
		status := strings.ToUpper(strings.TrimSpace(*input.CreditStatus))
		if !domain.CreditStatus(status).IsValid() {
			return db.Customer{}, ErrInvalidCreditStatus
		}
		creditStatus = status
	}

	priceListID := existing.PriceListID
	if input.PriceListID != nil {
		if *input.PriceListID == "" {
			priceListID = pgtype.UUID{}
		} else {
			if err := priceListID.Scan(*input.PriceListID); err != nil {
				return db.Customer{}, err
			}
		}
	}

	salesRepID := existing.SalesRepID
	if input.SalesRepID != nil {
		if *input.SalesRepID == "" {
			salesRepID = pgtype.UUID{}
		} else {
			if err := salesRepID.Scan(*input.SalesRepID); err != nil {
				return db.Customer{}, err
			}
		}
	}

	taxExempt := existing.TaxExempt
	if input.TaxExempt != nil {
		taxExempt = *input.TaxExempt
	}

	customerGroup := existing.CustomerGroup
	if input.CustomerGroup != nil {
		customerGroup = input.CustomerGroup
	}

	taxExemptionNumber := existing.TaxExemptionNumber
	if input.TaxExemptionNumber != nil {
		taxExemptionNumber = input.TaxExemptionNumber
	}

	params := db.UpdateCustomerParams{
		ID:                 customerID,
		CustomerGroup:      customerGroup,
		PriceListID:        priceListID,
		CreditLimit:        creditLimit,
		PaymentTermsDays:   paymentTerms,
		CreditStatus:       creditStatus,
		SalesRepID:         salesRepID,
		TaxExempt:          taxExempt,
		TaxExemptionNumber: taxExemptionNumber,
	}

	return s.customerRepository.UpdateCustomer(ctx, params)
}

func (s *CustomerService) CreateCustomerUserAccount(
	ctx context.Context,
	input domain.CreateCustomerUserAccountInput,
) (db.CustomerUserAccount, error) {
	var customerID pgtype.UUID
	if err := customerID.Scan(input.CustomerID); err != nil {
		return db.CustomerUserAccount{}, ErrInvalidCustomerID
	}

	var userID pgtype.UUID
	if err := userID.Scan(input.UserID); err != nil {
		return db.CustomerUserAccount{}, ErrInvalidUserID
	}

	// Ensure customer exists
	_, err := s.customerRepository.GetCustomerByID(ctx, customerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.CustomerUserAccount{}, ErrCustomerNotFound
		}
		return db.CustomerUserAccount{}, err
	}

	params := db.CreateCustomerUserAccountParams{
		CustomerID: customerID,
		UserID:     userID,
	}

	return s.customerRepository.CreateCustomerUserAccount(ctx, params)
}

func (s *CustomerService) ListCustomerUserAccounts(
	ctx context.Context,
	customerID pgtype.UUID,
) ([]db.CustomerUserAccount, error) {
	_, err := s.customerRepository.GetCustomerByID(ctx, customerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCustomerNotFound
		}
		return nil, err
	}

	return s.customerRepository.ListCustomerUserAccountsByCustomerID(ctx, customerID)
}

func (s *CustomerService) UpdateCustomerUserAccountStatus(
	ctx context.Context,
	id pgtype.UUID,
	isActive bool,
) (db.CustomerUserAccount, error) {
	params := db.UpdateCustomerUserAccountStatusParams{
		ID:       id,
		IsActive: isActive,
	}
	account, err := s.customerRepository.UpdateCustomerUserAccountStatus(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.CustomerUserAccount{}, ErrCustomerUserAccountNotFound
		}
		return db.CustomerUserAccount{}, err
	}
	return account, nil
}
