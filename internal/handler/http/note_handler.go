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

type NoteService interface {
	CreateNote(ctx context.Context, input domain.CreateNoteInput) (db.PartyNote, error)
	GetNoteByID(ctx context.Context, id pgtype.UUID) (db.PartyNote, error)
	ListNotesByPartyID(ctx context.Context, partyID pgtype.UUID) ([]db.PartyNote, error)
	DeleteNote(ctx context.Context, id pgtype.UUID) error
}

type NoteHandler struct {
	service NoteService
}

func NewNoteHandler(service NoteService) *NoteHandler {
	return &NoteHandler{service: service}
}

func (h *NoteHandler) CreateNote(c *fiber.Ctx) error {
	partyIDStr := c.Params("party_id")
	var partyID pgtype.UUID
	if err := partyID.Scan(partyIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid party id"})
	}

	var req dto.CreateNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	input := domain.CreateNoteInput{
		PartyID:   partyIDStr,
		NoteType:  req.NoteType,
		Content:   req.Content,
		CreatedBy: req.CreatedBy,
	}

	note, err := h.service.CreateNote(c.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPartyID):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid party id"})
		case errors.Is(err, service.ErrPartyNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "party not found"})
		case errors.Is(err, service.ErrEmptyNoteContent):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "content cannot be empty"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create note"})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(toNoteResponse(note))
}

func (h *NoteHandler) GetNoteByID(c *fiber.Ctx) error {
	var id pgtype.UUID
	if err := id.Scan(c.Params("note_id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid note id"})
	}

	note, err := h.service.GetNoteByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNoteNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "note not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get note"})
	}

	return c.JSON(toNoteResponse(note))
}

func (h *NoteHandler) ListNotesByPartyID(c *fiber.Ctx) error {
	var partyID pgtype.UUID
	if err := partyID.Scan(c.Params("party_id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid party id"})
	}

	notes, err := h.service.ListNotesByPartyID(c.Context(), partyID)
	if err != nil {
		if errors.Is(err, service.ErrPartyNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "party not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list notes"})
	}

	responses := make([]dto.NoteResponse, 0, len(notes))
	for _, n := range notes {
		responses = append(responses, toNoteResponse(n))
	}

	return c.JSON(responses)
}

func (h *NoteHandler) DeleteNote(c *fiber.Ctx) error {
	var id pgtype.UUID
	if err := id.Scan(c.Params("note_id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid note id"})
	}

	err := h.service.DeleteNote(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNoteNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "note not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete note"})
	}

	return c.JSON(fiber.Map{"message": "note deleted successfully"})
}

func toNoteResponse(n db.PartyNote) dto.NoteResponse {
	var createdBy *string
	if n.CreatedBy.Valid {
		val := n.CreatedBy.String()
		createdBy = &val
	}

	return dto.NoteResponse{
		ID:        n.ID.String(),
		PartyID:   n.PartyID.String(),
		NoteType:  n.NoteType,
		Content:   n.Content,
		CreatedBy: createdBy,
		CreatedAt: n.CreatedAt.Time.Format(time.RFC3339),
	}
}
