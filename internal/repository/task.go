package repository

import (
	"ToDoList/internal/db"
	"ToDoList/internal/errs"
	model "ToDoList/internal/model"
	"ToDoList/logger"
	"errors"
)

func ShowTasks(role string, userID int) ([]model.Tasks, error) {
	var (
		tasks []model.Tasks
		err   error
		query string
	)

	if role == "admin" {
		query = `
			SELECT 
				task_id,
				user_id,
				title, 
				description, 
				created_at,
				due_date, 
				updated_at,
				done
			FROM tasks
			WHERE deleted_at IS NULL
			ORDER BY task_id ASC
		`
		err = db.GetDBConn().Select(&tasks, query)
	} else {
		query = `
			SELECT 
				task_id,
				user_id,
				title, 
				description, 
				created_at,
				due_date, 
				updated_at,
				done
			FROM tasks
			WHERE deleted_at IS NULL AND user_id = $1
			ORDER BY task_id ASC
		`
		err = db.GetDBConn().Select(&tasks, query, userID)
	}

	if err != nil {
		logger.Error.Printf("[repository] ShowTasks(): error during getting all tasks from database: %s", err.Error())
		return nil, TranslateError(err)
	}

	return tasks, nil
}

func GetTaskByID(TaskID, userID int, role string) (model.Tasks, error) {
	task := model.Tasks{}
	query := `
		SELECT
			task_id,
			user_id,
			title,
			description,
			created_at,
			updated_at,
			due_date,
			done
		FROM tasks
		WHERE deleted_at IS NULL AND task_id = $1
	`
	err := db.GetDBConn().Get(&task, query, TaskID)
	if err != nil {
		logger.Error.Printf("[repository] GetTaskByID(): error during getting task by ID from database: %s", err.Error())
		return model.Tasks{}, TranslateError(err)
	}
	if role != "admin" && task.User_ID != userID {
		return model.Tasks{}, errs.ErrNotAccess
	}
	return task, nil
}

func DeleteTask(taskID, userID int, role string) error {
	var (
		query string
		err   error
	)

	if role != "admin" {
		query = `DELETE FROM tasks WHERE task_id = $1 AND user_id = $2`
		_, err = db.GetDBConn().Exec(query, taskID, userID)
	} else {
		query = `DELETE FROM tasks WHERE task_id = $1`
		_, err = db.GetDBConn().Exec(query, taskID)
	}

	if err != nil {
		logger.Error.Printf("[repository] DeleteTask(): error during deleting task from database: %s", err.Error())
		return TranslateError(err)
	}

	return nil
}

func GetUserByUsername(username, role string, currentUserID int) (model.User, error) {
	var user model.User

	query := `
		SELECT id, 
			   full_name, 
			   username, 
			   created_at
		FROM users 
		WHERE deleted_at IS NULL AND username = $1
	`

	err := db.GetDBConn().Get(&user, query, username)
	if err != nil {
		logger.Error.Printf("[repository] GetUserByUsername(): error during getting user by username from database: %s", err.Error())
		return model.User{}, TranslateError(err)
	}

	if role != "admin" && user.ID != currentUserID {
		return model.User{}, errs.ErrNotAccess
	}

	return user, nil
}

func CreateTask(t model.Tasks, role string, userID int) error {
	if role != "admin" && t.User_ID != userID {
		return errs.ErrNotAccess
	}

	query := `
		INSERT INTO tasks (user_id, title, description, done)
		VALUES ($1, $2, $3, $4)
	`

	_, err := db.GetDBConn().Exec(query, t.User_ID, t.Title, t.Description, t.Done)
	if err != nil {
		logger.Error.Printf("[repository] CreateTask(): error during creating task: %s", err.Error())
		return TranslateError(err)
	}

	return nil
}

func UpdateTask(d model.DoneTasks, taskID int, userID int, role string) error {
	var (
		query string
		err   error
	)

	if role != "admin" {
		query = `
			UPDATE tasks
			SET done = $1,
			    due_date = current_timestamp + make_interval(days => $2),
			    updated_at = NOW()
			WHERE task_id = $3 AND user_id = $4
		`
		_, err = db.GetDBConn().Exec(query, d.Done, d.DueDate, taskID, userID)
	} else {
		query = `
			UPDATE tasks
			SET done = $1,
			    due_date = current_timestamp + make_interval(days => $2),
			    updated_at = NOW()
			WHERE task_id = $3
		`
		_, err = db.GetDBConn().Exec(query, d.Done, d.DueDate, taskID)
	}

	if err != nil {
		logger.Error.Printf("[repository] UpdateTask(): error during updating task in database: %s", err.Error())
		return TranslateError(err)
	}

	return nil
}

func SearchTask(search, role string, userID int) ([]model.Tasks, error) {
	var tasks []model.Tasks
	var (
		query string
		args  []interface{}
		err   error
	)

	if role == "admin" {
		query = `
			SELECT 
				task_id,
				user_id,
				title, 
				description, 
				created_at,
				due_date, 
				done
			FROM tasks
			WHERE deleted_at IS NULL
			  AND (title ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%')
			ORDER BY task_id ASC
		`
		args = append(args, search)
	} else if role == "user" {
		query = `
			SELECT 
				task_id,
				user_id,
				title, 
				description, 
				created_at,
				due_date, 
				done
			FROM tasks
			WHERE deleted_at IS NULL AND user_id = $1
			  AND (title ILIKE '%' || $2 || '%' OR description ILIKE '%' || $2 || '%')
			ORDER BY task_id ASC
		`
		args = append(args, userID, search)
	} else {
		return nil, errors.New("invalid role")
	}

	err = db.GetDBConn().Select(&tasks, query, args...)
	if err != nil {
		logger.Error.Printf("[repository] SearchTask(): error during searching tasks from database: %s", err.Error())
		return nil, TranslateError(err)
	}

	return tasks, nil
}

func GetTasksByUserID(requestedUserID, currentUserID int, role string) ([]model.TaskWithUser, error) {
	if role != "admin" && requestedUserID != currentUserID {
		return nil, errs.ErrNotAccess
	}

	query := `
		SELECT 
			tasks.task_id,
			tasks.user_id,
			tasks.title,
			tasks.description,
			users.username
		FROM tasks
		JOIN users ON tasks.user_id = users.id
		WHERE tasks.deleted_at IS NULL AND tasks.user_id = $1
	`

	var tasks []model.TaskWithUser
	err := db.GetDBConn().Select(&tasks, query, requestedUserID)
	if err != nil {
		return nil, TranslateError(err)
	}

	return tasks, nil
}

func CheckTaskExists(ID int) error {
	var taskID int
	query := `SELECT task_id FROM tasks WHERE task_id = $1`
	err := db.GetDBConn().Get(&taskID, query, ID)
	if err != nil {
		return TranslateError(err) // Вернёт ErrNotFoundID, если нет
	}
	return nil
}
