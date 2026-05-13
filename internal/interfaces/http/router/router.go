package router

import (
	"github.com/gin-gonic/gin"

	"github.com/udevs/ai-chat/internal/interfaces/http/handler"
)

func New(h *handler.MessageHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	v1 := r.Group("/v1")
	{
		messages := v1.Group("/messages")
		messages.POST("", h.Create)
		messages.GET("", h.ListByChat)
		messages.GET("/:id", h.Get)
		messages.PUT("/:id", h.Update)
		messages.DELETE("/:id", h.Delete)
	}

	return r
}
