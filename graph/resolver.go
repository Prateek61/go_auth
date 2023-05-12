package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/Prateek61/go_auth/postgres"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	TodosRepo postgres.TodosRepo
	UsersRepo postgres.UsersRepo
}
