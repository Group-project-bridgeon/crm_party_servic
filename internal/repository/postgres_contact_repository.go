package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
)

type postgresContactRepository struct {
	queries *db.Queries
}

func NewPostgresContactRepository(queries *db.Queries) ContactRepository {
	return &postgresContactRepository{queries: queries}
}

// Direct contacts
func (r *postgresContactRepository) CreatePartyContact(
	ctx context.Context,
	params db.CreatePartyContactParams,
) (db.PartyContact, error) {
	return r.queries.CreatePartyContact(ctx, params)
}

func (r *postgresContactRepository) GetPartyContactByID(
	ctx context.Context,
	id pgtype.UUID,
) (db.PartyContact, error) {
	return r.queries.GetPartyContactByID(ctx, id)
}

func (r *postgresContactRepository) ListPartyContactsByPartyID(
	ctx context.Context,
	params db.ListPartyContactsByPartyIDParams,
) ([]db.PartyContact, error) {
	return r.queries.ListPartyContactsByPartyID(ctx, params)
}

func (r *postgresContactRepository) UpdatePartyContact(
	ctx context.Context,
	params db.UpdatePartyContactParams,
) (db.PartyContact, error) {
	return r.queries.UpdatePartyContact(ctx, params)
}

func (r *postgresContactRepository) ResetPrimaryPartyContact(
	ctx context.Context,
	params db.ResetPrimaryPartyContactParams,
) error {
	return r.queries.ResetPrimaryPartyContact(ctx, params)
}

func (r *postgresContactRepository) DeletePartyContact(
	ctx context.Context,
	id pgtype.UUID,
) error {
	return r.queries.DeletePartyContact(ctx, id)
}

// Contact persons
func (r *postgresContactRepository) CreatePartyContactPerson(
	ctx context.Context,
	params db.CreatePartyContactPersonParams,
) (db.PartyContactPerson, error) {
	return r.queries.CreatePartyContactPerson(ctx, params)
}

func (r *postgresContactRepository) GetPartyContactPersonByID(
	ctx context.Context,
	id pgtype.UUID,
) (db.PartyContactPerson, error) {
	return r.queries.GetPartyContactPersonByID(ctx, id)
}

func (r *postgresContactRepository) ListPartyContactPersonsByPartyID(
	ctx context.Context,
	partyID pgtype.UUID,
) ([]db.PartyContactPerson, error) {
	return r.queries.ListPartyContactPersonsByPartyID(ctx, partyID)
}

func (r *postgresContactRepository) UpdatePartyContactPerson(
	ctx context.Context,
	params db.UpdatePartyContactPersonParams,
) (db.PartyContactPerson, error) {
	return r.queries.UpdatePartyContactPerson(ctx, params)
}

func (r *postgresContactRepository) ResetPrimaryPartyContactPerson(
	ctx context.Context,
	partyID pgtype.UUID,
) error {
	return r.queries.ResetPrimaryPartyContactPerson(ctx, partyID)
}

func (r *postgresContactRepository) DeletePartyContactPerson(
	ctx context.Context,
	id pgtype.UUID,
) error {
	return r.queries.DeletePartyContactPerson(ctx, id)
}
