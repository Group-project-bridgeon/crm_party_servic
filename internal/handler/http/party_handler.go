package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/group-project/crm_party_service/internal/db"
	"github.com/group-project/crm_party_service/internal/domain"
	"github.com/group-project/crm_party_service/internal/handler/http/dto"
	"github.com/group-project/crm_party_service/internal/service"
	"github.com/jackc/pgx/v5/pgtype"
	"strconv"
)

type PartyHandler struct {
	service *service.PartyService
}

func NewPartyHandler(service *service.PartyService) *PartyHandler {
	return &PartyHandler{
		service: service,
	}
}

func (h *PartyHandler) CreateParty(c *fiber.Ctx) error {
	var req dto.CreatePartyRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	input := domain.CreatePartyInput{
		PartyType:   domain.PartyType(req.PartyType),
		DisplayName: req.DisplayName,
		LegalName:   req.LegalName,
		Email:       req.Email,
		Phone:       req.Phone,
		Website:     req.Website,
		TaxNumber:   req.TaxNumber,
		Notes:       req.Notes,
	}

	party, err := h.service.CreateParty(c.Context(), input)
	if err != nil {
		switch err {
		case service.ErrInvalidPartyType:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid party type",
			})

		case service.ErrEmptyDisplayName:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "display name is required",
			})

		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to create party",
			})
		}
	}

	response := dto.PartyResponse{
		ID:          party.ID.String(),
		PartyType:   party.PartyType,
		DisplayName: party.DisplayName,
		LegalName:   party.LegalName,
		Email:       party.Email,
		Phone:       party.Phone,
		Website:     party.Website,
		TaxNumber:   party.TaxNumber,
		Notes:       party.Notes,
		IsActive:    party.IsActive,
	}

	if party.CreatedAt.Valid {
		response.CreatedAt = party.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	if party.UpdatedAt.Valid {
		response.UpdatedAt = party.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *PartyHandler) GetPartyByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var partyID pgtype.UUID

	if err := partyID.Scan(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid party id",
		})
	}

	party, err := h.service.GetPartyByID(c.Context(), partyID)
	if err != nil {
		switch err {
		case service.ErrPartyNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "party not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to get party",
			})
		}
	}

	response := dto.PartyResponse{
		ID:          party.ID.String(),
		PartyType:   party.PartyType,
		DisplayName: party.DisplayName,
		LegalName:   party.LegalName,
		Email:       party.Email,
		Phone:       party.Phone,
		Website:     party.Website,
		TaxNumber:   party.TaxNumber,
		Notes:       party.Notes,
		IsActive:    party.IsActive,
	}

	if party.CreatedAt.Valid {
		response.CreatedAt = party.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	if party.UpdatedAt.Valid {
		response.UpdatedAt = party.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	return c.JSON(response)
}

func (h *PartyHandler) ListParties(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "offset cannot be negative",
		})
	}

	params := db.ListPartiesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	if partyType := c.Query("party_type"); partyType != "" {
		switch domain.PartyType(partyType) {
		case domain.PartyTypeOrganization, domain.PartyTypeIndividual:
			params.PartyType = &partyType
		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid party_type",
			})
		}
	}

	if rawIsActive := c.Context().QueryArgs().Peek("is_active"); rawIsActive != nil {
		value, err := strconv.ParseBool(string(rawIsActive))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "is_active must be true or false",
			})
		}

		params.IsActive = &value
	}

	if search := c.Query("search"); search != "" {
		params.Search = &search
	}

	parties, err := h.service.ListParties(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to list parties",
		})
	}

	response := dto.ListPartiesResponse{
		Data:   make([]dto.PartyResponse, 0, len(parties)),
		Limit:  limit,
		Offset: offset,
	}

	for _, party := range parties {
		item := dto.PartyResponse{
			ID:          party.ID.String(),
			PartyType:   party.PartyType,
			DisplayName: party.DisplayName,
			LegalName:   party.LegalName,
			Email:       party.Email,
			Phone:       party.Phone,
			Website:     party.Website,
			TaxNumber:   party.TaxNumber,
			Notes:       party.Notes,
			IsActive:    party.IsActive,
		}

		if party.CreatedAt.Valid {
			item.CreatedAt = party.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
		}

		if party.UpdatedAt.Valid {
			item.UpdatedAt = party.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
		}

		response.Data = append(response.Data, item)
	}

	return c.JSON(response)
}
