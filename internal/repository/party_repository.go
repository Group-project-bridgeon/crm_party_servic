package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
)

type PartyRepository interface {
	GetPartyByID(ctx context.Context, id pgtype.UUID) (db.Party, error)

	ListParties(ctx context.Context, params db.ListPartiesParams) ([]db.Party, error)

	CreateParty(ctx context.Context, params db.CreatePartyParams) (db.Party, error)
}
