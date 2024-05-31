// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: rbac.sql

package models

import (
	"context"
)

const getUserById = `-- name: GetUserById :one
SELECT id, created_at, deleted_at, first_name, last_name, username, password, email, description, account_locked, path_to_avatar, permissions FROM rbac.users WHERE id = $1 AND deleted_at IS NULL
`

func (q *Queries) GetUserById(ctx context.Context, id int32) (RbacUser, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i RbacUser
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.DeletedAt,
		&i.FirstName,
		&i.LastName,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.Description,
		&i.AccountLocked,
		&i.PathToAvatar,
		&i.Permissions,
	)
	return i, err
}

const insertUser = `-- name: InsertUser :one


SELECT rbac.insert_user($1, $2, $3)
`

// These are queries related to the built-in authentication/authorization system.
// User
func (q *Queries) InsertUser(ctx context.Context, pUsername string, pEmail string, pPassword string) (int32, error) {
	row := q.db.QueryRow(ctx, insertUser, pUsername, pEmail, pPassword)
	var insert_user int32
	err := row.Scan(&insert_user)
	return insert_user, err
}
