-- name: CreateUser :one
INSERT INTO users (
  email,
  full_name,
  password,
  role,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, $4, NOW(), NOW()
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET
  email = COALESCE(sqlc.narg(email), email),
  full_name = COALESCE(sqlc.narg(full_name), full_name),
  password = COALESCE(sqlc.narg(password), password),
  role = COALESCE(sqlc.narg(role), role),
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
