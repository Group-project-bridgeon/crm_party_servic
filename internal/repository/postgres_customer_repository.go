package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
)

type postgresCustomerRepository struct {
	queries *db.Queries
}

func NewPostgresCustomerRepository(queries *db.Queries) CustomerRepository {
	return &postgresCustomerRepository{
		queries: queries,
	}
}

func (r *postgresCustomerRepository) CreateCustomer(
	ctx context.Context,
	params db.CreateCustomerParams,
) (db.Customer, error) {
	return r.queries.CreateCustomer(ctx, params)
}

func (r *postgresCustomerRepository) GetCustomerByID(
	ctx context.Context,
	id pgtype.UUID,
) (db.Customer, error) {
	return r.queries.GetCustomerByID(ctx, id)
}

func (r *postgresCustomerRepository) GetCustomerByPartyID(
	ctx context.Context,
	partyID pgtype.UUID,
) (db.Customer, error) {
	return r.queries.GetCustomerByPartyID(ctx, partyID)
}

func (r *postgresCustomerRepository) GetCustomerByCode(
	ctx context.Context,
	customerCode string,
) (db.Customer, error) {
	return r.queries.GetCustomerByCode(ctx, customerCode)
}

func (r *postgresCustomerRepository) ListCustomers(
	ctx context.Context,
	params db.ListCustomersParams,
) ([]db.Customer, error) {
	return r.queries.ListCustomers(ctx, params)
}

func (r *postgresCustomerRepository) UpdateCustomer(
	ctx context.Context,
	params db.UpdateCustomerParams,
) (db.Customer, error) {
	return r.queries.UpdateCustomer(ctx, params)
}

func (r *postgresCustomerRepository) CreateCustomerUserAccount(
	ctx context.Context,
	params db.CreateCustomerUserAccountParams,
) (db.CustomerUserAccount, error) {
	return r.queries.CreateCustomerUserAccount(ctx, params)
}

func (r *postgresCustomerRepository) ListCustomerUserAccountsByCustomerID(
	ctx context.Context,
	customerID pgtype.UUID,
) ([]db.CustomerUserAccount, error) {
	return r.queries.ListCustomerUserAccountsByCustomerID(ctx, customerID)
}

func (r *postgresCustomerRepository) UpdateCustomerUserAccountStatus(
	ctx context.Context,
	params db.UpdateCustomerUserAccountStatusParams,
) (db.CustomerUserAccount, error) {
	return r.queries.UpdateCustomerUserAccountStatus(ctx, params)
}
