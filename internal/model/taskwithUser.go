package model

import "encoding/json"

type TaskWithUser struct {
	TaskID         int             `db:"task_id" json:"task_id"`
	UserID         int             `db:"user_id" json:"user_id"`
	Username       string          `db:"username" json:"username"`
	Title          string          `db:"title" json:"title"`
	Description    []string        `db:"-" json:"description"`
	DescriptionRaw json.RawMessage `db:"description" json:"-"`
	DoneTasks      bool            `db:"done" json:"done_tasks"`
}
