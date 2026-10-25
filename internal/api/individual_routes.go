package api

import (
	"github.com/gofiber/fiber/v2"
	httpHandler "github.com/group-project/crm_party_service/internal/handler/http"
)

// registerIndividualRoutes registers B2C individual profile endpoints.
func registerIndividualRoutes(router fiber.Router, indivHandler *httpHandler.IndividualHandler) {
	indivs := router.Group("/individuals")

	indivs.Post("", indivHandler.CreateIndividual)
	indivs.Get("/:party_id", indivHandler.GetIndividualByPartyID)
	indivs.Put("/:party_id", indivHandler.UpdateIndividual)
}
