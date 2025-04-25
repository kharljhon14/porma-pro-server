-- name: CreateSummary :one
INSERT INTO summaries(
    account_id,
    summary
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetSummary :one
SELECT * FROM summaries
WHERE id = $1;

-- name: UpdateSummary :one
UPDATE summaries
SET summary = $1
WHERE id = $2
RETURNING *;

-- name: DeleteSummary :exec
DELETE FROM summaries
WHERE id = $1;
