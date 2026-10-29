package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
)

type AddressRepository interface {
	CreatePartyAddress(ctx context.Context, params db.CreatePartyAddressParams) (db.PartyAddress, error)
	GetPartyAddressByID(ctx context.Context, id pgtype.UUID) (db.PartyAddress, error)
	ListPartyAddressesByPartyID(ctx context.Context, params db.ListPartyAddressesByPartyIDParams) ([]db.PartyAddress, error)
	UpdatePartyAddress(ctx context.Context, params db.UpdatePartyAddressParams) (db.PartyAddress, error)
	ResetPrimaryPartyAddress(ctx context.Context, params db.ResetPrimaryPartyAddressParams) error
	DeletePartyAddress(ctx context.Context, id pgtype.UUID) error
}
