package handler

import (
	"errors"
	"strconv"
	"users-api/internal/models"
	"users-api/internal/repository"
	"users-api/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	svc      service.UserService
	validate *validator.Validate
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc, validate: validator.New()}
}

func (h *UserHandler) parseID(c *fiber.Ctx) (int32, error) {
	id, err := strconv.Atoi(c.Params("id"))
	return int32(id), err
}

func errResponse(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(models.ErrorResponse{Error: msg})
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return errResponse(c, fiber.StatusBadRequest, "invalid request body")
	}
	if err := h.validate.Struct(req); err != nil {
		return errResponse(c, fiber.StatusBadRequest, err.Error())
	}
	user, err := h.svc.Create(c.Context(), req)
	if err != nil {
		return errResponse(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := h.parseID(c)
	if err != nil {
		return errResponse(c, fiber.StatusBadRequest, "invalid id")
	}
	user, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return errResponse(c, fiber.StatusNotFound, "user not found")
		}
		return errResponse(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := h.parseID(c)
	if err != nil {
		return errResponse(c, fiber.StatusBadRequest, "invalid id")
	}
	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return errResponse(c, fiber.StatusBadRequest, "invalid request body")
	}
	if err := h.validate.Struct(req); err != nil {
		return errResponse(c, fiber.StatusBadRequest, err.Error())
	}
	user, err := h.svc.Update(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return errResponse(c, fiber.StatusNotFound, "user not found")
		}
		return errResponse(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := h.parseID(c)
	if err != nil {
		return errResponse(c, fiber.StatusBadRequest, "invalid id")
	}
	if err := h.svc.Delete(c.Context(), id); err != nil {
		return errResponse(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	users, err := h.svc.List(c.Context(), c.QueryInt("page", 1), c.QueryInt("limit", 10))
	if err != nil {
		return errResponse(c, fiber.StatusInternalServerError, "internal server error")
	}
	return c.JSON(users)
}
