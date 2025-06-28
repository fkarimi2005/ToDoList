package model

import "time"

type User struct {
	ID        int       `json:"id" swaggerignore:"true"`
	FullName  string    `json:"full_name" db:"full_name" example:"Firuz Karimzoda"`
	Username  string    `json:"username" db:"username" example:"firuz"`
	UserRole  string    `json:"user_role" db:"user_role" example:"user"`
	Password  string    `json:"-" db:"password" example:"your_password"`
	CreatedAt time.Time `json:"created_at" db:"created_at" swaggerignore:"true"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" swaggerignore:"true"`
	DeletedAt time.Time `json:"-" db:"deleted_at"`
}
type UserSignUp struct {
	FullName string `json:"full_name" example:"Firuz Karimzoda"`
	Username string `json:"username" example:"firuz"`
	Password string `json:"password" example:"your_password"`
	UserRole string `json:"user_role" example:"user"`
}
type UserSignIn struct {
	Username string `json:"username" db:"username" example:"firuz"`
	Password string `json:"password" db:"password" example:"your_password"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
