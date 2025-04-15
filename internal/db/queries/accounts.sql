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
    updated_at = $2
WHERE id = $3
RETURNING *;

-- name: VerifyAccount :one
UPDATE accounts
SET is_verified = true
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;