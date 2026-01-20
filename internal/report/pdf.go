package report

import (
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

type PDFGenerator struct{}

func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{}
}

func (g *PDFGenerator) Generate(data *ReportData) ([]byte, error) {
	pdf := g.createPDF(data)

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

func (g *PDFGenerator) createPDF(data *ReportData) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 20)
	pdf.AddPage()
	pdf.SetFont("Arial", "", 10)

	g.renderHeader(pdf)
	g.renderProjectInfo(pdf, data)
	g.renderQualityGate(pdf, data)
	g.renderMetrics(pdf, data)
	g.renderIssues(pdf, data)
	g.renderHotspots(pdf, data)
	g.renderSummary(pdf, data)

	return pdf
}

func (g *PDFGenerator) renderHeader(pdf *gofpdf.Fpdf) {
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "SonarQube Analysis Report", "", 1, "C", false, 0, "")
	pdf.Ln(3)
	pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
	pdf.Ln(5)
}

func (g *PDFGenerator) renderProjectInfo(pdf *gofpdf.Fpdf, data *ReportData) {
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, "Project Information", "", 1, "L", false, 0, "")
	pdf.Ln(3)

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(45, 6, "Project Name:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, data.ProjectName, "", 1, "L", false, 0, "")

	pdf.CellFormat(45, 6, "Project Key:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, data.ProjectKey, "", 1, "L", false, 0, "")

	pdf.CellFormat(45, 6, "Branch:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, data.Branch, "", 1, "L", false, 0, "")

	pdf.CellFormat(45, 6, "Report Generated:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, formatTimeSimple(data.GeneratedAt), "", 1, "L", false, 0, "")

	pdf.Ln(5)
}

func (g *PDFGenerator) renderQualityGate(pdf *gofpdf.Fpdf, data *ReportData) {
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, "Quality Gate", "", 1, "L", false, 0, "")
	pdf.Ln(3)

	status := QualityGateText(data.QualityGateStatus)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 6, status, "", 1, "L", false, 0, "")
	pdf.Ln(3)

	if len(data.QualityGateConditions) > 0 {
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(0, 6, "Conditions:", "", 1, "L", false, 0, "")
		pdf.Ln(2)

		colW := []float64{60.0, 25.0, 40.0, 40.0}
		g.renderSimpleTable(pdf, []string{"Metric", "Status", "Value", "Threshold"}, []string{}, colW)

		for _, cond := range data.QualityGateConditions {
			statusIcon := "OK"
			if cond.Status != "OK" {
				statusIcon = "FAIL"
			}
			row := []string{cond.Metric, statusIcon, cond.ActualValue, cond.Comparator + " " + cond.ErrorThreshold}
			g.renderSimpleTable(pdf, []string{}, row, colW)
		}
	}

	pdf.Ln(5)
}

func (g *PDFGenerator) renderMetrics(pdf *gofpdf.Fpdf, data *ReportData) {
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, "Metrics Overview", "", 1, "L", false, 0, "")
	pdf.Ln(3)

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(50, 6, "Bugs:", "", 0, "L", false, 0, "")
	pdf.CellFormat(30, 6, data.Metrics.Bugs, "", 0, "L", false, 0, "")
	pdf.CellFormat(30, 6, "("+RatingToLetter(data.Metrics.ReliabilityRating)+")", "", 1, "L", false, 0, "")

	pdf.CellFormat(50, 6, "Vulnerabilities:", "", 0, "L", false, 0, "")
	pdf.CellFormat(30, 6, data.Metrics.Vulnerabilities, "", 0, "L", false, 0, "")
	pdf.CellFormat(30, 6, "("+RatingToLetter(data.Metrics.SecurityRating)+")", "", 1, "L", false, 0, "")

	pdf.CellFormat(50, 6, "Code Smells:", "", 0, "L", false, 0, "")
	pdf.CellFormat(30, 6, data.Metrics.CodeSmells, "", 0, "L", false, 0, "")
	pdf.CellFormat(30, 6, "("+RatingToLetter(data.Metrics.MaintainabilityRating)+")", "", 1, "L", false, 0, "")

	pdf.Ln(3)

	pdf.CellFormat(50, 6, "Lines of Code:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, data.Metrics.LinesOfCode, "", 1, "L", false, 0, "")

	pdf.CellFormat(50, 6, "Coverage:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, data.Metrics.Coverage, "", 1, "L", false, 0, "")

	pdf.CellFormat(50, 6, "Duplications:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, data.Metrics.DuplicatedLinesDensity, "", 1, "L", false, 0, "")

	pdf.CellFormat(50, 6, "Technical Debt:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, data.Metrics.TechnicalDebt, "", 1, "L", false, 0, "")

	pdf.Ln(5)
}

func (g *PDFGenerator) renderIssues(pdf *gofpdf.Fpdf, data *ReportData) {
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, "Issues Analysis", "", 1, "L", false, 0, "")
	pdf.Ln(3)

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, fmt.Sprintf("Total Issues: %d", data.TotalIssues), "", 1, "L", false, 0, "")
	pdf.Ln(3)

	colW := []float64{50.0, 30.0, 30.0}
	g.renderSimpleTable(pdf, []string{"Type", "Count", "%"}, []string{}, colW)

	for issueType, count := range data.IssuesByType {
		percentage := fmt.Sprintf("%.1f%%", float64(count)/float64(data.TotalIssues)*100)
		row := []string{issueType, fmt.Sprintf("%d", count), percentage}
		g.renderSimpleTable(pdf, []string{}, row, colW)
	}

	pdf.Ln(5)

	severities := []string{"BLOCKER", "CRITICAL", "MAJOR", "MINOR", "INFO"}
	for _, severity := range severities {
		issues := []IssueItem{}
		if data.IssuesBySeverity != nil {
			issues = data.IssuesBySeverity[severity]
		}
		if len(issues) == 0 {
			continue
		}

		pdf.SetFont("Arial", "B", 11)
		pdf.CellFormat(0, 7, fmt.Sprintf("%s Issues (%d)", severity, len(issues)), "", 1, "L", false, 0, "")
		pdf.Ln(2)

		pdf.SetFont("Arial", "", 9)
		for idx, issue := range issues {
			if idx >= 10 {
				break
			}

			pdf.CellFormat(0, 5, fmt.Sprintf("%d. %s", idx+1, truncateStr(issue.Message, 80)), "", 1, "L", false, 0, "")
			pdf.SetFont("Arial", "", 8)
			pdf.CellFormat(5, 4, "", "", 0, "L", false, 0, "")
			pdf.CellFormat(0, 4, fmt.Sprintf("File: %s | Line: %d | Rule: %s", truncateStr(issue.Component, 40), issue.Line, truncateStr(issue.Rule, 30)), "", 1, "L", false, 0, "")

			if issue.Effort != "" {
				pdf.CellFormat(5, 4, "", "", 0, "L", false, 0, "")
				pdf.CellFormat(0, 4, fmt.Sprintf("Effort: %s", issue.Effort), "", 1, "L", false, 0, "")
			}

			pdf.SetFont("Arial", "", 9)
			pdf.Ln(2)
		}

		pdf.Ln(3)
	}

	pdf.Ln(3)
}

func (g *PDFGenerator) renderHotspots(pdf *gofpdf.Fpdf, data *ReportData) {
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, "Security Hotspots", "", 1, "L", false, 0, "")
	pdf.Ln(3)

	if data.TotalHotspots == 0 {
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(0, 6, "No Security Hotspots Found", "", 1, "L", false, 0, "")
		pdf.Ln(5)
		return
	}

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, fmt.Sprintf("Total Hotspots: %d", data.TotalHotspots), "", 1, "L", false, 0, "")
	pdf.Ln(3)

	colW := []float64{40.0, 30.0}
	g.renderSimpleTable(pdf, []string{"Priority", "Count"}, []string{}, colW)

	for priority, count := range data.HotspotsByPriority {
		row := []string{priority, fmt.Sprintf("%d", count)}
		g.renderSimpleTable(pdf, []string{}, row, colW)
	}

	pdf.Ln(3)

	if len(data.Hotspots) > 0 {
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(0, 6, "Details:", "", 1, "L", false, 0, "")
		pdf.Ln(2)

		colW := []float64{10.0, 30.0, 40.0, 40.0}
		g.renderSimpleTable(pdf, []string{"#", "Priority", "Category", "Location"}, []string{}, colW)

		for idx, hotspot := range data.Hotspots {
			if idx >= 10 {
				break
			}
			row := []string{
				fmt.Sprintf("%d", idx+1),
				hotspot.VulnerabilityProbability,
				hotspot.SecurityCategory,
				fmt.Sprintf("%s:%d", truncateStr(hotspot.Component, 30), hotspot.Line),
			}
			g.renderSimpleTable(pdf, []string{}, row, colW)
		}
	}

	pdf.Ln(5)
}

func (g *PDFGenerator) renderSummary(pdf *gofpdf.Fpdf, data *ReportData) {
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, "Summary", "", 1, "L", false, 0, "")
	pdf.Ln(3)

	status := QualityGateText(data.QualityGateStatus)

	colW := []float64{40.0, 60.0}
	if data.QualityGateStatus == "OK" {
		g.renderSimpleTable(pdf, []string{"Quality Gate", status}, []string{}, colW)
		g.renderSimpleTable(pdf, []string{}, []string{"Bugs", data.Metrics.Bugs}, colW)
		g.renderSimpleTable(pdf, []string{}, []string{"Vulnerabilities", data.Metrics.Vulnerabilities}, colW)
		g.renderSimpleTable(pdf, []string{}, []string{"Code Smells", data.Metrics.CodeSmells}, colW)
	} else {
		g.renderSimpleTable(pdf, []string{}, []string{"Action Required", ""}, colW)
		g.renderSimpleTable(pdf, []string{}, []string{"Quality Gate", status}, colW)
	}

	pdf.Ln(5)

	pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
	pdf.Ln(3)

	pdf.SetFont("Arial", "I", 8)
	pdf.CellFormat(0, 5, "Report generated by SonarQube Report Generator", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, formatTimeSimple(data.GeneratedAt), "", 1, "C", false, 0, "")
}

func (g *PDFGenerator) renderSimpleTable(pdf *gofpdf.Fpdf, headers []string, row []string, colWidths []float64) {
	if len(headers) > 0 {
		pdf.SetFont("Arial", "B", 9)
		for i, h := range headers {
			pdf.CellFormat(colWidths[i], 6, h, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(6)
	}

	if len(row) > 0 {
		pdf.SetFont("Arial", "", 9)
		for i, cell := range row {
			border := "1"
			cellStr := truncateStr(cell, int(colWidths[i]/4))
			pdf.CellFormat(colWidths[i], 6, cellStr, border, 0, "L", false, 0, "")
		}
		pdf.Ln(6)
	}
}

func truncateStr(s string, maxLen int) string {
	if s == "" {
		return ""
	}
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func formatTimeSimple(t interface{}) string {
	switch v := t.(type) {
	case string:
		return v
	default:
		return fmt.Sprintf("%v", t)
	}
}
