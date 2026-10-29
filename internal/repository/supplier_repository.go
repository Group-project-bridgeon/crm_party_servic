package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
)

type SupplierRepository interface {
	CreateSupplier(ctx context.Context, params db.CreateSupplierParams) (db.Supplier, error)
	GetSupplierByID(ctx context.Context, id pgtype.UUID) (db.Supplier, error)
	GetSupplierByPartyID(ctx context.Context, partyID pgtype.UUID) (db.Supplier, error)
	GetSupplierByCode(ctx context.Context, supplierCode string) (db.Supplier, error)
	ListSuppliers(ctx context.Context, params db.ListSuppliersParams) ([]db.Supplier, error)
	UpdateSupplier(ctx context.Context, params db.UpdateSupplierParams) (db.Supplier, error)

	CreateSupplierUserAccount(ctx context.Context, params db.CreateSupplierUserAccountParams) (db.SupplierUserAccount, error)
	ListSupplierUserAccountsBySupplierID(ctx context.Context, supplierID pgtype.UUID) ([]db.SupplierUserAccount, error)
	UpdateSupplierUserAccountStatus(ctx context.Context, params db.UpdateSupplierUserAccountStatusParams) (db.SupplierUserAccount, error)
}
