// @title           AI Chat API
// @version         1.0
// @description     Chat backend: CRUD for messages plus AI chat sessions powered by the OpenAI Responses API.
// @host            ai-chat.leetcoders.uz
// @schemes         https http
// @BasePath        /
package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/udevs/ai-chat/docs"
	appchat "github.com/udevs/ai-chat/internal/application/chat"
	appmsg "github.com/udevs/ai-chat/internal/application/message"
	"github.com/udevs/ai-chat/internal/infrastructure/config"
	"github.com/udevs/ai-chat/internal/infrastructure/openai"
	"github.com/udevs/ai-chat/internal/infrastructure/postgres"
	"github.com/udevs/ai-chat/internal/interfaces/http/handler"
	"github.com/udevs/ai-chat/internal/interfaces/http/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		log.Fatalf("pgxpool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("postgres ping: %v", err)
	}

	msgRepo := postgres.NewMessageRepository(pool)
	msgSvc := appmsg.NewService(msgRepo)
	msgHandler := handler.NewMessageHandler(msgSvc)

	chatRepo := postgres.NewChatRepository(pool)
	aiClient := openai.New(cfg.OpenAIAPIKey)
	chatSvc := appchat.NewService(chatRepo, msgRepo, aiClient, cfg.OpenAIModel)
	chatHandler := handler.NewChatHandler(chatSvc)

	r := router.New(msgHandler, chatHandler)

	srv := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("http listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
