package api

import (
	"github.com/gofiber/fiber/v2"
	httpHandler "github.com/group-project/crm_party_service/internal/handler/http"
)

// registerContactRoutes registers party contact channels and contact persons (stakeholders).
func registerContactRoutes(router fiber.Router, contactHandler *httpHandler.ContactHandler) {
	// 1. Party Contact Channels (Email, Phone, Mobile, WhatsApp, etc.)
	router.Post("/parties/:party_id/contacts", contactHandler.CreateContact)
	router.Get("/parties/:party_id/contacts", contactHandler.ListContactsByPartyID)
	router.Get("/contacts/:contact_id", contactHandler.GetContactByID)
	router.Put("/contacts/:contact_id", contactHandler.UpdateContact)
	router.Delete("/contacts/:contact_id", contactHandler.DeleteContact)

	// 2. Party Contact Persons (Key Account Managers / Executive POCs)
	router.Post("/parties/:party_id/contact-persons", contactHandler.CreateContactPerson)
	router.Get("/parties/:party_id/contact-persons", contactHandler.ListContactPersonsByPartyID)
	router.Get("/contact-persons/:person_id", contactHandler.GetContactPersonByID)
	router.Put("/contact-persons/:person_id", contactHandler.UpdateContactPerson)
	router.Delete("/contact-persons/:person_id", contactHandler.DeleteContactPerson)
}
