package postgres

import (
	"context"
	"todo-api/internal/model"
)

func (db *DB) GetUsers(ctx context.Context) ([]model.User, error) {

	rows, err := db.conn.QueryContext(ctx, "SELECT * FROM users")

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var users []model.User

	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (db *DB) CreateUser(ctx context.Context, req model.CreateUserRequest) (*model.User, error) {

	var user model.User
	err := db.conn.QueryRowContext(ctx,
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name, email, created_at",
		req.Name, req.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
