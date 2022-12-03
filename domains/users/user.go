package users

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	_ UserServicer = (*UserService)(nil)
)

// UserService abstraction for
type UserService struct {
	db *gorm.DB
}

// New returns a new instance of UserService
func New(db *gorm.DB) (UserServicer, error) {
	if err := db.AutoMigrate(&User{}); err != nil {
		return nil, err
	}
	return &UserService{db: db}, nil
}

// Create creates a new user record
func (us *UserService) Create(ctx context.Context, user *User) error {
	user.ID = uuid.NewString()
	return us.db.WithContext(ctx).Create(&user).Error
}

// Delete remove a user's record
func (us *UserService) Delete(ctx context.Context, id string) error {
	return us.db.WithContext(ctx).Where("id = ?", id).Delete(&User{}).Error //Delete(&User{}, id).Error
}

// Find returns all the users in the system
func (us *UserService) Find(ctx context.Context) ([]User, error) {
	var result []User
	err := us.db.WithContext(ctx).Find(&result).Error
	return result, err
}

// FindByID finds a user by ID and returns it
func (us *UserService) FindByID(ctx context.Context, id string) (*User, error) {
	var user *User
	err := us.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return user, err
}

// Update updates a user's record with the given ID
func (us *UserService) Update(ctx context.Context, id string, user *User) error {
	existing, err := us.FindByID(ctx, id)
	if err != nil {
		return err
	}
	user.ID = id
	user.CreatedAt = existing.CreatedAt
	return us.db.WithContext(ctx).Save(&user).Error
}
