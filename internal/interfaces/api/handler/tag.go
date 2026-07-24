package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/agnathor/finances-go/internal/application/tag"
	"github.com/agnathor/finances-go/internal/domain"
	"github.com/agnathor/finances-go/internal/interfaces/api/dto"
	"github.com/agnathor/finances-go/internal/interfaces/api/middleware"
	"github.com/agnathor/finances-go/internal/interfaces/api/validator"
	"github.com/agnathor/finances-go/pkg/response"
)

type TagHandler struct {
	tagService tag.Service
}

func NewTagHandler(tagService tag.Service) *TagHandler {
	return &TagHandler{tagService: tagService}
}

func (h *TagHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req dto.CreateTagRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	tag, err := h.tagService.Create(c.Context(), userID, req.Name)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to create tag")
	}

	return response.Created(c, mapTagToResponse(tag))
}

func (h *TagHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tagID := c.Params("id")

	tag, err := h.tagService.GetByID(c.Context(), userID, tagID)
	if err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "tag not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to get tag")
	}

	return response.JSON(c, fiber.StatusOK, mapTagToResponse(tag))
}

func (h *TagHandler) GetAll(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	tags, err := h.tagService.GetAll(c.Context(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to get tags")
	}

	resp := make([]dto.TagResponse, len(tags))
	for i, t := range tags {
		resp[i] = mapTagToResponse(t)
	}

	return response.JSON(c, fiber.StatusOK, resp)
}

func (h *TagHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tagID := c.Params("id")

	var req dto.UpdateTagRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	tag, err := h.tagService.Update(c.Context(), userID, tagID, req.Name)
	if err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "tag not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to update tag")
	}

	return response.JSON(c, fiber.StatusOK, mapTagToResponse(tag))
}

func (h *TagHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	tagID := c.Params("id")

	if err := h.tagService.Delete(c.Context(), userID, tagID); err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "tag not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to delete tag")
	}

	return response.JSON(c, fiber.StatusOK, fiber.Map{
		"message": "tag deleted successfully",
	})
}

func mapTagToResponse(tag *domain.Tag) dto.TagResponse {
	return dto.TagResponse{
		ID:        tag.ID,
		UserID:    tag.UserID,
		Name:      tag.Name,
		CreatedAt: tag.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
