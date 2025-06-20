package model

type TaskWithUser struct {
	TaskID      int    `db:"task_id"`
	UserID      int    `db:"user_id"`
	Username    string `db:"username"`
	Title       string `db:"title"`
	Description string `db:"description"`
}
