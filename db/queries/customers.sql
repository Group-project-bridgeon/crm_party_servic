-- name: CreateCustomer :one
INSERT INTO customers (
    party_id,
    customer_code,
    customer_group,
    price_list_id,
    credit_limit,
    payment_terms_days,
    credit_status,
    sales_rep_id,
    tax_exempt,
    tax_exemption_number
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10
)
RETURNING
    id,
    party_id,
    customer_code,
    customer_group,
    price_list_id,
    credit_limit,
    payment_terms_days,
    credit_status,
    sales_rep_id,
    tax_exempt,
    tax_exemption_number,
    created_at,
    updated_at;

-- name: GetCustomerByID :one
SELECT
    id,
    party_id,
    customer_code,
    customer_group,
    price_list_id,
    credit_limit,
    payment_terms_days,
    credit_status,
    sales_rep_id,
    tax_exempt,
    tax_exemption_number,
    created_at,
    updated_at
FROM customers
WHERE id = $1;

-- name: GetCustomerByPartyID :one
SELECT
    id,
    party_id,
    customer_code,
    customer_group,
    price_list_id,
    credit_limit,
    payment_terms_days,
    credit_status,
    sales_rep_id,
    tax_exempt,
    tax_exemption_number,
    created_at,
    updated_at
FROM customers
WHERE party_id = $1;

-- name: GetCustomerByCode :one
SELECT
    id,
    party_id,
    customer_code,
    customer_group,
    price_list_id,
    credit_limit,
    payment_terms_days,
    credit_status,
    sales_rep_id,
    tax_exempt,
    tax_exemption_number,
    created_at,
    updated_at
FROM customers
WHERE customer_code = $1;

-- name: ListCustomers :many
SELECT
    c.id,
    c.party_id,
    c.customer_code,
    c.customer_group,
    c.price_list_id,
    c.credit_limit,
    c.payment_terms_days,
    c.credit_status,
    c.sales_rep_id,
    c.tax_exempt,
    c.tax_exemption_number,
    c.created_at,
    c.updated_at
FROM customers c
WHERE
    (
        sqlc.narg('customer_group')::VARCHAR IS NULL
        OR c.customer_group = sqlc.narg('customer_group')::VARCHAR
    )
    AND (
        sqlc.narg('credit_status')::VARCHAR IS NULL
        OR c.credit_status = sqlc.narg('credit_status')::VARCHAR
    )
    AND (
        sqlc.narg('search')::VARCHAR IS NULL
        OR c.customer_code ILIKE '%' || sqlc.narg('search')::VARCHAR || '%'
        OR c.tax_exemption_number ILIKE '%' || sqlc.narg('search')::VARCHAR || '%'
    )
ORDER BY c.created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: UpdateCustomer :one
UPDATE customers
SET
    customer_group = $2,
    price_list_id = $3,
    credit_limit = $4,
    payment_terms_days = $5,
    credit_status = $6,
    sales_rep_id = $7,
    tax_exempt = $8,
    tax_exemption_number = $9,
    updated_at = NOW()
WHERE id = $1
RETURNING
    id,
    party_id,
    customer_code,
    customer_group,
    price_list_id,
    credit_limit,
    payment_terms_days,
    credit_status,
    sales_rep_id,
    tax_exempt,
    tax_exemption_number,
    created_at,
    updated_at;

-- name: CreateCustomerUserAccount :one
INSERT INTO customer_user_accounts (
    customer_id,
    user_id
)
VALUES (
    $1,
    $2
)
RETURNING
    id,
    customer_id,
    user_id,
    is_active,
    created_at,
    updated_at;

-- name: ListCustomerUserAccountsByCustomerID :many
SELECT
    id,
    customer_id,
    user_id,
    is_active,
    created_at,
    updated_at
FROM customer_user_accounts
WHERE customer_id = $1
ORDER BY created_at DESC;

-- name: UpdateCustomerUserAccountStatus :one
UPDATE customer_user_accounts
SET
    is_active = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING
    id,
    customer_id,
    user_id,
    is_active,
    created_at,
    updated_at;
