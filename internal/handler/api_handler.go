package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"sonarqube-report-generator/internal/report"
	"sonarqube-report-generator/internal/sonarqube"
)

// APIHandler handles API endpoints
type APIHandler struct {
	sonarClient *sonarqube.Client
	generator   *report.Generator
	storage     *report.Storage
	mdGen       *report.MarkdownGenerator
	pdfGen      *report.PDFGenerator
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(client *sonarqube.Client, storage *report.Storage) *APIHandler {
	return &APIHandler{
		sonarClient: client,
		generator:   report.NewGenerator(client),
		storage:     storage,
		mdGen:       report.NewMarkdownGenerator(),
		pdfGen:      report.NewPDFGenerator(),
	}
}

// HealthCheck returns the health status
func (h *APIHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": "1.0.0",
	})
}

// GetProjects returns all projects
func (h *APIHandler) GetProjects(c *gin.Context) {
	projects, err := h.sonarClient.GetProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

// GetBranches returns branches for a project
func (h *APIHandler) GetBranches(c *gin.Context) {
	projectKey := c.Param("key")
	if projectKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project key is required"})
		return
	}

	branches, err := h.sonarClient.GetBranches(projectKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"branches": branches})
}

// GenerateRequest is the request body for generating reports
type GenerateRequest struct {
	ProjectKey          string `json:"projectKey" binding:"required"`
	Branch              string `json:"branch"`
	Format              string `json:"format"`              // md or pdf
	IncludeCodeSnippets *bool  `json:"includeCodeSnippets"` // include code snippets in report (default: true)
	IncludeHowToFix     *bool  `json:"includeHowToFix"`     // include how to fix in report (default: true)
}

// GenerateReport generates a report
func (h *APIHandler) GenerateReport(c *gin.Context) {
	var req GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default format
	if req.Format == "" {
		req.Format = "md"
	}

	// Validate format
	if req.Format != "md" && req.Format != "pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format must be 'md' or 'pdf'"})
		return
	}

	// Set default options
	options := report.GenerateOptions{
		IncludeCodeSnippets: true,
		IncludeHowToFix:     true,
	}
	if req.IncludeCodeSnippets != nil {
		options.IncludeCodeSnippets = *req.IncludeCodeSnippets
	}
	if req.IncludeHowToFix != nil {
		options.IncludeHowToFix = *req.IncludeHowToFix
	}

	// Generate report data
	data, err := h.generator.Generate(req.ProjectKey, req.Branch, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate content based on format
	var content []byte
	switch req.Format {
	case "pdf":
		content, err = h.pdfGen.Generate(data)
	default:
		content, err = h.mdGen.Generate(data)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save report
	record, err := h.storage.Save(data, content, req.Format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"report":  record,
		"data":    data,
	})
}

// GetHistory returns report history
func (h *APIHandler) GetHistory(c *gin.Context) {
	history := h.storage.GetHistory()
	c.JSON(http.StatusOK, gin.H{"reports": history})
}

// DownloadReport downloads a report file
func (h *APIHandler) DownloadReport(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "report id is required"})
		return
	}

	record, err := h.storage.GetRecord(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Determine content type
	var contentType string
	switch record.Format {
	case "pdf":
		contentType = "application/pdf"
	default:
		contentType = "text/markdown"
	}

	c.Header("Content-Disposition", "attachment; filename="+record.FileName)
	c.Header("Content-Type", contentType)
	c.File(record.FilePath)
}

// PreviewReport returns report content for preview
func (h *APIHandler) PreviewReport(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "report id is required"})
		return
	}

	record, err := h.storage.GetRecord(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// For markdown, read and return content
	if record.Format == "md" {
		content, err := filepath.Abs(record.FilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.File(content)
		return
	}

	// For PDF, just return metadata
	c.JSON(http.StatusOK, gin.H{
		"record":  record,
		"message": "PDF preview not supported, please download",
	})
}

// DeleteReport deletes a report
func (h *APIHandler) DeleteReport(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "report id is required"})
		return
	}

	if err := h.storage.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ClearHistory clears all report history
func (h *APIHandler) ClearHistory(c *gin.Context) {
	if err := h.storage.ClearAll(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// JenkinsTriggerRequest represents the request structure for Jenkins trigger
type JenkinsTriggerRequest struct {
	JenkinsURL   string `header:"X-Jenkins-Url"`
	JenkinsToken string `header:"X-Jenkins-Token"`
	JenkinsJob   string `header:"X-Jenkins-Job"`
	JenkinsUser  string `header:"X-Jenkins-User"`
}

// JenkinsCrumbResponse represents the crumb response from Jenkins
type JenkinsCrumbResponse struct {
	CrumbRequestField string `json:"crumbRequestField"`
	Crumb             string `json:"crumb"`
}

// TriggerJenkins triggers a Jenkins job build
func (h *APIHandler) TriggerJenkins(c *gin.Context) {
	// Parse headers
	var req JenkinsTriggerRequest
	if err := c.ShouldBindHeader(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid headers"})
		return
	}

	// Validate required headers
	if req.JenkinsURL == "" || req.JenkinsToken == "" || req.JenkinsJob == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Jenkins-Url, X-Jenkins-Token, and X-Jenkins-Job headers are required"})
		return
	}

	// Set default user if not provided
	user := "sonar"
	if req.JenkinsUser != "" {
		user = req.JenkinsUser
	}

	// Get crumb from Jenkins for CSRF protection
	crumbURL := fmt.Sprintf("%s/crumbIssuer/api/json", req.JenkinsURL)
	crumbReq, err := http.NewRequest("GET", crumbURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create crumb request: %v", err)})
		return
	}

	// Set basic auth for crumb request
	crumbReq.SetBasicAuth(user, req.JenkinsToken)

	crumbResp, err := http.DefaultClient.Do(crumbReq)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to connect to Jenkins: %v", err)
		if err.Error() == "dial tcp: lookup" || err.Error() == "connect: connection refused" || err.Error() == "no such host" {
			errorMsg = fmt.Sprintf("Cannot connect to Jenkins server at %s. Please check: 1) Jenkins is running, 2) Network connectivity, 3) Correct URL and port", req.JenkinsURL)
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}
	defer crumbResp.Body.Close()

	if crumbResp.StatusCode == http.StatusForbidden {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jenkins authentication failed: User does not have permission to access crumb issuer. Please check: 1) User has 'Overall/Read' permission, 2) CSRF protection is enabled, 3) User can access /crumbIssuer/api/json endpoint."})
		return
	}

	if crumbResp.StatusCode != http.StatusOK {
		// Try to read response body for more details
		errorBody, _ := io.ReadAll(crumbResp.Body)
		errorMsg := fmt.Sprintf("Failed to get Jenkins crumb: HTTP %d", crumbResp.StatusCode)
		if len(errorBody) > 0 {
			errorMsg += fmt.Sprintf(" - %s", string(errorBody))
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	// Parse crumb response
	var crumbData JenkinsCrumbResponse
	body, _ := io.ReadAll(crumbResp.Body)
	if err := json.Unmarshal(body, &crumbData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to parse Jenkins crumb response: %v", err)})
		return
	}

	// Build the Jenkins job URL
	jobURL := fmt.Sprintf("%s/job/%s/build", req.JenkinsURL, req.JenkinsJob)

	// Create request to trigger Jenkins job
	client := &http.Client{}
	buildReq, err := http.NewRequest("POST", jobURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create Jenkins build request: %v", err)})
		return
	}

	// Set headers for Jenkins authentication and CSRF protection
	buildReq.SetBasicAuth(user, req.JenkinsToken)
	buildReq.Header.Set(crumbData.CrumbRequestField, crumbData.Crumb)

	// Execute the request
	buildResp, err := client.Do(buildReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to trigger Jenkins job: %v", err)})
		return
	}
	defer buildResp.Body.Close()

	if buildResp.StatusCode == http.StatusForbidden {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jenkins job trigger failed: Invalid credentials or insufficient permissions. Please check your Jenkins token and user permissions."})
		return
	}

	if buildResp.StatusCode != http.StatusCreated && buildResp.StatusCode != http.StatusOK {
		// Read response body for error details
		errorBody, _ := io.ReadAll(buildResp.Body)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Jenkins job trigger failed: HTTP %d - %s", buildResp.StatusCode, string(errorBody))})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Jenkins job triggered successfully",
		"job":     req.JenkinsJob,
	})
}
