package repository

import (
	"ToDoList/db"
	"ToDoList/errs"
	"ToDoList/model"
	"ToDoList/utils"
)

func GetUserByUserNameAndPassword(username string, password string) (model.User, error) {
	var user model.User
	password = utils.GenerateHash(password)
	err := db.GetDBConn().Get(&user, `
        SELECT id, full_name, username, created_at, user_role
        FROM users
        WHERE deleted_at IS NULL AND username = $1 AND password = $2
    `, username, password)

	return user, err
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
	return user, err
}

func CreateUser(u model.User, role string) error {
	var err error
	if role == "admin" {
		_, err = db.GetDBConn().Exec(`INSERT INTO users (full_name, username, password, user_role)
	VALUES ($1, $2, $3, $4)`, u.FullName, u.Username, u.Password, u.UserRole)
	} else {
		return errs.ErrNotAccess
	}
	return err
}
func DeleteUser(userIDToDelete, requesterID int, role string) error {
	// Если это обычный пользователь — он может удалить только СЕБЯ
	if role != "admin" && userIDToDelete != requesterID {
		return errs.ErrNotAccess
	}

	// Удаление пользователя
	query := `DELETE FROM users WHERE id = $1`
	_, err := db.GetDBConn().Exec(query, userIDToDelete)
	if err != nil {
		return TranslateError(err)
	}

	return nil
}

func GetAllUsers(userID int, role string) ([]model.User, error) {
	var users []model.User
	var (
		err error
	)

	if role == "admin" {
		err = db.GetDBConn().Select(&users, `SELECT 
       id, 
      full_name,
      username,
      user_role,
      password, 
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
      password, 
      created_at, 
     updated_at
    from  users
    where id = $1
order by  id asc `, userID)
	}
	return users, err
}
func UpdateUser(user model.User, updateToID, userID int, role string) error {
	user.Password = utils.GenerateHash(user.Password)
	var err error
	if role == "admin" {
		_, err = db.GetDBConn().Exec(`
		UPDATE users
		SET full_name = $1,
		    username  = $2,
		    password  = $3,
		    user_role = $4,
		    updated_at= now()
		WHERE id = $5
	`, user.FullName, user.Username, user.Password, user.UserRole, updateToID)
	} else {
		if updateToID == user.ID {
			_, err = db.GetDBConn().Exec(`
		UPDATE users
		SET full_name = $1,
		    username  = $2,
		    password  = $3,
		    user_role = $4,
		    updated_at= now()
		WHERE id = $5 
	`, user.FullName, user.Username, user.Password, user.UserRole, updateToID)
		} else {
			return errs.ErrNotAccess
		}
	}

	return err
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

//func ChangeRoleUser(username, role string) error {
//	_, err := db.GetDBConn().Exec(`UPDATE users
//set user_role= role
//WHERE username = $1`, username)
//	return TranslateError(err)
//}
