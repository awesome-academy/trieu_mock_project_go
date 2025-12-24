package main

import (
	"fmt"
	"log"
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

	router.LoadHTMLGlob("templates/**/*.tmpl")
	router.Static("/static", "./static")

	store := cookie.NewStore([]byte(cfg.AppSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	})
	router.Use(sessions.Sessions("mysession", store))

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
