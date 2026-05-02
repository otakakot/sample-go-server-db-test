package gateway

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/otakakot/sample-go-server-db-test/internal/domain"
)

type Gateway struct {
	db *sql.DB
}

func New(
	db *sql.DB,
) *Gateway {
	return &Gateway{
		db: db,
	}
}

type CreateUserDAI struct {
	Name string
}

type CreateUserDAO struct {
	User domain.User
}

func (gw *Gateway) CreateUser(
	ctx context.Context,
	input CreateUserDAI,
) (*CreateUserDAO, error) {
	query := `INSERT INTO users (id, name) VALUES (?, ?)`

	id := uuid.New().String()

	if _, err := gw.db.ExecContext(ctx, query, id, input.Name); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &CreateUserDAO{
		User: domain.User{
			ID:   id,
			Name: input.Name,
		},
	}, nil
}

type ReadUserDAI struct {
	ID string
}

type ReadUserDAO struct {
	User domain.User
}

func (gw *Gateway) ReadUser(
	ctx context.Context,
	input ReadUserDAI,
) (*ReadUserDAO, error) {
	query := `SELECT id, name FROM users WHERE id = ?`

	var user domain.User

	if err := gw.db.QueryRowContext(ctx, query, input.ID).Scan(&user.ID, &user.Name); err != nil {
		return nil, fmt.Errorf("failed to read user: %w", err)
	}

	return &ReadUserDAO{
		User: user,
	}, nil
}

type UpdateUserDAI struct {
	ID   string
	Name string
}

type UpdateUserDAO struct{}

func (gw *Gateway) UpdateUser(
	ctx context.Context,
	input UpdateUserDAI,
) (*UpdateUserDAO, error) {
	query := `UPDATE users SET name = ? WHERE id = ?`

	if _, err := gw.db.ExecContext(ctx, query, input.Name, input.ID); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &UpdateUserDAO{}, nil
}

type DeleteUserDAI struct {
	ID string
}

type DeleteUserDAO struct{}

func (gw *Gateway) DeleteUser(
	ctx context.Context,
	input DeleteUserDAI,
) (*DeleteUserDAO, error) {
	query := `DELETE FROM users WHERE id = ?`

	if _, err := gw.db.ExecContext(ctx, query, input.ID); err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	return &DeleteUserDAO{}, nil
}
