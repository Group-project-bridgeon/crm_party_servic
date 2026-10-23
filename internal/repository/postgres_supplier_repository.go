package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
)

type postgresSupplierRepository struct {
	queries *db.Queries
}

func NewPostgresSupplierRepository(queries *db.Queries) SupplierRepository {
	return &postgresSupplierRepository{
		queries: queries,
	}
}

func (r *postgresSupplierRepository) CreateSupplier(
	ctx context.Context,
	params db.CreateSupplierParams,
) (db.Supplier, error) {
	return r.queries.CreateSupplier(ctx, params)
}

func (r *postgresSupplierRepository) GetSupplierByID(
	ctx context.Context,
	id pgtype.UUID,
) (db.Supplier, error) {
	return r.queries.GetSupplierByID(ctx, id)
}

func (r *postgresSupplierRepository) GetSupplierByPartyID(
	ctx context.Context,
	partyID pgtype.UUID,
) (db.Supplier, error) {
	return r.queries.GetSupplierByPartyID(ctx, partyID)
}

func (r *postgresSupplierRepository) GetSupplierByCode(
	ctx context.Context,
	supplierCode string,
) (db.Supplier, error) {
	return r.queries.GetSupplierByCode(ctx, supplierCode)
}

func (r *postgresSupplierRepository) ListSuppliers(
	ctx context.Context,
	params db.ListSuppliersParams,
) ([]db.Supplier, error) {
	return r.queries.ListSuppliers(ctx, params)
}

func (r *postgresSupplierRepository) UpdateSupplier(
	ctx context.Context,
	params db.UpdateSupplierParams,
) (db.Supplier, error) {
	return r.queries.UpdateSupplier(ctx, params)
}

func (r *postgresSupplierRepository) CreateSupplierUserAccount(
	ctx context.Context,
	params db.CreateSupplierUserAccountParams,
) (db.SupplierUserAccount, error) {
	return r.queries.CreateSupplierUserAccount(ctx, params)
}

func (r *postgresSupplierRepository) ListSupplierUserAccountsBySupplierID(
	ctx context.Context,
	supplierID pgtype.UUID,
) ([]db.SupplierUserAccount, error) {
	return r.queries.ListSupplierUserAccountsBySupplierID(ctx, supplierID)
}

func (r *postgresSupplierRepository) UpdateSupplierUserAccountStatus(
	ctx context.Context,
	params db.UpdateSupplierUserAccountStatusParams,
) (db.SupplierUserAccount, error) {
	return r.queries.UpdateSupplierUserAccountStatus(ctx, params)
}
