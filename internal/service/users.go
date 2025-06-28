package service

import (
	"ToDoList/internal/errs"
	"ToDoList/internal/model"
	"ToDoList/internal/repository"
	"ToDoList/utils"
	"errors"
	"fmt"
)

func CreateUser(u model.UserSignUp, userName string) error {
	_, err := repository.GetByUsername(u.Username)
	if err == nil {
		return errs.ErrUserAlreadyExists
	}

	u.Password = utils.GenerateHash(u.Password)

	if err = repository.CreateUser(u, userName); err != nil {
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
func GetUserByUserNameAndPassword(username string, password string) (user model.User, err error) {
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

func GetUserByUsername(username string, role string, currentUserID int) (model.User, error) {
	user, err := repository.GetUserByUsername(username, role, currentUserID)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return model.User{}, err
		}
	}
	return user, nil
}
func UpdateUserRole(user model.UserSignUp, role string, userID int) error {
	err := repository.UpdateUserRole(user, role, userID)
	if err != nil {
		return fmt.Errorf("could not update user role: %w", err)
	}
	return nil
}
