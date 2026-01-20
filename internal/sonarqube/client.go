package sonarqube

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is the SonarQube API client
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewClient creates a new SonarQube API client
func NewClient(baseURL, token string) *Client {
	// Remove trailing slash
	baseURL = strings.TrimSuffix(baseURL, "/")

	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(method, endpoint string, params url.Values) ([]byte, error) {
	reqURL := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	if params != nil && len(params) > 0 {
		reqURL = fmt.Sprintf("%s?%s", reqURL, params.Encode())
	}

	req, err := http.NewRequest(method, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// SonarQube uses token as username with empty password for basic auth
	auth := base64.StdEncoding.EncodeToString([]byte(c.token + ":"))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check if response body is empty
	if len(body) == 0 {
		return nil, fmt.Errorf("empty response from SonarQube API (status: %d)", resp.StatusCode)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// GetProjects returns all projects from SonarQube
func (c *Client) GetProjects() ([]Project, error) {
	var allProjects []Project
	page := 1
	pageSize := 100

	for {
		params := url.Values{}
		params.Set("ps", fmt.Sprintf("%d", pageSize))
		params.Set("p", fmt.Sprintf("%d", page))

		body, err := c.doRequest("GET", "/api/projects/search", params)
		if err != nil {
			return nil, err
		}

		var resp ProjectsResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse projects response: %w", err)
		}

		allProjects = append(allProjects, resp.Components...)

		if len(allProjects) >= resp.Paging.Total {
			break
		}
		page++
	}

	return allProjects, nil
}

// GetBranches returns all branches for a project
func (c *Client) GetBranches(projectKey string) ([]Branch, error) {
	params := url.Values{}
	params.Set("project", projectKey)

	body, err := c.doRequest("GET", "/api/project_branches/list", params)
	if err != nil {
		return nil, err
	}

	var resp BranchesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse branches response: %w", err)
	}

	return resp.Branches, nil
}

// GetQualityGateStatus returns the quality gate status for a project
func (c *Client) GetQualityGateStatus(projectKey, branch string) (*QualityGateStatus, error) {
	params := url.Values{}
	params.Set("projectKey", projectKey)
	if branch != "" {
		params.Set("branch", branch)
	}

	body, err := c.doRequest("GET", "/api/qualitygates/project_status", params)
	if err != nil {
		return nil, err
	}

	var resp QualityGateResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse quality gate response: %w", err)
	}

	return &resp.ProjectStatus, nil
}

// GetMeasures returns measures for a project
func (c *Client) GetMeasures(projectKey, branch string, metricKeys []string) ([]Measure, error) {
	params := url.Values{}
	params.Set("component", projectKey)
	params.Set("metricKeys", strings.Join(metricKeys, ","))
	if branch != "" {
		params.Set("branch", branch)
	}

	body, err := c.doRequest("GET", "/api/measures/component", params)
	if err != nil {
		return nil, err
	}

	var resp MeasuresResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse measures response: %w", err)
	}

	return resp.Component.Measures, nil
}

// GetIssues returns issues for a project
func (c *Client) GetIssues(projectKey, branch string, maxResults int) ([]Issue, int, error) {
	var allIssues []Issue
	page := 1
	pageSize := 100
	total := 0

	for {
		params := url.Values{}
		params.Set("componentKeys", projectKey)
		params.Set("ps", fmt.Sprintf("%d", pageSize))
		params.Set("p", fmt.Sprintf("%d", page))
		params.Set("resolved", "false")
		// Request additional fields for more accurate location info
		params.Set("additionalFields", "_all")
		if branch != "" {
			params.Set("branch", branch)
		}

		body, err := c.doRequest("GET", "/api/issues/search", params)
		if err != nil {
			return nil, 0, err
		}

		var resp IssuesResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, 0, fmt.Errorf("failed to parse issues response: %w", err)
		}

		total = resp.Total
		allIssues = append(allIssues, resp.Issues...)

		if len(allIssues) >= resp.Paging.Total || len(allIssues) >= maxResults {
			break
		}
		page++
	}

	// Limit to maxResults
	if len(allIssues) > maxResults {
		allIssues = allIssues[:maxResults]
	}

	return allIssues, total, nil
}

// GetHotspots returns security hotspots for a project
func (c *Client) GetHotspots(projectKey, branch string, maxResults int) ([]Hotspot, int, error) {
	var allHotspots []Hotspot
	page := 1
	pageSize := 100
	total := 0

	for {
		params := url.Values{}
		params.Set("projectKey", projectKey)
		params.Set("ps", fmt.Sprintf("%d", pageSize))
		params.Set("p", fmt.Sprintf("%d", page))
		if branch != "" {
			params.Set("branch", branch)
		}

		body, err := c.doRequest("GET", "/api/hotspots/search", params)
		if err != nil {
			return nil, 0, err
		}

		var resp HotspotsResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, 0, fmt.Errorf("failed to parse hotspots response: %w", err)
		}

		total = resp.Paging.Total
		allHotspots = append(allHotspots, resp.Hotspots...)

		if len(allHotspots) >= resp.Paging.Total || len(allHotspots) >= maxResults {
			break
		}
		page++
	}

	// Limit to maxResults
	if len(allHotspots) > maxResults {
		allHotspots = allHotspots[:maxResults]
	}

	return allHotspots, total, nil
}

// GetAnalyses returns analysis history for a project
func (c *Client) GetAnalyses(projectKey, branch string, limit int) ([]Analysis, error) {
	params := url.Values{}
	params.Set("project", projectKey)
	params.Set("ps", fmt.Sprintf("%d", limit))
	if branch != "" {
		params.Set("branch", branch)
	}

	body, err := c.doRequest("GET", "/api/project_analyses/search", params)
	if err != nil {
		return nil, err
	}

	var resp AnalysesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse analyses response: %w", err)
	}

	return resp.Analyses, nil
}

// Validate checks if the client can connect to SonarQube
func (c *Client) Validate() error {
	_, err := c.doRequest("GET", "/api/system/status", nil)
	return err
}

// GetSourceCode returns source code lines for a component
func (c *Client) GetSourceCode(componentKey string, fromLine, toLine int) ([]SourceLine, error) {
	// Use /api/sources/show first as it returns explicit line numbers
	sourceLines, err := c.getSourceCodeFromShow(componentKey, fromLine, toLine)
	if err == nil && len(sourceLines) > 0 {
		return sourceLines, nil
	}

	// Fallback to /api/sources/raw
	params := url.Values{}
	params.Set("key", componentKey)
	params.Set("from", fmt.Sprintf("%d", fromLine))
	params.Set("to", fmt.Sprintf("%d", toLine))

	body, err := c.doRequest("GET", "/api/sources/raw", params)
	if err != nil {
		return nil, err
	}

	// /api/sources/raw returns plain text, split by lines
	lines := strings.Split(string(body), "\n")
	var result []SourceLine
	for i, line := range lines {
		lineNum := fromLine + i
		if lineNum <= toLine {
			result = append(result, SourceLine{
				Line: lineNum,
				Code: line,
			})
		}
	}

	return result, nil
}

// getSourceCodeFromShow uses /api/sources/show which returns explicit line numbers
func (c *Client) getSourceCodeFromShow(componentKey string, fromLine, toLine int) ([]SourceLine, error) {
	params := url.Values{}
	params.Set("key", componentKey)
	params.Set("from", fmt.Sprintf("%d", fromLine))
	params.Set("to", fmt.Sprintf("%d", toLine))

	body, err := c.doRequest("GET", "/api/sources/show", params)
	if err != nil {
		return nil, err
	}

	var resp SourceResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse source response: %w", err)
	}

	var sourceLines []SourceLine
	for _, src := range resp.Sources {
		if len(src) >= 2 {
			lineNum := 0
			code := ""
			// First element is line number, second is code
			if num, ok := src[0].(float64); ok {
				lineNum = int(num)
			}
			if c, ok := src[1].(string); ok {
				code = c
			}
			sourceLines = append(sourceLines, SourceLine{
				Line: lineNum,
				Code: code,
			})
		}
	}

	return sourceLines, nil
}

// GetRule returns rule details including description and how to fix
func (c *Client) GetRule(ruleKey string) (*Rule, error) {
	params := url.Values{}
	params.Set("key", ruleKey)

	body, err := c.doRequest("GET", "/api/rules/show", params)
	if err != nil {
		return nil, err
	}

	var resp RuleResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse rule response: %w", err)
	}

	return &resp.Rule, nil
}

// DefaultMetricKeys returns the default metric keys to fetch
func DefaultMetricKeys() []string {
	return []string{
		"bugs",
		"vulnerabilities",
		"code_smells",
		"coverage",
		"duplicated_lines_density",
		"ncloc",
		"sqale_index",
		"sqale_rating",
		"reliability_rating",
		"security_rating",
		"security_hotspots",
		"new_bugs",
		"new_vulnerabilities",
		"new_code_smells",
		"new_coverage",
		"new_duplicated_lines_density",
	}
}
