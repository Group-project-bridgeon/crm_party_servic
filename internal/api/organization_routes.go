package api

import (
	"github.com/gofiber/fiber/v2"
	httpHandler "github.com/group-project/crm_party_service/internal/handler/http"
)

// registerOrganizationRoutes registers B2B organization profile endpoints.
func registerOrganizationRoutes(router fiber.Router, orgHandler *httpHandler.OrganizationHandler) {
	orgs := router.Group("/organizations")

	orgs.Post("", orgHandler.CreateOrganization)
	orgs.Get("/:party_id", orgHandler.GetOrganizationByPartyID)
	orgs.Put("/:party_id", orgHandler.UpdateOrganization)
}
