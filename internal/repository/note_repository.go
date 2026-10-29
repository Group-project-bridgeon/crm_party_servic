package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
)

type NoteRepository interface {
	CreatePartyNote(ctx context.Context, params db.CreatePartyNoteParams) (db.PartyNote, error)
	GetPartyNoteByID(ctx context.Context, id pgtype.UUID) (db.PartyNote, error)
	ListPartyNotesByPartyID(ctx context.Context, partyID pgtype.UUID) ([]db.PartyNote, error)
	DeletePartyNote(ctx context.Context, id pgtype.UUID) error
}
