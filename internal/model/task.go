package model

import "time"

type Tasks struct {
	TaskID      int       `json:"task_id" db:"task_id" swaggerignore:"true"`
	User_ID     int       `json:"user_id" db:"user_id" example:"0"`
	Title       string    `json:"title" binding:"required" db:"title"`
	Description string    `json:"description" binding:"required" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" example:"string" swaggerignore:"true"`
	DueDate     time.Time `json:"due_date" db:"due_date" swaggerignore:"true"`
	Done        bool      `json:"done" db:"done" swaggerignore:"true"`

	UpdatedAt time.Time  `json:"updated_at" db:"updated_at" swaggerignore:"true"`
	DeletedAt *time.Time `json:"-" db:"deleted_at" swaggerignore:"true"`
}
type DoneTasks struct {
	Done    bool `json:"done" db:"done"`
	DueDate int  `json:"due_date"`
}
