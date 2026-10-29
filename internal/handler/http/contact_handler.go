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

type ContactService interface {
	CreateContact(ctx context.Context, input domain.CreateContactInput) (db.PartyContact, error)
	GetContactByID(ctx context.Context, id pgtype.UUID) (db.PartyContact, error)
	ListContactsByPartyID(ctx context.Context, partyID pgtype.UUID, contactType *string) ([]db.PartyContact, error)
	UpdateContact(ctx context.Context, input domain.UpdateContactInput) (db.PartyContact, error)
	DeleteContact(ctx context.Context, id pgtype.UUID) error

	CreateContactPerson(ctx context.Context, input domain.CreateContactPersonInput) (db.PartyContactPerson, error)
	GetContactPersonByID(ctx context.Context, id pgtype.UUID) (db.PartyContactPerson, error)
	ListContactPersonsByPartyID(ctx context.Context, partyID pgtype.UUID) ([]db.PartyContactPerson, error)
	UpdateContactPerson(ctx context.Context, input domain.UpdateContactPersonInput) (db.PartyContactPerson, error)
	DeleteContactPerson(ctx context.Context, id pgtype.UUID) error
}

type ContactHandler struct {
	service ContactService
}

func NewContactHandler(service ContactService) *ContactHandler {
	return &ContactHandler{service: service}
}

// ----------------- Direct Contacts -----------------

func (h *ContactHandler) CreateContact(c *fiber.Ctx) error {
	partyIDStr := c.Params("party_id")
	var partyID pgtype.UUID
	if err := partyID.Scan(partyIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid party id"})
	}

	var req dto.CreateContactRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	input := domain.CreateContactInput{
		PartyID:      partyIDStr,
		ContactType:  req.ContactType,
		ContactValue: req.ContactValue,
		Label:        req.Label,
		IsPrimary:    req.IsPrimary,
		IsActive:     req.IsActive,
	}

	contact, err := h.service.CreateContact(c.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPartyID):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid party id"})
		case errors.Is(err, service.ErrPartyNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "party not found"})
		case errors.Is(err, service.ErrInvalidContactType):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid contact type"})
		case errors.Is(err, service.ErrEmptyContactValue):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "contact value is required"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create contact"})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(toContactResponse(contact))
}

func (h *ContactHandler) GetContactByID(c *fiber.Ctx) error {
	var id pgtype.UUID
	if err := id.Scan(c.Params("contact_id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid contact id"})
	}

	contact, err := h.service.GetContactByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrContactNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "contact not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get contact"})
	}

	return c.JSON(toContactResponse(contact))
}

func (h *ContactHandler) ListContactsByPartyID(c *fiber.Ctx) error {
	var partyID pgtype.UUID
	if err := partyID.Scan(c.Params("party_id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid party id"})
	}

	var contactType *string
	if ct := c.Query("contact_type"); ct != "" {
		contactType = &ct
	}

	contacts, err := h.service.ListContactsByPartyID(c.Context(), partyID, contactType)
	if err != nil {
		if errors.Is(err, service.ErrPartyNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "party not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list contacts"})
	}

	responses := make([]dto.ContactResponse, 0, len(contacts))
	for _, ct := range contacts {
		responses = append(responses, toContactResponse(ct))
	}

	return c.JSON(responses)
}

func (h *ContactHandler) UpdateContact(c *fiber.Ctx) error {
	idStr := c.Params("contact_id")
	var id pgtype.UUID
	if err := id.Scan(idStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid contact id"})
	}

	var req dto.UpdateContactRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	input := domain.UpdateContactInput{
		ID:           idStr,
		ContactType:  req.ContactType,
		ContactValue: req.ContactValue,
		Label:        req.Label,
		IsPrimary:    req.IsPrimary,
		IsActive:     req.IsActive,
	}

	contact, err := h.service.UpdateContact(c.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidContactID):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid contact id"})
		case errors.Is(err, service.ErrContactNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "contact not found"})
		case errors.Is(err, service.ErrInvalidContactType):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid contact type"})
		case errors.Is(err, service.ErrEmptyContactValue):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "contact value is required"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update contact"})
		}
	}

	return c.JSON(toContactResponse(contact))
}

func (h *ContactHandler) DeleteContact(c *fiber.Ctx) error {
	var id pgtype.UUID
	if err := id.Scan(c.Params("contact_id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid contact id"})
	}

	err := h.service.DeleteContact(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrContactNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "contact not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete contact"})
	}

	return c.JSON(fiber.Map{"message": "contact deleted successfully"})
}

// ----------------- Contact Persons -----------------

func (h *ContactHandler) CreateContactPerson(c *fiber.Ctx) error {
	partyIDStr := c.Params("party_id")
	var partyID pgtype.UUID
	if err := partyID.Scan(partyIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid party id"})
	}

	var req dto.CreateContactPersonRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	input := domain.CreateContactPersonInput{
		PartyID:     partyIDStr,
		FirstName:   req.FirstName,
		MiddleName:  req.MiddleName,
		LastName:    req.LastName,
		Designation: req.Designation,
		Department:  req.Department,
		Email:       req.Email,
		Phone:       req.Phone,
		Mobile:      req.Mobile,
		IsPrimary:   req.IsPrimary,
		IsActive:    req.IsActive,
		Notes:       req.Notes,
	}

	person, err := h.service.CreateContactPerson(c.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPartyID):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid party id"})
		case errors.Is(err, service.ErrPartyNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "party not found"})
		case errors.Is(err, service.ErrEmptyFirstName):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "first name is required"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create contact person"})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(toContactPersonResponse(person))
}

func (h *ContactHandler) GetContactPersonByID(c *fiber.Ctx) error {
	var id pgtype.UUID
	if err := id.Scan(c.Params("person_id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid person id"})
	}

	person, err := h.service.GetContactPersonByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrContactPersonNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "contact person not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get contact person"})
	}

	return c.JSON(toContactPersonResponse(person))
}

func (h *ContactHandler) ListContactPersonsByPartyID(c *fiber.Ctx) error {
	var partyID pgtype.UUID
	if err := partyID.Scan(c.Params("party_id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid party id"})
	}

	persons, err := h.service.ListContactPersonsByPartyID(c.Context(), partyID)
	if err != nil {
		if errors.Is(err, service.ErrPartyNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "party not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list contact persons"})
	}

	responses := make([]dto.ContactPersonResponse, 0, len(persons))
	for _, p := range persons {
		responses = append(responses, toContactPersonResponse(p))
	}

	return c.JSON(responses)
}

func (h *ContactHandler) UpdateContactPerson(c *fiber.Ctx) error {
	idStr := c.Params("person_id")
	var id pgtype.UUID
	if err := id.Scan(idStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid person id"})
	}

	var req dto.UpdateContactPersonRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	input := domain.UpdateContactPersonInput{
		ID:          idStr,
		FirstName:   req.FirstName,
		MiddleName:  req.MiddleName,
		LastName:    req.LastName,
		Designation: req.Designation,
		Department:  req.Department,
		Email:       req.Email,
		Phone:       req.Phone,
		Mobile:      req.Mobile,
		IsPrimary:   req.IsPrimary,
		IsActive:    req.IsActive,
		Notes:       req.Notes,
	}

	person, err := h.service.UpdateContactPerson(c.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidContactPersonID):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid person id"})
		case errors.Is(err, service.ErrContactPersonNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "contact person not found"})
		case errors.Is(err, service.ErrEmptyFirstName):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "first name is required"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update contact person"})
		}
	}

	return c.JSON(toContactPersonResponse(person))
}

func (h *ContactHandler) DeleteContactPerson(c *fiber.Ctx) error {
	var id pgtype.UUID
	if err := id.Scan(c.Params("person_id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid person id"})
	}

	err := h.service.DeleteContactPerson(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrContactPersonNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "contact person not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete contact person"})
	}

	return c.JSON(fiber.Map{"message": "contact person deleted successfully"})
}

func toContactResponse(c db.PartyContact) dto.ContactResponse {
	return dto.ContactResponse{
		ID:           c.ID.String(),
		PartyID:      c.PartyID.String(),
		ContactType:  c.ContactType,
		ContactValue: c.ContactValue,
		Label:        c.Label,
		IsPrimary:    c.IsPrimary,
		IsActive:     c.IsActive,
		CreatedAt:    c.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:    c.UpdatedAt.Time.Format(time.RFC3339),
	}
}

func toContactPersonResponse(p db.PartyContactPerson) dto.ContactPersonResponse {
	return dto.ContactPersonResponse{
		ID:          p.ID.String(),
		PartyID:     p.PartyID.String(),
		FirstName:   p.FirstName,
		MiddleName:  p.MiddleName,
		LastName:    p.LastName,
		Designation: p.Designation,
		Department:  p.Department,
		Email:       p.Email,
		Phone:       p.Phone,
		Mobile:      p.Mobile,
		IsPrimary:   p.IsPrimary,
		IsActive:    p.IsActive,
		Notes:       p.Notes,
		CreatedAt:   p.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.Time.Format(time.RFC3339),
	}
}
