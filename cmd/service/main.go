package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/group-project/crm_party_service/internal/config"
	"github.com/group-project/crm_party_service/internal/db"
	httpHandler "github.com/group-project/crm_party_service/internal/handler/http"
	"github.com/group-project/crm_party_service/internal/repository"
	"github.com/group-project/crm_party_service/internal/service"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to create database pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	queries := db.New(pool)

	// Party
	partyRepository := repository.NewPostgresPartyRepository(queries)
	partyService := service.NewPartyService(partyRepository)
	partyHandler := httpHandler.NewPartyHandler(partyService)

	// Organization
	organizationRepository := repository.NewPostgresOrganizationRepository(queries)
	organizationService := service.NewOrganizationService(
		partyRepository,
		organizationRepository,
	)
	organizationHandler := httpHandler.NewOrganizationHandler(organizationService)

	// Individual
	individualRepository := repository.NewPostgresIndividualRepository(queries)
	individualService := service.NewIndividualService(
		partyRepository,
		individualRepository,
	)
	individualHandler := httpHandler.NewIndividualHandler(individualService)

	// Customer
	customerRepository := repository.NewPostgresCustomerRepository(queries)
	customerService := service.NewCustomerService(
		partyRepository,
		customerRepository,
	)
	customerHandler := httpHandler.NewCustomerHandler(customerService)

	// Supplier
	supplierRepository := repository.NewPostgresSupplierRepository(queries)
	supplierService := service.NewSupplierService(
		partyRepository,
		supplierRepository,
	)
	supplierHandler := httpHandler.NewSupplierHandler(supplierService)

	// Address
	addressRepository := repository.NewPostgresAddressRepository(queries)
	addressService := service.NewAddressService(partyRepository, addressRepository)
	addressHandler := httpHandler.NewAddressHandler(addressService)

	// Contact & Contact Person
	contactRepository := repository.NewPostgresContactRepository(queries)
	contactService := service.NewContactService(partyRepository, contactRepository)
	contactHandler := httpHandler.NewContactHandler(contactService)

	// Note
	noteRepository := repository.NewPostgresNoteRepository(queries)
	noteService := service.NewNoteService(partyRepository, noteRepository)
	noteHandler := httpHandler.NewNoteHandler(noteService)

	// Composite & Transaction Manager
	txManager := repository.NewPostgresTransactionManager(pool)
	compositeService := service.NewCompositeService(
		txManager,
		partyRepository,
		organizationRepository,
		individualRepository,
		customerRepository,
		supplierRepository,
		addressRepository,
		contactRepository,
		noteRepository,
	)
	compositeHandler := httpHandler.NewCompositeHandler(compositeService)

	app := fiber.New(fiber.Config{
		AppName: "CRM Party Service",
	})




	// Middlewares
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "UP",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Routes: Parties
	app.Post("/api/v1/parties", partyHandler.CreateParty)
	app.Get("/api/v1/parties", partyHandler.ListParties)
	app.Get("/api/v1/parties/:id", partyHandler.GetPartyByID)
	app.Put("/api/v1/parties/:id", partyHandler.UpdateParty)
	app.Patch("/api/v1/parties/:id/status", partyHandler.UpdatePartyStatus)

	// Routes: Organizations
	app.Post("/api/v1/organizations", organizationHandler.CreateOrganization)
	app.Get("/api/v1/organizations/:party_id", organizationHandler.GetOrganizationByPartyID)
	app.Put("/api/v1/organizations/:party_id", organizationHandler.UpdateOrganization)

	// Routes: Individuals
	app.Post("/api/v1/individuals", individualHandler.CreateIndividual)
	app.Get("/api/v1/individuals/:party_id", individualHandler.GetIndividualByPartyID)
	app.Put("/api/v1/individuals/:party_id", individualHandler.UpdateIndividual)

	// Routes: Customers
	app.Post("/api/v1/customers", customerHandler.CreateCustomer)
	app.Get("/api/v1/customers", customerHandler.ListCustomers)
	app.Get("/api/v1/customers/:id", customerHandler.GetCustomerByID)
	app.Get("/api/v1/customers/party/:party_id", customerHandler.GetCustomerByPartyID)
	app.Put("/api/v1/customers/:id", customerHandler.UpdateCustomer)
	app.Post("/api/v1/customers/:id/user-accounts", customerHandler.CreateCustomerUserAccount)
	app.Get("/api/v1/customers/:id/user-accounts", customerHandler.ListCustomerUserAccounts)
	app.Patch("/api/v1/customers/user-accounts/:account_id/status", customerHandler.UpdateCustomerUserAccountStatus)

	// Routes: Suppliers
	app.Post("/api/v1/suppliers", supplierHandler.CreateSupplier)
	app.Get("/api/v1/suppliers", supplierHandler.ListSuppliers)
	app.Get("/api/v1/suppliers/:id", supplierHandler.GetSupplierByID)
	app.Get("/api/v1/suppliers/party/:party_id", supplierHandler.GetSupplierByPartyID)
	app.Put("/api/v1/suppliers/:id", supplierHandler.UpdateSupplier)
	app.Post("/api/v1/suppliers/:id/user-accounts", supplierHandler.CreateSupplierUserAccount)
	app.Get("/api/v1/suppliers/:id/user-accounts", supplierHandler.ListSupplierUserAccounts)
	app.Patch("/api/v1/suppliers/user-accounts/:account_id/status", supplierHandler.UpdateSupplierUserAccountStatus)

	// Routes: Addresses
	app.Post("/api/v1/parties/:party_id/addresses", addressHandler.CreateAddress)
	app.Get("/api/v1/parties/:party_id/addresses", addressHandler.ListAddressesByPartyID)
	app.Get("/api/v1/addresses/:address_id", addressHandler.GetAddressByID)
	app.Put("/api/v1/addresses/:address_id", addressHandler.UpdateAddress)
	app.Delete("/api/v1/addresses/:address_id", addressHandler.DeleteAddress)

	// Routes: Contacts
	app.Post("/api/v1/parties/:party_id/contacts", contactHandler.CreateContact)
	app.Get("/api/v1/parties/:party_id/contacts", contactHandler.ListContactsByPartyID)
	app.Get("/api/v1/contacts/:contact_id", contactHandler.GetContactByID)
	app.Put("/api/v1/contacts/:contact_id", contactHandler.UpdateContact)
	app.Delete("/api/v1/contacts/:contact_id", contactHandler.DeleteContact)

	// Routes: Contact Persons
	app.Post("/api/v1/parties/:party_id/contact-persons", contactHandler.CreateContactPerson)
	app.Get("/api/v1/parties/:party_id/contact-persons", contactHandler.ListContactPersonsByPartyID)
	app.Get("/api/v1/contact-persons/:person_id", contactHandler.GetContactPersonByID)
	app.Put("/api/v1/contact-persons/:person_id", contactHandler.UpdateContactPerson)
	app.Delete("/api/v1/contact-persons/:person_id", contactHandler.DeleteContactPerson)

	// Routes: Notes
	app.Post("/api/v1/parties/:party_id/notes", noteHandler.CreateNote)
	app.Get("/api/v1/parties/:party_id/notes", noteHandler.ListNotesByPartyID)
	app.Get("/api/v1/notes/:note_id", noteHandler.GetNoteByID)
	app.Delete("/api/v1/notes/:note_id", noteHandler.DeleteNote)

	// Routes: Composite & Party 360
	app.Get("/api/v1/parties/:id/360", compositeHandler.GetParty360)
	app.Post("/api/v1/customers/composite", compositeHandler.CreateCompositeCustomer)
	app.Post("/api/v1/suppliers/composite", compositeHandler.CreateCompositeSupplier)




	// Graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdownChan
		log.Println("gracefully shutting down server...")
		if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
			log.Printf("error during server shutdown: %v", err)
		}
	}()

	log.Printf("party service running on port %s", cfg.Port)
	if err := app.Listen(cfg.Port); err != nil && err != fiber.ErrBadGateway {
		log.Printf("server terminated: %v", err)
	}
}
