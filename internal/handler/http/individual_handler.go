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

type IndividualService interface {
	CreateIndividual(
		ctx context.Context,
		input domain.CreateIndividualInput,
	) (db.Individual, error)

	GetIndividualByPartyID(
		ctx context.Context,
		partyID pgtype.UUID,
	) (db.Individual, error)

	UpdateIndividual(
		ctx context.Context,
		input domain.UpdateIndividualInput,
	) (db.Individual, error)
}

type IndividualHandler struct {
	service IndividualService
}

func NewIndividualHandler(
	service IndividualService,
) *IndividualHandler {
	return &IndividualHandler{
		service: service,
	}
}

func (h *IndividualHandler) CreateIndividual(c *fiber.Ctx) error {
	var req dto.CreateIndividualRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	input := domain.CreateIndividualInput{
		PartyID:     req.PartyID,
		FirstName:   req.FirstName,
		MiddleName:  req.MiddleName,
		LastName:    req.LastName,
		DateOfBirth: req.DateOfBirth,
		Gender:      req.Gender,
		JobTitle:    req.JobTitle,
		NationalID:  req.NationalID,
	}

	individual, err := h.service.CreateIndividual(c.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPartyID):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid party id",
			})
		case errors.Is(err, service.ErrEmptyFirstName):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "first name is required",
			})
		case errors.Is(err, service.ErrPartyNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "party not found",
			})
		case errors.Is(err, service.ErrPartyNotIndividual):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "party is not an individual",
			})
		case errors.Is(err, service.ErrIndividualAlreadyExists):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "individual already exists",
			})
		case errors.Is(err, service.ErrInvalidDateOfBirth):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "date_of_birth must be in YYYY-MM-DD format",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(toIndividualResponse(individual))
}

func (h *IndividualHandler) GetIndividualByPartyID(c *fiber.Ctx) error {
	var partyID pgtype.UUID
	if err := partyID.Scan(c.Params("party_id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid party id",
		})
	}

	individual, err := h.service.GetIndividualByPartyID(c.Context(), partyID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrIndividualNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "individual not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(toIndividualResponse(individual))
}

func (h *IndividualHandler) UpdateIndividual(c *fiber.Ctx) error {
	partyIDStr := c.Params("party_id")
	var partyID pgtype.UUID
	if err := partyID.Scan(partyIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid party id",
		})
	}

	var req dto.UpdateIndividualRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	input := domain.UpdateIndividualInput{
		PartyID:     partyIDStr,
		FirstName:   req.FirstName,
		MiddleName:  req.MiddleName,
		LastName:    req.LastName,
		DateOfBirth: req.DateOfBirth,
		Gender:      req.Gender,
		JobTitle:    req.JobTitle,
		NationalID:  req.NationalID,
	}

	individual, err := h.service.UpdateIndividual(c.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPartyID):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid party id",
			})
		case errors.Is(err, service.ErrEmptyFirstName):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "first name is required",
			})
		case errors.Is(err, service.ErrPartyNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "party not found",
			})
		case errors.Is(err, service.ErrPartyNotIndividual):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "party is not an individual",
			})
		case errors.Is(err, service.ErrIndividualNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "individual not found",
			})
		case errors.Is(err, service.ErrInvalidDateOfBirth):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "date_of_birth must be in YYYY-MM-DD format",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(toIndividualResponse(individual))
}

func toIndividualResponse(individual db.Individual) dto.IndividualResponse {
	var dob *string
	if individual.DateOfBirth.Valid {
		val := individual.DateOfBirth.Time.Format("2006-01-02")
		dob = &val
	}

	return dto.IndividualResponse{
		ID:          individual.ID.String(),
		PartyID:     individual.PartyID.String(),
		FirstName:   individual.FirstName,
		MiddleName:  individual.MiddleName,
		LastName:    individual.LastName,
		DateOfBirth: dob,
		Gender:      individual.Gender,
		JobTitle:    individual.JobTitle,
		NationalID:  individual.NationalID,
		CreatedAt:   individual.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:   individual.UpdatedAt.Time.Format(time.RFC3339),
	}
}
