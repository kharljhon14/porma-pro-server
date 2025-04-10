-- name: CreateAccount :one
INSERT INTO accounts(
    email,
    password_hash,
    full_name
) VALUES(
 $1, $2, $3
) RETURNING *;
-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1;

-- name: UpdateAccount :one
UPDATE accounts
SET full_name = $1,
    is_verified = $2,
    updated_at = $3
WHERE id = $4
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;