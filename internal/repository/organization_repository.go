package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
)

type OrganizationRepository interface {
	CreateOrganization(
		ctx context.Context,
		params db.CreateOrganizationParams,
	) (db.Organization, error)

	GetOrganizationByPartyID(
		ctx context.Context,
		partyID pgtype.UUID,
	) (db.Organization, error)
}
