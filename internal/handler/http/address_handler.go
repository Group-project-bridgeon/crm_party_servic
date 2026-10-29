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

type AddressService interface {
	CreateAddress(ctx context.Context, input domain.CreateAddressInput) (db.PartyAddress, error)
	GetAddressByID(ctx context.Context, id pgtype.UUID) (db.PartyAddress, error)
	ListAddressesByPartyID(ctx context.Context, partyID pgtype.UUID, addressType *string) ([]db.PartyAddress, error)
	UpdateAddress(ctx context.Context, input domain.UpdateAddressInput) (db.PartyAddress, error)
	DeleteAddress(ctx context.Context, id pgtype.UUID) error
}

type AddressHandler struct {
	service AddressService
}

func NewAddressHandler(service AddressService) *AddressHandler {
	return &AddressHandler{service: service}
}

func (h *AddressHandler) CreateAddress(c *fiber.Ctx) error {
	partyIDStr := c.Params("party_id")
	var partyID pgtype.UUID
	if err := partyID.Scan(partyIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid party id"})
	}

	var req dto.CreateAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	input := domain.CreateAddressInput{
		PartyID:      partyIDStr,
		AddressType:  req.AddressType,
		AddressName:  req.AddressName,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		Landmark:     req.Landmark,
		City:         req.City,
		State:        req.State,
		PostalCode:   req.PostalCode,
		CountryCode:  req.CountryCode,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		IsPrimary:    req.IsPrimary,
		IsActive:     req.IsActive,
	}

	address, err := h.service.CreateAddress(c.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPartyID):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid party id"})
		case errors.Is(err, service.ErrPartyNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "party not found"})
		case errors.Is(err, service.ErrInvalidAddressType):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid address type"})
		case errors.Is(err, service.ErrEmptyAddressLine1):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "address_line1 is required"})
		case errors.Is(err, service.ErrEmptyCity):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "city is required"})
		case errors.Is(err, service.ErrInvalidCountryCode):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "country_code must be a 2-letter ISO code"})
		case errors.Is(err, service.ErrInvalidCoordinates):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "latitude must be between -90 and 90, longitude between -180 and 180"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create address"})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(toAddressResponse(address))
}

func (h *AddressHandler) GetAddressByID(c *fiber.Ctx) error {
	var id pgtype.UUID
	if err := id.Scan(c.Params("address_id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid address id"})
	}

	address, err := h.service.GetAddressByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrAddressNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "address not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get address"})
	}

	return c.JSON(toAddressResponse(address))
}

func (h *AddressHandler) ListAddressesByPartyID(c *fiber.Ctx) error {
	var partyID pgtype.UUID
	if err := partyID.Scan(c.Params("party_id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid party id"})
	}

	var addressType *string
	if at := c.Query("address_type"); at != "" {
		addressType = &at
	}

	addresses, err := h.service.ListAddressesByPartyID(c.Context(), partyID, addressType)
	if err != nil {
		if errors.Is(err, service.ErrPartyNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "party not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list addresses"})
	}

	responses := make([]dto.AddressResponse, 0, len(addresses))
	for _, a := range addresses {
		responses = append(responses, toAddressResponse(a))
	}

	return c.JSON(responses)
}

func (h *AddressHandler) UpdateAddress(c *fiber.Ctx) error {
	idStr := c.Params("address_id")
	var id pgtype.UUID
	if err := id.Scan(idStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid address id"})
	}

	var req dto.UpdateAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	input := domain.UpdateAddressInput{
		ID:           idStr,
		AddressType:  req.AddressType,
		AddressName:  req.AddressName,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		Landmark:     req.Landmark,
		City:         req.City,
		State:        req.State,
		PostalCode:   req.PostalCode,
		CountryCode:  req.CountryCode,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		IsPrimary:    req.IsPrimary,
		IsActive:     req.IsActive,
	}

	address, err := h.service.UpdateAddress(c.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidAddressID):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid address id"})
		case errors.Is(err, service.ErrAddressNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "address not found"})
		case errors.Is(err, service.ErrInvalidAddressType):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid address type"})
		case errors.Is(err, service.ErrEmptyAddressLine1):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "address_line1 is required"})
		case errors.Is(err, service.ErrEmptyCity):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "city is required"})
		case errors.Is(err, service.ErrInvalidCountryCode):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "country_code must be a 2-letter ISO code"})
		case errors.Is(err, service.ErrInvalidCoordinates):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "latitude must be between -90 and 90, longitude between -180 and 180"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update address"})
		}
	}

	return c.JSON(toAddressResponse(address))
}

func (h *AddressHandler) DeleteAddress(c *fiber.Ctx) error {
	var id pgtype.UUID
	if err := id.Scan(c.Params("address_id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid address id"})
	}

	err := h.service.DeleteAddress(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrAddressNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "address not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete address"})
	}

	return c.JSON(fiber.Map{"message": "address deleted successfully"})
}

func toAddressResponse(a db.PartyAddress) dto.AddressResponse {
	var lat *float64
	if a.Latitude.Valid {
		f, _ := a.Latitude.Float64Value()
		if f.Valid {
			lat = &f.Float64
		}
	}

	var lon *float64
	if a.Longitude.Valid {
		f, _ := a.Longitude.Float64Value()
		if f.Valid {
			lon = &f.Float64
		}
	}

	return dto.AddressResponse{
		ID:           a.ID.String(),
		PartyID:      a.PartyID.String(),
		AddressType:  a.AddressType,
		AddressName:  a.AddressName,
		AddressLine1: a.AddressLine1,
		AddressLine2: a.AddressLine2,
		Landmark:     a.Landmark,
		City:         a.City,
		State:        a.State,
		PostalCode:   a.PostalCode,
		CountryCode:  a.CountryCode,
		Latitude:     lat,
		Longitude:    lon,
		IsPrimary:    a.IsPrimary,
		IsActive:     a.IsActive,
		CreatedAt:    a.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:    a.UpdatedAt.Time.Format(time.RFC3339),
	}
}
