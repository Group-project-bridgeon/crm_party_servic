package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
)

type IndividualRepository interface {
	CreateIndividual(
		ctx context.Context,
		params db.CreateIndividualParams,
	) (db.Individual, error)

	GetIndividualByPartyID(
		ctx context.Context,
		partyID pgtype.UUID,
	) (db.Individual, error)

	UpdateIndividual(
		ctx context.Context,
		params db.UpdateIndividualParams,
	) (db.Individual, error)
}
