package handler

import (
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
