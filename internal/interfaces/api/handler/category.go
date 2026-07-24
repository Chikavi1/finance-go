package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/agnathor/finances-go/internal/application/category"
	"github.com/agnathor/finances-go/internal/domain"
	"github.com/agnathor/finances-go/internal/interfaces/api/dto"
	"github.com/agnathor/finances-go/internal/interfaces/api/middleware"
	"github.com/agnathor/finances-go/internal/interfaces/api/validator"
	"github.com/agnathor/finances-go/pkg/response"
)

type CategoryHandler struct {
	categoryService category.Service
}

func NewCategoryHandler(categoryService category.Service) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req dto.CreateCategoryRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	category, err := h.categoryService.Create(c.Context(), userID, req.Name, req.Type, req.Color, req.Icon)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to create category")
	}

	return response.Created(c, mapCategoryToResponse(category))
}

func (h *CategoryHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	categoryID := c.Params("id")

	category, err := h.categoryService.GetByID(c.Context(), userID, categoryID)
	if err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "category not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to get category")
	}

	return response.JSON(c, fiber.StatusOK, mapCategoryToResponse(category))
}

func (h *CategoryHandler) GetAll(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	categories, err := h.categoryService.GetAll(c.Context(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to get categories")
	}

	resp := make([]dto.CategoryResponse, len(categories))
	for i, a := range categories {
		resp[i] = mapCategoryToResponse(a)
	}

	return response.JSON(c, fiber.StatusOK, resp)
}

func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	categoryID := c.Params("id")

	var req dto.UpdateCategoryRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	category, err := h.categoryService.Update(c.Context(), userID, categoryID, req.Name, req.Type, req.Color, req.Icon)
	if err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "category not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to update category")
	}

	return response.JSON(c, fiber.StatusOK, mapCategoryToResponse(category))
}

func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	categoryID := c.Params("id")

	if err := h.categoryService.Delete(c.Context(), userID, categoryID); err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "category not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to delete category")
	}

	return response.JSON(c, fiber.StatusOK, fiber.Map{
		"message": "category deleted successfully",
	})
}

func mapCategoryToResponse(category *domain.Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:        category.ID,
		UserID:    category.UserID,
		Name:      category.Name,
		Type:      category.Type,
		Color:     category.Color,
		Icon:      category.Icon,
		CreatedAt: category.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
