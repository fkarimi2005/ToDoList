package repository

import (
	"ToDoList/errs"
	"database/sql"
	"errors"
)

func TranslateError(err error) error {
	if err == nil {
		return nil
	} else if errors.Is(err, sql.ErrNoRows) {
		return errs.ErrNotFoundID
	} else {
		return err
	}
}
