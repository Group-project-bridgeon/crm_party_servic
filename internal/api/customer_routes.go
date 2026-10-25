package api

import (
	"github.com/gofiber/fiber/v2"
	httpHandler "github.com/group-project/crm_party_service/internal/handler/http"
)

// registerCustomerRoutes registers customer, user account, and composite customer onboarding endpoints.
func registerCustomerRoutes(router fiber.Router, custHandler *httpHandler.CustomerHandler, compositeHandler *httpHandler.CompositeHandler) {
	customers := router.Group("/customers")

	// Base Customer CRUD & Search
	customers.Post("", custHandler.CreateCustomer)
	customers.Get("", custHandler.ListCustomers)
	customers.Get("/:id", custHandler.GetCustomerByID)
	customers.Get("/party/:party_id", custHandler.GetCustomerByPartyID)
	customers.Put("/:id", custHandler.UpdateCustomer)

	// Composite Customer Atomic Onboarding
	customers.Post("/composite", compositeHandler.CreateCompositeCustomer)

	// Customer Portal User Accounts
	customers.Post("/:id/user-accounts", custHandler.CreateCustomerUserAccount)
	customers.Get("/:id/user-accounts", custHandler.ListCustomerUserAccounts)
	customers.Patch("/user-accounts/:account_id/status", custHandler.UpdateCustomerUserAccountStatus)
}
