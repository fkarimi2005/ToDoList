package service

import (
	model "ToDoList/internal/model"
	"ToDoList/internal/repository"
)

func ShowTask(role string, userID int) ([]model.Tasks, error) {
	task, err := repository.ShowTasks(role, userID)
	if err != nil {
		return nil, err
	}
	return task, nil
}
func GetTaskByID(TaskID, userID int, role string) (model.Tasks, error) {
	task, err := repository.GetTaskByID(TaskID, userID, role)
	if err != nil {
		return model.Tasks{}, err
	}
	return task, nil
}
func DeleteTask(TaskID, userID int, role string) error {
	err := repository.CheckTaskExists(TaskID)
	if err != nil {
		return err
	}
	err = repository.DeleteTask(TaskID, userID, role)
	if err != nil {
		return err
	}
	return nil
}

func CreateTask(t model.Tasks, role string, userID int) error {
	if err := repository.CreateTask(t, role, userID); err != nil {
		return err
	}
	return nil
}
func UpdateTask(t model.DoneTasks, tasksID, userID int, role string) error {
	if err := repository.UpdateTask(t, tasksID, userID, role); err != nil {
		return err
	}
	return nil

}
func SearchTask(search, role string, userID int) ([]model.Tasks, error) {
	tasks, err := repository.SearchTask(search, role, userID)
	if err != nil {
		return nil, err
	}
	return tasks, nil

}
func GetTasksByUserID(RequestID, UserID int, role string) ([]model.TaskWithUser, error) {

	tasks, err := repository.GetTasksByUserID(RequestID, UserID, role)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
func GetCompletedTasks(role string, userID int) ([]model.Tasks, error) {
	task, err := repository.GetCompletedTasks(role, userID)
	if err != nil {
		return nil, err
	}
	return task, nil
}
func GetInCompletedTasks(role string, userID int) ([]model.Tasks, error) {
	task, err := repository.GetInCompletedTasks(role, userID)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func GetPendingTasks(role string, userID int) ([]model.Tasks, error) {
	task, err := repository.GetPendingTasks(role, userID)
	if err != nil {
		return nil, err
	}
	return task, nil
}
func GetTaskByPriority(role string, UserID int) ([]model.Tasks, error) {
	tasks, err := repository.GetTasksByPriority(role, UserID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
