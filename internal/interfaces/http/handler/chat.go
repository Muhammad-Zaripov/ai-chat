package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	appchat "github.com/udevs/ai-chat/internal/application/chat"
	domain "github.com/udevs/ai-chat/internal/domain/chat"
	"github.com/udevs/ai-chat/internal/interfaces/http/dto"
)

type ChatHandler struct {
	svc *appchat.Service
}

func NewChatHandler(svc *appchat.Service) *ChatHandler {
	return &ChatHandler{svc: svc}
}

// Create godoc
// @Summary      Create a new chat session
// @Description  Allocates a chat with its own conversation state (server-side, via OpenAI Responses API).
// @Tags         chats
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateChatRequest  true  "Chat options"
// @Success      201   {object}  dto.ChatResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /v1/chats [post]
func (h *ChatHandler) Create(c *gin.Context) {
	var req dto.CreateChatRequest
	// Body is optional — accept empty POST as "use defaults".
	_ = c.ShouldBindJSON(&req)

	out, err := h.svc.Create(c.Request.Context(), appchat.CreateInput{
		Title: req.Title,
		Model: req.Model,
	})
	if err != nil {
		writeChatError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.ChatFromDomain(out))
}

// List godoc
// @Summary      List chats
// @Tags         chats
// @Produce      json
// @Param        limit   query     int  false  "Page size (1-200, default 50)"
// @Param        offset  query     int  false  "Page offset (default 0)"
// @Success      200     {object}  dto.ListChatsResponse
// @Failure      500     {object}  dto.ErrorResponse
// @Router       /v1/chats [get]
func (h *ChatHandler) List(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 32)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 32)

	items, err := h.svc.List(c.Request.Context(), int32(limit), int32(offset))
	if err != nil {
		writeChatError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.ListChatsResponse{Items: dto.ChatsFromDomain(items)})
}

// Get godoc
// @Summary      Get a chat by id
// @Tags         chats
// @Produce      json
// @Param        id   path      string  true  "Chat ID (UUID)"
// @Success      200  {object}  dto.ChatResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /v1/chats/{id} [get]
func (h *ChatHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	out, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		writeChatError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.ChatFromDomain(out))
}

// Send godoc
// @Summary      Send a message and get the AI's reply
// @Description  Appends a user turn to the chat. The previous response_id is sent to OpenAI so the model has full conversation context.
// @Tags         chats
// @Accept       json
// @Produce      json
// @Param        id    path      string                  true  "Chat ID (UUID)"
// @Param        body  body      dto.SendMessageRequest  true  "User message"
// @Success      200   {object}  dto.SendMessageResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      404   {object}  dto.ErrorResponse
// @Failure      502   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /v1/chats/{id}/messages [post]
func (h *ChatHandler) Send(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	senderID := uuid.Nil
	if req.SenderID != nil {
		senderID = *req.SenderID
	}
	out, err := h.svc.Send(c.Request.Context(), id, appchat.SendInput{
		Input:    req.Input,
		SenderID: senderID,
	})
	if err != nil {
		writeChatError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.SendMessageResponse{
		Chat:  dto.ChatFromDomain(out.Chat),
		Reply: out.Reply,
	})
}

// Delete godoc
// @Summary      Delete a chat
// @Tags         chats
// @Param        id   path      string  true  "Chat ID (UUID)"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /v1/chats/{id} [delete]
func (h *ChatHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		writeChatError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// DeleteAll godoc
// @Summary      Delete all chats
// @Tags         chats
// @Success      204
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /v1/chats [delete]
func (h *ChatHandler) DeleteAll(c *gin.Context) {
	if err := h.svc.DeleteAll(c.Request.Context()); err != nil {
		writeChatError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func writeChatError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrEmptyInput):
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}
}
