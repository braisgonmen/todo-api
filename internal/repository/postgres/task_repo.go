package postgres

import (
	"context"
	"database/sql"
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
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.CreatedAt, &t.UserId); err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (db *DB) FindTaskByID(ctx context.Context, id int) (*model.Task, error) {

	var t model.Task

	err := db.conn.QueryRow("SELECT * FROM tasks WHERE id = $1", id).Scan(&t.ID, &t.Title, &t.Description, &t.CreatedAt, &t.UserId)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (db *DB) CreateTask(ctx context.Context, req model.CreateTaskRequest) (*model.Task, error) {

	var task model.Task
	err := db.conn.QueryRowContext(ctx,
		"INSERT INTO tasks (title, description, user_id) VALUES ($1, $2, $3) RETURNING id, title, description, created_at, user_id",
		req.Title, req.Description, req.UserId,
	).Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAt, &task.UserId)

	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (db *DB) UpdateTask(ctx context.Context, id int, req model.CreateTaskRequest) (*model.Task, error) {

	var updatedTask model.Task

	err := db.conn.QueryRowContext(ctx,
		"UPDATE tasks SET title = $1, description = $2, user_id = $3 WHERE id = $4 RETURNING id, title, description, created_at, user_id",
		req.Title, req.Description, req.UserId, id,
	).Scan(&updatedTask.ID, &updatedTask.Title, &updatedTask.Description, &updatedTask.CreatedAt, &updatedTask.UserId)

	if err != nil {
		return nil, err
	}

	return &updatedTask, nil
}

func (db *DB) DeleteTask(ctx context.Context, id int) error {

	result, err := db.conn.ExecContext(ctx,
		"DELETE FROM tasks WHERE id = $1",
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
