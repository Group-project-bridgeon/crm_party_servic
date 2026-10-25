package api

import (
	"github.com/gofiber/fiber/v2"
	httpHandler "github.com/group-project/crm_party_service/internal/handler/http"
)

// registerSupplierRoutes registers supplier, vendor account, and composite supplier onboarding endpoints.
func registerSupplierRoutes(router fiber.Router, suppHandler *httpHandler.SupplierHandler, compositeHandler *httpHandler.CompositeHandler) {
	suppliers := router.Group("/suppliers")

	// Base Supplier CRUD & Search
	suppliers.Post("", suppHandler.CreateSupplier)
	suppliers.Get("", suppHandler.ListSuppliers)
	suppliers.Get("/:id", suppHandler.GetSupplierByID)
	suppliers.Get("/party/:party_id", suppHandler.GetSupplierByPartyID)
	suppliers.Put("/:id", suppHandler.UpdateSupplier)

	// Composite Supplier Atomic Onboarding
	suppliers.Post("/composite", compositeHandler.CreateCompositeSupplier)

	// Vendor Portal User Accounts
	suppliers.Post("/:id/user-accounts", suppHandler.CreateSupplierUserAccount)
	suppliers.Get("/:id/user-accounts", suppHandler.ListSupplierUserAccounts)
	suppliers.Patch("/user-accounts/:account_id/status", suppHandler.UpdateSupplierUserAccountStatus)
}
