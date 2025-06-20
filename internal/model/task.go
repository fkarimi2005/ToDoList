package model

import "time"

type Tasks struct {
	TaskID      int       `json:"task_id" db:"task_id"`
	User_ID     int       `json:"user_id" db:"user_id"`
	Title       string    `json:"title" binding:"required" db:"title"`
	Description string    `json:"description" binding:"required" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	DueDate     time.Time `json:"due_date" db:"due_date"`
	Done        bool      `json:"done" db:"done"`

	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}
type DoneTasks struct {
	Done    bool `json:"done" db:"done"`
	DueDate int  `json:"due_date"`
}
