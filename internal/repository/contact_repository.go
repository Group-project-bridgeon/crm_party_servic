package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
)

type ContactRepository interface {
	// Direct contacts
	CreatePartyContact(ctx context.Context, params db.CreatePartyContactParams) (db.PartyContact, error)
	GetPartyContactByID(ctx context.Context, id pgtype.UUID) (db.PartyContact, error)
	ListPartyContactsByPartyID(ctx context.Context, params db.ListPartyContactsByPartyIDParams) ([]db.PartyContact, error)
	UpdatePartyContact(ctx context.Context, params db.UpdatePartyContactParams) (db.PartyContact, error)
	ResetPrimaryPartyContact(ctx context.Context, params db.ResetPrimaryPartyContactParams) error
	DeletePartyContact(ctx context.Context, id pgtype.UUID) error

	// Contact persons
	CreatePartyContactPerson(ctx context.Context, params db.CreatePartyContactPersonParams) (db.PartyContactPerson, error)
	GetPartyContactPersonByID(ctx context.Context, id pgtype.UUID) (db.PartyContactPerson, error)
	ListPartyContactPersonsByPartyID(ctx context.Context, partyID pgtype.UUID) ([]db.PartyContactPerson, error)
	UpdatePartyContactPerson(ctx context.Context, params db.UpdatePartyContactPersonParams) (db.PartyContactPerson, error)
	ResetPrimaryPartyContactPerson(ctx context.Context, partyID pgtype.UUID) error
	DeletePartyContactPerson(ctx context.Context, id pgtype.UUID) error
}
