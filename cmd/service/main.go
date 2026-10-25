package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/group-project/crm_party_service/internal/api"
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
		AppName:      "CRM Party Service",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	// Middlewares
	api.SetupMiddlewares(app)

	// Register Categorized API Routes
	api.RegisterRoutes(app, api.Handlers{
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

	// Channel to catch server startup / runtime errors
	serverErrors := make(chan error, 1)

	// Start server in background goroutine
	go func() {
		log.Printf("party service running on port %s (env: %s)", cfg.Port, cfg.Environment)
		if err := app.Listen(cfg.Port); err != nil {
			serverErrors <- err
		}
	}()

	// Channel to listen for OS interrupt and termination signals
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Block until either a fatal server error or an OS termination signal
	select {
	case err := <-serverErrors:
		log.Fatalf("server terminated unexpectedly: %v", err)

	case sig := <-shutdownChan:
		log.Printf("received signal '%s', initiating graceful shutdown...", sig)

		// 1. Gracefully shutdown Fiber HTTP server (draining in-flight requests)
		if err := app.ShutdownWithTimeout(cfg.ShutdownTimeout); err != nil {
			log.Printf("error during server shutdown: %v", err)
		} else {
			log.Println("HTTP server closed gracefully")
		}

		// 2. Close PostgreSQL connection pool cleanly
		log.Println("closing database connection pool...")
		pool.Close()
		log.Println("database pool closed successfully")

		log.Println("graceful shutdown completed successfully")
	}
}
