package api

import (
	"github.com/gofiber/fiber/v2"
	httpHandler "github.com/group-project/crm_party_service/internal/handler/http"
)

// registerNoteRoutes registers CRM interaction history and note endpoints.
func registerNoteRoutes(router fiber.Router, noteHandler *httpHandler.NoteHandler) {
	router.Post("/parties/:party_id/notes", noteHandler.CreateNote)
	router.Get("/parties/:party_id/notes", noteHandler.ListNotesByPartyID)
	router.Get("/notes/:note_id", noteHandler.GetNoteByID)
	router.Delete("/notes/:note_id", noteHandler.DeleteNote)
}
