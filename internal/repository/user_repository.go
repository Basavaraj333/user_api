package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
	db "users-api/db/sqlc"
)

var ErrNotFound = errors.New("user not found")

type UserRepo interface {
	Create(ctx context.Context, name string, dob time.Time) (db.User, error)
	GetByID(ctx context.Context, id int32) (db.User, error)
	Update(ctx context.Context, id int32, name string, dob time.Time) (db.User, error)
	Delete(ctx context.Context, id int32) error
	List(ctx context.Context, limit, offset int32) ([]db.User, error)
	Count(ctx context.Context) (int64, error)
}

type userRepo struct{ q *db.Queries }

func NewUserRepo(q *db.Queries) UserRepo { return &userRepo{q} }

func (r *userRepo) Create(ctx context.Context, name string, dob time.Time) (db.User, error) {
	return r.q.CreateUser(ctx, db.CreateUserParams{Name: name, Dob: dob})
}

func (r *userRepo) GetByID(ctx context.Context, id int32) (db.User, error) {
	u, err := r.q.GetUserByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return db.User{}, ErrNotFound
	}
	return u, err
}

func (r *userRepo) Update(ctx context.Context, id int32, name string, dob time.Time) (db.User, error) {
	u, err := r.q.UpdateUser(ctx, db.UpdateUserParams{ID: id, Name: name, Dob: dob})
	if errors.Is(err, sql.ErrNoRows) {
		return db.User{}, ErrNotFound
	}
	return u, err
}

func (r *userRepo) Delete(ctx context.Context, id int32) error {
	return r.q.DeleteUser(ctx, id)
}

func (r *userRepo) List(ctx context.Context, limit, offset int32) ([]db.User, error) {
	return r.q.ListUsers(ctx, db.ListUsersParams{Limit: limit, Offset: offset})
}

func (r *userRepo) Count(ctx context.Context) (int64, error) {
	return r.q.CountUsers(ctx)
}
