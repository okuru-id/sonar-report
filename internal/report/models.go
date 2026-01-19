package report

import "time"

// ReportData contains all data needed for report generation
type ReportData struct {
	// Project info
	ProjectKey   string    `json:"projectKey"`
	ProjectName  string    `json:"projectName"`
	Branch       string    `json:"branch"`
	GeneratedAt  time.Time `json:"generatedAt"`
	AnalysisDate string    `json:"analysisDate,omitempty"`

	// Quality Gate
	QualityGateStatus     string            `json:"qualityGateStatus"` // PASSED, FAILED, WARNING
	QualityGateConditions []ConditionResult `json:"qualityGateConditions,omitempty"`

	// Metrics
	Metrics MetricsSummary `json:"metrics"`

	// Issues
	TotalIssues      int                    `json:"totalIssues"`
	IssuesByType     map[string]int         `json:"issuesByType"`
	IssuesBySeverity map[string][]IssueItem `json:"issuesBySeverity"`

	// Hotspots
	TotalHotspots      int            `json:"totalHotspots"`
	Hotspots           []HotspotItem  `json:"hotspots"`
	HotspotsByPriority map[string]int `json:"hotspotsByPriority"`
}

// ConditionResult represents a quality gate condition result
type ConditionResult struct {
	Metric         string `json:"metric"`
	Status         string `json:"status"`
	ActualValue    string `json:"actualValue"`
	ErrorThreshold string `json:"errorThreshold"`
	Comparator     string `json:"comparator"`
}

// MetricsSummary contains summarized metrics
type MetricsSummary struct {
	Bugs                   string `json:"bugs"`
	Vulnerabilities        string `json:"vulnerabilities"`
	CodeSmells             string `json:"codeSmells"`
	Coverage               string `json:"coverage"`
	DuplicatedLinesDensity string `json:"duplicatedLinesDensity"`
	LinesOfCode            string `json:"linesOfCode"`
	TechnicalDebt          string `json:"technicalDebt"`

	// Ratings (A-E)
	ReliabilityRating     string `json:"reliabilityRating"`
	SecurityRating        string `json:"securityRating"`
	MaintainabilityRating string `json:"maintainabilityRating"`

	// New code metrics
	NewBugs            string `json:"newBugs,omitempty"`
	NewVulnerabilities string `json:"newVulnerabilities,omitempty"`
	NewCodeSmells      string `json:"newCodeSmells,omitempty"`
	NewCoverage        string `json:"newCoverage,omitempty"`
	NewDuplicatedLines string `json:"newDuplicatedLines,omitempty"`
}

// IssueItem represents an issue for display
type IssueItem struct {
	Key         string `json:"key"`
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Message     string `json:"message"`
	Component   string `json:"component"`
	Line        int    `json:"line,omitempty"`
	EndLine     int    `json:"endLine,omitempty"`
	Effort      string `json:"effort,omitempty"`
	Rule        string `json:"rule"`
	CodeSnippet string `json:"codeSnippet,omitempty"` // Source code snippet
	HowToFix    string `json:"howToFix,omitempty"`    // Rule description / how to fix
	Language    string `json:"language,omitempty"`    // Programming language for syntax highlighting
}

// HotspotItem represents a security hotspot for display
type HotspotItem struct {
	Key                      string `json:"key"`
	SecurityCategory         string `json:"securityCategory"`
	VulnerabilityProbability string `json:"vulnerabilityProbability"`
	Status                   string `json:"status"`
	Message                  string `json:"message"`
	Component                string `json:"component"`
	Line                     int    `json:"line,omitempty"`
}

// ReportRecord represents a saved report record
type ReportRecord struct {
	ID          string    `json:"id"`
	ProjectKey  string    `json:"projectKey"`
	ProjectName string    `json:"projectName"`
	Branch      string    `json:"branch"`
	Format      string    `json:"format"` // md, pdf
	FileName    string    `json:"fileName"`
	FilePath    string    `json:"filePath"`
	FileSize    int64     `json:"fileSize"`
	GeneratedAt time.Time `json:"generatedAt"`
}

// GenerateRequest represents a report generation request
type GenerateRequest struct {
	ProjectKey string `json:"projectKey" form:"projectKey" binding:"required"`
	Branch     string `json:"branch" form:"branch"`
	Format     string `json:"format" form:"format"` // md, pdf
}

// RatingToLetter converts a numeric rating to letter grade
func RatingToLetter(rating string) string {
	switch rating {
	case "1", "1.0":
		return "A"
	case "2", "2.0":
		return "B"
	case "3", "3.0":
		return "C"
	case "4", "4.0":
		return "D"
	case "5", "5.0":
		return "E"
	default:
		return rating
	}
}

// SeverityOrder returns the order for severity sorting
func SeverityOrder(severity string) int {
	switch severity {
	case "BLOCKER":
		return 0
	case "CRITICAL":
		return 1
	case "MAJOR":
		return 2
	case "MINOR":
		return 3
	case "INFO":
		return 4
	default:
		return 5
	}
}

// SeverityEmoji returns an emoji for severity
func SeverityEmoji(severity string) string {
	switch severity {
	case "BLOCKER":
		return "🔴"
	case "CRITICAL":
		return "🟠"
	case "MAJOR":
		return "🟡"
	case "MINOR":
		return "🔵"
	case "INFO":
		return "⚪"
	default:
		return "⚫"
	}
}

// QualityGateEmoji returns an emoji for quality gate status
func QualityGateEmoji(status string) string {
	switch status {
	case "OK":
		return "✅"
	case "WARN":
		return "⚠️"
	case "ERROR":
		return "❌"
	default:
		return "❓"
	}
}

// QualityGateText returns human-readable text for quality gate status
func QualityGateText(status string) string {
	switch status {
	case "OK":
		return "PASSED"
	case "WARN":
		return "WARNING"
	case "ERROR":
		return "FAILED"
	default:
		return status
	}
}
