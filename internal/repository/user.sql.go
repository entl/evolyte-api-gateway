// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: user.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
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
RETURNING id, email, full_name, password, created_at, updated_at, role
`

type CreateUserParams struct {
	Email    string
	FullName string
	Password pgtype.Text
	Role     Roles
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Email,
		arg.FullName,
		arg.Password,
		arg.Role,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FullName,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Role,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, full_name, password, created_at, updated_at, role FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FullName,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Role,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, email, full_name, password, created_at, updated_at, role FROM users
WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FullName,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Role,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, email, full_name, password, created_at, updated_at, role FROM users
ORDER BY id
LIMIT $1 OFFSET $2
`

type ListUsersParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.FullName,
			&i.Password,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Role,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET
  email = COALESCE($2, email),
  full_name = COALESCE($3, full_name),
  password = COALESCE($4, password),
  role = COALESCE($5, role),
  updated_at = NOW()
WHERE id = $1
RETURNING id, email, full_name, password, created_at, updated_at, role
`

type UpdateUserParams struct {
	ID       int32
	Email    pgtype.Text
	FullName pgtype.Text
	Password pgtype.Text
	Role     NullRoles
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.ID,
		arg.Email,
		arg.FullName,
		arg.Password,
		arg.Role,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FullName,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Role,
	)
	return i, err
}
