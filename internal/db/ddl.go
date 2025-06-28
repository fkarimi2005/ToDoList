package db

import (
	"ToDoList/utils"
	"fmt"
	"os"
)

func InitMigrations() error {
	userTableQuery := `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    full_name  VARCHAR(255) NOT NULL,
    username   VARCHAR(255) NOT NULL UNIQUE,
    user_role varchar(255) NOT NULL default 'user',
    password   VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);`
	if _, err := db.Exec(userTableQuery); err != nil {
		return err
	}

	taskTableQuery := `
    CREATE TABLE IF NOT EXISTS tasks (
    task_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    description TEXT[] NOT NULL,
    done bool NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    due_date TIMESTAMP DEFAULT (CURRENT_TIMESTAMP + interval '1 day'),
    due_in_days integer NOT NULL,
    updated_at TIMESTAMP Default CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);`
	if _, err := db.Exec(taskTableQuery); err != nil {
		return err
	}

	if err := SeedUser(); err != nil {
		return err
	}
	return nil

}
func SeedUser() error {
	adminPassword := os.Getenv("Password_admin")
	if adminPassword == "" {
		return fmt.Errorf("Password_admin environment variable not set")
	}
	var exists bool
	err := GetDBConn().QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM users WHERE username = $1
		)
	`, "firuz7").Scan(&exists)

	if err != nil {
		return fmt.Errorf("could not check if user exists: %w", err)
	}

	if exists {
		return nil
	}

	password := utils.GenerateHash(adminPassword)
	_, err = GetDBConn().Exec(`
		INSERT INTO users (full_name, username, password, user_role)
		VALUES ($1, $2, $3, $4)
	`, "Firuz Karimzoda", "firuz7", password, "admin")
	if err != nil {
		return fmt.Errorf("could not insert user: %w", err)
	}

	return nil
}
