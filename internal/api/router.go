package api

import (
	"github.com/gofiber/fiber/v2"
	httpHandler "github.com/group-project/crm_party_service/internal/handler/http"
)

// Handlers encapsulates all domain HTTP handlers for route registration.
type Handlers struct {
	Party        *httpHandler.PartyHandler
	Organization *httpHandler.OrganizationHandler
	Individual   *httpHandler.IndividualHandler
	Customer     *httpHandler.CustomerHandler
	Supplier     *httpHandler.SupplierHandler
	Address      *httpHandler.AddressHandler
	Contact      *httpHandler.ContactHandler
	Note         *httpHandler.NoteHandler
	Composite    *httpHandler.CompositeHandler
}

// RegisterRoutes registers all categorized API routes on the Fiber application.
func RegisterRoutes(app *fiber.App, h Handlers) {
	// 1. Health & system status routes
	registerHealthRoutes(app)

	// 2. Base API v1 group
	v1 := app.Group("/api/v1")

	// 3. Categorized domain routes
	registerPartyRoutes(v1, h.Party, h.Composite)
	registerOrganizationRoutes(v1, h.Organization)
	registerIndividualRoutes(v1, h.Individual)
	registerCustomerRoutes(v1, h.Customer, h.Composite)
	registerSupplierRoutes(v1, h.Supplier, h.Composite)
	registerAddressRoutes(v1, h.Address)
	registerContactRoutes(v1, h.Contact)
	registerNoteRoutes(v1, h.Note)
}
