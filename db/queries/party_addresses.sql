-- name: CreatePartyAddress :one
INSERT INTO party_addresses (
    party_id,
    address_type,
    address_name,
    address_line1,
    address_line2,
    landmark,
    city,
    state,
    postal_code,
    country_code,
    latitude,
    longitude,
    is_primary,
    is_active
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
    $10,
    $11,
    $12,
    $13,
    $14
)
RETURNING
    id,
    party_id,
    address_type,
    address_name,
    address_line1,
    address_line2,
    landmark,
    city,
    state,
    postal_code,
    country_code,
    latitude,
    longitude,
    is_primary,
    is_active,
    created_at,
    updated_at;

-- name: GetPartyAddressByID :one
SELECT
    id,
    party_id,
    address_type,
    address_name,
    address_line1,
    address_line2,
    landmark,
    city,
    state,
    postal_code,
    country_code,
    latitude,
    longitude,
    is_primary,
    is_active,
    created_at,
    updated_at
FROM party_addresses
WHERE id = $1;

-- name: ListPartyAddressesByPartyID :many
SELECT
    id,
    party_id,
    address_type,
    address_name,
    address_line1,
    address_line2,
    landmark,
    city,
    state,
    postal_code,
    country_code,
    latitude,
    longitude,
    is_primary,
    is_active,
    created_at,
    updated_at
FROM party_addresses
WHERE
    party_id = $1
    AND (
        sqlc.narg('address_type')::VARCHAR IS NULL
        OR address_type = sqlc.narg('address_type')::VARCHAR
    )
ORDER BY is_primary DESC, created_at DESC;

-- name: UpdatePartyAddress :one
UPDATE party_addresses
SET
    address_type = $2,
    address_name = $3,
    address_line1 = $4,
    address_line2 = $5,
    landmark = $6,
    city = $7,
    state = $8,
    postal_code = $9,
    country_code = $10,
    latitude = $11,
    longitude = $12,
    is_primary = $13,
    is_active = $14,
    updated_at = NOW()
WHERE id = $1
RETURNING
    id,
    party_id,
    address_type,
    address_name,
    address_line1,
    address_line2,
    landmark,
    city,
    state,
    postal_code,
    country_code,
    latitude,
    longitude,
    is_primary,
    is_active,
    created_at,
    updated_at;

-- name: ResetPrimaryPartyAddress :exec
UPDATE party_addresses
SET is_primary = FALSE, updated_at = NOW()
WHERE party_id = $1 AND address_type = $2 AND is_primary = TRUE;

-- name: DeletePartyAddress :exec
DELETE FROM party_addresses
WHERE id = $1;
