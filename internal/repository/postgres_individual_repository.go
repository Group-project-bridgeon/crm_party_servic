package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
)

type postgresIndividualRepository struct {
	queries *db.Queries
}

func NewPostgresIndividualRepository(
	queries *db.Queries,
) IndividualRepository {
	return &postgresIndividualRepository{
		queries: queries,
	}
}

func (r *postgresIndividualRepository) CreateIndividual(
	ctx context.Context,
	params db.CreateIndividualParams,
) (db.Individual, error) {
	return r.queries.CreateIndividual(ctx, params)
}

func (r *postgresIndividualRepository) GetIndividualByPartyID(
	ctx context.Context,
	partyID pgtype.UUID,
) (db.Individual, error) {
	return r.queries.GetIndividualByPartyID(ctx, partyID)
}

func (r *postgresIndividualRepository) UpdateIndividual(
	ctx context.Context,
	params db.UpdateIndividualParams,
) (db.Individual, error) {
	return r.queries.UpdateIndividual(ctx, params)
}
