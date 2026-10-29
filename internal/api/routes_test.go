package api_test

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/group-project/crm_party_service/internal/api"
	httpHandler "github.com/group-project/crm_party_service/internal/handler/http"
)

func TestRegisterRoutes_AllCategoriesMounted(t *testing.T) {
	app := fiber.New()

	api.SetupMiddlewares(app)

	// Register with empty/nil handlers to verify router path registration
	api.RegisterRoutes(app, api.Handlers{
		Party:        httpHandler.NewPartyHandler(nil),
		Organization: httpHandler.NewOrganizationHandler(nil),
		Individual:   httpHandler.NewIndividualHandler(nil),
		Customer:     httpHandler.NewCustomerHandler(nil),
		Supplier:     httpHandler.NewSupplierHandler(nil),
		Address:      httpHandler.NewAddressHandler(nil),
		Contact:      httpHandler.NewContactHandler(nil),
		Note:         httpHandler.NewNoteHandler(nil),
		Composite:    httpHandler.NewCompositeHandler(nil),
	})

	expectedPaths := []string{
		"/health",
		"/metrics",
		"/api/v1/parties",
		"/api/v1/parties/:id",
		"/api/v1/parties/:id/status",
		"/api/v1/parties/:id/360",
		"/api/v1/organizations",
		"/api/v1/organizations/:party_id",
		"/api/v1/individuals",
		"/api/v1/individuals/:party_id",
		"/api/v1/customers",
		"/api/v1/customers/composite",
		"/api/v1/customers/:id",
		"/api/v1/customers/party/:party_id",
		"/api/v1/customers/:id/user-accounts",
		"/api/v1/customers/user-accounts/:account_id/status",
		"/api/v1/suppliers",
		"/api/v1/suppliers/composite",
		"/api/v1/suppliers/:id",
		"/api/v1/suppliers/party/:party_id",
		"/api/v1/suppliers/:id/user-accounts",
		"/api/v1/suppliers/user-accounts/:account_id/status",
		"/api/v1/parties/:party_id/addresses",
		"/api/v1/addresses/:address_id",
		"/api/v1/parties/:party_id/contacts",
		"/api/v1/contacts/:contact_id",
		"/api/v1/parties/:party_id/contact-persons",
		"/api/v1/contact-persons/:person_id",
		"/api/v1/parties/:party_id/notes",
		"/api/v1/notes/:note_id",
	}

	routes := app.GetRoutes()
	registeredPaths := make(map[string]bool)
	for _, r := range routes {
		registeredPaths[r.Path] = true
	}

	for _, expected := range expectedPaths {
		if !registeredPaths[expected] {
			t.Errorf("expected route path %q was not registered", expected)
		}
	}
}
