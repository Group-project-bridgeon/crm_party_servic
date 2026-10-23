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

	// Register Clean Architecture routes
	httpHandler.RegisterRoutes(app, httpHandler.Handlers{
		Party:        partyHandler,
		Organization: organizationHandler,
		Individual:   individualHandler,
		Customer:     customerHandler,
		Supplier:     supplierHandler,
		Address:      addressHandler,
		Contact:      contactHandler,
		Note:         noteHandler,
		Composite:    compositeHandler,
	})






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
