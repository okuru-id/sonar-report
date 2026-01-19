package report

import (
	"fmt"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

// PDFGenerator generates PDF reports
type PDFGenerator struct{}

// NewPDFGenerator creates a new PDF generator
func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{}
}

// Color constants
var (
	colorPrimary    = []int{41, 128, 185}  // Blue
	colorSuccess    = []int{39, 174, 96}   // Green
	colorWarning    = []int{241, 196, 15}  // Yellow
	colorDanger     = []int{231, 76, 60}   // Red
	colorDark       = []int{44, 62, 80}    // Dark
	colorLight      = []int{236, 240, 241} // Light Gray
	colorWhite      = []int{255, 255, 255} // White
	colorMuted      = []int{127, 140, 141} // Muted
	colorRatingA    = []int{0, 170, 0}     // Green
	colorRatingB    = []int{133, 187, 43}  // Light Green
	colorRatingC    = []int{237, 186, 0}   // Yellow
	colorRatingD    = []int{237, 127, 15}  // Orange
	colorRatingE    = []int{208, 0, 0}     // Red
	colorCodeBg     = []int{248, 249, 250} // Code background
	colorCodeBorder = []int{200, 200, 200} // Code border
)

// Generate generates a PDF report
func (g *PDFGenerator) Generate(data *ReportData) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 20)

	// Add first page with header
	pdf.AddPage()
	g.addHeader(pdf, data)
	g.addProjectInfo(pdf, data)
	g.addQualityGate(pdf, data)
	g.addMetricsOverview(pdf, data)

	// Add issues page if needed
	if data.TotalIssues > 0 {
		pdf.AddPage()
		g.addIssuesSection(pdf, data)
	}

	// Add hotspots page if needed
	if data.TotalHotspots > 0 {
		g.addHotspotsSection(pdf, data)
	}

	// Add footer to all pages
	g.addFooterToAllPages(pdf, data)

	// Output to bytes
	var buf strings.Builder
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return []byte(buf.String()), nil
}

func (g *PDFGenerator) addHeader(pdf *gofpdf.Fpdf, data *ReportData) {
	// Header background
	pdf.SetFillColor(colorPrimary[0], colorPrimary[1], colorPrimary[2])
	pdf.Rect(0, 0, 210, 45, "F")

	// Title
	pdf.SetFont("Arial", "B", 24)
	pdf.SetTextColor(colorWhite[0], colorWhite[1], colorWhite[2])
	pdf.SetXY(15, 12)
	pdf.CellFormat(0, 10, "SonarQube Analysis Report", "", 1, "C", false, 0, "")

	// Subtitle
	pdf.SetFont("Arial", "", 12)
	pdf.SetXY(15, 25)
	pdf.CellFormat(0, 8, "Code Quality Analysis Summary", "", 1, "C", false, 0, "")

	pdf.Ln(25)
}

func (g *PDFGenerator) addProjectInfo(pdf *gofpdf.Fpdf, data *ReportData) {
	pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])

	// Section title
	g.addSectionTitle(pdf, "Project Information")

	// Info box
	startY := pdf.GetY()
	pdf.SetFillColor(colorLight[0], colorLight[1], colorLight[2])
	pdf.RoundedRect(15, startY, 180, 32, 3, "1234", "F")

	pdf.SetXY(20, startY+5)
	g.addInfoRow(pdf, "Project Name:", data.ProjectName)
	g.addInfoRow(pdf, "Project Key:", data.ProjectKey)
	g.addInfoRow(pdf, "Branch:", data.Branch)
	g.addInfoRow(pdf, "Generated:", formatTime(data.GeneratedAt))

	pdf.Ln(10)
}

func (g *PDFGenerator) addInfoRow(pdf *gofpdf.Fpdf, label, value string) {
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(colorMuted[0], colorMuted[1], colorMuted[2])
	pdf.CellFormat(35, 6, label, "", 0, "L", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
	pdf.CellFormat(55, 6, value, "", 0, "L", false, 0, "")
}

func (g *PDFGenerator) addQualityGate(pdf *gofpdf.Fpdf, data *ReportData) {
	g.addSectionTitle(pdf, "Quality Gate Status")

	// Quality Gate badge
	var bgColor, textColor []int
	var statusText string

	switch data.QualityGateStatus {
	case "OK":
		bgColor = colorSuccess
		textColor = colorWhite
		statusText = "PASSED"
	case "WARN":
		bgColor = colorWarning
		textColor = colorDark
		statusText = "WARNING"
	default:
		bgColor = colorDanger
		textColor = colorWhite
		statusText = "FAILED"
	}

	// Draw badge
	badgeWidth := float64(50)
	badgeX := (210 - badgeWidth) / 2
	pdf.SetFillColor(bgColor[0], bgColor[1], bgColor[2])
	pdf.RoundedRect(badgeX, pdf.GetY(), badgeWidth, 12, 3, "1234", "F")

	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(textColor[0], textColor[1], textColor[2])
	pdf.SetXY(badgeX, pdf.GetY()+2)
	pdf.CellFormat(badgeWidth, 8, statusText, "", 1, "C", false, 0, "")

	pdf.Ln(8)

	// Quality Gate Conditions
	if len(data.QualityGateConditions) > 0 {
		pdf.SetFont("Arial", "B", 10)
		pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
		pdf.CellFormat(0, 6, "Quality Gate Conditions:", "", 1, "L", false, 0, "")
		pdf.Ln(2)

		// Table header
		pdf.SetFillColor(colorDark[0], colorDark[1], colorDark[2])
		pdf.SetTextColor(colorWhite[0], colorWhite[1], colorWhite[2])
		pdf.SetFont("Arial", "B", 9)

		colWidths := []float64{60, 25, 35, 50}
		headers := []string{"Metric", "Status", "Actual", "Threshold"}

		for i, header := range headers {
			pdf.CellFormat(colWidths[i], 7, header, "", 0, "C", true, 0, "")
		}
		pdf.Ln(-1)

		// Table rows
		pdf.SetFont("Arial", "", 9)
		for i, cond := range data.QualityGateConditions {
			if i%2 == 0 {
				pdf.SetFillColor(colorLight[0], colorLight[1], colorLight[2])
			} else {
				pdf.SetFillColor(colorWhite[0], colorWhite[1], colorWhite[2])
			}

			pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
			pdf.CellFormat(colWidths[0], 6, cond.Metric, "LR", 0, "L", true, 0, "")

			// Status with color
			if cond.Status == "OK" {
				pdf.SetTextColor(colorSuccess[0], colorSuccess[1], colorSuccess[2])
				pdf.CellFormat(colWidths[1], 6, "Pass", "LR", 0, "C", true, 0, "")
			} else {
				pdf.SetTextColor(colorDanger[0], colorDanger[1], colorDanger[2])
				pdf.CellFormat(colWidths[1], 6, "Fail", "LR", 0, "C", true, 0, "")
			}

			pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
			pdf.CellFormat(colWidths[2], 6, cond.ActualValue, "LR", 0, "C", true, 0, "")
			pdf.CellFormat(colWidths[3], 6, cond.Comparator+" "+cond.ErrorThreshold, "LR", 0, "L", true, 0, "")
			pdf.Ln(-1)
		}

		// Bottom border
		pdf.SetDrawColor(colorDark[0], colorDark[1], colorDark[2])
		pdf.Line(15, pdf.GetY(), 185, pdf.GetY())
	}

	pdf.Ln(8)
}

func (g *PDFGenerator) addMetricsOverview(pdf *gofpdf.Fpdf, data *ReportData) {
	g.addSectionTitle(pdf, "Metrics Overview")

	// Metrics cards
	cardWidth := float64(55)
	cardHeight := float64(35)
	startX := float64(17)
	startY := pdf.GetY()

	// Bug card
	g.drawMetricCard(pdf, startX, startY, cardWidth, cardHeight,
		"Bugs", data.Metrics.Bugs, data.Metrics.ReliabilityRating, colorDanger)

	// Vulnerability card
	g.drawMetricCard(pdf, startX+cardWidth+5, startY, cardWidth, cardHeight,
		"Vulnerabilities", data.Metrics.Vulnerabilities, data.Metrics.SecurityRating, colorWarning)

	// Code Smells card
	g.drawMetricCard(pdf, startX+2*(cardWidth+5), startY, cardWidth, cardHeight,
		"Code Smells", data.Metrics.CodeSmells, data.Metrics.MaintainabilityRating, colorPrimary)

	pdf.SetY(startY + cardHeight + 8)

	// Additional metrics table
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
	pdf.CellFormat(0, 6, "Additional Metrics:", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// Table
	metrics := [][]string{
		{"Lines of Code", data.Metrics.LinesOfCode},
		{"Coverage", data.Metrics.Coverage},
		{"Duplications", data.Metrics.DuplicatedLinesDensity},
		{"Technical Debt", data.Metrics.TechnicalDebt},
	}

	pdf.SetFont("Arial", "", 9)
	for i, row := range metrics {
		if i%2 == 0 {
			pdf.SetFillColor(colorLight[0], colorLight[1], colorLight[2])
		} else {
			pdf.SetFillColor(colorWhite[0], colorWhite[1], colorWhite[2])
		}

		pdf.SetTextColor(colorMuted[0], colorMuted[1], colorMuted[2])
		pdf.CellFormat(60, 6, row[0], "", 0, "L", true, 0, "")

		pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
		pdf.SetFont("Arial", "B", 9)
		pdf.CellFormat(50, 6, row[1], "", 1, "L", true, 0, "")
		pdf.SetFont("Arial", "", 9)
	}

	pdf.Ln(5)
}

func (g *PDFGenerator) drawMetricCard(pdf *gofpdf.Fpdf, x, y, w, h float64, title, value, rating string, baseColor []int) {
	// Card background
	pdf.SetFillColor(colorWhite[0], colorWhite[1], colorWhite[2])
	pdf.SetDrawColor(colorLight[0], colorLight[1], colorLight[2])
	pdf.RoundedRect(x, y, w, h, 3, "1234", "FD")

	// Title
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(colorMuted[0], colorMuted[1], colorMuted[2])
	pdf.SetXY(x, y+3)
	pdf.CellFormat(w, 5, title, "", 1, "C", false, 0, "")

	// Value
	pdf.SetFont("Arial", "B", 20)
	pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
	pdf.SetXY(x, y+10)
	pdf.CellFormat(w, 10, value, "", 1, "C", false, 0, "")

	// Rating badge
	ratingColor := g.getRatingColor(rating)
	badgeWidth := float64(25)
	badgeX := x + (w-badgeWidth)/2

	pdf.SetFillColor(ratingColor[0], ratingColor[1], ratingColor[2])
	pdf.RoundedRect(badgeX, y+h-10, badgeWidth, 7, 2, "1234", "F")

	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(colorWhite[0], colorWhite[1], colorWhite[2])
	pdf.SetXY(badgeX, y+h-9)
	pdf.CellFormat(badgeWidth, 5, "Rating: "+rating, "", 0, "C", false, 0, "")
}

func (g *PDFGenerator) getRatingColor(rating string) []int {
	switch rating {
	case "A":
		return colorRatingA
	case "B":
		return colorRatingB
	case "C":
		return colorRatingC
	case "D":
		return colorRatingD
	case "E":
		return colorRatingE
	default:
		return colorMuted
	}
}

func (g *PDFGenerator) addIssuesSection(pdf *gofpdf.Fpdf, data *ReportData) {
	g.addSectionTitle(pdf, "Issues Analysis")

	// Total issues badge
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
	pdf.CellFormat(0, 8, fmt.Sprintf("Total Issues: %d", data.TotalIssues), "", 1, "C", false, 0, "")
	pdf.Ln(3)

	// Issues by severity
	severities := GetSortedSeverities(data.IssuesBySeverity)

	// Severity summary table
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 6, "Issues by Severity:", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	pdf.SetFont("Arial", "", 9)
	for _, sev := range severities {
		issues := data.IssuesBySeverity[sev]
		count := len(issues)

		// Severity color
		sevColor := g.getSeverityColor(sev)
		pdf.SetFillColor(sevColor[0], sevColor[1], sevColor[2])
		pdf.RoundedRect(15, pdf.GetY(), 4, 4, 1, "1234", "F")

		pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
		pdf.SetX(22)
		pdf.CellFormat(40, 5, sev, "", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "B", 9)
		pdf.CellFormat(20, 5, fmt.Sprintf("%d", count), "", 1, "L", false, 0, "")
		pdf.SetFont("Arial", "", 9)
	}

	pdf.Ln(5)

	// Top issues table
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 6, "Top Issues:", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// Table header
	pdf.SetFillColor(colorDark[0], colorDark[1], colorDark[2])
	pdf.SetTextColor(colorWhite[0], colorWhite[1], colorWhite[2])
	pdf.SetFont("Arial", "B", 8)

	colWidths := []float64{20, 45, 12, 93}
	headers := []string{"Severity", "File", "Line", "Message"}

	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 6, header, "", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Table rows
	pdf.SetFont("Arial", "", 7)
	count := 0
	maxIssues := 25

	for _, sev := range severities {
		issues := data.IssuesBySeverity[sev]
		for _, issue := range issues {
			if count >= maxIssues {
				break
			}

			if count%2 == 0 {
				pdf.SetFillColor(colorLight[0], colorLight[1], colorLight[2])
			} else {
				pdf.SetFillColor(colorWhite[0], colorWhite[1], colorWhite[2])
			}

			// Severity with color indicator
			sevColor := g.getSeverityColor(issue.Severity)
			pdf.SetTextColor(sevColor[0], sevColor[1], sevColor[2])
			pdf.CellFormat(colWidths[0], 5, issue.Severity, "LR", 0, "L", true, 0, "")

			pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
			pdf.CellFormat(colWidths[1], 5, truncateString(issue.Component, 28), "LR", 0, "L", true, 0, "")
			pdf.CellFormat(colWidths[2], 5, fmt.Sprintf("%d", issue.Line), "LR", 0, "C", true, 0, "")
			pdf.CellFormat(colWidths[3], 5, truncateString(issue.Message, 60), "LR", 0, "L", true, 0, "")
			pdf.Ln(-1)
			count++
		}
		if count >= maxIssues {
			break
		}
	}

	// Bottom border
	pdf.SetDrawColor(colorDark[0], colorDark[1], colorDark[2])
	pdf.Line(15, pdf.GetY(), 185, pdf.GetY())

	if data.TotalIssues > maxIssues {
		pdf.Ln(3)
		pdf.SetFont("Arial", "I", 8)
		pdf.SetTextColor(colorMuted[0], colorMuted[1], colorMuted[2])
		pdf.CellFormat(0, 5, fmt.Sprintf("Showing %d of %d issues. See full report for complete list.", maxIssues, data.TotalIssues), "", 1, "C", false, 0, "")
	}

	// Add detailed issues with code snippets
	g.addDetailedIssues(pdf, data)
}

// addDetailedIssues adds detailed issue information with code snippets and how to fix
func (g *PDFGenerator) addDetailedIssues(pdf *gofpdf.Fpdf, data *ReportData) {
	severities := GetSortedSeverities(data.IssuesBySeverity)
	detailedCount := 0
	maxDetailed := 10 // Show detailed view for top 10 issues

	for _, sev := range severities {
		issues := data.IssuesBySeverity[sev]
		for _, issue := range issues {
			if detailedCount >= maxDetailed {
				return
			}

			// Skip if no code snippet
			if issue.CodeSnippet == "" && issue.HowToFix == "" {
				continue
			}

			// Check if we need a new page
			if pdf.GetY() > 220 {
				pdf.AddPage()
			}

			detailedCount++

			// Issue header
			pdf.Ln(5)
			g.addSubSectionTitle(pdf, fmt.Sprintf("Issue #%d: %s", detailedCount, truncateString(issue.Message, 60)))

			// Issue info table
			pdf.SetFont("Arial", "", 8)
			pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])

			// File info
			pdf.SetFont("Arial", "B", 8)
			pdf.CellFormat(25, 5, "File:", "", 0, "L", false, 0, "")
			pdf.SetFont("Courier", "", 8)
			pdf.CellFormat(0, 5, truncateString(issue.Component, 70), "", 1, "L", false, 0, "")

			// Line info
			pdf.SetFont("Arial", "B", 8)
			pdf.CellFormat(25, 5, "Line:", "", 0, "L", false, 0, "")
			pdf.SetFont("Arial", "", 8)
			lineInfo := fmt.Sprintf("%d", issue.Line)
			if issue.EndLine != 0 && issue.EndLine != issue.Line {
				lineInfo = fmt.Sprintf("%d - %d", issue.Line, issue.EndLine)
			}
			pdf.CellFormat(25, 5, lineInfo, "", 0, "L", false, 0, "")

			// Severity
			pdf.SetFont("Arial", "B", 8)
			pdf.CellFormat(25, 5, "Severity:", "", 0, "L", false, 0, "")
			sevColor := g.getSeverityColor(issue.Severity)
			pdf.SetTextColor(sevColor[0], sevColor[1], sevColor[2])
			pdf.SetFont("Arial", "B", 8)
			pdf.CellFormat(25, 5, issue.Severity, "", 0, "L", false, 0, "")

			// Type
			pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
			pdf.SetFont("Arial", "B", 8)
			pdf.CellFormat(15, 5, "Type:", "", 0, "L", false, 0, "")
			pdf.SetFont("Arial", "", 8)
			pdf.CellFormat(0, 5, issue.Type, "", 1, "L", false, 0, "")

			// Rule
			pdf.SetFont("Arial", "B", 8)
			pdf.CellFormat(25, 5, "Rule:", "", 0, "L", false, 0, "")
			pdf.SetFont("Courier", "", 7)
			pdf.CellFormat(0, 5, issue.Rule, "", 1, "L", false, 0, "")

			pdf.Ln(3)

			// Code snippet
			if issue.CodeSnippet != "" {
				g.addCodeBlock(pdf, "Problematic Code:", issue.CodeSnippet)
				pdf.Ln(3)
			}

			// How to fix
			if issue.HowToFix != "" {
				g.addHowToFix(pdf, issue.HowToFix)
			}

			// Separator
			pdf.Ln(3)
			pdf.SetDrawColor(colorLight[0], colorLight[1], colorLight[2])
			pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
		}
	}
}

// addSubSectionTitle adds a smaller section title
func (g *PDFGenerator) addSubSectionTitle(pdf *gofpdf.Fpdf, title string) {
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
	pdf.CellFormat(0, 6, title, "", 1, "L", false, 0, "")
	pdf.Ln(2)
}

// addCodeBlock adds a code block with monospace font
func (g *PDFGenerator) addCodeBlock(pdf *gofpdf.Fpdf, title, code string) {
	// Title
	pdf.SetFont("Arial", "B", 8)
	pdf.SetTextColor(colorPrimary[0], colorPrimary[1], colorPrimary[2])
	pdf.CellFormat(0, 5, title, "", 1, "L", false, 0, "")

	// Code background
	startY := pdf.GetY()
	lines := strings.Split(code, "\n")
	lineHeight := float64(3.5)
	boxHeight := float64(len(lines))*lineHeight + 4

	// Limit height to prevent overflow
	if boxHeight > 60 {
		boxHeight = 60
		maxLines := int((boxHeight - 4) / lineHeight)
		if maxLines < len(lines) {
			lines = lines[:maxLines]
			lines = append(lines, "... (truncated)")
		}
	}

	// Draw code box
	pdf.SetFillColor(colorCodeBg[0], colorCodeBg[1], colorCodeBg[2])
	pdf.SetDrawColor(colorCodeBorder[0], colorCodeBorder[1], colorCodeBorder[2])
	pdf.RoundedRect(15, startY, 180, boxHeight, 2, "1234", "FD")

	// Code text
	pdf.SetFont("Courier", "", 7)
	pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
	pdf.SetXY(17, startY+2)

	for _, line := range lines {
		// Highlight problematic lines (starting with "> ")
		if strings.HasPrefix(line, "> ") {
			pdf.SetTextColor(colorDanger[0], colorDanger[1], colorDanger[2])
			pdf.SetFont("Courier", "B", 7)
		} else {
			pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
			pdf.SetFont("Courier", "", 7)
		}

		// Truncate long lines
		displayLine := line
		if len(displayLine) > 100 {
			displayLine = displayLine[:97] + "..."
		}

		pdf.CellFormat(176, lineHeight, displayLine, "", 1, "L", false, 0, "")
		pdf.SetX(17)
	}

	pdf.SetY(startY + boxHeight + 2)
}

// addHowToFix adds the "How to fix" section
func (g *PDFGenerator) addHowToFix(pdf *gofpdf.Fpdf, howToFix string) {
	// Title
	pdf.SetFont("Arial", "B", 8)
	pdf.SetTextColor(colorSuccess[0], colorSuccess[1], colorSuccess[2])
	pdf.CellFormat(0, 5, "How to Fix:", "", 1, "L", false, 0, "")

	// Content box
	startY := pdf.GetY()

	// Truncate if too long
	displayText := howToFix
	if len(displayText) > 400 {
		displayText = displayText[:397] + "..."
	}

	// Calculate box height based on text
	pdf.SetFont("Arial", "", 7)
	lines := pdf.SplitLines([]byte(displayText), 174)
	lineHeight := float64(3.5)
	boxHeight := float64(len(lines))*lineHeight + 4

	if boxHeight > 40 {
		boxHeight = 40
	}

	// Draw box with light green background
	pdf.SetFillColor(240, 255, 240) // Very light green
	pdf.SetDrawColor(colorSuccess[0], colorSuccess[1], colorSuccess[2])
	pdf.RoundedRect(15, startY, 180, boxHeight, 2, "1234", "FD")

	// Text
	pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
	pdf.SetXY(17, startY+2)

	for i, line := range lines {
		if float64(i)*lineHeight > boxHeight-6 {
			pdf.CellFormat(174, lineHeight, "...", "", 1, "L", false, 0, "")
			break
		}
		pdf.CellFormat(174, lineHeight, string(line), "", 1, "L", false, 0, "")
		pdf.SetX(17)
	}

	pdf.SetY(startY + boxHeight + 2)
}

func (g *PDFGenerator) getSeverityColor(severity string) []int {
	switch severity {
	case "BLOCKER":
		return colorDanger
	case "CRITICAL":
		return []int{230, 126, 34} // Orange
	case "MAJOR":
		return colorWarning
	case "MINOR":
		return colorPrimary
	case "INFO":
		return colorMuted
	default:
		return colorMuted
	}
}

func (g *PDFGenerator) addHotspotsSection(pdf *gofpdf.Fpdf, data *ReportData) {
	// Check if we need a new page
	if pdf.GetY() > 200 {
		pdf.AddPage()
	}

	g.addSectionTitle(pdf, "Security Hotspots")

	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
	pdf.CellFormat(0, 8, fmt.Sprintf("Total Hotspots: %d", data.TotalHotspots), "", 1, "C", false, 0, "")
	pdf.Ln(3)

	// Hotspots by priority
	if len(data.HotspotsByPriority) > 0 {
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(0, 6, "Hotspots by Priority:", "", 1, "L", false, 0, "")
		pdf.Ln(2)

		pdf.SetFont("Arial", "", 9)
		for priority, count := range data.HotspotsByPriority {
			prioColor := g.getPriorityColor(priority)
			pdf.SetFillColor(prioColor[0], prioColor[1], prioColor[2])
			pdf.RoundedRect(15, pdf.GetY(), 4, 4, 1, "1234", "F")

			pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
			pdf.SetX(22)
			pdf.CellFormat(40, 5, priority, "", 0, "L", false, 0, "")
			pdf.SetFont("Arial", "B", 9)
			pdf.CellFormat(20, 5, fmt.Sprintf("%d", count), "", 1, "L", false, 0, "")
			pdf.SetFont("Arial", "", 9)
		}
	}

	pdf.Ln(5)

	// Hotspot details table
	if len(data.Hotspots) > 0 {
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(0, 6, "Hotspot Details:", "", 1, "L", false, 0, "")
		pdf.Ln(2)

		// Table header
		pdf.SetFillColor(colorDark[0], colorDark[1], colorDark[2])
		pdf.SetTextColor(colorWhite[0], colorWhite[1], colorWhite[2])
		pdf.SetFont("Arial", "B", 8)

		colWidths := []float64{22, 40, 70, 25}
		headers := []string{"Priority", "Category", "Location", "Status"}

		for i, header := range headers {
			pdf.CellFormat(colWidths[i], 6, header, "", 0, "C", true, 0, "")
		}
		pdf.Ln(-1)

		// Table rows
		pdf.SetFont("Arial", "", 7)
		maxHotspots := 15

		for i, hotspot := range data.Hotspots {
			if i >= maxHotspots {
				break
			}

			if i%2 == 0 {
				pdf.SetFillColor(colorLight[0], colorLight[1], colorLight[2])
			} else {
				pdf.SetFillColor(colorWhite[0], colorWhite[1], colorWhite[2])
			}

			prioColor := g.getPriorityColor(hotspot.VulnerabilityProbability)
			pdf.SetTextColor(prioColor[0], prioColor[1], prioColor[2])
			pdf.CellFormat(colWidths[0], 5, hotspot.VulnerabilityProbability, "LR", 0, "L", true, 0, "")

			pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
			pdf.CellFormat(colWidths[1], 5, truncateString(hotspot.SecurityCategory, 24), "LR", 0, "L", true, 0, "")
			location := fmt.Sprintf("%s:%d", truncateString(hotspot.Component, 40), hotspot.Line)
			pdf.CellFormat(colWidths[2], 5, location, "LR", 0, "L", true, 0, "")
			pdf.CellFormat(colWidths[3], 5, hotspot.Status, "LR", 0, "C", true, 0, "")
			pdf.Ln(-1)
		}

		// Bottom border
		pdf.SetDrawColor(colorDark[0], colorDark[1], colorDark[2])
		pdf.Line(15, pdf.GetY(), 172, pdf.GetY())
	}
}

func (g *PDFGenerator) getPriorityColor(priority string) []int {
	switch priority {
	case "HIGH":
		return colorDanger
	case "MEDIUM":
		return colorWarning
	case "LOW":
		return colorSuccess
	default:
		return colorMuted
	}
}

func (g *PDFGenerator) addSectionTitle(pdf *gofpdf.Fpdf, title string) {
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(colorPrimary[0], colorPrimary[1], colorPrimary[2])
	pdf.CellFormat(0, 10, title, "", 1, "L", false, 0, "")

	// Underline
	pdf.SetDrawColor(colorPrimary[0], colorPrimary[1], colorPrimary[2])
	pdf.SetLineWidth(0.5)
	pdf.Line(15, pdf.GetY(), 60, pdf.GetY())
	pdf.SetLineWidth(0.2)

	pdf.Ln(5)
}

func (g *PDFGenerator) addFooterToAllPages(pdf *gofpdf.Fpdf, data *ReportData) {
	totalPages := pdf.PageCount()

	for i := 1; i <= totalPages; i++ {
		pdf.SetPage(i)

		// Footer line
		pdf.SetDrawColor(colorLight[0], colorLight[1], colorLight[2])
		pdf.Line(15, 280, 195, 280)

		// Footer text
		pdf.SetFont("Arial", "I", 8)
		pdf.SetTextColor(colorMuted[0], colorMuted[1], colorMuted[2])
		pdf.SetXY(15, 282)
		pdf.CellFormat(90, 5, "Generated by SonarQube Report Generator", "", 0, "L", false, 0, "")
		pdf.CellFormat(90, 5, fmt.Sprintf("Page %d of %d", i, totalPages), "", 0, "R", false, 0, "")
	}
}
