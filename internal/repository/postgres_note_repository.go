package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
)

type postgresNoteRepository struct {
	queries *db.Queries
}

func NewPostgresNoteRepository(queries *db.Queries) NoteRepository {
	return &postgresNoteRepository{queries: queries}
}

func (r *postgresNoteRepository) CreatePartyNote(
	ctx context.Context,
	params db.CreatePartyNoteParams,
) (db.PartyNote, error) {
	return r.queries.CreatePartyNote(ctx, params)
}

func (r *postgresNoteRepository) GetPartyNoteByID(
	ctx context.Context,
	id pgtype.UUID,
) (db.PartyNote, error) {
	return r.queries.GetPartyNoteByID(ctx, id)
}

func (r *postgresNoteRepository) ListPartyNotesByPartyID(
	ctx context.Context,
	partyID pgtype.UUID,
) ([]db.PartyNote, error) {
	return r.queries.ListPartyNotesByPartyID(ctx, partyID)
}

func (r *postgresNoteRepository) DeletePartyNote(
	ctx context.Context,
	id pgtype.UUID,
) error {
	return r.queries.DeletePartyNote(ctx, id)
}
