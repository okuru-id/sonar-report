package handler

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"

	"sonarqube-report-generator/internal/auth"
	"sonarqube-report-generator/internal/sonarqube"
)

// WebHandler handles web UI endpoints
type WebHandler struct {
	authenticator *auth.Authenticator
	sonarClient   *sonarqube.Client
	templates     *template.Template
}

// NewWebHandler creates a new web handler
func NewWebHandler(authenticator *auth.Authenticator, client *sonarqube.Client) *WebHandler {
	return &WebHandler{
		authenticator: authenticator,
		sonarClient:   client,
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
