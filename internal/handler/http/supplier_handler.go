package http

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
	"github.com/group-project/crm_party_service/internal/domain"
	"github.com/group-project/crm_party_service/internal/handler/http/dto"
	"github.com/group-project/crm_party_service/internal/service"
)

type SupplierService interface {
	CreateSupplier(ctx context.Context, input domain.CreateSupplierInput) (db.Supplier, error)
	GetSupplierByID(ctx context.Context, id pgtype.UUID) (db.Supplier, error)
	GetSupplierByPartyID(ctx context.Context, partyID pgtype.UUID) (db.Supplier, error)
	ListSuppliers(ctx context.Context, params db.ListSuppliersParams) ([]db.Supplier, error)
	UpdateSupplier(ctx context.Context, input domain.UpdateSupplierInput) (db.Supplier, error)

	CreateSupplierUserAccount(ctx context.Context, input domain.CreateSupplierUserAccountInput) (db.SupplierUserAccount, error)
	ListSupplierUserAccounts(ctx context.Context, supplierID pgtype.UUID) ([]db.SupplierUserAccount, error)
	UpdateSupplierUserAccountStatus(ctx context.Context, id pgtype.UUID, isActive bool) (db.SupplierUserAccount, error)
}

type SupplierHandler struct {
	service SupplierService
}

func NewSupplierHandler(service SupplierService) *SupplierHandler {
	return &SupplierHandler{service: service}
}

func (h *SupplierHandler) CreateSupplier(c *fiber.Ctx) error {
	var req dto.CreateSupplierRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	input := domain.CreateSupplierInput{
		PartyID:            req.PartyID,
		SupplierCode:       req.SupplierCode,
		SupplierGroup:      req.SupplierGroup,
		PaymentTermsDays:   req.PaymentTermsDays,
		PaymentMethod:      req.PaymentMethod,
		CurrencyCode:       req.CurrencyCode,
		PurchaseRepID:      req.PurchaseRepID,
		TaxExempt:          req.TaxExempt,
		TaxExemptionNumber: req.TaxExemptionNumber,
		SupplierRating:     req.SupplierRating,
	}

	supplier, err := h.service.CreateSupplier(c.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPartyID):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid party id"})
		case errors.Is(err, service.ErrEmptySupplierCode):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "supplier code is required"})
		case errors.Is(err, service.ErrPartyNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "party not found"})
		case errors.Is(err, service.ErrSupplierAlreadyExists):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "party is already registered as a supplier"})
		case errors.Is(err, service.ErrDuplicateSupplierCode):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "supplier code is already in use"})
		case errors.Is(err, service.ErrInvalidPaymentTerms):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "payment terms days must be zero or greater"})
		case errors.Is(err, service.ErrInvalidCurrencyCode):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "currency code must be 3 uppercase letters"})
		case errors.Is(err, service.ErrInvalidSupplierRating):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "supplier rating must be between 0 and 5"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create supplier"})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(toSupplierResponse(supplier))
}

func (h *SupplierHandler) GetSupplierByID(c *fiber.Ctx) error {
	var id pgtype.UUID
	if err := id.Scan(c.Params("id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid supplier id"})
	}

	supplier, err := h.service.GetSupplierByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrSupplierNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "supplier not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get supplier"})
	}

	return c.JSON(toSupplierResponse(supplier))
}

func (h *SupplierHandler) GetSupplierByPartyID(c *fiber.Ctx) error {
	var partyID pgtype.UUID
	if err := partyID.Scan(c.Params("party_id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid party id"})
	}

	supplier, err := h.service.GetSupplierByPartyID(c.Context(), partyID)
	if err != nil {
		if errors.Is(err, service.ErrSupplierNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "supplier not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get supplier"})
	}

	return c.JSON(toSupplierResponse(supplier))
}

func (h *SupplierHandler) ListSuppliers(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "offset cannot be negative"})
	}

	params := db.ListSuppliersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	if group := c.Query("supplier_group"); group != "" {
		params.SupplierGroup = &group
	}
	if currency := c.Query("currency_code"); currency != "" {
		params.CurrencyCode = &currency
	}
	if search := c.Query("search"); search != "" {
		params.Search = &search
	}

	suppliers, err := h.service.ListSuppliers(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list suppliers"})
	}

	response := dto.ListSuppliersResponse{
		Data:   make([]dto.SupplierResponse, 0, len(suppliers)),
		Limit:  limit,
		Offset: offset,
	}

	for _, s := range suppliers {
		response.Data = append(response.Data, toSupplierResponse(s))
	}

	return c.JSON(response)
}

func (h *SupplierHandler) UpdateSupplier(c *fiber.Ctx) error {
	idStr := c.Params("id")
	var id pgtype.UUID
	if err := id.Scan(idStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid supplier id"})
	}

	var req dto.UpdateSupplierRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	input := domain.UpdateSupplierInput{
		ID:                 idStr,
		SupplierGroup:      req.SupplierGroup,
		PaymentTermsDays:   req.PaymentTermsDays,
		PaymentMethod:      req.PaymentMethod,
		CurrencyCode:       req.CurrencyCode,
		PurchaseRepID:      req.PurchaseRepID,
		TaxExempt:          req.TaxExempt,
		TaxExemptionNumber: req.TaxExemptionNumber,
		SupplierRating:     req.SupplierRating,
	}

	supplier, err := h.service.UpdateSupplier(c.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidSupplierID):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid supplier id"})
		case errors.Is(err, service.ErrSupplierNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "supplier not found"})
		case errors.Is(err, service.ErrInvalidPaymentTerms):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "payment terms days must be zero or greater"})
		case errors.Is(err, service.ErrInvalidCurrencyCode):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "currency code must be 3 uppercase letters"})
		case errors.Is(err, service.ErrInvalidSupplierRating):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "supplier rating must be between 0 and 5"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update supplier"})
		}
	}

	return c.JSON(toSupplierResponse(supplier))
}

func (h *SupplierHandler) CreateSupplierUserAccount(c *fiber.Ctx) error {
	supplierIDStr := c.Params("id")
	var supplierID pgtype.UUID
	if err := supplierID.Scan(supplierIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid supplier id"})
	}

	var req dto.CreateSupplierUserAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	input := domain.CreateSupplierUserAccountInput{
		SupplierID: supplierIDStr,
		UserID:     req.UserID,
	}

	account, err := h.service.CreateSupplierUserAccount(c.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidSupplierID):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid supplier id"})
		case errors.Is(err, service.ErrInvalidUserID):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
		case errors.Is(err, service.ErrSupplierNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "supplier not found"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to link user account"})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(toSupplierUserAccountResponse(account))
}

func (h *SupplierHandler) ListSupplierUserAccounts(c *fiber.Ctx) error {
	var supplierID pgtype.UUID
	if err := supplierID.Scan(c.Params("id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid supplier id"})
	}

	accounts, err := h.service.ListSupplierUserAccounts(c.Context(), supplierID)
	if err != nil {
		if errors.Is(err, service.ErrSupplierNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "supplier not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list user accounts"})
	}

	responses := make([]dto.SupplierUserAccountResponse, 0, len(accounts))
	for _, acc := range accounts {
		responses = append(responses, toSupplierUserAccountResponse(acc))
	}

	return c.JSON(responses)
}

func (h *SupplierHandler) UpdateSupplierUserAccountStatus(c *fiber.Ctx) error {
	var accountID pgtype.UUID
	if err := accountID.Scan(c.Params("account_id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid account id"})
	}

	var req dto.UpdateUserAccountStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	account, err := h.service.UpdateSupplierUserAccountStatus(c.Context(), accountID, req.IsActive)
	if err != nil {
		if errors.Is(err, service.ErrSupplierUserAccountNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "supplier user account not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update account status"})
	}

	return c.JSON(toSupplierUserAccountResponse(account))
}

func toSupplierResponse(s db.Supplier) dto.SupplierResponse {
	var purchaseRepID *string
	if s.PurchaseRepID.Valid {
		val := s.PurchaseRepID.String()
		purchaseRepID = &val
	}

	var supplierRating *float64
	floatVal, _ := s.SupplierRating.Float64Value()
	if floatVal.Valid {
		val := floatVal.Float64
		supplierRating = &val
	}

	return dto.SupplierResponse{
		ID:                 s.ID.String(),
		PartyID:            s.PartyID.String(),
		SupplierCode:       s.SupplierCode,
		SupplierGroup:      s.SupplierGroup,
		PaymentTermsDays:   s.PaymentTermsDays,
		PaymentMethod:      s.PaymentMethod,
		CurrencyCode:       s.CurrencyCode,
		PurchaseRepID:      purchaseRepID,
		TaxExempt:          s.TaxExempt,
		TaxExemptionNumber: s.TaxExemptionNumber,
		SupplierRating:     supplierRating,
		CreatedAt:          s.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:          s.UpdatedAt.Time.Format(time.RFC3339),
	}
}

func toSupplierUserAccountResponse(a db.SupplierUserAccount) dto.SupplierUserAccountResponse {
	return dto.SupplierUserAccountResponse{
		ID:         a.ID.String(),
		SupplierID: a.SupplierID.String(),
		UserID:     a.UserID.String(),
		IsActive:   a.IsActive,
		CreatedAt:  a.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:  a.UpdatedAt.Time.Format(time.RFC3339),
	}
}
