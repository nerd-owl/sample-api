// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package db

import (
	"context"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) error
	DeactivateUser(ctx context.Context, phone string) error
	ListUser(ctx context.Context) ([]Kuser, error)
}

var _ Querier = (*Queries)(nil)