-- name: GetIndividualByPartyID :one
SELECT
    id,
    party_id,
    first_name,
    middle_name,
    last_name,
    date_of_birth,
    gender,
    job_title,
    national_id,
    created_at,
    updated_at
FROM individuals
WHERE party_id = $1;

-- name: CreateIndividual :one
INSERT INTO individuals (
    party_id,
    first_name,
    middle_name,
    last_name,
    date_of_birth,
    gender,
    job_title,
    national_id
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
    party_id,
    first_name,
    middle_name,
    last_name,
    date_of_birth,
    gender,
    job_title,
    national_id,
    created_at,
    updated_at;

-- name: UpdateIndividual :one
UPDATE individuals
SET
    first_name = $2,
    middle_name = $3,
    last_name = $4,
    date_of_birth = $5,
    gender = $6,
    job_title = $7,
    national_id = $8,
    updated_at = NOW()
WHERE party_id = $1
RETURNING
    id,
    party_id,
    first_name,
    middle_name,
    last_name,
    date_of_birth,
    gender,
    job_title,
    national_id,
    created_at,
    updated_at;
