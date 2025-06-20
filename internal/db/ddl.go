package db

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
    description TEXT NOT NULL,
    done bool NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    due_date TIMESTAMP DEFAULT (CURRENT_TIMESTAMP + interval '1 day'),
    updated_at TIMESTAMP Default CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);`
	if _, err := db.Exec(taskTableQuery); err != nil {
		return err
	}
	//if err := Seed(); err != nil {
	//	return err
	//}

	return nil
}

//func Seed() error {
//	password := utils.GenerateHash("mypassword")
//	_, err := db.Exec(
//		`INSERT INTO users (full_name, username, password)
//         VALUES ($1, $2, $3)`,
//		"Salmon Fors",
//		"salmon7745",
//		password,
//	)
//	return err
//}
