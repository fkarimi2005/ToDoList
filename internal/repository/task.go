package repository

import (
	"ToDoList/internal/db"
	"ToDoList/internal/errs"
	model "ToDoList/internal/model"
	"ToDoList/logger"
	"encoding/json"
	"errors"
	"time"
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
				due_date, 
				priority,
				done,
				created_at,
				updated_at
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
				priority,
				due_date, 
				updated_at,
				done
			FROM tasks
			WHERE deleted_at IS NULL AND user_id = $1
			ORDER BY task_id ASC
		`
		err = db.GetDBConn().Select(&tasks, query, userID)
	}
	for i := range tasks {
		err := json.Unmarshal(tasks[i].DescriptionRaw, &tasks[i].Description)
		if err != nil {
			logger.Error.Printf("[repository] ShowTasks(): error unmarshaling description JSON: %s", err.Error())
			return nil, err
		}
	}
	if tasks == nil {
		return []model.Tasks{}, errors.New("no tasks found")
	}
	return tasks, err
}

func GetCompletedTasks(role string, userID int) ([]model.Tasks, error) {
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
				priority,
				due_date, 
				updated_at,
				done
			FROM tasks
			WHERE deleted_at IS NULL
			AND done=true
			order by 
			CASE LOWER(priority)
					WHEN 'high' THEN 1
					WHEN 'medium' THEN 2
					WHEN 'low' THEN 3
					ELSE 4
				END ASC,
				due_date ASC
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
				priority,
				due_date, 
				updated_at,
				done
			FROM tasks
			WHERE deleted_at IS NULL AND user_id = $1 and done=true
			order by 
			CASE LOWER(priority)
					WHEN 'high' THEN 1
					WHEN 'medium' THEN 2
					WHEN 'low' THEN 3
					ELSE 4
				END ASC,
				due_date ASC
		`
		err = db.GetDBConn().Select(&tasks, query, userID)
	}

	for i := range tasks {
		err = json.Unmarshal(tasks[i].DescriptionRaw, &tasks[i].Description)
		if err != nil {
			logger.Error.Printf("[repository] GetCompletedTasks(): error unmarshaling description JSON: %s", err.Error())
			return nil, err
		}
	}
	if tasks == nil {
		return []model.Tasks{}, errs.ErrTaskNotFound
	}

	return tasks, nil
}

func GetInCompletedTasks(role string, userID int) ([]model.Tasks, error) {
	var (
		tasks []model.Tasks
		err   error
		query string
	)

	if role == "admin" || role == "superadmin" {
		query = `
			SELECT 
				task_id,
				user_id,
				title, 
				description, 
				created_at,
				priority,
				due_date, 
				updated_at,
				done
			FROM tasks
			WHERE deleted_at IS NULL
			AND due_date < Now() and done = false
			order by 
			CASE LOWER(priority)
					WHEN 'high' THEN 1
					WHEN 'medium' THEN 2
					WHEN 'low' THEN 3
					ELSE 4
				END ASC,
				due_date ASC
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
				priority,
				due_date, 
				updated_at,
				done
			FROM tasks
			WHERE deleted_at IS NULL AND user_id = $1 
			and due_date < Now() and done=false
			order by 
			CASE LOWER(priority)
					WHEN 'high' THEN 1
					WHEN 'medium' THEN 2
					WHEN 'low' THEN 3
					ELSE 4
				END ASC,
				due_date ASC
		`
		err = db.GetDBConn().Select(&tasks, query, userID)
	}

	for i := range tasks {
		err = json.Unmarshal(tasks[i].DescriptionRaw, &tasks[i].Description)
		if err != nil {
			logger.Error.Printf("[repository] GetInCompletedTasks(): error unmarshaling description JSON: %s", err.Error())
			return nil, err
		}
	}
	if tasks == nil {
		return []model.Tasks{}, errs.ErrTaskNotFound
	}

	return tasks, nil
}
func GetPendingTasks(role string, userID int) ([]model.Tasks, error) {
	var (
		tasks []model.Tasks
		err   error
		query string
	)

	if role == "admin" || role == "superadmin" {
		query = `
			SELECT 
				task_id,
				user_id,
				title, 
				description, 
				created_at,
				priority,
				due_date, 
				updated_at,
				done
			FROM tasks
			WHERE deleted_at IS NULL
			  AND due_date > now()
			  AND done = false
			  order by 
			CASE LOWER(priority)
					WHEN 'high' THEN 1
					WHEN 'medium' THEN 2
					WHEN 'low' THEN 3
					ELSE 4
				END ASC,
				due_date ASC
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
				priority,
				due_date, 
				updated_at,
				done
			FROM tasks
			WHERE deleted_at IS NULL
			  AND user_id = $1
			  AND due_date > now()
			  AND done = false
			  order by 
			CASE LOWER(priority)
					WHEN 'high' THEN 1
					WHEN 'medium' THEN 2
					WHEN 'low' THEN 3
					ELSE 4
				END ASC,
				due_date ASC
		`
		err = db.GetDBConn().Select(&tasks, query, userID)
	}
	if err != nil {
		logger.Error.Printf("[repository] GetPendingTasks(): error selecting tasks: %s", err.Error())
		return nil, TranslateError(err)
	}

	for i := range tasks {
		if err = json.Unmarshal(tasks[i].DescriptionRaw, &tasks[i].Description); err != nil {
			logger.Error.Printf("[repository] GetPendingTasks(): error unmarshaling description JSON: %s", err.Error())
			return nil, err
		}
	}
	if tasks == nil {
		return []model.Tasks{}, errs.ErrTaskNotFound
	}

	return tasks, nil
}
func GetTaskByID(taskID, userID int, role string) (model.Tasks, error) {
	task := model.Tasks{}
	var (
		query string
		args  []any
	)

	if role == "user" {
		query = `
			SELECT
				task_id,
				user_id,
				title,
				description,
				created_at,
				priority,
				updated_at,
				due_date,
				done
			FROM tasks
			WHERE deleted_at IS NULL
			  AND task_id = $1
			  AND user_id = $2
			ORDER BY 
				CASE LOWER(priority)
					WHEN 'high' THEN 1
					WHEN 'medium' THEN 2
					WHEN 'low' THEN 3
					ELSE 4
				END ASC,
				due_date ASC
		`
		args = []any{taskID, userID}
	} else if role == "admin" || role == "superadmin" {
		query = `
			SELECT
				task_id,
				user_id,
				title,
				description,
				created_at,
				priority,
				updated_at,
				due_date,
				done
			FROM tasks
			WHERE deleted_at IS NULL
			  AND task_id = $1
			ORDER BY 
				CASE LOWER(priority)
					WHEN 'high' THEN 1
					WHEN 'medium' THEN 2
					WHEN 'low' THEN 3
					ELSE 4
				END ASC,
				due_date ASC
		`
		args = []any{taskID}
	} else {
		return model.Tasks{}, errs.ErrForbidden
	}

	err := db.GetDBConn().Get(&task, query, args...)
	if err != nil {
		logger.Error.Printf("[repository] GetTaskByID(): error fetching task: %s", err.Error())
		return model.Tasks{}, TranslateError(err)
	}

	if len(task.DescriptionRaw) > 0 {
		err = json.Unmarshal(task.DescriptionRaw, &task.Description)
		if err != nil {
			logger.Error.Printf("[repository] GetTaskByID(): error unmarshaling description JSON: %s", err.Error())
			return model.Tasks{}, err
		}
	}

	return task, nil
}

func DeleteTask(taskID, userID int, role string) error {
	var (
		query string
		args  []any
	)

	if role == "superadmin" {
		role = "admin"
	}

	if role == "user" {
		query = `DELETE FROM tasks WHERE task_id = $1 AND user_id = $2`
		args = []any{taskID, userID}
	} else if role == "admin" {
		query = `DELETE FROM tasks WHERE task_id = $1`
		args = []any{taskID}
	} else {
		return errs.ErrForbidden
	}

	res, err := db.GetDBConn().Exec(query, args...)
	if err != nil {
		logger.Error.Printf("[repository] DeleteTask(): error during deleting task from database: %s", err.Error())
		return TranslateError(err)
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errs.ErrNotFoud
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
	if role == "superadmin" {
		role = "admin"
	}
	if role != "admin" && user.ID != currentUserID {
		return model.User{}, errs.ErrNotAccess
	}

	return user, nil
}

func CreateTask(t model.Tasks, role string, userID int) error {
	if role == "superadmin" {
		role = "admin"
	}

	err := CheckUsersExists(t.User_ID)
	if err != nil {
		return errs.ErrUserNotFound
	}
	var dueInDays int
	if t.DueInDays != nil {
		dueInDays = *t.DueInDays
	} else {
		dueInDays = 1
	}

	t.DueDate = time.Now().AddDate(0, 0, dueInDays)

	var descJSON []byte
	if t.Description != nil {
		var err error
		descJSON, err = json.Marshal(t.Description)
		if err != nil {
			return err
		}
	} else {
		descJSON = []byte("null")
	}

	query := `
        INSERT INTO tasks (user_id, title, description, done, due_date, due_in_days, priority)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	if role == "admin" {
		_, err = db.GetDBConn().Exec(query,
			t.User_ID,
			t.Title,
			descJSON,
			t.Done,
			t.DueDate,
			dueInDays,
			t.Priority,
		)
	} else if role == "user" {
		_, err = db.GetDBConn().Exec(query,
			userID,
			t.Title,
			descJSON,
			t.Done,
			t.DueDate,
			dueInDays,
			t.Priority,
		)
	}

	if err != nil {
		logger.Error.Printf("[repository] CreateTask(): error during creating task: %s", err.Error())
		return TranslateError(err)
	}

	return nil
}

func UpdateTask(d model.DoneTasks, taskID int, userID int, role string) error {
	var (
		query string
		args  []any
	)

	if role == "superadmin" {
		role = "admin"
	}

	if role == "user" {
		query = `
			UPDATE tasks
			SET done = $1,
			    due_date = current_timestamp + make_interval(days => $2),
			    updated_at = NOW()
			WHERE task_id = $3 AND user_id = $4
		`
		args = []any{d.Done, d.DueDate, taskID, userID}
	} else if role == "admin" {
		query = `
			UPDATE tasks
			SET done = $1,
			    due_date = current_timestamp + make_interval(days => $2),
			    updated_at = NOW()
			WHERE task_id = $3
		`
		args = []any{d.Done, d.DueDate, taskID}
	} else {
		return errs.ErrForbidden
	}

	res, err := db.GetDBConn().Exec(query, args...)
	if err != nil {
		logger.Error.Printf("[repository] UpdateTask(): error during updating task in database: %s", err.Error())
		return TranslateError(err)
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errs.ErrNotFoud
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

	if role == "admin" || role == "superadmin" {
		query = `
			SELECT 
				task_id,
				user_id,
				title, 
				description, 
				created_at,
				priority,
				due_date, 
				done
			FROM tasks
			WHERE deleted_at IS NULL
			  AND (title ILIKE '%' || $1 || '%' OR description::text ILIKE '%' || $1 || '%')
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
			  AND (title ILIKE '%' || $2 || '%' OR description::text ILIKE '%' || $2 || '%')
			ORDER BY task_id ASC
		`
		args = append(args, userID, search)
	} else {
		return nil, errors.New("invalid role")
	}
	err = db.GetDBConn().Select(&tasks, query, args...)
	for i := range tasks {
		if err = json.Unmarshal(tasks[i].DescriptionRaw, &tasks[i].Description); err != nil {
			logger.Error.Printf("[repository] SearchTask(): error unmarshaling description JSON: %s", err.Error())
			return nil, err
		}
	}

	return tasks, nil
}

func GetTasksByUserID(requestedUserID, currentUserID int, role string) ([]model.TaskWithUser, error) {
	var err error
	var tasks []model.TaskWithUser
	if role != "admin" && requestedUserID != currentUserID {
		return nil, errs.ErrForbidden
	}
	if role == "admin" || role == "superadmin" {
		query := `
		SELECT 
			tasks.task_id,
			tasks.user_id,
			tasks.title,
			tasks.description,
			users.username,
		    tasks.done
		FROM tasks
		JOIN users ON tasks.user_id = users.id
		WHERE tasks.deleted_at IS NULL AND tasks.user_id = $1
	`

		err = db.GetDBConn().Select(&tasks, query, requestedUserID)
	} else {
		query := `
		SELECT 
			tasks.task_id,
			tasks.user_id,
			tasks.title,
			tasks.description,
			users.username,
		    tasks.done
		FROM tasks
		JOIN users ON tasks.user_id = users.id
		WHERE tasks.deleted_at IS NULL AND tasks.user_id = $1 and  tasks.user_id = $2
	`

		err = db.GetDBConn().Select(&tasks, query, requestedUserID, currentUserID)
	}

	for i := range tasks {
		err = json.Unmarshal(tasks[i].DescriptionRaw, &tasks[i].Description)
		if err != nil {
			logger.Error.Printf("[repository] GetTasksByUserID(): error unmarshaling description JSON: %s", err.Error())
			return nil, err
		}
	}

	return tasks, nil
}
func GetTasksByPriority(role string, userID int) ([]model.Tasks, error) {
	var (
		tasks []model.Tasks
		err   error
		query string
	)

	if role == "admin" || role == "superadmin" {
		query = `
			SELECT 
				task_id,
				user_id,
				title, 
				description, 
				due_date, 
				priority,
				done,
				created_at,
				updated_at
			FROM tasks
			WHERE deleted_at IS NULL AND done = false
			ORDER BY 
				CASE LOWER(priority)
					WHEN 'high' THEN 1
					WHEN 'medium' THEN 2
					ELSE 3
				END ASC,
				due_date ASC
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
				priority,
				due_date, 
				updated_at,
				done
			FROM tasks
			WHERE deleted_at IS NULL AND user_id = $1 AND done = false
			ORDER BY 
				CASE LOWER(priority)
					WHEN 'high' THEN 1
					WHEN 'medium' THEN 2
					ELSE 3
				END ASC,
				due_date ASC
		`
		err = db.GetDBConn().Select(&tasks, query, userID)
	}

	if err != nil {
		logger.Error.Printf("[repository] GetTasksByPriority(): error selecting tasks: %s", err.Error())
		return nil, err
	}

	for i := range tasks {
		if err := json.Unmarshal(tasks[i].DescriptionRaw, &tasks[i].Description); err != nil {
			logger.Error.Printf("[repository] GetTasksByPriority(): error unmarshaling description JSON: %s", err.Error())
			return nil, err
		}
	}
	if tasks == nil {
		return []model.Tasks{}, errors.New("no tasks found")
	}

	return tasks, nil
}

func CheckTaskExists(ID int) error {
	var taskID int
	query := `SELECT task_id FROM tasks WHERE task_id = $1`
	err := db.GetDBConn().Get(&taskID, query, ID)
	if err != nil {
		return TranslateError(err)
	}
	return nil
}
