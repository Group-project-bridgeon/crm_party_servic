package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/group-project/crm_party_service/internal/db"
	httpHandler "github.com/group-project/crm_party_service/internal/handler/http"
	"github.com/group-project/crm_party_service/internal/repository"
	"github.com/group-project/crm_party_service/internal/service"
)

func main() {
	ctx := context.Background()

	databaseURL := "postgres://postgres:postgres@localhost:5432/party?sslmode=disable"

	pool, err := pgxpool.New(ctx, databaseURL)
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

	app := fiber.New()

	app.Post("/api/v1/parties", partyHandler.CreateParty)
	app.Get("/api/v1/parties", partyHandler.ListParties)
	app.Get("/api/v1/parties/:id", partyHandler.GetPartyByID)

	app.Post("/api/v1/organizations", organizationHandler.CreateOrganization)

	app.Get(
    "/api/v1/organizations/:party_id",
    organizationHandler.GetOrganizationByPartyID,
)

	log.Println("party service running on :8080")

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
