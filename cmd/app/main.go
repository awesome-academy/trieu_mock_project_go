package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"trieu_mock_project_go/internal/bootstrap"
	"trieu_mock_project_go/internal/config"
	"trieu_mock_project_go/internal/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Initialize database
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
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

	// Setup routes
	routes.SetupRoutes(router, appContainer)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
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
	store := cookie.NewStore([]byte(cfg.SessionConfig.Secret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   cfg.SessionConfig.MaxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   cfg.SessionConfig.Secure,
	})
	router.Use(sessions.Sessions("trieu_mock_project_session", store))
}
