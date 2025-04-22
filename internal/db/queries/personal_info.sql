-- name: CreatePersonalInfo :one
INSERT INTO personal_infos (
    account_id,
    full_name,
    email,
    phone_number,
    linkedin_url,
    personal_url,
    country,
    state,
    city
) VALUES (
 $1, $2, $3,
 $4, $5, $6,
 $7, $8, $9
) RETURNING *;

-- name: GetPersonalInfo :one
SELECT * FROM personal_infos
WHERE id = $1;

-- name: UpdatePersonalInfo :one
UPDATE personal_infos
SET full_name = $1,
    email = $2,
    phone_number = $3,
    linkedin_url = $4,
    personal_url = $5,
    country = $6,
    state = $7,
    city = $8
WHERE id = $9
RETURNING *;

-- name: DeletePersonalInfo :exec
DELETE FROM personal_infos
WHERE id = $1;