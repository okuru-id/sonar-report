package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"sonarqube-report-generator/internal/auth"
	"sonarqube-report-generator/internal/config"
	"sonarqube-report-generator/internal/handler"
	"sonarqube-report-generator/internal/report"
	"sonarqube-report-generator/internal/sonarqube"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Validate required config
	if cfg.SonarQubeURL == "" {
		log.Fatal("SONARQUBE_URL is required")
	}
	if cfg.SonarQubeToken == "" {
		log.Fatal("SONARQUBE_TOKEN is required")
	}

	// Initialize SonarQube client
	sonarClient := sonarqube.NewClient(cfg.SonarQubeURL, cfg.SonarQubeToken)

	// Validate connection
	if err := sonarClient.Validate(); err != nil {
		log.Printf("Warning: Failed to connect to SonarQube: %v", err)
	} else {
		log.Printf("Connected to SonarQube: %s", cfg.SonarQubeURL)
	}

	// Initialize storage
	storage, err := report.NewStorage(cfg.ReportStoragePath)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Clean old reports
	if err := storage.CleanOld(cfg.ReportRetentionDays); err != nil {
		log.Printf("Warning: Failed to clean old reports: %v", err)
	}

	// Initialize authenticator
	authenticator := auth.NewAuthenticator(cfg.AdminUsername, cfg.AdminPassword)

	// Initialize handlers
	apiHandler := handler.NewAPIHandler(sonarClient, storage)
	webHandler := handler.NewWebHandler(authenticator, sonarClient, storage)

	// Setup Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.RecoveryWithWriter(gin.DefaultWriter))
	r.Use(gin.Logger())

	// Custom recovery that returns JSON
	r.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatusJSON(500, gin.H{
					"error": fmt.Sprintf("Internal server error: %v", err),
				})
			}
		}()
		c.Next()
	})

	// Load templates
	r.LoadHTMLGlob("web/templates/*.html")

	// Static files
	r.Static("/static", "web/static")

	// Public routes
	r.GET("/", webHandler.Index)
	r.GET("/login", webHandler.LoginPage)
	r.POST("/login", webHandler.Login)
	r.GET("/logout", webHandler.Logout)
	r.GET("/share/:id", webHandler.ShareReportPage)

	// Health check (public)
	r.GET("/api/v1/health", apiHandler.HealthCheck)

	// Public API routes (no authentication)
	r.POST("/api/v1/jenkins/trigger", apiHandler.TriggerJenkins)

	// Protected web routes
	protected := r.Group("/")
	protected.Use(authenticator.AuthMiddleware())
	{
		protected.GET("/dashboard", webHandler.Dashboard)
		protected.GET("/reports/:id/preview", webHandler.PreviewReportPage)
	}

	// Protected API routes
	api := r.Group("/api/v1")
	api.Use(authenticator.AuthMiddleware())
	{
		// Projects
		api.GET("/projects", apiHandler.GetProjects)
		api.GET("/projects/:key/branches", apiHandler.GetBranches)

		// Reports
		api.POST("/reports/generate", apiHandler.GenerateReport)
		api.GET("/reports/history", apiHandler.GetHistory)
		api.GET("/reports/:id/download", apiHandler.DownloadReport)
		api.GET("/reports/:id/preview", apiHandler.PreviewReport)
		api.DELETE("/reports/:id", apiHandler.DeleteReport)
		api.DELETE("/reports/history", apiHandler.ClearHistory)
	}

	// Start server
	addr := ":" + cfg.ServerPort
	log.Printf("Starting server on %s", addr)
	log.Printf("Admin dashboard: http://localhost%s", addr)
	log.Printf("API endpoint: http://localhost%s/api/v1", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
