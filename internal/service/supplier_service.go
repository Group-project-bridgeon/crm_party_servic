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

type SupplierService struct {
	partyRepository    repository.PartyRepository
	supplierRepository repository.SupplierRepository
}

func NewSupplierService(
	partyRepository repository.PartyRepository,
	supplierRepository repository.SupplierRepository,
) *SupplierService {
	return &SupplierService{
		partyRepository:    partyRepository,
		supplierRepository: supplierRepository,
	}
}

func (s *SupplierService) CreateSupplier(
	ctx context.Context,
	input domain.CreateSupplierInput,
) (db.Supplier, error) {
	// 1. Parse party ID.
	var partyID pgtype.UUID
	if err := partyID.Scan(input.PartyID); err != nil {
		return db.Supplier{}, ErrInvalidPartyID
	}

	// 2. Validate supplier code.
	input.SupplierCode = strings.TrimSpace(input.SupplierCode)
	if input.SupplierCode == "" {
		return db.Supplier{}, ErrEmptySupplierCode
	}

	// 3. Ensure party exists.
	_, err := s.partyRepository.GetPartyByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Supplier{}, ErrPartyNotFound
		}
		return db.Supplier{}, err
	}

	// 4. Check if party is already a supplier.
	_, err = s.supplierRepository.GetSupplierByPartyID(ctx, partyID)
	if err == nil {
		return db.Supplier{}, ErrSupplierAlreadyExists
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return db.Supplier{}, err
	}

	// 5. Check if supplier code is already in use.
	_, err = s.supplierRepository.GetSupplierByCode(ctx, input.SupplierCode)
	if err == nil {
		return db.Supplier{}, ErrDuplicateSupplierCode
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return db.Supplier{}, err
	}

	// 6. Payment terms validation.
	var paymentTerms int32 = 0
	if input.PaymentTermsDays != nil {
		if *input.PaymentTermsDays < 0 {
			return db.Supplier{}, ErrInvalidPaymentTerms
		}
		paymentTerms = int32(*input.PaymentTermsDays)
	}

	// 7. Currency code validation.
	var currencyCode *string
	if input.CurrencyCode != nil && strings.TrimSpace(*input.CurrencyCode) != "" {
		c := strings.ToUpper(strings.TrimSpace(*input.CurrencyCode))
		if len(c) != 3 {
			return db.Supplier{}, ErrInvalidCurrencyCode
		}
		currencyCode = &c
	}

	// 8. Supplier rating validation & conversion.
	var supplierRating pgtype.Numeric
	if input.SupplierRating != nil {
		if *input.SupplierRating < 0 || *input.SupplierRating > 5 {
			return db.Supplier{}, ErrInvalidSupplierRating
		}
		if err := supplierRating.Scan(fmt.Sprintf("%.2f", *input.SupplierRating)); err != nil {
			return db.Supplier{}, ErrInvalidSupplierRating
		}
	}

	var purchaseRepID pgtype.UUID
	if input.PurchaseRepID != nil && *input.PurchaseRepID != "" {
		if err := purchaseRepID.Scan(*input.PurchaseRepID); err != nil {
			return db.Supplier{}, err
		}
	}

	taxExempt := false
	if input.TaxExempt != nil {
		taxExempt = *input.TaxExempt
	}

	params := db.CreateSupplierParams{
		PartyID:            partyID,
		SupplierCode:       input.SupplierCode,
		SupplierGroup:      input.SupplierGroup,
		PaymentTermsDays:   paymentTerms,
		PaymentMethod:      input.PaymentMethod,
		CurrencyCode:       currencyCode,
		PurchaseRepID:      purchaseRepID,
		TaxExempt:          taxExempt,
		TaxExemptionNumber: input.TaxExemptionNumber,
		SupplierRating:     supplierRating,
	}

	return s.supplierRepository.CreateSupplier(ctx, params)
}

func (s *SupplierService) GetSupplierByID(
	ctx context.Context,
	id pgtype.UUID,
) (db.Supplier, error) {
	supplier, err := s.supplierRepository.GetSupplierByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Supplier{}, ErrSupplierNotFound
		}
		return db.Supplier{}, err
	}
	return supplier, nil
}

func (s *SupplierService) GetSupplierByPartyID(
	ctx context.Context,
	partyID pgtype.UUID,
) (db.Supplier, error) {
	supplier, err := s.supplierRepository.GetSupplierByPartyID(ctx, partyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Supplier{}, ErrSupplierNotFound
		}
		return db.Supplier{}, err
	}
	return supplier, nil
}

func (s *SupplierService) ListSuppliers(
	ctx context.Context,
	params db.ListSuppliersParams,
) ([]db.Supplier, error) {
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

	return s.supplierRepository.ListSuppliers(ctx, params)
}

func (s *SupplierService) UpdateSupplier(
	ctx context.Context,
	input domain.UpdateSupplierInput,
) (db.Supplier, error) {
	var supplierID pgtype.UUID
	if err := supplierID.Scan(input.ID); err != nil {
		return db.Supplier{}, ErrInvalidSupplierID
	}

	existing, err := s.supplierRepository.GetSupplierByID(ctx, supplierID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Supplier{}, ErrSupplierNotFound
		}
		return db.Supplier{}, err
	}

	paymentTerms := existing.PaymentTermsDays
	if input.PaymentTermsDays != nil {
		if *input.PaymentTermsDays < 0 {
			return db.Supplier{}, ErrInvalidPaymentTerms
		}
		paymentTerms = int32(*input.PaymentTermsDays)
	}

	currencyCode := existing.CurrencyCode
	if input.CurrencyCode != nil {
		c := strings.ToUpper(strings.TrimSpace(*input.CurrencyCode))
		if c == "" {
			currencyCode = nil
		} else if len(c) != 3 {
			return db.Supplier{}, ErrInvalidCurrencyCode
		} else {
			currencyCode = &c
		}
	}

	supplierRating := existing.SupplierRating
	if input.SupplierRating != nil {
		if *input.SupplierRating < 0 || *input.SupplierRating > 5 {
			return db.Supplier{}, ErrInvalidSupplierRating
		}
		if err := supplierRating.Scan(fmt.Sprintf("%.2f", *input.SupplierRating)); err != nil {
			return db.Supplier{}, ErrInvalidSupplierRating
		}
	}

	purchaseRepID := existing.PurchaseRepID
	if input.PurchaseRepID != nil {
		if *input.PurchaseRepID == "" {
			purchaseRepID = pgtype.UUID{}
		} else {
			if err := purchaseRepID.Scan(*input.PurchaseRepID); err != nil {
				return db.Supplier{}, err
			}
		}
	}

	taxExempt := existing.TaxExempt
	if input.TaxExempt != nil {
		taxExempt = *input.TaxExempt
	}

	supplierGroup := existing.SupplierGroup
	if input.SupplierGroup != nil {
		supplierGroup = input.SupplierGroup
	}

	paymentMethod := existing.PaymentMethod
	if input.PaymentMethod != nil {
		paymentMethod = input.PaymentMethod
	}

	taxExemptionNumber := existing.TaxExemptionNumber
	if input.TaxExemptionNumber != nil {
		taxExemptionNumber = input.TaxExemptionNumber
	}

	params := db.UpdateSupplierParams{
		ID:                 supplierID,
		SupplierGroup:      supplierGroup,
		PaymentTermsDays:   paymentTerms,
		PaymentMethod:      paymentMethod,
		CurrencyCode:       currencyCode,
		PurchaseRepID:      purchaseRepID,
		TaxExempt:          taxExempt,
		TaxExemptionNumber: taxExemptionNumber,
		SupplierRating:     supplierRating,
	}

	return s.supplierRepository.UpdateSupplier(ctx, params)
}

func (s *SupplierService) CreateSupplierUserAccount(
	ctx context.Context,
	input domain.CreateSupplierUserAccountInput,
) (db.SupplierUserAccount, error) {
	var supplierID pgtype.UUID
	if err := supplierID.Scan(input.SupplierID); err != nil {
		return db.SupplierUserAccount{}, ErrInvalidSupplierID
	}

	var userID pgtype.UUID
	if err := userID.Scan(input.UserID); err != nil {
		return db.SupplierUserAccount{}, ErrInvalidUserID
	}

	_, err := s.supplierRepository.GetSupplierByID(ctx, supplierID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.SupplierUserAccount{}, ErrSupplierNotFound
		}
		return db.SupplierUserAccount{}, err
	}

	params := db.CreateSupplierUserAccountParams{
		SupplierID: supplierID,
		UserID:     userID,
	}

	return s.supplierRepository.CreateSupplierUserAccount(ctx, params)
}

func (s *SupplierService) ListSupplierUserAccounts(
	ctx context.Context,
	supplierID pgtype.UUID,
) ([]db.SupplierUserAccount, error) {
	_, err := s.supplierRepository.GetSupplierByID(ctx, supplierID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSupplierNotFound
		}
		return nil, err
	}

	return s.supplierRepository.ListSupplierUserAccountsBySupplierID(ctx, supplierID)
}

func (s *SupplierService) UpdateSupplierUserAccountStatus(
	ctx context.Context,
	id pgtype.UUID,
	isActive bool,
) (db.SupplierUserAccount, error) {
	params := db.UpdateSupplierUserAccountStatusParams{
		ID:       id,
		IsActive: isActive,
	}
	account, err := s.supplierRepository.UpdateSupplierUserAccountStatus(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.SupplierUserAccount{}, ErrSupplierUserAccountNotFound
		}
		return db.SupplierUserAccount{}, err
	}
	return account, nil
}
