package http

import "github.com/gofiber/fiber/v2"

// Handlers encapsulates all HTTP handlers for route registration.
type Handlers struct {
	Party        *PartyHandler
	Organization *OrganizationHandler
	Individual   *IndividualHandler
	Customer     *CustomerHandler
	Supplier     *SupplierHandler
	Address      *AddressHandler
	Contact      *ContactHandler
	Note         *NoteHandler
	Composite    *CompositeHandler
}

// RegisterRoutes registers all API routes according to the Clean Architecture delivery layer.
func RegisterRoutes(app *fiber.App, h Handlers) {
	api := app.Group("/api/v1")

	// --- 1. Parties ---
	parties := api.Group("/parties")
	parties.Post("", h.Party.CreateParty)
	parties.Get("", h.Party.ListParties)
	parties.Get("/:id", h.Party.GetPartyByID)
	parties.Put("/:id", h.Party.UpdateParty)
	parties.Patch("/:id/status", h.Party.UpdatePartyStatus)
	parties.Delete("/:id", h.Party.DeleteParty)
	parties.Get("/:id/360", h.Composite.GetParty360)

	// --- 2. Organizations (B2B Profile) ---
	orgs := api.Group("/organizations")
	orgs.Post("", h.Organization.CreateOrganization)
	orgs.Get("/:party_id", h.Organization.GetOrganizationByPartyID)
	orgs.Put("/:party_id", h.Organization.UpdateOrganization)

	// --- 3. Individuals (B2C Profile) ---
	indivs := api.Group("/individuals")
	indivs.Post("", h.Individual.CreateIndividual)
	indivs.Get("/:party_id", h.Individual.GetIndividualByPartyID)
	indivs.Put("/:party_id", h.Individual.UpdateIndividual)

	// --- 4. Customers & Accounts ---
	customers := api.Group("/customers")
	customers.Post("", h.Customer.CreateCustomer)
	customers.Get("", h.Customer.ListCustomers)
	customers.Post("/composite", h.Composite.CreateCompositeCustomer)
	customers.Get("/:id", h.Customer.GetCustomerByID)
	customers.Get("/party/:party_id", h.Customer.GetCustomerByPartyID)
	customers.Put("/:id", h.Customer.UpdateCustomer)
	customers.Post("/:id/user-accounts", h.Customer.CreateCustomerUserAccount)
	customers.Get("/:id/user-accounts", h.Customer.ListCustomerUserAccounts)
	customers.Patch("/user-accounts/:account_id/status", h.Customer.UpdateCustomerUserAccountStatus)

	// --- 5. Suppliers & Accounts ---
	suppliers := api.Group("/suppliers")
	suppliers.Post("", h.Supplier.CreateSupplier)
	suppliers.Get("", h.Supplier.ListSuppliers)
	suppliers.Post("/composite", h.Composite.CreateCompositeSupplier)
	suppliers.Get("/:id", h.Supplier.GetSupplierByID)
	suppliers.Get("/party/:party_id", h.Supplier.GetSupplierByPartyID)
	suppliers.Put("/:id", h.Supplier.UpdateSupplier)
	suppliers.Post("/:id/user-accounts", h.Supplier.CreateSupplierUserAccount)
	suppliers.Get("/:id/user-accounts", h.Supplier.ListSupplierUserAccounts)
	suppliers.Patch("/user-accounts/:account_id/status", h.Supplier.UpdateSupplierUserAccountStatus)

	// --- 6. Addresses (Multi-address & Primary reset) ---
	api.Post("/parties/:party_id/addresses", h.Address.CreateAddress)
	api.Get("/parties/:party_id/addresses", h.Address.ListAddressesByPartyID)
	api.Get("/addresses/:address_id", h.Address.GetAddressByID)
	api.Put("/addresses/:address_id", h.Address.UpdateAddress)
	api.Delete("/addresses/:address_id", h.Address.DeleteAddress)

	// --- 7. Contacts (Multi-channel & Primary reset) ---
	api.Post("/parties/:party_id/contacts", h.Contact.CreateContact)
	api.Get("/parties/:party_id/contacts", h.Contact.ListContactsByPartyID)
	api.Get("/contacts/:contact_id", h.Contact.GetContactByID)
	api.Put("/contacts/:contact_id", h.Contact.UpdateContact)
	api.Delete("/contacts/:contact_id", h.Contact.DeleteContact)

	// --- 8. Contact Persons (Key Account Managers / Stakeholders) ---
	api.Post("/parties/:party_id/contact-persons", h.Contact.CreateContactPerson)
	api.Get("/parties/:party_id/contact-persons", h.Contact.ListContactPersonsByPartyID)
	api.Get("/contact-persons/:person_id", h.Contact.GetContactPersonByID)
	api.Put("/contact-persons/:person_id", h.Contact.UpdateContactPerson)
	api.Delete("/contact-persons/:person_id", h.Contact.DeleteContactPerson)

	// --- 9. CRM Notes ---
	api.Post("/parties/:party_id/notes", h.Note.CreateNote)
	api.Get("/parties/:party_id/notes", h.Note.ListNotesByPartyID)
	api.Get("/notes/:note_id", h.Note.GetNoteByID)
	api.Delete("/notes/:note_id", h.Note.DeleteNote)
}
