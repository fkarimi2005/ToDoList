package model

import "time"

type User struct {
	ID        int       `json:"id"`
	FullName  string    `json:"full_name" db:"full_name"`
	Username  string    `json:"username" db:"username"`
	UserRole  string    `json:"user_role" db:"user_role"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt time.Time `json:"-" db:"deleted_at"`
}
type UserSignIn struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}
