package api

import (
	"github.com/gofiber/fiber/v2"
	httpHandler "github.com/group-project/crm_party_service/internal/handler/http"
)

// registerAddressRoutes registers party address management endpoints.
func registerAddressRoutes(router fiber.Router, addrHandler *httpHandler.AddressHandler) {
	// Party-scoped Address Endpoints
	router.Post("/parties/:party_id/addresses", addrHandler.CreateAddress)
	router.Get("/parties/:party_id/addresses", addrHandler.ListAddressesByPartyID)

	// Direct Address Resource Endpoints
	router.Get("/addresses/:address_id", addrHandler.GetAddressByID)
	router.Put("/addresses/:address_id", addrHandler.UpdateAddress)
	router.Delete("/addresses/:address_id", addrHandler.DeleteAddress)
}
