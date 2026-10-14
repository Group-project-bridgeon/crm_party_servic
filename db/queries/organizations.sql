-- name: GetOrganizationByPartyID :one
SELECT
    id,
    party_id,
    registration_number,
    industry,
    organization_type,
    founded_date,
    employee_count,
    description,
    created_at,
    updated_at
FROM organizations
WHERE party_id = $1;

-- name: CreateOrganization :one
INSERT INTO organizations (
    party_id,
    registration_number,
    industry,
    organization_type,
    founded_date,
    employee_count,
    description
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING
    id,
    party_id,
    registration_number,
    industry,
    organization_type,
    founded_date,
    employee_count,
    description,
    created_at,
    updated_at;