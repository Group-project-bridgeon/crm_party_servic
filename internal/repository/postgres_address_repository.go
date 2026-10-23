package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
)

type postgresAddressRepository struct {
	queries *db.Queries
}

func NewPostgresAddressRepository(queries *db.Queries) AddressRepository {
	return &postgresAddressRepository{queries: queries}
}

func (r *postgresAddressRepository) CreatePartyAddress(
	ctx context.Context,
	params db.CreatePartyAddressParams,
) (db.PartyAddress, error) {
	return r.queries.CreatePartyAddress(ctx, params)
}

func (r *postgresAddressRepository) GetPartyAddressByID(
	ctx context.Context,
	id pgtype.UUID,
) (db.PartyAddress, error) {
	return r.queries.GetPartyAddressByID(ctx, id)
}

func (r *postgresAddressRepository) ListPartyAddressesByPartyID(
	ctx context.Context,
	params db.ListPartyAddressesByPartyIDParams,
) ([]db.PartyAddress, error) {
	return r.queries.ListPartyAddressesByPartyID(ctx, params)
}

func (r *postgresAddressRepository) UpdatePartyAddress(
	ctx context.Context,
	params db.UpdatePartyAddressParams,
) (db.PartyAddress, error) {
	return r.queries.UpdatePartyAddress(ctx, params)
}

func (r *postgresAddressRepository) ResetPrimaryPartyAddress(
	ctx context.Context,
	params db.ResetPrimaryPartyAddressParams,
) error {
	return r.queries.ResetPrimaryPartyAddress(ctx, params)
}

func (r *postgresAddressRepository) DeletePartyAddress(
	ctx context.Context,
	id pgtype.UUID,
) error {
	return r.queries.DeletePartyAddress(ctx, id)
}
