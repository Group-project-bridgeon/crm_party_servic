-- name: CreatePartyContactPerson :one
INSERT INTO party_contact_persons (
    party_id,
    first_name,
    middle_name,
    last_name,
    designation,
    department,
    email,
    phone,
    mobile,
    is_primary,
    is_active,
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
    $8,
    $9,
    $10,
    $11,
    $12
)
RETURNING
    id,
    party_id,
    first_name,
    middle_name,
    last_name,
    designation,
    department,
    email,
    phone,
    mobile,
    is_primary,
    is_active,
    notes,
    created_at,
    updated_at;

-- name: GetPartyContactPersonByID :one
SELECT
    id,
    party_id,
    first_name,
    middle_name,
    last_name,
    designation,
    department,
    email,
    phone,
    mobile,
    is_primary,
    is_active,
    notes,
    created_at,
    updated_at
FROM party_contact_persons
WHERE id = $1;

-- name: ListPartyContactPersonsByPartyID :many
SELECT
    id,
    party_id,
    first_name,
    middle_name,
    last_name,
    designation,
    department,
    email,
    phone,
    mobile,
    is_primary,
    is_active,
    notes,
    created_at,
    updated_at
FROM party_contact_persons
WHERE party_id = $1
ORDER BY is_primary DESC, created_at DESC;

-- name: UpdatePartyContactPerson :one
UPDATE party_contact_persons
SET
    first_name = $2,
    middle_name = $3,
    last_name = $4,
    designation = $5,
    department = $6,
    email = $7,
    phone = $8,
    mobile = $9,
    is_primary = $10,
    is_active = $11,
    notes = $12,
    updated_at = NOW()
WHERE id = $1
RETURNING
    id,
    party_id,
    first_name,
    middle_name,
    last_name,
    designation,
    department,
    email,
    phone,
    mobile,
    is_primary,
    is_active,
    notes,
    created_at,
    updated_at;

-- name: ResetPrimaryPartyContactPerson :exec
UPDATE party_contact_persons
SET is_primary = FALSE, updated_at = NOW()
WHERE party_id = $1 AND is_primary = TRUE;

-- name: DeletePartyContactPerson :exec
DELETE FROM party_contact_persons
WHERE id = $1;
