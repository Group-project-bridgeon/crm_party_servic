-- name: CreatePartyNote :one
INSERT INTO party_notes (
    party_id,
    note_type,
    content,
    created_by
)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING
    id,
    party_id,
    note_type,
    content,
    created_by,
    created_at;

-- name: GetPartyNoteByID :one
SELECT
    id,
    party_id,
    note_type,
    content,
    created_by,
    created_at
FROM party_notes
WHERE id = $1;

-- name: ListPartyNotesByPartyID :many
SELECT
    id,
    party_id,
    note_type,
    content,
    created_by,
    created_at
FROM party_notes
WHERE party_id = $1
ORDER BY created_at DESC;

-- name: DeletePartyNote :exec
DELETE FROM party_notes
WHERE id = $1;
