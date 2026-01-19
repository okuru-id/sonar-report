# SonarQube Report Generator

A Go-based web application that automates code quality report generation from SonarQube. It produces Markdown and PDF reports featuring severity-categorized issues, accurate code snippets, and fix guidance. With parallel fetching optimization, reports generate in ~1.5 seconds—25x faster than sequential processing.

## Features

- **Multi-format Reports**: Generate reports in Markdown (.md) and PDF formats
- **Comprehensive Issue Tracking**: Issues categorized by severity (Blocker, Critical, Major, Minor, Info)
- **Code Snippets**: Accurate source code snippets highlighting problematic lines
- **How to Fix Guidance**: Extracted fix recommendations from SonarQube rule documentation
- **Security Hotspots**: Track security vulnerabilities with priority levels
- **Quality Metrics**: Coverage, duplications, technical debt, and quality ratings
- **Parallel Processing**: Optimized fetching with worker pool pattern (~1.5s vs ~40s)
- **Web Dashboard**: Interactive UI with configurable report options
- **Docker Ready**: Easy deployment with Docker Compose
- **Session-based Auth**: Secure admin authentication

## Tech Stack

- **Backend**: Go 1.22+ with Gin Framework
- **PDF Generation**: go-pdf/fpdf
- **Containerization**: Docker & Docker Compose
- **Database**: PostgreSQL (for SonarQube)

## Project Structure

```
sonarqube-report-generator/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── auth/
│   │   └── auth.go              # Authentication middleware
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── handler/
│   │   ├── api_handler.go       # REST API handlers
│   │   └── web_handler.go       # Web UI handlers
│   ├── report/
│   │   ├── generator.go         # Report generation logic
│   │   ├── markdown.go          # Markdown template & rendering
│   │   ├── pdf.go               # PDF generation
│   │   ├── models.go            # Report data models
│   │   └── storage.go           # Report file storage
│   └── sonarqube/
│       ├── client.go            # SonarQube API client
│       └── models.go            # SonarQube data models
├── web/
│   └── templates/
│       ├── dashboard.html       # Main dashboard UI
│       └── login.html           # Login page
├── docker-compose.yml           # Docker services configuration
├── Dockerfile                   # Multi-stage Docker build
├── .env.example                 # Environment variables template
├── go.mod                       # Go module definition
└── go.sum                       # Go dependencies checksum
```

## Quick Start

### Prerequisites

- Docker & Docker Compose
- SonarQube instance with API token
- Go 1.22+ (for local development)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/sonarqube-report-generator.git
   cd sonarqube-report-generator
   ```

2. **Configure environment variables**
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` with your settings:
   ```env
   # SonarQube Configuration
   SONARQUBE_URL=https://your-sonarqube-instance.com
   SONARQUBE_TOKEN=sqp_your_token_here
   
   # Admin Authentication
   ADMIN_USERNAME=admin
   ADMIN_PASSWORD=your-secure-password
   
   # Session Configuration
   SESSION_SECRET=your-random-secret-key-min-32-chars
   ```

3. **Start the application**
   ```bash
   docker-compose up -d --build report-generator
   ```

4. **Access the dashboard**
   
   Open [http://localhost:8080](http://localhost:8080) in your browser

### Running with Local SonarQube

To run with a local SonarQube instance:

```bash
docker-compose up -d
```

This starts:
- SonarQube at `http://localhost:9002`
- PostgreSQL database
- Report Generator at `http://localhost:8080`

## Usage

### Web Dashboard

1. Login with your admin credentials
2. Enter the SonarQube project key
3. Select report format (Markdown or PDF)
4. Configure advanced options:
   - Include/Exclude Code Snippets
   - Include/Exclude How to Fix guidance
5. Click "Generate Report"

### REST API

#### Health Check
```bash
curl http://localhost:8080/api/v1/health
```

#### Get Projects List
```bash
curl -b cookies.txt http://localhost:8080/api/v1/projects
```

#### Generate Report
```bash
# Login first
curl -c cookies.txt -X POST "http://localhost:8080/login" \
  -d "username=admin&password=your-password"

# Generate Markdown report
curl -b cookies.txt -X POST "http://localhost:8080/api/v1/reports/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "projectKey": "your-project-key",
    "format": "md",
    "includeCodeSnippets": true,
    "includeHowToFix": true
  }'

# Generate PDF report
curl -b cookies.txt -X POST "http://localhost:8080/api/v1/reports/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "projectKey": "your-project-key",
    "format": "pdf"
  }' --output report.pdf
```

#### API Response (JSON format for Markdown)
```json
{
  "success": true,
  "data": {
    "projectName": "My Project",
    "generatedAt": "2026-01-19T10:30:00Z",
    "qualityGateStatus": "OK",
    "metrics": {
      "bugs": "5",
      "vulnerabilities": "2",
      "codeSmells": "150",
      "coverage": "75.5%",
      "duplicatedLinesDensity": "3.2%"
    },
    "issuesBySeverity": {
      "BLOCKER": [...],
      "CRITICAL": [...],
      "MAJOR": [...],
      "MINOR": [...],
      "INFO": [...]
    }
  }
}
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SONARQUBE_URL` | SonarQube server URL | `http://sonarqube:9000` |
| `SONARQUBE_TOKEN` | SonarQube API token | - |
| `SERVER_PORT` | Application port | `8080` |
| `ADMIN_USERNAME` | Dashboard login username | `admin` |
| `ADMIN_PASSWORD` | Dashboard login password | - |
| `SESSION_SECRET` | Session encryption key (min 32 chars) | - |
| `REPORT_STORAGE_PATH` | Report files storage path | `./reports` |
| `REPORT_RETENTION_DAYS` | Days to keep generated reports | `30` |

### Getting SonarQube Token

1. Login to your SonarQube instance
2. Go to **My Account** > **Security**
3. Generate a new token with `Read` permissions
4. Copy the token to your `.env` file

## Report Contents

### Metrics Summary
- Bugs, Vulnerabilities, Code Smells
- Coverage & Duplication percentage
- Technical Debt estimation
- Reliability, Security, Maintainability ratings (A-E)

### Issues by Severity
Each issue includes:
- Rule ID and description
- File path and line number
- Code snippet with highlighted problematic line
- How to fix guidance (from SonarQube rules)

### Security Hotspots
- Vulnerability probability (HIGH, MEDIUM, LOW)
- Security category
- Review status

## Development

### Local Development

```bash
# Install dependencies
go mod download

# Run locally
go run cmd/server/main.go

# Build binary
go build -o sonar-report ./cmd/server
```

### Running Tests

```bash
go test ./...
```

### Building Docker Image

```bash
docker build -t sonarqube-report-generator .
```

## Performance

| Metric | Before Optimization | After Optimization |
|--------|--------------------|--------------------|
| Report Generation | ~40 seconds | ~1.5 seconds |
| Code Snippet Fetching | Sequential | Parallel (5 workers) |
| Improvement | - | **25x faster** |

## Troubleshooting

### Common Issues

1. **Cannot connect to SonarQube**
   - Verify `SONARQUBE_URL` is correct
   - Check if token has proper permissions
   - Ensure network connectivity between containers

2. **Empty code snippets**
   - Verify SonarQube has source code indexed
   - Check if `includeCodeSnippets` is enabled

3. **Login failed**
   - Verify `ADMIN_USERNAME` and `ADMIN_PASSWORD` in `.env`
   - Clear browser cookies and try again

### Logs

```bash
# View container logs
docker logs sonar-report-generator

# Follow logs in real-time
docker logs -f sonar-report-generator
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Acknowledgments

- [SonarQube](https://www.sonarqube.org/) - Code quality platform
- [Gin Framework](https://gin-gonic.com/) - HTTP web framework
- [go-pdf/fpdf](https://github.com/go-pdf/fpdf) - PDF generation library
