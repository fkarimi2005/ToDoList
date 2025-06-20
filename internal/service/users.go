package service

import (
	"ToDoList/internal/errs"
	"ToDoList/internal/model"
	"ToDoList/internal/repository"
	"ToDoList/utils"

	"errors"
	"fmt"
)

func CreateUser(u model.User, role string) error {

	_, err := repository.GetByUsername(u.Username)
	if err == nil {

		return errs.ErrUserAlreadyExists
	}
	u.Password = utils.GenerateHash(u.Password)

	if err := repository.CreateUser(u, role); err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}

	return nil
}
func GetAllUsers(userID int, role string) ([]model.User, error) {
	users, err := repository.GetAllUsers(userID, role)
	if err != nil {
		return nil, fmt.Errorf("could not get all users: %w", err)
	}
	return users, nil
}
func UpdateUser(u model.User, upID, userID int, role string) error {

	if err := repository.UpdateUser(u, upID, userID, role); err != nil {
		return fmt.Errorf("could not update user: %w", err)
	}

	return nil
}
func DeleteUsers(UserID, requestID int, role string) error {
	err := repository.CheckUsersExists(UserID)
	if err != nil {
		return err
	}
	err = repository.DeleteUser(UserID, requestID, role)
	if err != nil {
		return err
	}
	return nil
}
func GetUserByUsernameAndPassword(username string, password string) (user model.User, err error) {
	password = utils.GenerateHash(password)

	user, err = repository.GetUserByUserNameAndPassword(username, password)
	if err != nil {
		if errors.Is(err, errs.ErrNotFoud) {
			return model.User{}, err
		}
		return model.User{}, err
	}
	return user, nil
}

//func ChangeRoleUser(username, role string) error {
//	_, err := repository.GetByUsername(username)
//	if err != nil {
//		return err
//	}
//	err = repository.ChangeRoleUser(username, role)
//	if err != nil {
//		return err
//	}
//	return nil
//
//}
