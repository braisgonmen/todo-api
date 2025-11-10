package postgres

import (
	"context"
	"todo-api/internal/model"
)

func (db *DB) GetAllTask(ctx context.Context) ([]model.Task, error) {

	rows, err := db.conn.QueryContext(ctx, "SELECT * FROM tasks")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []model.Task

	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.CreatedAt); err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (db *DB) GetTaskByID(ctx context.Context, id int) (*model.Task, error) {

	var t model.Task

	err := db.conn.QueryRow("SELECT * FROM tasks WHERE id = $1", id).Scan(&t.ID, &t.Title, &t.Description, &t.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (db *DB) CreateTask(ctx context.Context, req model.CreateTaskRequest) (*model.Task, error) {

	var task model.Task
	err := db.conn.QueryRowContext(ctx,
		"INSERT INTO tasks (title, description) VALUES ($1, $2) RETURNING id, title, email, description, user_id",
		req.Title, req.Description,
	).Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAt, &task.UserId)

	if err != nil {
		return nil, err
	}
	return &task, nil
}
