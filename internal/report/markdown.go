package report

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

// MarkdownGenerator generates markdown reports
type MarkdownGenerator struct{}

// NewMarkdownGenerator creates a new markdown generator
func NewMarkdownGenerator() *MarkdownGenerator {
	return &MarkdownGenerator{}
}

// Generate generates a markdown report
func (g *MarkdownGenerator) Generate(data *ReportData) ([]byte, error) {
	tmpl, err := template.New("report").Funcs(template.FuncMap{
		"severityIcon":    severityIcon,
		"qualityGateIcon": qualityGateIcon,
		"qualityGateText": QualityGateText,
		"formatTime":      formatTime,
		"formatDate":      formatDate,
		"getSortedSeverities": func(m map[string][]IssueItem) []string {
			return GetSortedSeverities(m)
		},
		"truncate":     truncateString,
		"ratingIcon":   ratingIcon,
		"priorityIcon": priorityIcon,
		"icon":         icon,
		"issueCount": func(m map[string][]IssueItem, sev string) int {
			return len(m[sev])
		},
		"hasCodeSnippet": func(s string) bool {
			return s != ""
		},
		"add": func(a, b int) int {
			return a + b
		},
		"mul": func(a, b float64) float64 {
			return a * b
		},
		"div": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"float64": func(i int) float64 {
			return float64(i)
		},
	}).Parse(markdownTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.Bytes(), nil
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func formatDate(t time.Time) string {
	return t.Format("January 02, 2006")
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func icon(name string, color string) string {
	colors := map[string]string{
		"success": "#22c55e",
		"warning": "#f59e0b",
		"danger":  "#ef4444",
		"info":    "#3b82f6",
		"gray":    "#94a3b8",
	}
	hexColor := colors[color]
	if hexColor == "" {
		hexColor = color
	}

	icons := map[string]string{
		"chart-bar":      `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="` + hexColor + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M18 20V10"/><path d="M12 20V4"/><path d="M6 20v-6"/></svg>`,
		"info-circle":    `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="` + hexColor + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><circle cx="12" cy="12" r="10"/><path d="M12 16v-4"/><path d="M12 8h.01"/></svg>`,
		"activity":       `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="` + hexColor + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M22 12h-4l-3 9L9 3l-3 9H2"/></svg>`,
		"check-circle":   `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="` + hexColor + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><path d="M22 4L12 14.01l-3-3"/></svg>`,
		"alert-triangle": `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="` + hexColor + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><path d="M12 9v4"/><path d="M12 17h.01"/></svg>`,
		"circle-x":       `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="` + hexColor + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><circle cx="12" cy="12" r="10"/><path d="m15 9-6 6"/><path d="m9 9 6 6"/></svg>`,
		"bug":            `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="` + hexColor + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><rect width="8" height="14" x="8" y="6" rx="4"/><path d="m19 7-3 2"/><path d="m5 7 3 2"/><path d="m19 19-3-2"/><path d="m5 19 3-2"/><path d="M20 13h-4"/><path d="M4 13h4"/><path d="m10 4 1 2"/></svg>`,
		"shield":         `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="` + hexColor + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>`,
		"broom":          `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="` + hexColor + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="m16 5-6.4 11.2a2.1 2.1 0 0 1-1.2 1.3 2.1 2.1 0 0 1-1.2-.3L4.6 16"/><path d="m10 5 1.7 3.4a2.1 2.1 0 0 1 .4 1.5L9 16"/><path d="m14 5 6.4 11.2a2.1 2.1 0 0 1 1.2 1.3 2.1 2.1 0 0 1-1.2.3l-2.4-2"/></svg>`,
		"search":         `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="` + hexColor + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><circle cx="10" cy="10" r="7"/><path d="m21 21-4.3-4.3"/></svg>`,
		"code":           `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="` + hexColor + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><polyline points="16 18 22 12 16 6"/><polyline points="8 6 2 12 8 18"/></svg>`,
		"bulb":           `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="` + hexColor + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M9 18h6"/><path d="M10 22h4"/><path d="M15.09 14c.18-.98.65-1.74 1.41-2.5A4.65 4.65 0 0 0 18 8 6 6 0 0 0 6 8c0 1 .23 2.23 1.5 3.5A4.61 4.61 0 0 1 8.91 14"/></svg>`,
		"ruler":          `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="` + hexColor + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M21.21 15.89A10 10 0 1 1 8 2.83"/><path d="M22 12A10 10 0 0 0 12 2v10z"/></svg>`,
		"copy":           `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="` + hexColor + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><rect width="14" height="14" x="8" y="8" rx="2" ry="2"/><path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"/></svg>`,
		"clock":          `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="` + hexColor + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>`,
		"sparkles":       `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="` + hexColor + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="m12 3-1.912 5.813a2 2 0 0 1-1.275 1.275L3 12l5.813 1.912a2 2 0 0 1 1.275 1.275L12 21l1.912-5.813a2 2 0 0 1 1.275-1.275L21 12l-5.813-1.912a2 2 0 0 1-1.275-1.275L12 3Z"/><path d="M5 3v4"/><path d="M9 3v4"/><path d="M1 7h4"/><path d="M3 5h4"/><path d="M3 7h4"/><path d="M1 11h4"/></svg>`,
		"list":           `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="` + hexColor + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><line x1="8" x2="21" y1="6" y2="6"/><line x1="8" x2="21" y1="12" y2="12"/><line x1="8" x2="21" y1="18" y2="18"/><line x1="3" x2="3.01" y1="6" y2="6"/><line x1="3" x2="3.01" y1="12" y2="12"/><line x1="3" x2="3.01" y1="18" y2="18"/></svg>`,
		"trending-up":    `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="` + hexColor + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><polyline points="23 6 13.5 15.5 8.5 10.5 1 18"/><polyline points="17 6 23 6 23 12"/></svg>`,
		"folder":         `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="` + hexColor + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>`,
		"circle-filled":  `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="` + hexColor + `" stroke="` + hexColor + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><circle cx="12" cy="12" r="10"/></svg>`,
		"circle":         `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="` + hexColor + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><circle cx="12" cy="12" r="10"/></svg>`,
	}

	if i, ok := icons[name]; ok {
		return i
	}
	return ""
}

func ratingIcon(rating string) string {
	switch rating {
	case "A":
		return icon("circle", "#22c55e")
	case "B":
		return icon("circle", "#84cc16")
	case "C":
		return icon("circle", "#f59e0b")
	case "D":
		return icon("circle", "#f97316")
	case "E":
		return icon("circle", "#ef4444")
	default:
		return icon("circle", "#94a3b8")
	}
}

func severityIcon(severity string) string {
	switch severity {
	case "BLOCKER":
		return icon("circle", "#ef4444") + " BLOCKER"
	case "CRITICAL":
		return icon("circle", "#f97316") + " CRITICAL"
	case "MAJOR":
		return icon("circle", "#f59e0b") + " MAJOR"
	case "MINOR":
		return icon("circle", "#3b82f6") + " MINOR"
	case "INFO":
		return icon("circle", "#94a3b8") + " INFO"
	default:
		return severity
	}
}

func priorityIcon(priority string) string {
	switch priority {
	case "HIGH":
		return icon("circle", "#ef4444")
	case "MEDIUM":
		return icon("circle", "#f59e0b")
	case "LOW":
		return icon("circle", "#22c55e")
	default:
		return icon("circle", "#94a3b8")
	}
}

func qualityGateIcon(status string) string {
	switch status {
	case "OK":
		return icon("check-circle", "#22c55e")
	case "WARN":
		return icon("alert-triangle", "#f59e0b")
	default:
		return icon("circle-x", "#ef4444")
	}
}

const markdownTemplate = `# <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M18 20V10"/><path d="M12 20V4"/><path d="M6 20v-6"/></svg> SonarQube Analysis Report

### Code Quality Analysis Summary

---

## <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><circle cx="12" cy="12" r="10"/><path d="M12 16v-4"/><path d="M12 8h.01"/></svg> Project Information

| | |
|---|---|
| **Project Name** | {{ .ProjectName }} |
| **Project Key** | ` + "`{{ .ProjectKey }}`" + ` |
| **Branch** | ` + "`{{ .Branch }}`" + ` |
| **Report Generated** | {{ formatTime .GeneratedAt }} |
{{- if .AnalysisDate }}
| **Last Analysis** | {{ .AnalysisDate }} |
{{- end }}

---

## <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M22 12h-4l-3 9L9 3l-3 9H2"/></svg> Quality Gate

{{- if eq .QualityGateStatus "OK" }}

### <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#22c55e" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><path d="M22 4L12 14.01l-3-3"/></svg> PASSED

> **Congratulations!** Your code meets all quality standards.

{{- else if eq .QualityGateStatus "WARN" }}

### <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#f59e0b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><path d="M12 9v4"/><path d="M12 17h.01"/></svg> WARNING

> **Attention needed.** Some quality thresholds are close to failing.

{{- else }}

### <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#ef4444" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><circle cx="12" cy="12" r="10"/><path d="m15 9-6 6"/><path d="m9 9 6 6"/></svg> FAILED

> **Action required!** Your code does not meet quality standards.

{{- end }}

{{- if .QualityGateConditions }}

### Quality Gate Conditions

| Metric | Status | Actual Value | Threshold |
|:-------|:------:|:------------:|:----------|
{{- range .QualityGateConditions }}
| {{ .Metric }} | {{ if eq .Status "OK" }}<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#22c55e" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><path d="M22 4L12 14.01l-3-3"/></svg>{{ else if eq .Status "WARN" }}<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#f59e0b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><path d="M12 9v4"/><path d="M12 17h.01"/></svg>{{ else }}<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#ef4444" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><circle cx="12" cy="12" r="10"/><path d="m15 9-6 6"/><path d="m9 9 6 6"/></svg>{{ end }} | **{{ .ActualValue }}** | {{ .Comparator }} {{ .ErrorThreshold }} |
{{- end }}
{{- end }}

---

## <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><polyline points="23 6 13.5 15.5 8.5 10.5 1 18"/><polyline points="17 6 23 6 23 12"/></svg> Metrics Overview

### Code Health Dashboard

| <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#ef4444" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><rect width="8" height="14" x="8" y="6" rx="4"/><path d="m19 7-3 2"/><path d="m5 7 3 2"/><path d="m19 19-3-2"/><path d="m5 19 3-2"/><path d="M20 13h-4"/><path d="M4 13h4"/><path d="m10 4 1 2"/></svg> Bugs | <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#f59e0b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg> Vulnerabilities | <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="m16 5-6.4 11.2a2.1 2.1 0 0 1-1.2 1.3 2.1 2.1 0 0 1-1.2.3L4.6 16"/><path d="m10 5 1.7 3.4a2.1 2.1 0 0 1 .4 1.5L9 16"/><path d="m14 5 6.4 11.2a2.1 2.1 0 0 1 1.2 1.3 2.1 2.1 0 0 1-1.2.3l-2.4-2"/></svg> Code Smells |
|:---|:---|:---|
| **{{ .Metrics.Bugs }}** | **{{ .Metrics.Vulnerabilities }}** | **{{ .Metrics.CodeSmells }}** |
| Rating: {{ .Metrics.ReliabilityRating }} | Rating: {{ .Metrics.SecurityRating }} | Rating: {{ .Metrics.MaintainabilityRating }} |

### Additional Metrics

| Metric | Value | Description |
|:-------|:-----:|:------------|
| <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M21.21 15.89A10 10 0 1 1 8 2.83"/><path d="M22 12A10 10 0 0 0 12 2v10z"/></svg> **Lines of Code** | {{ .Metrics.LinesOfCode }} | Total lines analyzed |
| <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M18 20V10"/><path d="M12 20V4"/><path d="M6 20v-6"/></svg> **Coverage** | {{ .Metrics.Coverage }} | Test coverage percentage |
| <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><rect width="14" height="14" x="8" y="8" rx="2" ry="2"/><path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"/></svg> **Duplications** | {{ .Metrics.DuplicatedLinesDensity }} | Duplicated code percentage |
| <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg> **Technical Debt** | {{ .Metrics.TechnicalDebt }} | Estimated time to fix all issues |

{{- if or .Metrics.NewBugs .Metrics.NewVulnerabilities .Metrics.NewCodeSmells }}

### <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#f59e0b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="m12 3-1.912 5.813a2 2 0 0 1-1.275 1.275L3 12l5.813 1.912a2 2 0 0 1 1.275 1.275L12 21l1.912-5.813a2 2 0 0 1 1.275-1.275L21 12l-5.813-1.912a2 2 0 0 1-1.275-1.275L12 3Z"/><path d="M5 3v4"/><path d="M9 3v4"/><path d="M1 7h4"/><path d="M3 5h4"/><path d="M3 7h4"/><path d="M1 11h4"/></svg> New Code Analysis

| Metric | Value |
|:-------|:-----:|
{{- if .Metrics.NewBugs }}
| <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#ef4444" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><rect width="8" height="14" x="8" y="6" rx="4"/><path d="m19 7-3 2"/><path d="m5 7 3 2"/><path d="m19 19-3-2"/><path d="m5 19 3-2"/><path d="M20 13h-4"/><path d="M4 13h4"/><path d="m10 4 1 2"/></svg> New Bugs | **{{ .Metrics.NewBugs }}** |
{{- end }}
{{- if .Metrics.NewVulnerabilities }}
| <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#f59e0b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg> New Vulnerabilities | **{{ .Metrics.NewVulnerabilities }}** |
{{- end }}
{{- if .Metrics.NewCodeSmells }}
| <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="m16 5-6.4 11.2a2.1 2.1 0 0 1-1.2 1.3 2.1 2.1 0 0 1-1.2.3L4.6 16"/><path d="m10 5 1.7 3.4a2.1 2.1 0 0 1 .4 1.5L9 16"/><path d="m14 5 6.4 11.2a2.1 2.1 0 0 1 1.2 1.3 2.1 2.1 0 0 1-1.2.3l-2.4-2"/></svg> New Code Smells | **{{ .Metrics.NewCodeSmells }}** |
{{- end }}
{{- if .Metrics.NewCoverage }}
| <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M18 20V10"/><path d="M12 20V4"/><path d="M6 20v-6"/></svg> New Coverage | **{{ .Metrics.NewCoverage }}** |
{{- end }}
{{- if .Metrics.NewDuplicatedLines }}
| <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><rect width="14" height="14" x="8" y="8" rx="2" ry="2"/><path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"/></svg> New Duplications | **{{ .Metrics.NewDuplicatedLines }}** |
{{- end }}
{{- end }}

---

## <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><circle cx="10" cy="10" r="7"/><path d="m21 21-4.3-4.3"/></svg> Issues Analysis

### Total Issues: **{{ .TotalIssues }}**

### Issues by Type

| Type | Count | Percentage |
|:-----|:-----:|:----------:|
{{- range $type, $count := .IssuesByType }}
| {{ if eq $type "BUG" }}<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#ef4444" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><rect width="8" height="14" x="8" y="6" rx="4"/><path d="m19 7-3 2"/><path d="m5 7 3 2"/><path d="m19 19-3-2"/><path d="m5 19 3-2"/><path d="M20 13h-4"/><path d="M4 13h4"/><path d="m10 4 1 2"/></svg>{{ else if eq $type "VULNERABILITY" }}<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#f59e0b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>{{ else }}<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="m16 5-6.4 11.2a2.1 2.1 0 0 1-1.2 1.3 2.1 2.1 0 0 1-1.2.3L4.6 16"/><path d="m10 5 1.7 3.4a2.1 2.1 0 0 1 .4 1.5L9 16"/><path d="m14 5 6.4 11.2a2.1 2.1 0 0 1 1.2 1.3 2.1 2.1 0 0 1-1.2.3l-2.4-2"/></svg>{{ end }} {{ $type }} | **{{ $count }}** | {{ if $.TotalIssues }}{{ printf "%.1f%%" (mul (div (float64 $count) (float64 $.TotalIssues)) 100) }}{{ else }}0%{{ end }} |
{{- end }}

### Issues by Severity

| Severity | Count |
|:---------|:-----:|
{{- $severities := getSortedSeverities .IssuesBySeverity }}
{{- range $sev := $severities }}
| {{ severityIcon $sev }} | **{{ issueCount $.IssuesBySeverity $sev }}** |
{{- end }}

{{- $severities := getSortedSeverities .IssuesBySeverity }}
{{- range $sev := $severities }}
{{- $issues := index $.IssuesBySeverity $sev }}
{{- if $issues }}

---

### {{ severityIcon $sev }} Issues ({{ len $issues }})

{{- if gt (len $issues) 10 }}

**{{ severityIcon $sev }} Issues List:**

| # | File | Line | Message |
|:-:|:-----|:----:|:--------|
{{- range $idx, $issue := $issues }}
{{- if and (ge $idx 10) (lt $idx 25) }}
| {{ add $idx 1 }} | ` + "`{{ .Component }}`" + ` | {{ .Line }} | {{ truncate .Message 60 }} |
{{- end }}
{{- end }}
{{- if gt (len $issues) 25 }}

> Showing 25 of {{ len $issues }} {{ $sev }} issues. See SonarQube for full list.
{{- end }}

{{- end }}

<details>
<summary>Click to expand {{ $sev }} issues with details and code</summary>

{{- range $idx, $issue := $issues }}
{{- if lt $idx 10 }}

#### {{ add $idx 1 }}. {{ .Message }}

| Property | Value |
|:---------|:------|
| **File** | ` + "`{{ .Component }}`" + ` |
| **Line** | {{ .Line }}{{ if and .EndLine (ne .EndLine .Line) }} - {{ .EndLine }}{{ end }} |
| **Type** | {{ .Type }} |
| **Rule** | ` + "`{{ .Rule }}`" + ` |
{{- if .Effort }}
| **Effort** | {{ .Effort }} |
{{- end }}

{{- if hasCodeSnippet .CodeSnippet }}

**<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><polyline points="16 18 22 12 16 6"/><polyline points="8 6 2 12 8 18"/></svg> Problematic Code:**

` + "```{{ .Language }}" + `
{{ .CodeSnippet }}
` + "```" + `

{{- end }}

{{- if hasCodeSnippet .HowToFix }}

**<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#f59e0b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M9 18h6"/><path d="M10 22h4"/><path d="M15.09 14c.18-.98.65-1.74 1.41-2.5A4.65 4.65 0 0 0 18 8 6 6 0 0 0 6 8c0 1 .23 2.23 1.5 3.5A4.61 4.61 0 0 1 8.91 14"/></svg> How to Fix:**

> {{ truncate .HowToFix 500 }}

{{- end }}

---

{{- end }}
{{- end }}

{{- if gt (len $issues) 10 }}

> **Note:** Showing detailed view for first 10 of {{ len $issues }} {{ $sev }} issues.

{{- end }}

</details>

{{- end }}
{{- end }}

---

## <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg> Security Hotspots

{{- if gt .TotalHotspots 0 }}

### Total Hotspots: **{{ .TotalHotspots }}**

### Hotspots by Priority

| Priority | Count |
|:---------|:-----:|
{{- range $priority, $count := .HotspotsByPriority }}
| {{ priorityIcon $priority }} {{ $priority }} | **{{ $count }}** |
{{- end }}

{{- if .Hotspots }}

### Hotspot Details

<details>
<summary>Click to expand security hotspots</summary>

| # | Priority | Category | Location | Status |
|:-:|:--------:|:---------|:---------|:------:|
{{- range $idx, $hotspot := .Hotspots }}
{{- if lt $idx 20 }}
| {{ add $idx 1 }} | {{ priorityIcon .VulnerabilityProbability }} {{ .VulnerabilityProbability }} | {{ .SecurityCategory }} | ` + "`{{ .Component }}:{{ .Line }}`" + ` | {{ .Status }} |
{{- end }}
{{- end }}
{{- if gt (len .Hotspots) 20 }}

> **Note:** Showing first 20 of {{ len .Hotspots }} hotspots.
{{- end }}

</details>
{{- end }}

{{- else }}

### <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#22c55e" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><path d="M22 4L12 14.01l-3-3"/></svg> No Security Hotspots Found

> Great job! No security hotspots detected in this analysis.

{{- end }}

---

## <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><line x1="8" x2="21" y1="6" y2="6"/><line x1="8" x2="21" y1="12" y2="12"/><line x1="8" x2="21" y1="18" y2="18"/><line x1="3" x2="3.01" y1="6" y2="6"/><line x1="3" x2="3.01" y1="12" y2="12"/><line x1="3" x2="3.01" y1="18" y2="18"/></svg> Summary

{{- if eq .QualityGateStatus "OK" }}

| Status | Result |
|:------:|:------:|
| Quality Gate | <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#22c55e" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><path d="M22 4L12 14.01l-3-3"/></svg> **PASSED** |
| Bugs | {{ .Metrics.Bugs }} ({{ .Metrics.ReliabilityRating }}) |
| Vulnerabilities | {{ .Metrics.Vulnerabilities }} ({{ .Metrics.SecurityRating }}) |
| Code Smells | {{ .Metrics.CodeSmells }} ({{ .Metrics.MaintainabilityRating }}) |

{{- else }}

| <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#f59e0b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><path d="M12 9v4"/><path d="M12 17h.01"/></svg> **Action Required** |
|:----------------------:|
| Quality Gate: **{{ qualityGateText .QualityGateStatus }}** |
| Please review and fix the issues above |

{{- end }}

---

*Report generated by **SonarQube Report Generator***  
*{{ formatTime .GeneratedAt }}*
`
