-- name: CreateWorkExperience :one
INSERT INTO work_experiences (
    account_id,
    role,
    company,
    location,
    summary,
    start_date,
    end_date
) VALUES (
    $1, $2, $3,
    $4, $5, $6,
    $7
) RETURNING *;

-- name: GetWorkExperience :one
SELECT * FROM work_experiences
WHERE id = $1;

-- name: GetWorkExperiences :many
SELECT * FROM work_experiences
WHERE account_id = $1;

-- name: UpdateWorkExperience :one
UPDATE work_experiences
SET role = $1,
    company = $2,
    location = $3,
    summary = $4,
    start_date = $5,
    end_date = $6
WHERE id = $7
RETURNING *;

-- name: DeleteWorkExperience :exec
DELETE FROM work_experiences
WHERE id = $1;


