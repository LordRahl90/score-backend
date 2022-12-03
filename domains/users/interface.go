package users

import "context"

// UserServicer interface for user management
type UserServicer interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, id string, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	Find(ctx context.Context) ([]User, error)
	Delete(ctx context.Context, id string) error
}
