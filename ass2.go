package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "flora.db.elephantsql.com"
	port     = 5432
	user     = "muolkfz" // User & Default database
	password = "1234"
	dbname   = "ymuolkfz" // User & Default database
)

type Task struct {
	ID        int
	Name      string
	Completed bool
}

func main() {
	// Connect to the PostgreSQL database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new task
	taskName := "Sample Task"
	err = createTask(db, taskName)
	if err != nil {
		log.Fatal(err)
	}

	// List all tasks
	tasks, err := listTasks(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Tasks:")
	for _, task := range tasks {
		fmt.Printf("- %s (ID: %d, Completed: %t)\n", task.Name, task.ID, task.Completed)
	}

	// Update a task (mark as completed)
	taskID := 1 // Assuming the task ID to update
	err = updateTask(db, taskID, true)
	if err != nil {
		log.Fatal(err)
	}

	// Delete a task
	err = deleteTask(db, taskID)
	if err != nil {
		log.Fatal(err)
	}
}

func createTask(db *sql.DB, name string) error {
	_, err := db.Exec("INSERT INTO tasks (name) VALUES ($1)", name)
	return err
}

func listTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query("SELECT id, name, completed FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Completed); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func updateTask(db *sql.DB, id int, completed bool) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec("UPDATE tasks SET completed = $1 WHERE id = $2", completed, id)
	return err
}

func deleteTask(db *sql.DB, id int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec("DELETE FROM tasks WHERE id = $1", id)
	return err
}
