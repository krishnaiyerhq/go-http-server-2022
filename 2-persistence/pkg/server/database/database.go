package database

import "context"

// User represents the JSON value that's sent as a response to a request.
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// Database represents the operations that are done on a Database.
// This interface abstracts the underlying implementation.
type Database interface {
	Create(ctx context.Context, user User) error
	Get(ctx context.Context, name string) *User
	// Update updates a given user (if found) and returns the updated value.
	Update(ctx context.Context, user User) (*User, error)
	Delete(ctx context.Context, name string) error
}
