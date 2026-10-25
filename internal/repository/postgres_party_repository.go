package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
)

type postgresPartyRepository struct {
	queries *db.Queries
}

func NewPostgresPartyRepository(queries *db.Queries) PartyRepository {
	return &postgresPartyRepository{
		queries: queries,
	}
}

func (r *postgresPartyRepository) GetPartyByID(
	ctx context.Context,
	id pgtype.UUID,
) (db.Party, error) {
	return r.queries.GetPartyByID(ctx, id)
}

func (r *postgresPartyRepository) ListParties(
	ctx context.Context,
	params db.ListPartiesParams,
) ([]db.Party, error) {
	return r.queries.ListParties(ctx, params)
}

func (r *postgresPartyRepository) CreateParty(
	ctx context.Context,
	params db.CreatePartyParams,
) (db.Party, error) {
	return r.queries.CreateParty(ctx, params)
}

func (r *postgresPartyRepository) UpdateParty(
	ctx context.Context,
	params db.UpdatePartyParams,
) (db.Party, error) {
	return r.queries.UpdateParty(ctx, params)
}

func (r *postgresPartyRepository) UpdatePartyStatus(
	ctx context.Context,
	params db.UpdatePartyStatusParams,
) (db.Party, error) {
	return r.queries.UpdatePartyStatus(ctx, params)
}
