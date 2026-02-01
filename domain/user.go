package domain

import (
	"context"
	"time"
)

type User struct {
	Id          int        `json:"id" db:"id"`
	GoogleId    string     `json:"google_id" db:"google_id"`
	AvatarEmoji int        `json:"avatar_emoji" db:"avatar_emoji"`
	Name        string     `json:"name" db:"name"`
	Password    string     `json:"password" db:"password"`
	Email       string     `json:"email" db:"email"`
	Phone       string     `json:"phone" db:"phone"`
	IsVerified  bool       `json:"is_verified" db:"is_verified"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at" db:"updated_at"`
}

type UserResponse struct {
	Id          int       `json:"id" db:"id"`
	GoogleId    string    `json:"google_id" db:"google_id"`
	AvatarEmoji int       `json:"avatar_emoji" db:"avatar_emoji"`
	Name        string    `json:"name" db:"name"`
	Email       string    `json:"email" db:"email"`
	Phone       string    `json:"phone" db:"phone"`
	IsVerified  bool      `json:"is_verified" db:"is_verified"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// UserUpdateRequest represents the request body for updating user profile
type UserUpdateRequest struct {
	Name        string `json:"name,omitempty" example:"John Doe"`
	Email       string `json:"email,omitempty" example:"john@example.com"`
	Password    string `json:"password,omitempty" example:"newpassword123"`
	Phone       string `json:"phone,omitempty" example:"+1234567890"`
	AvatarEmoji *int   `json:"avatar_emoji,omitempty" example:"5"`
}

type UserUseCase interface {
	GetUserById(c context.Context, id int) (*UserResponse, error)
	GetUsers(c context.Context) ([]*UserResponse, error)
	UpdateUser(c context.Context, user *User) error
	DeleteUser(c context.Context, id int) error
}
