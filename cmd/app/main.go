package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trieu_mock_project_go/internal/bootstrap"
	"trieu_mock_project_go/internal/config"
	"trieu_mock_project_go/internal/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Initialize database
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Redis
	if err := config.InitRedis(); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	// Create Gin router
	router := gin.Default()

	setupHtmlTemplate(router)

	router.Static("/static", "./static")

	setupSessionConfiguration(router, cfg)

	// Register custom validations
	bootstrap.RegisterCustomValidations()

	// Initialize app container
	appContainer := bootstrap.NewAppContainer()

	// Initialize Application Services (Redis sub, RabbitMQ sub)
	if err := appContainer.InitializeApp(); err != nil {
		log.Fatalf("Failed to initialize application services: %v", err)
	}

	// Setup routes
	routes.SetupRoutes(router, appContainer)

	startServerWithGracefulShutdown(router, appContainer)
}

func setupHtmlTemplate(router *gin.Engine) {
	router.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"min": func(a, b int64) int64 {
			if a < b {
				return a
			}
			return b
		},
		"int64": func(v int) int64 {
			return int64(v)
		},
		"iterate": func(start, end int) []int {
			var items []int
			for i := start; i <= end; i++ {
				items = append(items, i)
			}
			return items
		},
	})

	router.LoadHTMLGlob("templates/**/*")
}

func setupSessionConfiguration(router *gin.Engine, cfg *config.Config) {
	store, err := redis.NewStore(
		10,
		"tcp",
		fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		cfg.Redis.Username,
		cfg.Redis.Password,
		[]byte(cfg.SessionConfig.Secret),
	)
	if err != nil {
		log.Fatalf("Failed to initialize redis session store: %v", err)
	}
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   cfg.SessionConfig.MaxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   cfg.SessionConfig.Secure,
	})
	router.Use(sessions.Sessions("trieu_mock_project_session", store))
}

func startServerWithGracefulShutdown(router http.Handler, appContainer *bootstrap.AppContainer) {
	cfg := config.LoadConfig()

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		log.Printf("Starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Shutdown application services (RabbitMQ, etc.)
	appContainer.Shutdown()

	log.Println("Server exiting")
}
