package api

import (
	"github.com/gofiber/fiber/v2"
	httpHandler "github.com/group-project/crm_party_service/internal/handler/http"
)

// registerPartyRoutes registers all Party lifecycle & 360 overview endpoints.
func registerPartyRoutes(router fiber.Router, partyHandler *httpHandler.PartyHandler, compositeHandler *httpHandler.CompositeHandler) {
	parties := router.Group("/parties")

	parties.Post("", partyHandler.CreateParty)
	parties.Get("", partyHandler.ListParties)
	parties.Get("/:id", partyHandler.GetPartyByID)
	parties.Put("/:id", partyHandler.UpdateParty)
	parties.Patch("/:id/status", partyHandler.UpdatePartyStatus)
	parties.Delete("/:id", partyHandler.DeleteParty)

	// Party 360 Unified Overview
	parties.Get("/:id/360", compositeHandler.GetParty360)
}
