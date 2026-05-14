package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	appmsg "github.com/udevs/ai-chat/internal/application/message"
	domain "github.com/udevs/ai-chat/internal/domain/message"
	"github.com/udevs/ai-chat/internal/interfaces/http/dto"
)

type MessageHandler struct {
	svc *appmsg.Service
}

func NewMessageHandler(svc *appmsg.Service) *MessageHandler {
	return &MessageHandler{svc: svc}
}

// Create godoc
// @Summary      Create a message
// @Description  Stores a new chat message. `message` is an arbitrary JSON object.
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateMessageRequest  true  "Message to create"
// @Success      201   {object}  dto.MessageResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /v1/messages [post]
func (h *MessageHandler) Create(c *gin.Context) {
	var req dto.CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m, err := h.svc.Create(c.Request.Context(), appmsg.CreateInput{
		ChatID:   req.ChatID,
		SenderID: req.SenderID,
		Body:     req.Message,
	})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.FromDomain(m))
}

// Get godoc
// @Summary      Get a message by id
// @Tags         messages
// @Produce      json
// @Param        id   path      string  true  "Message ID (UUID)"
// @Success      200  {object}  dto.MessageResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /v1/messages/{id} [get]
func (h *MessageHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	m, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.FromDomain(m))
}

// ListByChat godoc
// @Summary      List messages for a chat
// @Tags         messages
// @Produce      json
// @Param        chat_id  query     string  true   "Chat ID (UUID)"
// @Param        limit    query     int     false  "Page size (1-200, default 50)"
// @Param        offset   query     int     false  "Page offset (default 0)"
// @Success      200      {object}  dto.ListMessagesResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /v1/messages [get]
func (h *MessageHandler) ListByChat(c *gin.Context) {
	chatID, err := uuid.Parse(c.Query("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat_id"})
		return
	}
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 32)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 32)

	items, err := h.svc.ListByChat(c.Request.Context(), appmsg.ListInput{
		ChatID: chatID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.ListMessagesResponse{Items: dto.FromDomainList(items)})
}

// ListByChatID godoc
// @Summary      List messages for a chat
// @Tags         chats
// @Produce      json
// @Param        id      path      string  true   "Chat ID (UUID)"
// @Param        limit   query     int     false  "Page size (1-200, default 50)"
// @Param        offset  query     int     false  "Page offset (default 0)"
// @Success      200     {object}  dto.ListMessagesResponse
// @Failure      400     {object}  dto.ErrorResponse
// @Failure      500     {object}  dto.ErrorResponse
// @Router       /v1/chats/{id}/messages [get]
func (h *MessageHandler) ListByChatID(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 32)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 32)

	items, err := h.svc.ListByChat(c.Request.Context(), appmsg.ListInput{
		ChatID: chatID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.ListMessagesResponse{Items: dto.FromDomainList(items)})
}

// Update godoc
// @Summary      Update a message body
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        id    path      string                    true  "Message ID (UUID)"
// @Param        body  body      dto.UpdateMessageRequest  true  "New message body"
// @Success      200   {object}  dto.MessageResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      404   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /v1/messages/{id} [put]
func (h *MessageHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req dto.UpdateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m, err := h.svc.Update(c.Request.Context(), appmsg.UpdateInput{
		ID:   id,
		Body: req.Message,
	})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.FromDomain(m))
}

// Delete godoc
// @Summary      Delete a message
// @Tags         messages
// @Param        id   path  string  true  "Message ID (UUID)"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /v1/messages/{id} [delete]
func (h *MessageHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrInvalidChatID),
		errors.Is(err, domain.ErrInvalidSenderID),
		errors.Is(err, domain.ErrInvalidBody):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
