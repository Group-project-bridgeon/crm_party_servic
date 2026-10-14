package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
)

type postgresOrganizationRepository struct {
	queries *db.Queries
}

func NewPostgresOrganizationRepository(
	queries *db.Queries,
) OrganizationRepository {
	return &postgresOrganizationRepository{
		queries: queries,
	}
}

func (r *postgresOrganizationRepository) CreateOrganization(
	ctx context.Context,
	params db.CreateOrganizationParams,
) (db.Organization, error) {
	return r.queries.CreateOrganization(ctx, params)
}

func (r *postgresOrganizationRepository) GetOrganizationByPartyID(
	ctx context.Context,
	partyID pgtype.UUID,
) (db.Organization, error) {
	return r.queries.GetOrganizationByPartyID(ctx, partyID)
}
