package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.31

import (
	"context"
	"errors"
	"log"

	"github.com/Prateek61/go_auth/graph/model"
	"github.com/Prateek61/go_auth/middleware"
)

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, input model.TodoInput) (*model.Todo, error) {
	currentUser, err := middleware.GetCurrentUserFromCTX(ctx)
	if err != nil {
		return nil, ErrUnauthenticated
	}

	// Create new todo
	newTodo := &model.Todo{
		Text:   input.Text,
		UserID: currentUser.ID,
	}

	err = r.TodosRepo.CreateTodo(newTodo)
	return newTodo, err
}

// Register is the resolver for the register field.
// It creates a new user and returns the user and auth token.
func (r *mutationResolver) Register(ctx context.Context, input *model.RegisterInput) (*model.AuthResponse, error) {
	// Check if user with input.Email exists
	_, err := r.UsersRepo.GetUserByEmail(input.Email)
	if err == nil {
		return nil, errors.New("email already in use")
	}

	// Check if user with input.Username exists
	_, err = r.UsersRepo.GetUserByUsername(input.Username)
	if err == nil {
		return nil, errors.New("username already in use")
	}

	// Create new user
	newUser := &model.User{
		Username:  input.Username,
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		DeletedAt: nil,
	}

	// Hash password
	err = newUser.HashPassword(input.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, errors.New("something went wrong")
	}

	// Start transaction
	transaction, err := r.UsersRepo.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return nil, errors.New("something went wrong")
	}
	defer transaction.Rollback()

	// Add user to transaction
	_, err = r.UsersRepo.CreateUser(transaction, newUser)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, errors.New("something went wrong")
	}

	// Commit transaction
	err = transaction.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, errors.New("something went wrong")
	}

	// Generate token for user
	token, err := newUser.GenerateToken()
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return nil, errors.New("something went wrong")
	}

	return &model.AuthResponse{
		User:      newUser,
		AuthToken: token,
	}, nil
}

// Login is the resolver for the login field.
// It checks if the user exists and if the password is correct and returns the user and auth token.
func (r *mutationResolver) Login(ctx context.Context, input *model.LoginInput) (*model.AuthResponse, error) {
	// Check if user with input.Username exists
	user, err := r.UsersRepo.GetUserByUsername(input.Username)
	if err != nil {
		return nil, ErrBadCredentials
	}

	// Check if password is correct
	err = user.CheckPassword(input.Password)
	if err != nil {
		return nil, ErrBadCredentials
	}

	// Generate token for user
	token, err := user.GenerateToken()
	if err != nil {
		return nil, errors.New("something went wrong")
	}

	return &model.AuthResponse{
		User:      user,
		AuthToken: token,
	}, nil
}

// Todos is the resolver for the todos field.
// It returns all todos for the current user.
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	// Get current user
	currentUser, err := middleware.GetCurrentUserFromCTX(ctx)
	if err != nil {
		return nil, ErrUnauthenticated
	}

	// Get todos for current user
	return r.TodosRepo.GetTodosByUserID(currentUser.ID)
}

// User is the resolver for the user field.
// It returns the current user.
func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	currentUser, err := middleware.GetCurrentUserFromCTX(ctx)
	if err != nil {
		return nil, ErrUnauthenticated
	}

	return currentUser, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
