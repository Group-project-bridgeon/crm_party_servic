-- name: CreateSupplier :one
INSERT INTO suppliers (
    party_id,
    supplier_code,
    supplier_group,
    payment_terms_days,
    payment_method,
    currency_code,
    purchase_rep_id,
    tax_exempt,
    tax_exemption_number,
    supplier_rating
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
    supplier_code,
    supplier_group,
    payment_terms_days,
    payment_method,
    currency_code,
    purchase_rep_id,
    tax_exempt,
    tax_exemption_number,
    supplier_rating,
    created_at,
    updated_at;

-- name: GetSupplierByID :one
SELECT
    id,
    party_id,
    supplier_code,
    supplier_group,
    payment_terms_days,
    payment_method,
    currency_code,
    purchase_rep_id,
    tax_exempt,
    tax_exemption_number,
    supplier_rating,
    created_at,
    updated_at
FROM suppliers
WHERE id = $1;

-- name: GetSupplierByPartyID :one
SELECT
    id,
    party_id,
    supplier_code,
    supplier_group,
    payment_terms_days,
    payment_method,
    currency_code,
    purchase_rep_id,
    tax_exempt,
    tax_exemption_number,
    supplier_rating,
    created_at,
    updated_at
FROM suppliers
WHERE party_id = $1;

-- name: GetSupplierByCode :one
SELECT
    id,
    party_id,
    supplier_code,
    supplier_group,
    payment_terms_days,
    payment_method,
    currency_code,
    purchase_rep_id,
    tax_exempt,
    tax_exemption_number,
    supplier_rating,
    created_at,
    updated_at
FROM suppliers
WHERE supplier_code = $1;

-- name: ListSuppliers :many
SELECT
    s.id,
    s.party_id,
    s.supplier_code,
    s.supplier_group,
    s.payment_terms_days,
    s.payment_method,
    s.currency_code,
    s.purchase_rep_id,
    s.tax_exempt,
    s.tax_exemption_number,
    s.supplier_rating,
    s.created_at,
    s.updated_at
FROM suppliers s
WHERE
    (
        sqlc.narg('supplier_group')::VARCHAR IS NULL
        OR s.supplier_group = sqlc.narg('supplier_group')::VARCHAR
    )
    AND (
        sqlc.narg('currency_code')::VARCHAR IS NULL
        OR s.currency_code = sqlc.narg('currency_code')::VARCHAR
    )
    AND (
        sqlc.narg('search')::VARCHAR IS NULL
        OR s.supplier_code ILIKE '%' || sqlc.narg('search')::VARCHAR || '%'
        OR s.tax_exemption_number ILIKE '%' || sqlc.narg('search')::VARCHAR || '%'
    )
ORDER BY s.created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: UpdateSupplier :one
UPDATE suppliers
SET
    supplier_group = $2,
    payment_terms_days = $3,
    payment_method = $4,
    currency_code = $5,
    purchase_rep_id = $6,
    tax_exempt = $7,
    tax_exemption_number = $8,
    supplier_rating = $9,
    updated_at = NOW()
WHERE id = $1
RETURNING
    id,
    party_id,
    supplier_code,
    supplier_group,
    payment_terms_days,
    payment_method,
    currency_code,
    purchase_rep_id,
    tax_exempt,
    tax_exemption_number,
    supplier_rating,
    created_at,
    updated_at;

-- name: CreateSupplierUserAccount :one
INSERT INTO supplier_user_accounts (
    supplier_id,
    user_id
)
VALUES (
    $1,
    $2
)
RETURNING
    id,
    supplier_id,
    user_id,
    is_active,
    created_at,
    updated_at;

-- name: ListSupplierUserAccountsBySupplierID :many
SELECT
    id,
    supplier_id,
    user_id,
    is_active,
    created_at,
    updated_at
FROM supplier_user_accounts
WHERE supplier_id = $1
ORDER BY created_at DESC;

-- name: UpdateSupplierUserAccountStatus :one
UPDATE supplier_user_accounts
SET
    is_active = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING
    id,
    supplier_id,
    user_id,
    is_active,
    created_at,
    updated_at;
