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

type OrganizationService interface {
    CreateOrganization(
        ctx context.Context,
        input domain.CreateOrganizationInput,
    ) (db.Organization, error)

    GetOrganizationByPartyID(
        ctx context.Context,
        partyID pgtype.UUID,
    ) (db.Organization, error)
}


type OrganizationHandler struct {
	service OrganizationService
}

func NewOrganizationHandler(
	service OrganizationService,
) *OrganizationHandler {
	return &OrganizationHandler{
		service: service,
	}
}

func (h *OrganizationHandler) CreateOrganization(c *fiber.Ctx) error {
	var req dto.CreateOrganizationRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	input := domain.CreateOrganizationInput{
		PartyID:            req.PartyID,
		RegistrationNumber: req.RegistrationNumber,
		Industry:           req.Industry,
		OrganizationType:   req.OrganizationType,
		FoundedDate:        req.FoundedDate,
		EmployeeCount:      req.EmployeeCount,
		Description:        req.Description,
	}

	organization, err := h.service.CreateOrganization(c.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPartyID):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid party id",
			})

		case errors.Is(err, service.ErrPartyNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "party not found",
			})

		case errors.Is(err, service.ErrPartyNotOrganization):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "party is not an organization",
			})

		case errors.Is(err, service.ErrOrganizationAlreadyExists):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "organization already exists",
			})

		case errors.Is(err, service.ErrInvalidEmployeeCount):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "employee_count must be zero or greater",
			})

		case errors.Is(err, service.ErrInvalidFoundedDate):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "founded_date must be in YYYY-MM-DD format",
			})

		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(toOrganizationResponse(organization))
}

func toOrganizationResponse(
	organization db.Organization,
) dto.OrganizationResponse {
	var foundedDate *string

	if organization.FoundedDate.Valid {
		value := organization.FoundedDate.Time.Format("2006-01-02")
		foundedDate = &value
	}

	var employeeCount *int32
	if organization.EmployeeCount != nil {
		value := *organization.EmployeeCount
		employeeCount = &value
	}

	return dto.OrganizationResponse{
		ID:                 organization.ID.String(),
		PartyID:            organization.PartyID.String(),
		RegistrationNumber: organization.RegistrationNumber,
		Industry:           organization.Industry,
		OrganizationType:   organization.OrganizationType,
		FoundedDate:        foundedDate,
		EmployeeCount:      employeeCount,
		Description:        organization.Description,
		CreatedAt:          organization.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:          organization.UpdatedAt.Time.Format(time.RFC3339),
	}
}


func (h *OrganizationHandler) GetOrganizationByPartyID(c *fiber.Ctx) error {
    var partyID pgtype.UUID

    if err := partyID.Scan(c.Params("party_id")); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid party id",
        })
    }

    organization, err := h.service.GetOrganizationByPartyID(
        c.Context(),
        partyID,
    )
    if err != nil {
        switch {
        case errors.Is(err, service.ErrOrganizationNotFound):
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": "organization not found",
            })

        default:
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "internal server error",
            })
        }
    }

    return c.Status(fiber.StatusOK).JSON(
        toOrganizationResponse(organization),
    )
}