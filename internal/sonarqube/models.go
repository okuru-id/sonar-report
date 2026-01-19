package sonarqube

// Project represents a SonarQube project
type Project struct {
	Key          string `json:"key"`
	Name         string `json:"name"`
	Qualifier    string `json:"qualifier"`
	LastAnalysis string `json:"lastAnalysisDate,omitempty"`
}

// ProjectsResponse from /api/projects/search
type ProjectsResponse struct {
	Paging struct {
		PageIndex int `json:"pageIndex"`
		PageSize  int `json:"pageSize"`
		Total     int `json:"total"`
	} `json:"paging"`
	Components []Project `json:"components"`
}

// Branch represents a project branch
type Branch struct {
	Name         string `json:"name"`
	IsMain       bool   `json:"isMain"`
	Type         string `json:"type"`
	AnalysisDate string `json:"analysisDate,omitempty"`
}

// BranchesResponse from /api/project_branches/list
type BranchesResponse struct {
	Branches []Branch `json:"branches"`
}

// QualityGateStatus represents quality gate status
type QualityGateStatus struct {
	Status     string      `json:"status"` // OK, WARN, ERROR
	Conditions []Condition `json:"conditions"`
}

// Condition represents a quality gate condition
type Condition struct {
	Status         string `json:"status"`
	MetricKey      string `json:"metricKey"`
	Comparator     string `json:"comparator"`
	ErrorThreshold string `json:"errorThreshold"`
	ActualValue    string `json:"actualValue"`
}

// QualityGateResponse from /api/qualitygates/project_status
type QualityGateResponse struct {
	ProjectStatus QualityGateStatus `json:"projectStatus"`
}

// Measure represents a metric measure
type Measure struct {
	Metric    string  `json:"metric"`
	Value     string  `json:"value"`
	BestValue bool    `json:"bestValue,omitempty"`
	Period    *Period `json:"period,omitempty"`
}

// Period represents a period for measures
type Period struct {
	Value     string `json:"value"`
	BestValue bool   `json:"bestValue,omitempty"`
}

// Component with measures
type ComponentWithMeasures struct {
	Key      string    `json:"key"`
	Name     string    `json:"name"`
	Measures []Measure `json:"measures"`
}

// MeasuresResponse from /api/measures/component
type MeasuresResponse struct {
	Component ComponentWithMeasures `json:"component"`
}

// TextRange represents the location of code in a file
type TextRange struct {
	StartLine   int `json:"startLine"`
	EndLine     int `json:"endLine"`
	StartOffset int `json:"startOffset"`
	EndOffset   int `json:"endOffset"`
}

// FlowLocation represents a location within a flow
type FlowLocation struct {
	Component string     `json:"component"`
	TextRange *TextRange `json:"textRange,omitempty"`
	Msg       string     `json:"msg,omitempty"`
}

// Flow represents a data flow path for an issue
type Flow struct {
	Locations []FlowLocation `json:"locations"`
}

// Issue represents a SonarQube issue
type Issue struct {
	Key          string     `json:"key"`
	Rule         string     `json:"rule"`
	Severity     string     `json:"severity"` // BLOCKER, CRITICAL, MAJOR, MINOR, INFO
	Component    string     `json:"component"`
	Project      string     `json:"project"`
	Line         int        `json:"line,omitempty"`
	TextRange    *TextRange `json:"textRange,omitempty"`
	Message      string     `json:"message"`
	Type         string     `json:"type"` // BUG, VULNERABILITY, CODE_SMELL
	Effort       string     `json:"effort,omitempty"`
	CreationDate string     `json:"creationDate"`
	Status       string     `json:"status"`
	Tags         []string   `json:"tags,omitempty"`
	Flows        []Flow     `json:"flows,omitempty"` // Additional location info
}

// IssuesResponse from /api/issues/search
type IssuesResponse struct {
	Total  int     `json:"total"`
	Paging Paging  `json:"paging"`
	Issues []Issue `json:"issues"`
}

// Paging represents pagination info
type Paging struct {
	PageIndex int `json:"pageIndex"`
	PageSize  int `json:"pageSize"`
	Total     int `json:"total"`
}

// Hotspot represents a security hotspot
type Hotspot struct {
	Key                      string `json:"key"`
	Component                string `json:"component"`
	Project                  string `json:"project"`
	SecurityCategory         string `json:"securityCategory"`
	VulnerabilityProbability string `json:"vulnerabilityProbability"` // HIGH, MEDIUM, LOW
	Status                   string `json:"status"`
	Line                     int    `json:"line,omitempty"`
	Message                  string `json:"message"`
	CreationDate             string `json:"creationDate"`
}

// HotspotsResponse from /api/hotspots/search
type HotspotsResponse struct {
	Paging   Paging    `json:"paging"`
	Hotspots []Hotspot `json:"hotspots"`
}

// Analysis represents a project analysis
type Analysis struct {
	Key            string `json:"key"`
	Date           string `json:"date"`
	ProjectVersion string `json:"projectVersion,omitempty"`
}

// AnalysesResponse from /api/project_analyses/search
type AnalysesResponse struct {
	Paging   Paging     `json:"paging"`
	Analyses []Analysis `json:"analyses"`
}

// SourceLine represents a line of source code
type SourceLine struct {
	Line int    `json:"line"`
	Code string `json:"code"`
}

// SourceResponse from /api/sources/show
type SourceResponse struct {
	Sources [][]interface{} `json:"sources"`
}

// Rule represents a SonarQube rule
type Rule struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	HtmlDesc string `json:"htmlDesc"`
	MdDesc   string `json:"mdDesc"`
	Severity string `json:"severity"`
	Type     string `json:"type"`
	Lang     string `json:"lang"`
	LangName string `json:"langName"`
}

// RuleResponse from /api/rules/show
type RuleResponse struct {
	Rule Rule `json:"rule"`
}
