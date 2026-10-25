package main

import (
	"context"
	"log/slog"
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
	"github.com/group-project/crm_party_service/internal/logger"
	"github.com/group-project/crm_party_service/internal/repository"
	"github.com/group-project/crm_party_service/internal/service"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	// Initialize structured slog logger
	logger.Init(logger.Config{
		Environment: cfg.Environment,
		Level:       cfg.LogLevel,
		ServiceName: "crm_party_service",
	})

	slog.Info("starting crm party service",
		"env", cfg.Environment,
		"port", cfg.Port,
		"log_level", cfg.LogLevel,
	)

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to create database pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
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
		slog.Info("party service http server listening", "port", cfg.Port, "env", cfg.Environment)
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
		slog.Error("server terminated unexpectedly", "error", err)
		os.Exit(1)

	case sig := <-shutdownChan:
		slog.Info("received termination signal, initiating graceful shutdown", "signal", sig.String())

		// 1. Gracefully shutdown Fiber HTTP server (draining in-flight requests)
		if err := app.ShutdownWithTimeout(cfg.ShutdownTimeout); err != nil {
			slog.Error("error during http server shutdown", "error", err)
		} else {
			slog.Info("http server closed gracefully")
		}

		// 2. Close PostgreSQL connection pool cleanly
		slog.Info("closing database connection pool")
		pool.Close()
		slog.Info("database pool closed successfully")

		slog.Info("graceful shutdown completed successfully")
	}
}
