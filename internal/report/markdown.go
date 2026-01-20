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
		"severityEmoji":    SeverityEmoji,
		"qualityGateEmoji": QualityGateEmoji,
		"qualityGateText":  QualityGateText,
		"formatTime":       formatTime,
		"formatDate":       formatDate,
		"getSortedSeverities": func(m map[string][]IssueItem) []string {
			return GetSortedSeverities(m)
		},
		"truncate":      truncateString,
		"ratingEmoji":   ratingEmoji,
		"ratingColor":   ratingColor,
		"severityBadge": severityBadge,
		"priorityEmoji": priorityEmoji,
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

func ratingEmoji(rating string) string {
	switch rating {
	case "A":
		return "🟢"
	case "B":
		return "🟡"
	case "C":
		return "🟠"
	case "D":
		return "🔴"
	case "E":
		return "⛔"
	default:
		return "⚪"
	}
}

func ratingColor(rating string) string {
	switch rating {
	case "A":
		return "brightgreen"
	case "B":
		return "green"
	case "C":
		return "yellow"
	case "D":
		return "orange"
	case "E":
		return "red"
	default:
		return "lightgrey"
	}
}

func severityBadge(severity string) string {
	switch severity {
	case "BLOCKER":
		return "🔴 BLOCKER"
	case "CRITICAL":
		return "🟠 CRITICAL"
	case "MAJOR":
		return "🟡 MAJOR"
	case "MINOR":
		return "🔵 MINOR"
	case "INFO":
		return "⚪ INFO"
	default:
		return severity
	}
}

func priorityEmoji(priority string) string {
	switch priority {
	case "HIGH":
		return "🔴"
	case "MEDIUM":
		return "🟡"
	case "LOW":
		return "🟢"
	default:
		return "⚪"
	}
}

const markdownTemplate = `# 📊 SonarQube Analysis Report

### Code Quality Analysis Summary

---

## 📋 Project Information

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

## 🚦 Quality Gate

{{- if eq .QualityGateStatus "OK" }}

### ✅ PASSED

> **Congratulations!** Your code meets all quality standards.

{{- else if eq .QualityGateStatus "WARN" }}

### ⚠️ WARNING

> **Attention needed.** Some quality thresholds are close to failing.

{{- else }}

### ❌ FAILED

> **Action required!** Your code does not meet quality standards.

{{- end }}

{{- if .QualityGateConditions }}

### Quality Gate Conditions

| Metric | Status | Actual Value | Threshold |
|:-------|:------:|:------------:|:----------|
{{- range .QualityGateConditions }}
| {{ .Metric }} | {{ if eq .Status "OK" }}✅{{ else if eq .Status "WARN" }}⚠️{{ else }}❌{{ end }} | **{{ .ActualValue }}** | {{ .Comparator }} {{ .ErrorThreshold }} |
{{- end }}
{{- end }}

---

## 📈 Metrics Overview

### Code Health Dashboard

<table>
<tr>
<td align="center" width="33%">

### 🐛 Bugs
# {{ .Metrics.Bugs }}
{{ ratingEmoji .Metrics.ReliabilityRating }} Rating: **{{ .Metrics.ReliabilityRating }}**

</td>
<td align="center" width="33%">

### 🔓 Vulnerabilities  
# {{ .Metrics.Vulnerabilities }}
{{ ratingEmoji .Metrics.SecurityRating }} Rating: **{{ .Metrics.SecurityRating }}**

</td>
<td align="center" width="33%">

### 🧹 Code Smells
# {{ .Metrics.CodeSmells }}
{{ ratingEmoji .Metrics.MaintainabilityRating }} Rating: **{{ .Metrics.MaintainabilityRating }}**

</td>
</tr>
</table>

### Additional Metrics

| Metric | Value | Description |
|:-------|:-----:|:------------|
| 📏 **Lines of Code** | {{ .Metrics.LinesOfCode }} | Total lines analyzed |
| 📊 **Coverage** | {{ .Metrics.Coverage }} | Test coverage percentage |
| 📑 **Duplications** | {{ .Metrics.DuplicatedLinesDensity }} | Duplicated code percentage |
| ⏱️ **Technical Debt** | {{ .Metrics.TechnicalDebt }} | Estimated time to fix all issues |

{{- if or .Metrics.NewBugs .Metrics.NewVulnerabilities .Metrics.NewCodeSmells }}

### 🆕 New Code Analysis

| Metric | Value |
|:-------|:-----:|
{{- if .Metrics.NewBugs }}
| 🐛 New Bugs | **{{ .Metrics.NewBugs }}** |
{{- end }}
{{- if .Metrics.NewVulnerabilities }}
| 🔓 New Vulnerabilities | **{{ .Metrics.NewVulnerabilities }}** |
{{- end }}
{{- if .Metrics.NewCodeSmells }}
| 🧹 New Code Smells | **{{ .Metrics.NewCodeSmells }}** |
{{- end }}
{{- if .Metrics.NewCoverage }}
| 📊 New Coverage | **{{ .Metrics.NewCoverage }}** |
{{- end }}
{{- if .Metrics.NewDuplicatedLines }}
| 📑 New Duplications | **{{ .Metrics.NewDuplicatedLines }}** |
{{- end }}
{{- end }}

---

## 🔍 Issues Analysis

### Total Issues: **{{ .TotalIssues }}**

### Issues by Type

| Type | Count | Percentage |
|:-----|:-----:|:----------:|
{{- range $type, $count := .IssuesByType }}
| {{ if eq $type "BUG" }}🐛{{ else if eq $type "VULNERABILITY" }}🔓{{ else }}🧹{{ end }} {{ $type }} | **{{ $count }}** | {{ if $.TotalIssues }}{{ printf "%.1f%%" (mul (div (float64 $count) (float64 $.TotalIssues)) 100) }}{{ else }}0%{{ end }} |
{{- end }}

### Issues by Severity

| Severity | Count |
|:---------|:-----:|
{{- $severities := getSortedSeverities .IssuesBySeverity }}
{{- range $sev := $severities }}
| {{ severityBadge $sev }} | **{{ issueCount $.IssuesBySeverity $sev }}** |
{{- end }}

{{- $severities := getSortedSeverities .IssuesBySeverity }}
{{- range $sev := $severities }}
{{- $issues := index $.IssuesBySeverity $sev }}
{{- if $issues }}

---

### {{ severityBadge $sev }} Issues ({{ len $issues }})

<details>
<summary>Click to expand {{ $sev }} issues</summary>

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

**📝 Problematic Code:**

` + "```{{ .Language }}" + `
{{ .CodeSnippet }}
` + "```" + `

{{- end }}

{{- if hasCodeSnippet .HowToFix }}

**💡 How to Fix:**

> {{ truncate .HowToFix 500 }}

{{- end }}

---

{{- end }}
{{- end }}

{{- if gt (len $issues) 10 }}

> **Note:** Showing detailed view for first 10 of {{ len $issues }} {{ $sev }} issues.

**Additional {{ $sev }} Issues:**

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

</details>
{{- end }}
{{- end }}

---

## 🛡️ Security Hotspots

{{- if gt .TotalHotspots 0 }}

### Total Hotspots: **{{ .TotalHotspots }}**

### Hotspots by Priority

| Priority | Count |
|:---------|:-----:|
{{- range $priority, $count := .HotspotsByPriority }}
| {{ priorityEmoji $priority }} {{ $priority }} | **{{ $count }}** |
{{- end }}

{{- if .Hotspots }}

### Hotspot Details

<details>
<summary>Click to expand security hotspots</summary>

| # | Priority | Category | Location | Status |
|:-:|:--------:|:---------|:---------|:------:|
{{- range $idx, $hotspot := .Hotspots }}
{{- if lt $idx 20 }}
| {{ add $idx 1 }} | {{ priorityEmoji .VulnerabilityProbability }} {{ .VulnerabilityProbability }} | {{ .SecurityCategory }} | ` + "`{{ .Component }}:{{ .Line }}`" + ` | {{ .Status }} |
{{- end }}
{{- end }}
{{- if gt (len .Hotspots) 20 }}

> **Note:** Showing first 20 of {{ len .Hotspots }} hotspots.
{{- end }}

</details>
{{- end }}

{{- else }}

### ✅ No Security Hotspots Found

> Great job! No security hotspots detected in this analysis.

{{- end }}

---

## 📝 Summary

{{- if eq .QualityGateStatus "OK" }}

| Status | Result |
|:------:|:------:|
| Quality Gate | ✅ **PASSED** |
| Bugs | {{ .Metrics.Bugs }} ({{ .Metrics.ReliabilityRating }}) |
| Vulnerabilities | {{ .Metrics.Vulnerabilities }} ({{ .Metrics.SecurityRating }}) |
| Code Smells | {{ .Metrics.CodeSmells }} ({{ .Metrics.MaintainabilityRating }}) |

{{- else }}

| ⚠️ **Action Required** |
|:----------------------:|
| Quality Gate: **{{ qualityGateText .QualityGateStatus }}** |
| Please review and fix the issues above |

{{- end }}

---

*Report generated by **SonarQube Report Generator***  
*{{ formatTime .GeneratedAt }}*
`
