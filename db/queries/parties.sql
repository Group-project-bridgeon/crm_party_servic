-- name: GetPartyByID :one
SELECT
    id,
    party_type,
    display_name,
    legal_name,
    email,
    phone,
    website,
    tax_number,
    notes,
    is_active,
    created_at,
    updated_at
FROM parties
WHERE id = $1;



-- name: ListParties :many
SELECT
    id,
    party_type,
    display_name,
    legal_name,
    email,
    phone,
    website,
    tax_number,
    notes,
    is_active,
    created_at,
    updated_at
FROM parties
WHERE
    (
        sqlc.narg('party_type')::VARCHAR IS NULL
        OR party_type = sqlc.narg('party_type')::VARCHAR
    )
    AND (
        sqlc.narg('is_active')::BOOLEAN IS NULL
        OR is_active = sqlc.narg('is_active')::BOOLEAN
    )
    AND (
        sqlc.narg('search')::VARCHAR IS NULL
        OR display_name ILIKE '%' || sqlc.narg('search')::VARCHAR || '%'
        OR legal_name ILIKE '%' || sqlc.narg('search')::VARCHAR || '%'
        OR email ILIKE '%' || sqlc.narg('search')::VARCHAR || '%'
        OR phone ILIKE '%' || sqlc.narg('search')::VARCHAR || '%'
    )
ORDER BY created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');



-- name: CreateParty :one
INSERT INTO parties (
    party_type,
    display_name,
    legal_name,
    email,
    phone,
    website,
    tax_number,
    notes
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING
    id,
    party_type,
    display_name,
    legal_name,
    email,
    phone,
    website,
    tax_number,
    notes,
    is_active,
    created_at,
    updated_at;

-- name: UpdateParty :one
UPDATE parties
SET
    display_name = $2,
    legal_name = $3,
    email = $4,
    phone = $5,
    website = $6,
    tax_number = $7,
    notes = $8,
    updated_at = NOW()
WHERE id = $1
RETURNING
    id,
    party_type,
    display_name,
    legal_name,
    email,
    phone,
    website,
    tax_number,
    notes,
    is_active,
    created_at,
    updated_at;

-- name: UpdatePartyStatus :one
UPDATE parties
SET
    is_active = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING
    id,
    party_type,
    display_name,
    legal_name,
    email,
    phone,
    website,
    tax_number,
    notes,
    is_active,
    created_at,
    updated_at;
