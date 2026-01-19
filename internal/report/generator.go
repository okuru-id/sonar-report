package report

import (
	"fmt"
	"html"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"sonarqube-report-generator/internal/sonarqube"
)

// stripHTMLTags removes HTML tags and decodes HTML entities from code
func stripHTMLTags(s string) string {
	// Remove HTML tags (like <span class="...">...</span>)
	tagRegex := regexp.MustCompile(`<[^>]*>`)
	result := tagRegex.ReplaceAllString(s, "")
	// Decode HTML entities (like &lt; &gt; &amp; &quot;)
	result = html.UnescapeString(result)
	return result
}

// Generator generates reports from SonarQube data
type Generator struct {
	client *sonarqube.Client
}

// GenerateOptions contains options for report generation
type GenerateOptions struct {
	IncludeCodeSnippets bool // Include code snippets in issues (default: true)
	IncludeHowToFix     bool // Include how to fix info from rules (default: true)
}

// NewGenerator creates a new report generator
func NewGenerator(client *sonarqube.Client) *Generator {
	return &Generator{client: client}
}

// Generate generates a report for a project
func (g *Generator) Generate(projectKey, branch string, options GenerateOptions) (*ReportData, error) {
	// Get project info
	projects, err := g.client.GetProjects()
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	var projectName string
	for _, p := range projects {
		if p.Key == projectKey {
			projectName = p.Name
			break
		}
	}
	if projectName == "" {
		projectName = projectKey
	}

	// If no branch specified, get main branch
	if branch == "" {
		branches, err := g.client.GetBranches(projectKey)
		if err == nil {
			for _, b := range branches {
				if b.IsMain {
					branch = b.Name
					break
				}
			}
		}
	}

	// Get quality gate status
	qgStatus, err := g.client.GetQualityGateStatus(projectKey, branch)
	if err != nil {
		return nil, fmt.Errorf("failed to get quality gate status: %w", err)
	}

	// Get measures
	measures, err := g.client.GetMeasures(projectKey, branch, sonarqube.DefaultMetricKeys())
	if err != nil {
		return nil, fmt.Errorf("failed to get measures: %w", err)
	}

	// Get issues
	issues, totalIssues, err := g.client.GetIssues(projectKey, branch, 500)
	if err != nil {
		return nil, fmt.Errorf("failed to get issues: %w", err)
	}

	// Get hotspots
	hotspots, totalHotspots, err := g.client.GetHotspots(projectKey, branch, 100)
	if err != nil {
		// Hotspots API might not be available in all editions
		hotspots = []sonarqube.Hotspot{}
		totalHotspots = 0
	}

	// Get latest analysis date
	analyses, err := g.client.GetAnalyses(projectKey, branch, 1)
	var analysisDate string
	if err == nil && len(analyses) > 0 {
		analysisDate = analyses[0].Date
	}

	// Build report data
	reportData := &ReportData{
		ProjectKey:   projectKey,
		ProjectName:  projectName,
		Branch:       branch,
		GeneratedAt:  time.Now(),
		AnalysisDate: analysisDate,
	}

	// Quality gate
	reportData.QualityGateStatus = qgStatus.Status
	for _, cond := range qgStatus.Conditions {
		reportData.QualityGateConditions = append(reportData.QualityGateConditions, ConditionResult{
			Metric:         cond.MetricKey,
			Status:         cond.Status,
			ActualValue:    cond.ActualValue,
			ErrorThreshold: cond.ErrorThreshold,
			Comparator:     cond.Comparator,
		})
	}

	// Metrics
	reportData.Metrics = buildMetricsSummary(measures)

	// Issues
	reportData.TotalIssues = totalIssues
	reportData.IssuesByType = make(map[string]int)
	reportData.IssuesBySeverity = make(map[string][]IssueItem)

	// Cache for rule descriptions to avoid duplicate API calls
	ruleCache := make(map[string]string)
	var ruleCacheMu sync.Mutex

	// Limit code snippet fetching to top issues per severity (to avoid too many API calls)
	maxCodeSnippetsPerSeverity := 10

	// First pass: create all issue items and count issues for snippet fetching
	type issueWithIndex struct {
		issue sonarqube.Issue
		index int
	}
	issueItems := make([]IssueItem, len(issues))
	issuesToFetch := make([]issueWithIndex, 0)
	severityCount := make(map[string]int)

	for i, issue := range issues {
		// Count by type
		reportData.IssuesByType[issue.Type]++

		// Determine end line from TextRange
		endLine := issue.Line
		if issue.TextRange != nil {
			endLine = issue.TextRange.EndLine
		}

		// Create issue item
		issueItems[i] = IssueItem{
			Key:       issue.Key,
			Type:      issue.Type,
			Severity:  issue.Severity,
			Message:   issue.Message,
			Component: extractFileName(issue.Component),
			Line:      issue.Line,
			EndLine:   endLine,
			Effort:    issue.Effort,
			Rule:      issue.Rule,
			Language:  getLanguageFromFile(issue.Component),
		}

		// Track which issues need code snippet fetching (only if enabled)
		if options.IncludeCodeSnippets || options.IncludeHowToFix {
			currentCount := severityCount[issue.Severity]
			if currentCount < maxCodeSnippetsPerSeverity {
				issuesToFetch = append(issuesToFetch, issueWithIndex{issue: issue, index: i})
				severityCount[issue.Severity]++
			}
		}
	}

	// Parallel fetch code snippets and how to fix using worker pool (only if enabled)
	if len(issuesToFetch) > 0 && (options.IncludeCodeSnippets || options.IncludeHowToFix) {
		const numWorkers = 5 // Limit concurrent API calls
		jobs := make(chan issueWithIndex, len(issuesToFetch))
		var wg sync.WaitGroup

		// Start workers
		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for job := range jobs {
					// Fetch code snippet if enabled
					if options.IncludeCodeSnippets {
						issueItems[job.index].CodeSnippet = g.fetchCodeSnippet(job.issue)
					}

					// Fetch how to fix from rule if enabled
					if options.IncludeHowToFix {
						issueItems[job.index].HowToFix = g.fetchHowToFix(job.issue.Rule, ruleCache, &ruleCacheMu)
					}
				}
			}()
		}

		// Send jobs to workers
		for _, job := range issuesToFetch {
			jobs <- job
		}
		close(jobs)

		// Wait for all workers to complete
		wg.Wait()
	}

	// Group issues by severity
	for _, item := range issueItems {
		reportData.IssuesBySeverity[item.Severity] = append(reportData.IssuesBySeverity[item.Severity], item)
	}

	// Hotspots
	reportData.TotalHotspots = totalHotspots
	reportData.HotspotsByPriority = make(map[string]int)

	for _, hotspot := range hotspots {
		reportData.HotspotsByPriority[hotspot.VulnerabilityProbability]++
		reportData.Hotspots = append(reportData.Hotspots, HotspotItem{
			Key:                      hotspot.Key,
			SecurityCategory:         hotspot.SecurityCategory,
			VulnerabilityProbability: hotspot.VulnerabilityProbability,
			Status:                   hotspot.Status,
			Message:                  hotspot.Message,
			Component:                extractFileName(hotspot.Component),
			Line:                     hotspot.Line,
		})
	}

	return reportData, nil
}

func buildMetricsSummary(measures []sonarqube.Measure) MetricsSummary {
	summary := MetricsSummary{}

	metricsMap := make(map[string]string)
	for _, m := range measures {
		metricsMap[m.Metric] = m.Value
	}

	summary.Bugs = getMetricValue(metricsMap, "bugs", "0")
	summary.Vulnerabilities = getMetricValue(metricsMap, "vulnerabilities", "0")
	summary.CodeSmells = getMetricValue(metricsMap, "code_smells", "0")
	summary.Coverage = formatPercentage(getMetricValue(metricsMap, "coverage", "0"))
	summary.DuplicatedLinesDensity = formatPercentage(getMetricValue(metricsMap, "duplicated_lines_density", "0"))
	summary.LinesOfCode = getMetricValue(metricsMap, "ncloc", "0")
	summary.TechnicalDebt = formatDebt(getMetricValue(metricsMap, "sqale_index", "0"))

	summary.ReliabilityRating = RatingToLetter(getMetricValue(metricsMap, "reliability_rating", "1"))
	summary.SecurityRating = RatingToLetter(getMetricValue(metricsMap, "security_rating", "1"))
	summary.MaintainabilityRating = RatingToLetter(getMetricValue(metricsMap, "sqale_rating", "1"))

	// New code metrics
	summary.NewBugs = getMetricValue(metricsMap, "new_bugs", "")
	summary.NewVulnerabilities = getMetricValue(metricsMap, "new_vulnerabilities", "")
	summary.NewCodeSmells = getMetricValue(metricsMap, "new_code_smells", "")
	summary.NewCoverage = formatPercentage(getMetricValue(metricsMap, "new_coverage", ""))
	summary.NewDuplicatedLines = formatPercentage(getMetricValue(metricsMap, "new_duplicated_lines_density", ""))

	return summary
}

func getMetricValue(metrics map[string]string, key, defaultValue string) string {
	if val, ok := metrics[key]; ok {
		return val
	}
	return defaultValue
}

func formatPercentage(value string) string {
	if value == "" {
		return ""
	}
	return value + "%"
}

func formatDebt(minutes string) string {
	if minutes == "" || minutes == "0" {
		return "0min"
	}

	var m int
	fmt.Sscanf(minutes, "%d", &m)

	if m < 60 {
		return fmt.Sprintf("%dmin", m)
	}

	hours := m / 60
	mins := m % 60

	if hours < 24 {
		if mins > 0 {
			return fmt.Sprintf("%dh %dmin", hours, mins)
		}
		return fmt.Sprintf("%dh", hours)
	}

	days := hours / 8 // 8 hour work day
	remainingHours := hours % 8

	if remainingHours > 0 {
		return fmt.Sprintf("%dd %dh", days, remainingHours)
	}
	return fmt.Sprintf("%dd", days)
}

func extractFileName(component string) string {
	// Component is usually in format "project:src/path/file.go"
	parts := strings.SplitN(component, ":", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return filepath.Base(component)
}

// GetSortedSeverities returns severities in order
func GetSortedSeverities(issuesBySeverity map[string][]IssueItem) []string {
	severities := make([]string, 0, len(issuesBySeverity))
	for sev := range issuesBySeverity {
		severities = append(severities, sev)
	}
	sort.Slice(severities, func(i, j int) bool {
		return SeverityOrder(severities[i]) < SeverityOrder(severities[j])
	})
	return severities
}

// fetchCodeSnippet fetches source code for an issue
func (g *Generator) fetchCodeSnippet(issue sonarqube.Issue) string {
	// Use the issue's reported location as the primary source
	issueStartLine := issue.Line
	issueEndLine := issue.Line
	component := issue.Component

	if issue.TextRange != nil {
		issueStartLine = issue.TextRange.StartLine
		issueEndLine = issue.TextRange.EndLine
	}

	// Only use flows for issues where the reported line is at file level (line <= 5)
	// and flows provide a more specific location within the actual code
	// This helps for issues like Cognitive Complexity where the issue is on a function
	// but reported line might be 0 or file header
	if issueStartLine <= 5 && len(issue.Flows) > 0 && len(issue.Flows[0].Locations) > 0 {
		firstLoc := issue.Flows[0].Locations[0]
		if firstLoc.TextRange != nil && firstLoc.TextRange.StartLine > 5 {
			issueStartLine = firstLoc.TextRange.StartLine
			issueEndLine = firstLoc.TextRange.EndLine
			if firstLoc.Component != "" {
				component = firstLoc.Component
			}
		}
	}

	// Handle special case: no line number specified
	if issueStartLine == 0 {
		// For issues like "Add a new line at the end of file", show beginning of file
		sourceLines, err := g.client.GetSourceCode(component, 1, 10)
		if err != nil || len(sourceLines) == 0 {
			return ""
		}

		var codeBuilder strings.Builder
		for _, sl := range sourceLines {
			cleanCode := stripHTMLTags(sl.Code)
			codeBuilder.WriteString(fmt.Sprintf("  %d: %s\n", sl.Line, cleanCode))
		}
		return strings.TrimSuffix(codeBuilder.String(), "\n")
	}

	// Determine line range (3 lines before and after for context)
	startLine := issueStartLine - 3
	if startLine < 1 {
		startLine = 1
	}
	endLine := issueEndLine + 3

	sourceLines, err := g.client.GetSourceCode(component, startLine, endLine)
	if err != nil {
		return ""
	}

	var codeBuilder strings.Builder
	for _, sl := range sourceLines {
		// Mark the problematic line(s)
		prefix := "  "
		if sl.Line >= issueStartLine && sl.Line <= issueEndLine {
			prefix = "> "
		}
		cleanCode := stripHTMLTags(sl.Code)
		codeBuilder.WriteString(fmt.Sprintf("%s%d: %s\n", prefix, sl.Line, cleanCode))
	}

	return strings.TrimSuffix(codeBuilder.String(), "\n")
}

// fetchHowToFix fetches rule description
func (g *Generator) fetchHowToFix(ruleKey string, ruleCache map[string]string, mu *sync.Mutex) string {
	// Check cache first
	mu.Lock()
	if desc, ok := ruleCache[ruleKey]; ok {
		mu.Unlock()
		return desc
	}
	mu.Unlock()

	rule, err := g.client.GetRule(ruleKey)
	if err != nil {
		return ""
	}

	// Extract "How to fix" or use description
	description := rule.MdDesc
	if description == "" {
		description = stripHTML(rule.HtmlDesc)
	}

	// Try to extract "How to fix" section
	howToFix := extractHowToFix(description)

	// Cache the result
	mu.Lock()
	ruleCache[ruleKey] = howToFix
	mu.Unlock()

	return howToFix
}

// stripHTML removes HTML tags from a string
func stripHTML(html string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	text := re.ReplaceAllString(html, "")
	// Clean up extra whitespace
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	return strings.TrimSpace(text)
}

// extractHowToFix tries to extract the "How to fix" section from rule description
func extractHowToFix(description string) string {
	// Common patterns for "how to fix" sections in SonarQube rules
	patterns := []string{
		`(?i)how\s+to\s+fix[:\s]+(.+?)(?:\n\n|$)`,
		`(?i)compliant\s+solution[:\s]+(.+?)(?:\n\n|$)`,
		`(?i)noncompliant\s+code\s+example.+?compliant\s+solution[:\s]+(.+?)(?:\n\n|$)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(description)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	// If no specific "how to fix" found, return truncated description
	if len(description) > 500 {
		return description[:500] + "..."
	}
	return description
}

// getLanguageFromFile determines programming language from file extension
func getLanguageFromFile(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".go":
		return "go"
	case ".java":
		return "java"
	case ".js":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".py":
		return "python"
	case ".rb":
		return "ruby"
	case ".php":
		return "php"
	case ".cs":
		return "csharp"
	case ".cpp", ".cc", ".cxx":
		return "cpp"
	case ".c", ".h":
		return "c"
	case ".swift":
		return "swift"
	case ".kt":
		return "kotlin"
	case ".rs":
		return "rust"
	case ".scala":
		return "scala"
	case ".vue":
		return "vue"
	case ".jsx":
		return "jsx"
	case ".tsx":
		return "tsx"
	case ".html":
		return "html"
	case ".css":
		return "css"
	case ".scss":
		return "scss"
	case ".sql":
		return "sql"
	case ".xml":
		return "xml"
	case ".yaml", ".yml":
		return "yaml"
	case ".json":
		return "json"
	case ".sh":
		return "bash"
	default:
		return ""
	}
}
