package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"sonarqube-report-generator/internal/auth"
	"sonarqube-report-generator/internal/report"
	"sonarqube-report-generator/internal/sonarqube"
)

// WebHandler handles web UI endpoints
type WebHandler struct {
	authenticator *auth.Authenticator
	sonarClient   *sonarqube.Client
	storage       *report.Storage
	templates     *template.Template
}

// NewWebHandler creates a new web handler
func NewWebHandler(authenticator *auth.Authenticator, client *sonarqube.Client, storage *report.Storage) *WebHandler {
	return &WebHandler{
		authenticator: authenticator,
		sonarClient:   client,
		storage:       storage,
	}
}

// LoadTemplates loads HTML templates
func (h *WebHandler) LoadTemplates(templatesDir string) error {
	tmpl, err := template.ParseGlob(templatesDir + "/*.html")
	if err != nil {
		return err
	}
	h.templates = tmpl
	return nil
}

// LoginPage renders the login page
func (h *WebHandler) LoginPage(c *gin.Context) {
	// Check if already logged in
	if session := h.authenticator.GetSession(c); session != nil {
		c.Redirect(http.StatusFound, "/dashboard")
		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"error": c.Query("error"),
	})
}

// LoginRequest is the login form data
type LoginRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// Login processes the login form
func (h *WebHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Redirect(http.StatusFound, "/login?error=Username+and+password+are+required")
		return
	}

	if !h.authenticator.Login(c, req.Username, req.Password) {
		c.Redirect(http.StatusFound, "/login?error=Invalid+username+or+password")
		return
	}

	c.Redirect(http.StatusFound, "/dashboard")
}

// Logout handles user logout
func (h *WebHandler) Logout(c *gin.Context) {
	h.authenticator.Logout(c)
	c.Redirect(http.StatusFound, "/login")
}

// Dashboard renders the main dashboard
func (h *WebHandler) Dashboard(c *gin.Context) {
	session := h.authenticator.GetSession(c)

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"user": session.Username,
	})
}

// Index redirects to dashboard or login
func (h *WebHandler) Index(c *gin.Context) {
	if session := h.authenticator.GetSession(c); session != nil {
		c.Redirect(http.StatusFound, "/dashboard")
		return
	}
	c.Redirect(http.StatusFound, "/login")
}

// PreviewReportPage renders the report preview page in a new tab
func (h *WebHandler) PreviewReportPage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.HTML(http.StatusBadRequest, "preview.html", gin.H{
			"error": "Report ID is required",
		})
		return
	}

	// Get report record
	record, err := h.storage.GetRecord(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "preview.html", gin.H{
			"error": "Report not found: " + err.Error(),
		})
		return
	}

	// Only support markdown preview
	if record.Format != "md" {
		c.HTML(http.StatusBadRequest, "preview.html", gin.H{
			"error": "Preview is only available for Markdown reports. Please download the PDF instead.",
		})
		return
	}

	// Read markdown content
	content, err := os.ReadFile(record.FilePath)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "preview.html", gin.H{
			"error": "Failed to read report file: " + err.Error(),
		})
		return
	}

	// Convert to string and escape for JavaScript
	markdownStr := string(content)

	// Marshal to JSON to properly escape special characters
	markdownJSON, err := json.Marshal(markdownStr)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "preview.html", gin.H{
			"error": "Failed to process markdown content: " + err.Error(),
		})
		return
	}

	// Format generated date
	generatedAt := record.GeneratedAt.Format("Jan 2, 2006 at 3:04 PM")

	// Render preview template
	// markdownJSON is already a JSON string (with quotes), so we use it directly
	c.HTML(http.StatusOK, "preview.html", gin.H{
		"reportID":        record.ID,
		"projectKey":      record.ProjectKey,
		"projectName":     record.ProjectName,
		"branch":          record.Branch,
		"generatedAt":     generatedAt,
		"markdownContent": template.JS(markdownJSON),
	})
}

// ShareReportPage renders the report share page (public, no login required)
func (h *WebHandler) ShareReportPage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.HTML(http.StatusBadRequest, "share.html", gin.H{
			"error": "Report ID is required",
		})
		return
	}

	// Get report record
	record, err := h.storage.GetRecord(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "share.html", gin.H{
			"error": "Report not found: " + err.Error(),
		})
		return
	}

	// Only support markdown preview
	if record.Format != "md" {
		c.HTML(http.StatusBadRequest, "share.html", gin.H{
			"error": "Share is only available for Markdown reports. Please download the PDF instead.",
		})
		return
	}

	// Read markdown content
	content, err := os.ReadFile(record.FilePath)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "share.html", gin.H{
			"error": "Failed to read report file: " + err.Error(),
		})
		return
	}

	// Convert to string and escape for JavaScript
	markdownStr := string(content)

	// Marshal to JSON to properly escape special characters
	markdownJSON, err := json.Marshal(markdownStr)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "share.html", gin.H{
			"error": "Failed to process markdown content: " + err.Error(),
		})
		return
	}

	// Format generated date
	generatedAt := record.GeneratedAt.Format("Jan 2, 2006 at 3:04 PM")

	// Render share template
	c.HTML(http.StatusOK, "share.html", gin.H{
		"reportID":        record.ID,
		"projectKey":      record.ProjectKey,
		"projectName":     record.ProjectName,
		"branch":          record.Branch,
		"generatedAt":     generatedAt,
		"markdownContent": template.JS(markdownJSON),
	})
}
