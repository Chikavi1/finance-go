package handler

import (
	"path/filepath"

	"github.com/gofiber/fiber/v2"

	"github.com/agnathor/finances-go/internal/application/attachment"
	"github.com/agnathor/finances-go/internal/domain"
	"github.com/agnathor/finances-go/internal/interfaces/api/middleware"
	"github.com/agnathor/finances-go/pkg/response"
)

type AttachmentHandler struct {
	attachmentService attachment.Service
}

func NewAttachmentHandler(attachmentService attachment.Service) *AttachmentHandler {
	return &AttachmentHandler{attachmentService: attachmentService}
}

type uploadRequest struct {
	TransactionID *string `json:"transaction_id,omitempty" form:"transaction_id"`
}

func (h *AttachmentHandler) Upload(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	file, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "file is required")
	}

	transactionID := c.FormValue("transaction_id")
	var txID *string
	if transactionID != "" {
		txID = &transactionID
	}

	src, err := file.Open()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to read file")
	}
	defer src.Close()

	data := make([]byte, file.Size)
	if _, err := src.Read(data); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to read file data")
	}

	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		ext := filepath.Ext(file.Filename)
		switch ext {
		case ".pdf":
			mimeType = "application/pdf"
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".png":
			mimeType = "image/png"
		case ".gif":
			mimeType = "image/gif"
		default:
			mimeType = "application/octet-stream"
		}
	}

	att, err := h.attachmentService.Create(c.Context(), userID, txID, file.Filename, mimeType, data)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to upload file")
	}

	return response.Created(c, mapAttachmentToResponse(att))
}

func (h *AttachmentHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id := c.Params("id")

	att, err := h.attachmentService.GetByID(c.Context(), userID, id)
	if err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "attachment not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to get attachment")
	}

	return response.JSON(c, fiber.StatusOK, mapAttachmentToResponse(att))
}

func (h *AttachmentHandler) GetByTransactionID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	transactionID := c.Params("transactionId")

	attachments, err := h.attachmentService.GetByTransactionID(c.Context(), userID, transactionID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to get attachments")
	}

	responses := make([]AttachmentResponse, len(attachments))
	for i, a := range attachments {
		responses[i] = mapAttachmentToResponse(a)
	}

	return response.JSON(c, fiber.StatusOK, responses)
}

func (h *AttachmentHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id := c.Params("id")

	if err := h.attachmentService.Delete(c.Context(), userID, id); err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "attachment not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to delete attachment")
	}

	return response.JSON(c, fiber.StatusOK, fiber.Map{
		"message": "attachment deleted successfully",
	})
}

type AttachmentResponse struct {
	ID            string  `json:"id"`
	TransactionID *string `json:"transaction_id,omitempty"`
	Filename      string  `json:"filename"`
	OriginalName  string  `json:"original_name"`
	MimeType      string  `json:"mime_type"`
	Size          int64   `json:"size"`
	URL           string  `json:"url"`
	CreatedAt     string  `json:"created_at"`
}

func mapAttachmentToResponse(att *domain.Attachment) AttachmentResponse {
	return AttachmentResponse{
		ID:            att.ID,
		TransactionID: att.TransactionID,
		Filename:      att.Filename,
		OriginalName:  att.OriginalName,
		MimeType:      att.MimeType,
		Size:          att.Size,
		URL:           att.URL,
		CreatedAt:     att.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
