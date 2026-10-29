package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/group-project/crm_party_service/internal/db"
)

type CustomerRepository interface {
	CreateCustomer(ctx context.Context, params db.CreateCustomerParams) (db.Customer, error)
	GetCustomerByID(ctx context.Context, id pgtype.UUID) (db.Customer, error)
	GetCustomerByPartyID(ctx context.Context, partyID pgtype.UUID) (db.Customer, error)
	GetCustomerByCode(ctx context.Context, customerCode string) (db.Customer, error)
	ListCustomers(ctx context.Context, params db.ListCustomersParams) ([]db.Customer, error)
	UpdateCustomer(ctx context.Context, params db.UpdateCustomerParams) (db.Customer, error)

	CreateCustomerUserAccount(ctx context.Context, params db.CreateCustomerUserAccountParams) (db.CustomerUserAccount, error)
	ListCustomerUserAccountsByCustomerID(ctx context.Context, customerID pgtype.UUID) ([]db.CustomerUserAccount, error)
	UpdateCustomerUserAccountStatus(ctx context.Context, params db.UpdateCustomerUserAccountStatusParams) (db.CustomerUserAccount, error)
}
