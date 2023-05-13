package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"errors"
	"github.com/Prateek61/go_auth/postgres"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	TodosRepo postgres.TodosRepo
	UsersRepo postgres.UsersRepo
}

// Error definitions
var (
	// ErrUnauthenticated is returned when the user is not authenticated.
	ErrUnauthenticated = errors.New("unauthenticated")
	// ErrBadCredentials is returned when the user provides invalid credentials.
	ErrBadCredentials = errors.New("username or password is incorrect")
)
