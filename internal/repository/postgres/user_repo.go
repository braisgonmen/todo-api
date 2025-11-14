package postgres

import (
	"context"
	"log"
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

func (db *DB) FindByEmail(ctx context.Context, email string) (*model.User, error) {

	var u model.User

	log.Printf("FindByEmail(repo): querying email=%s", email)
	err := db.conn.QueryRow("SELECT * FROM users WHERE email = $1", email).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	if err != nil {
		log.Printf("FindByEmail(repo): email=%s error=%v", email, err)
		return nil, err
	}

	return &u, nil
}

func (db *DB) FindUserByID(ctx context.Context, id int) (*model.User, error) {
	var u model.User

	log.Printf("FindUserByID(repo): querying id=%d", id)
	err := db.conn.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	if err != nil {
		log.Printf("FindUserByID(repo): id=%d error=%v", id, err)
		return nil, err
	}

	return &u, nil
}
