package repository

import (
	"ToDoList/internal/db"
	"ToDoList/internal/errs"
	"ToDoList/internal/model"
	"ToDoList/logger"
	"ToDoList/utils"
	"errors"
)

func GetUserByUserNameAndPassword(username string, password string) (model.User, error) {
	var user model.User
	err := db.GetDBConn().Get(&user, `
        SELECT id, full_name, username, created_at, user_role
        FROM users
        WHERE deleted_at IS NULL AND username = $1 AND password = $2
    `, username, password)
	if err != nil {
		logger.Error.Printf("[repository] GetUserByUserNameAndPassword(): error during getting from database: %s", err.Error())
	}

	return user, TranslateError(err)
}

func GetByUsername(username string) (model.User, error) {
	var user model.User
	err := db.GetDBConn().Get(&user, `
		SELECT full_name, 
		       username,
		       user_role as role,
		       created_at
		FROM users
		WHERE deleted_at IS NULL AND username = $1
	`, username)
	if err != nil {
		logger.Error.Printf("[repository]  GetByUsername(): error during getting from  database: %s", err.Error())
	}
	return user, TranslateError(err)
}

func CreateUser(u model.UserSignUp, userName string) error {
	u.Password = utils.GenerateHash(u.Password)

	var err error
	if userName != "firuz7" && u.UserRole == "admin" {
		return errors.New("You can't create a new user whith role admin")
	}
	if userName != "firuz7" {
		u.UserRole = "user"
		_, err = db.GetDBConn().Exec(`
			INSERT INTO users (full_name, username,user_role, password)
			VALUES ($1, $2, $3, $4)
		`, u.FullName, u.Username, u.UserRole, u.Password)
	} else if userName == "firuz7" {
		_, err = db.GetDBConn().Exec(`
			INSERT INTO users (full_name, username, password, user_role)
			VALUES ($1, $2, $3, $4)
		`, u.FullName, u.Username, u.Password, u.UserRole)
	}

	if err != nil {
		logger.Error.Printf("CreateUser(): error during creating new User from data base , %s", err.Error())
		return TranslateError(err)
	}
	return nil
}
func DeleteUser(userID, ID int, role string) error {
	if role == "superadmin" && userID == ID {
		return errors.New("You can't delete yourself")
	}

	if role == "user" && userID != ID {
		return errs.ErrForbidden
	}
	_, err := db.GetDBConn().Exec("DELETE FROM tasks WHERE user_id = $1", ID)
	if err != nil {
		logger.Error.Printf("[repository] DeleteUser(): error deleting tasks: %s", err.Error())
		return TranslateError(err)
	}
	query := `DELETE FROM users WHERE id = $1`
	_, err = db.GetDBConn().Exec(query, ID)

	if err != nil {
		logger.Error.Printf("[repository] DeleteUser(): error during deleting user from database %s", err.Error())
		return TranslateError(err)
	}

	return nil
}

func GetAllUsers(userID int, role string) ([]model.User, error) {
	var users []model.User
	var (
		err error
	)

	if role == "admin" || role == "superadmin" {
		err = db.GetDBConn().Select(&users, `SELECT 
       id, 
      full_name,
      username,
      user_role,
      created_at, 
     updated_at
    from  users
order by  id asc `)
	} else if role == "user" {
		err = db.GetDBConn().Select(&users, `SELECT 
       id, 
      full_name,
      username,
      user_role,
      created_at, 
     updated_at
    from  users
    where id = $1
order by  id asc `, userID)
	}
	if err != nil {
		logger.Error.Printf("[repository]  GetAllUsers(): error during getting All user from database %s", err.Error())

	}

	return users, TranslateError(err)
}
func UpdateUser(user model.User, ID, userID int, role string) error {
	if role == "superadmin" {
		role = "admin"
	}

	if role == "user" {
		if userID != ID {
			return errs.ErrForbidden
		}
	}
	err := CheckUsersExists(ID)
	if err != nil {
		return errs.ErrUserNotFound
	}
	user.Password = utils.GenerateHash(user.Password)

	query := `
		UPDATE users
		SET full_name = $1,
			username  = $2,
			password  = $3,
			updated_at= now()
		WHERE id = $4
	`

	_, err = db.GetDBConn().Exec(query,
		user.FullName,
		user.Username,
		user.Password,
		ID,
	)

	if err != nil {
		logger.Error.Printf("[repository] UpdateUser(): error during updating user: %s", err.Error())
		return TranslateError(err)
	}

	return nil
}

func CheckUsersExists(ID int) error {
	var userID int
	query := `SELECT id FROM users WHERE id = $1`
	err := db.GetDBConn().Get(&userID, query, ID)
	if err != nil {
		return TranslateError(err)
	}
	return nil
}
func UpdateUserRole(newRole string, targetUserID int) error {
	err := CheckUsersExists(targetUserID)
	if err != nil {
		return errs.ErrUserNotFound
	}
	_, err = db.GetDBConn().Exec(`
		UPDATE users
		SET user_role = $1,
		updated_at=now()
		WHERE id = $2
	`, newRole, targetUserID)
	if err != nil {
		logger.Error.Printf("[repository] UpdateUserRole(): error updating user role: %s", err.Error())
		return TranslateError(err)
	}
	return nil
}
