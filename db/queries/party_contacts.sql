-- name: CreatePartyContact :one
INSERT INTO party_contacts (
    party_id,
    contact_type,
    contact_value,
    label,
    is_primary,
    is_active
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING
    id,
    party_id,
    contact_type,
    contact_value,
    label,
    is_primary,
    is_active,
    created_at,
    updated_at;

-- name: GetPartyContactByID :one
SELECT
    id,
    party_id,
    contact_type,
    contact_value,
    label,
    is_primary,
    is_active,
    created_at,
    updated_at
FROM party_contacts
WHERE id = $1;

-- name: ListPartyContactsByPartyID :many
SELECT
    id,
    party_id,
    contact_type,
    contact_value,
    label,
    is_primary,
    is_active,
    created_at,
    updated_at
FROM party_contacts
WHERE
    party_id = $1
    AND (
        sqlc.narg('contact_type')::VARCHAR IS NULL
        OR contact_type = sqlc.narg('contact_type')::VARCHAR
    )
ORDER BY is_primary DESC, created_at DESC;

-- name: UpdatePartyContact :one
UPDATE party_contacts
SET
    contact_type = $2,
    contact_value = $3,
    label = $4,
    is_primary = $5,
    is_active = $6,
    updated_at = NOW()
WHERE id = $1
RETURNING
    id,
    party_id,
    contact_type,
    contact_value,
    label,
    is_primary,
    is_active,
    created_at,
    updated_at;

-- name: ResetPrimaryPartyContact :exec
UPDATE party_contacts
SET is_primary = FALSE, updated_at = NOW()
WHERE party_id = $1 AND contact_type = $2 AND is_primary = TRUE;

-- name: DeletePartyContact :exec
DELETE FROM party_contacts
WHERE id = $1;
