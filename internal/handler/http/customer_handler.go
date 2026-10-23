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

type CustomerService interface {
	CreateCustomer(ctx context.Context, input domain.CreateCustomerInput) (db.Customer, error)
	GetCustomerByID(ctx context.Context, id pgtype.UUID) (db.Customer, error)
	GetCustomerByPartyID(ctx context.Context, partyID pgtype.UUID) (db.Customer, error)
	ListCustomers(ctx context.Context, params db.ListCustomersParams) ([]db.Customer, error)
	UpdateCustomer(ctx context.Context, input domain.UpdateCustomerInput) (db.Customer, error)

	CreateCustomerUserAccount(ctx context.Context, input domain.CreateCustomerUserAccountInput) (db.CustomerUserAccount, error)
	ListCustomerUserAccounts(ctx context.Context, customerID pgtype.UUID) ([]db.CustomerUserAccount, error)
	UpdateCustomerUserAccountStatus(ctx context.Context, id pgtype.UUID, isActive bool) (db.CustomerUserAccount, error)
}

type CustomerHandler struct {
	service CustomerService
}

func NewCustomerHandler(service CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func (h *CustomerHandler) CreateCustomer(c *fiber.Ctx) error {
	var req dto.CreateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	input := domain.CreateCustomerInput{
		PartyID:            req.PartyID,
		CustomerCode:       req.CustomerCode,
		CustomerGroup:      req.CustomerGroup,
		PriceListID:        req.PriceListID,
		CreditLimit:        req.CreditLimit,
		PaymentTermsDays:   req.PaymentTermsDays,
		CreditStatus:       req.CreditStatus,
		SalesRepID:         req.SalesRepID,
		TaxExempt:          req.TaxExempt,
		TaxExemptionNumber: req.TaxExemptionNumber,
	}

	customer, err := h.service.CreateCustomer(c.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPartyID):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid party id"})
		case errors.Is(err, service.ErrEmptyCustomerCode):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "customer code is required"})
		case errors.Is(err, service.ErrPartyNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "party not found"})
		case errors.Is(err, service.ErrCustomerAlreadyExists):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "party is already registered as a customer"})
		case errors.Is(err, service.ErrDuplicateCustomerCode):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "customer code is already in use"})
		case errors.Is(err, service.ErrInvalidCreditLimit):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "credit limit must be zero or greater"})
		case errors.Is(err, service.ErrInvalidPaymentTerms):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "payment terms days must be zero or greater"})
		case errors.Is(err, service.ErrInvalidCreditStatus):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "credit status must be ACTIVE, BLOCKED, or HOLD"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create customer"})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(toCustomerResponse(customer))
}

func (h *CustomerHandler) GetCustomerByID(c *fiber.Ctx) error {
	var id pgtype.UUID
	if err := id.Scan(c.Params("id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid customer id"})
	}

	customer, err := h.service.GetCustomerByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrCustomerNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "customer not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get customer"})
	}

	return c.JSON(toCustomerResponse(customer))
}

func (h *CustomerHandler) GetCustomerByPartyID(c *fiber.Ctx) error {
	var partyID pgtype.UUID
	if err := partyID.Scan(c.Params("party_id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid party id"})
	}

	customer, err := h.service.GetCustomerByPartyID(c.Context(), partyID)
	if err != nil {
		if errors.Is(err, service.ErrCustomerNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "customer not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get customer"})
	}

	return c.JSON(toCustomerResponse(customer))
}

func (h *CustomerHandler) ListCustomers(c *fiber.Ctx) error {
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

	params := db.ListCustomersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	if group := c.Query("customer_group"); group != "" {
		params.CustomerGroup = &group
	}
	if status := c.Query("credit_status"); status != "" {
		params.CreditStatus = &status
	}
	if search := c.Query("search"); search != "" {
		params.Search = &search
	}

	customers, err := h.service.ListCustomers(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list customers"})
	}

	response := dto.ListCustomersResponse{
		Data:   make([]dto.CustomerResponse, 0, len(customers)),
		Limit:  limit,
		Offset: offset,
	}

	for _, cust := range customers {
		response.Data = append(response.Data, toCustomerResponse(cust))
	}

	return c.JSON(response)
}

func (h *CustomerHandler) UpdateCustomer(c *fiber.Ctx) error {
	idStr := c.Params("id")
	var id pgtype.UUID
	if err := id.Scan(idStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid customer id"})
	}

	var req dto.UpdateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	input := domain.UpdateCustomerInput{
		ID:                 idStr,
		CustomerGroup:      req.CustomerGroup,
		PriceListID:        req.PriceListID,
		CreditLimit:        req.CreditLimit,
		PaymentTermsDays:   req.PaymentTermsDays,
		CreditStatus:       req.CreditStatus,
		SalesRepID:         req.SalesRepID,
		TaxExempt:          req.TaxExempt,
		TaxExemptionNumber: req.TaxExemptionNumber,
	}

	customer, err := h.service.UpdateCustomer(c.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCustomerID):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid customer id"})
		case errors.Is(err, service.ErrCustomerNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "customer not found"})
		case errors.Is(err, service.ErrInvalidCreditLimit):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "credit limit must be zero or greater"})
		case errors.Is(err, service.ErrInvalidPaymentTerms):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "payment terms days must be zero or greater"})
		case errors.Is(err, service.ErrInvalidCreditStatus):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "credit status must be ACTIVE, BLOCKED, or HOLD"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update customer"})
		}
	}

	return c.JSON(toCustomerResponse(customer))
}

func (h *CustomerHandler) CreateCustomerUserAccount(c *fiber.Ctx) error {
	customerIDStr := c.Params("id")
	var customerID pgtype.UUID
	if err := customerID.Scan(customerIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid customer id"})
	}

	var req dto.CreateCustomerUserAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	input := domain.CreateCustomerUserAccountInput{
		CustomerID: customerIDStr,
		UserID:     req.UserID,
	}

	account, err := h.service.CreateCustomerUserAccount(c.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCustomerID):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid customer id"})
		case errors.Is(err, service.ErrInvalidUserID):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
		case errors.Is(err, service.ErrCustomerNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "customer not found"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to link user account"})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(toCustomerUserAccountResponse(account))
}

func (h *CustomerHandler) ListCustomerUserAccounts(c *fiber.Ctx) error {
	var customerID pgtype.UUID
	if err := customerID.Scan(c.Params("id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid customer id"})
	}

	accounts, err := h.service.ListCustomerUserAccounts(c.Context(), customerID)
	if err != nil {
		if errors.Is(err, service.ErrCustomerNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "customer not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list user accounts"})
	}

	responses := make([]dto.CustomerUserAccountResponse, 0, len(accounts))
	for _, acc := range accounts {
		responses = append(responses, toCustomerUserAccountResponse(acc))
	}

	return c.JSON(responses)
}

func (h *CustomerHandler) UpdateCustomerUserAccountStatus(c *fiber.Ctx) error {
	var accountID pgtype.UUID
	if err := accountID.Scan(c.Params("account_id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid account id"})
	}

	var req dto.UpdateUserAccountStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	account, err := h.service.UpdateCustomerUserAccountStatus(c.Context(), accountID, req.IsActive)
	if err != nil {
		if errors.Is(err, service.ErrCustomerUserAccountNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "customer user account not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update account status"})
	}

	return c.JSON(toCustomerUserAccountResponse(account))
}

func toCustomerResponse(c db.Customer) dto.CustomerResponse {
	var priceListID *string
	if c.PriceListID.Valid {
		val := c.PriceListID.String()
		priceListID = &val
	}

	var salesRepID *string
	if c.SalesRepID.Valid {
		val := c.SalesRepID.String()
		salesRepID = &val
	}

	var creditLimit float64
	floatVal, _ := c.CreditLimit.Float64Value()
	if floatVal.Valid {
		creditLimit = floatVal.Float64
	}

	return dto.CustomerResponse{
		ID:                 c.ID.String(),
		PartyID:            c.PartyID.String(),
		CustomerCode:       c.CustomerCode,
		CustomerGroup:      c.CustomerGroup,
		PriceListID:        priceListID,
		CreditLimit:        creditLimit,
		PaymentTermsDays:   c.PaymentTermsDays,
		CreditStatus:       c.CreditStatus,
		SalesRepID:         salesRepID,
		TaxExempt:          c.TaxExempt,
		TaxExemptionNumber: c.TaxExemptionNumber,
		CreatedAt:          c.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:          c.UpdatedAt.Time.Format(time.RFC3339),
	}
}

func toCustomerUserAccountResponse(a db.CustomerUserAccount) dto.CustomerUserAccountResponse {
	return dto.CustomerUserAccountResponse{
		ID:         a.ID.String(),
		CustomerID: a.CustomerID.String(),
		UserID:     a.UserID.String(),
		IsActive:   a.IsActive,
		CreatedAt:  a.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:  a.UpdatedAt.Time.Format(time.RFC3339),
	}
}
