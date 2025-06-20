package service

import (
	model2 "ToDoList/internal/model"
	"ToDoList/internal/repository"
)

func ShowTask(role string, userID int) ([]model2.Tasks, error) {
	task, err := repository.ShowTasks(role, userID)
	if err != nil {
		return nil, err
	}
	return task, nil
}
func GetTaskByID(TaskID, userID int, role string) (model2.Tasks, error) {
	task, err := repository.GetTaskByID(TaskID, userID, role)
	if err != nil {
		return model2.Tasks{}, err
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

func CreateTask(t model2.Tasks, role string, userID int) error {
	if err := repository.CreateTask(t, role, userID); err != nil {
		return err
	}
	return nil
}
func UpdateTask(t model2.DoneTasks, tasksID, userID int, role string) error {
	if err := repository.UpdateTask(t, tasksID, userID, role); err != nil {
		return err
	}
	return nil

}
func SearchTask(search, role string, userID int) ([]model2.Tasks, error) {
	tasks, err := repository.SearchTask(search, role, userID)
	if err != nil {
		return nil, err
	}
	return tasks, nil

}
func GetTasksByUserID(RequestID, UserID int, role string) ([]model2.TaskWithUser, error) {

	tasks, err := repository.GetTasksByUserID(RequestID, UserID, role)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
