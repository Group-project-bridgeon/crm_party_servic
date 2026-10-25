package http

import (
	"context"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/group-project/crm_party_service/internal/db"
	"github.com/group-project/crm_party_service/internal/domain"
	"github.com/group-project/crm_party_service/internal/handler/http/dto"
	"github.com/group-project/crm_party_service/internal/service"
	"github.com/jackc/pgx/v5/pgtype"
)

type PartyService interface {
	CreateParty(ctx context.Context, input domain.CreatePartyInput) (db.Party, error)
	GetPartyByID(ctx context.Context, id pgtype.UUID) (db.Party, error)
	ListParties(ctx context.Context, params db.ListPartiesParams) ([]db.Party, error)
	UpdateParty(ctx context.Context, input domain.UpdatePartyInput) (db.Party, error)
	UpdatePartyStatus(ctx context.Context, id pgtype.UUID, isActive bool) (db.Party, error)
}

type PartyHandler struct {
	service PartyService
}

func NewPartyHandler(service PartyService) *PartyHandler {
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

	return c.Status(fiber.StatusCreated).JSON(toPartyResponse(party))
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

	return c.JSON(toPartyResponse(party))
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
		response.Data = append(response.Data, toPartyResponse(party))
	}

	return c.JSON(response)
}

func (h *PartyHandler) UpdateParty(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdatePartyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	input := domain.UpdatePartyInput{
		ID:          id,
		DisplayName: req.DisplayName,
		LegalName:   req.LegalName,
		Email:       req.Email,
		Phone:       req.Phone,
		Website:     req.Website,
		TaxNumber:   req.TaxNumber,
		Notes:       req.Notes,
	}

	party, err := h.service.UpdateParty(c.Context(), input)
	if err != nil {
		switch err {
		case service.ErrInvalidPartyID:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid party id",
			})
		case service.ErrEmptyDisplayName:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "display name is required",
			})
		case service.ErrPartyNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "party not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to update party",
			})
		}
	}

	return c.JSON(toPartyResponse(party))
}

func (h *PartyHandler) UpdatePartyStatus(c *fiber.Ctx) error {
	id := c.Params("id")

	var partyID pgtype.UUID
	if err := partyID.Scan(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid party id",
		})
	}

	var req dto.UpdatePartyStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	party, err := h.service.UpdatePartyStatus(c.Context(), partyID, req.IsActive)
	if err != nil {
		switch err {
		case service.ErrPartyNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "party not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to update party status",
			})
		}
	}

	return c.JSON(toPartyResponse(party))
}

func (h *PartyHandler) DeleteParty(c *fiber.Ctx) error {
	partyIDStr := c.Params("id")
	var partyID pgtype.UUID
	if err := partyID.Scan(partyIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid party id",
		})
	}

	_, err := h.service.UpdatePartyStatus(c.Context(), partyID, false)
	if err != nil {
		if errors.Is(err, service.ErrPartyNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "party not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete party",
		})
	}

	return c.JSON(fiber.Map{
		"message": "party deactivated successfully",
	})
}

func toPartyResponse(party db.Party) dto.PartyResponse {
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

	return response
}
