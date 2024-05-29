-- These are queries related to the built-in authentication/authorization system.

-- User

-- name: InsertUser :one
SELECT rbac.insert_user($1, $2, $3);

-- name: GetUserById :one
SELECT * FROM rbac.users WHERE id = $1 AND deleted_at IS NULL;