package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/udevs/ai-chat/internal/interfaces/http/handler"
)

func New(messages *handler.MessageHandler, chats *handler.ChatHandler, images *handler.ImageHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/v1")
	{
		m := v1.Group("/messages")
		m.POST("", messages.Create)
		m.GET("", messages.ListByChat)
		m.GET("/:id", messages.Get)
		m.PUT("/:id", messages.Update)
		m.DELETE("/:id", messages.Delete)

		ch := v1.Group("/chats")
		ch.POST("", chats.Create)
		ch.GET("", chats.List)
		ch.GET("/:id", chats.Get)
		ch.GET("/:id/messages", messages.ListByChatID)
		ch.POST("/:id/messages", chats.Send)
		ch.DELETE("", chats.DeleteAll)
		ch.DELETE("/:id", chats.Delete)

		img := v1.Group("/images")
		img.POST("/generate", images.Generate)
	}

	return r
}
