package service

import (
	"context"
	"math"
	"time"
	db "users-api/db/sqlc"
	"users-api/internal/logger"
	"users-api/internal/models"
	"users-api/internal/repository"

	"go.uber.org/zap"
)

type UserService interface {
	Create(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error)
	GetByID(ctx context.Context, id int32) (models.UserWithAge, error)
	Update(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserResponse, error)
	Delete(ctx context.Context, id int32) error
	List(ctx context.Context, page, limit int) (models.PaginatedResponse, error)
}

type userService struct{ repo repository.UserRepo }

func NewUserService(repo repository.UserRepo) UserService { return &userService{repo} }

const dobLayout = "2006-01-02"

func toResponse(u db.User) models.UserResponse {
	return models.UserResponse{ID: u.ID, Name: u.Name, Dob: u.Dob.Format(dobLayout)}
}

func toResponseWithAge(u db.User) models.UserWithAge {
	return models.UserWithAge{ID: u.ID, Name: u.Name, Dob: u.Dob.Format(dobLayout), Age: CalculateAge(u.Dob)}
}

func (s *userService) Create(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error) {
	dob, _ := time.Parse(dobLayout, req.Dob) // validator already confirmed format
	u, err := s.repo.Create(ctx, req.Name, dob)
	if err != nil {
		logger.Log.Error("create user", zap.Error(err))
		return models.UserResponse{}, err
	}
	logger.Log.Info("user created", zap.Int32("id", u.ID))
	return toResponse(u), nil
}

func (s *userService) GetByID(ctx context.Context, id int32) (models.UserWithAge, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return models.UserWithAge{}, err
	}
	return toResponseWithAge(u), nil
}

func (s *userService) Update(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserResponse, error) {
	dob, _ := time.Parse(dobLayout, req.Dob)
	u, err := s.repo.Update(ctx, id, req.Name, dob)
	if err != nil {
		return models.UserResponse{}, err
	}
	logger.Log.Info("user updated", zap.Int32("id", u.ID))
	return toResponse(u), nil
}

func (s *userService) Delete(ctx context.Context, id int32) error {
	err := s.repo.Delete(ctx, id)
	if err == nil {
		logger.Log.Info("user deleted", zap.Int32("id", id))
	}
	return err
}

func (s *userService) List(ctx context.Context, page, limit int) (models.PaginatedResponse, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 10 }

	offset := (page - 1) * limit
	users, err := s.repo.List(ctx, int32(limit), int32(offset))
	if err != nil {
		return models.PaginatedResponse{}, err
	}
	total, err := s.repo.Count(ctx)
	if err != nil {
		return models.PaginatedResponse{}, err
	}

	data := make([]models.UserWithAge, len(users))
	for i, u := range users {
		data[i] = toResponseWithAge(u)
	}

	return models.PaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}, nil
}

// CalculateAge returns completed years since dob.
// Exported so it can be unit tested independently.
func CalculateAge(dob time.Time) int {
	now := time.Now()
	years := now.Year() - dob.Year()
	// Subtract 1 if birthday hasn't occurred yet this year
	if now.YearDay() < dob.YearDay() {
		years--
	}
	return years
}
