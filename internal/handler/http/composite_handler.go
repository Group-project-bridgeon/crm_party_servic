package http

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/domain"
	"github.com/group-project/crm_party_service/internal/handler/http/dto"
	"github.com/group-project/crm_party_service/internal/service"
)

type CompositeService interface {
	GetParty360(ctx context.Context, partyID pgtype.UUID) (*domain.Party360, error)
	CreateCompositeCustomer(ctx context.Context, input domain.CompositeCustomerInput) (*domain.CompositeCustomerResult, error)
	CreateCompositeSupplier(ctx context.Context, input domain.CompositeSupplierInput) (*domain.CompositeSupplierResult, error)
}

type CompositeHandler struct {
	service CompositeService
}

func NewCompositeHandler(service CompositeService) *CompositeHandler {
	return &CompositeHandler{service: service}
}

func (h *CompositeHandler) GetParty360(c *fiber.Ctx) error {
	var partyID pgtype.UUID
	if err := partyID.Scan(c.Params("id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid party id"})
	}

	result, err := h.service.GetParty360(c.Context(), partyID)
	if err != nil {
		if errors.Is(err, service.ErrPartyNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "party not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get party 360 view"})
	}

	resp := dto.Party360Response{
		Party:          toPartyResponse(result.Party),
		Addresses:      make([]dto.AddressResponse, 0, len(result.Addresses)),
		Contacts:       make([]dto.ContactResponse, 0, len(result.Contacts)),
		ContactPersons: make([]dto.ContactPersonResponse, 0, len(result.ContactPersons)),
		Notes:          make([]dto.NoteResponse, 0, len(result.Notes)),
	}

	if result.Organization != nil {
		orgResp := toOrganizationResponse(*result.Organization)
		resp.Organization = &orgResp
	}

	if result.Individual != nil {
		indivResp := toIndividualResponse(*result.Individual)
		resp.Individual = &indivResp
	}

	if result.Customer != nil {
		custResp := toCustomerResponse(*result.Customer)
		resp.Customer = &custResp
	}

	if result.Supplier != nil {
		suppResp := toSupplierResponse(*result.Supplier)
		resp.Supplier = &suppResp
	}

	for _, a := range result.Addresses {
		resp.Addresses = append(resp.Addresses, toAddressResponse(a))
	}
	for _, ct := range result.Contacts {
		resp.Contacts = append(resp.Contacts, toContactResponse(ct))
	}
	for _, p := range result.ContactPersons {
		resp.ContactPersons = append(resp.ContactPersons, toContactPersonResponse(p))
	}
	for _, n := range result.Notes {
		resp.Notes = append(resp.Notes, toNoteResponse(n))
	}

	return c.JSON(resp)
}

func (h *CompositeHandler) CreateCompositeCustomer(c *fiber.Ctx) error {
	var req dto.CreateCompositeCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	var addrInput *domain.CreateAddressInput
	if req.Address != nil {
		addrInput = &domain.CreateAddressInput{
			AddressType:  req.Address.AddressType,
			AddressName:  req.Address.AddressName,
			AddressLine1: req.Address.AddressLine1,
			AddressLine2: req.Address.AddressLine2,
			Landmark:     req.Address.Landmark,
			City:         req.Address.City,
			State:        req.Address.State,
			PostalCode:   req.Address.PostalCode,
			CountryCode:  req.Address.CountryCode,
			Latitude:     req.Address.Latitude,
			Longitude:    req.Address.Longitude,
			IsPrimary:    true,
		}
	}

	input := domain.CompositeCustomerInput{
		PartyType:          domain.PartyType(req.PartyType),
		DisplayName:        req.DisplayName,
		LegalName:          req.LegalName,
		Email:              req.Email,
		Phone:              req.Phone,
		Website:            req.Website,
		TaxNumber:          req.TaxNumber,
		Notes:              req.Notes,
		RegistrationNumber: req.RegistrationNumber,
		Industry:           req.Industry,
		OrganizationType:   req.OrganizationType,
		FoundedDate:        req.FoundedDate,
		EmployeeCount:      req.EmployeeCount,
		Description:        req.Description,
		FirstName:          req.FirstName,
		MiddleName:         req.MiddleName,
		LastName:           req.LastName,
		DateOfBirth:        req.DateOfBirth,
		Gender:             req.Gender,
		JobTitle:           req.JobTitle,
		NationalID:         req.NationalID,
		CustomerCode:       req.CustomerCode,
		CustomerGroup:      req.CustomerGroup,
		PriceListID:        req.PriceListID,
		CreditLimit:        req.CreditLimit,
		PaymentTermsDays:   req.PaymentTermsDays,
		CreditStatus:       req.CreditStatus,
		SalesRepID:         req.SalesRepID,
		TaxExempt:          req.TaxExempt,
		TaxExemptionNumber: req.TaxExemptionNumber,
		Address:            addrInput,
	}

	result, err := h.service.CreateCompositeCustomer(c.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPartyType):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid party type"})
		case errors.Is(err, service.ErrEmptyDisplayName):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "display name is required"})
		case errors.Is(err, service.ErrEmptyCustomerCode):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "customer code is required"})
		case errors.Is(err, service.ErrDuplicateCustomerCode):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "customer code already exists"})
		case errors.Is(err, service.ErrEmptyFirstName):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "first name is required for individual customer"})
		case errors.Is(err, service.ErrInvalidCreditLimit):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "credit limit must be zero or greater"})
		case errors.Is(err, service.ErrInvalidPaymentTerms):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "payment terms days must be zero or greater"})
		case errors.Is(err, service.ErrInvalidCreditStatus):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "credit status must be ACTIVE, BLOCKED, or HOLD"})
		case errors.Is(err, service.ErrEmptyAddressLine1):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "address line 1 is required"})
		case errors.Is(err, service.ErrEmptyCity):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "city is required"})
		case errors.Is(err, service.ErrInvalidCountryCode):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "country code must be a 2-letter ISO code"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to onboard customer: " + err.Error()})
		}
	}

	resp := dto.CompositeCustomerResponse{
		Party:    toPartyResponse(result.Party),
		Customer: toCustomerResponse(result.Customer),
	}
	if result.Organization != nil {
		orgResp := toOrganizationResponse(*result.Organization)
		resp.Organization = &orgResp
	}
	if result.Individual != nil {
		indivResp := toIndividualResponse(*result.Individual)
		resp.Individual = &indivResp
	}
	if result.Address != nil {
		addrResp := toAddressResponse(*result.Address)
		resp.Address = &addrResp
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *CompositeHandler) CreateCompositeSupplier(c *fiber.Ctx) error {
	var req dto.CreateCompositeSupplierRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	var addrInput *domain.CreateAddressInput
	if req.Address != nil {
		addrInput = &domain.CreateAddressInput{
			AddressType:  req.Address.AddressType,
			AddressName:  req.Address.AddressName,
			AddressLine1: req.Address.AddressLine1,
			AddressLine2: req.Address.AddressLine2,
			Landmark:     req.Address.Landmark,
			City:         req.Address.City,
			State:        req.Address.State,
			PostalCode:   req.Address.PostalCode,
			CountryCode:  req.Address.CountryCode,
			Latitude:     req.Address.Latitude,
			Longitude:    req.Address.Longitude,
			IsPrimary:    true,
		}
	}

	input := domain.CompositeSupplierInput{
		PartyType:          domain.PartyType(req.PartyType),
		DisplayName:        req.DisplayName,
		LegalName:          req.LegalName,
		Email:              req.Email,
		Phone:              req.Phone,
		Website:            req.Website,
		TaxNumber:          req.TaxNumber,
		Notes:              req.Notes,
		RegistrationNumber: req.RegistrationNumber,
		Industry:           req.Industry,
		OrganizationType:   req.OrganizationType,
		FoundedDate:        req.FoundedDate,
		EmployeeCount:      req.EmployeeCount,
		Description:        req.Description,
		FirstName:          req.FirstName,
		MiddleName:         req.MiddleName,
		LastName:           req.LastName,
		DateOfBirth:        req.DateOfBirth,
		Gender:             req.Gender,
		JobTitle:           req.JobTitle,
		NationalID:         req.NationalID,
		SupplierCode:       req.SupplierCode,
		SupplierGroup:      req.SupplierGroup,
		PaymentTermsDays:   req.PaymentTermsDays,
		PaymentMethod:      req.PaymentMethod,
		CurrencyCode:       req.CurrencyCode,
		PurchaseRepID:      req.PurchaseRepID,
		TaxExempt:          req.TaxExempt,
		TaxExemptionNumber: req.TaxExemptionNumber,
		SupplierRating:     req.SupplierRating,
		Address:            addrInput,
	}

	result, err := h.service.CreateCompositeSupplier(c.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPartyType):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid party type"})
		case errors.Is(err, service.ErrEmptyDisplayName):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "display name is required"})
		case errors.Is(err, service.ErrEmptySupplierCode):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "supplier code is required"})
		case errors.Is(err, service.ErrDuplicateSupplierCode):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "supplier code already exists"})
		case errors.Is(err, service.ErrEmptyFirstName):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "first name is required for individual supplier"})
		case errors.Is(err, service.ErrInvalidPaymentTerms):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "payment terms days must be zero or greater"})
		case errors.Is(err, service.ErrInvalidCurrencyCode):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "currency code must be 3 uppercase letters"})
		case errors.Is(err, service.ErrInvalidSupplierRating):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "supplier rating must be between 0 and 5"})
		case errors.Is(err, service.ErrEmptyAddressLine1):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "address line 1 is required"})
		case errors.Is(err, service.ErrEmptyCity):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "city is required"})
		case errors.Is(err, service.ErrInvalidCountryCode):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "country code must be a 2-letter ISO code"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to onboard supplier: " + err.Error()})
		}
	}

	resp := dto.CompositeSupplierResponse{
		Party:    toPartyResponse(result.Party),
		Supplier: toSupplierResponse(result.Supplier),
	}
	if result.Organization != nil {
		orgResp := toOrganizationResponse(*result.Organization)
		resp.Organization = &orgResp
	}
	if result.Individual != nil {
		indivResp := toIndividualResponse(*result.Individual)
		resp.Individual = &indivResp
	}
	if result.Address != nil {
		addrResp := toAddressResponse(*result.Address)
		resp.Address = &addrResp
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}
