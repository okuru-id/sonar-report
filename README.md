# SonarQube Report Generator with Jenkins Trigger

## Overview

This application provides SonarQube report generation and Jenkins job triggering via API.

## Quick Start

```bash
# Build
go build -o sonarqube-report ./cmd/server/main.go

# Run
./sonarqube-report
```

## Jenkins Trigger API

### Endpoint
```
POST /api/v1/jenkins/trigger
```

### Headers
```
X-Jenkins-Url: https://jenkins.example.com
X-Jenkins-Token: your-jenkins-token
X-Jenkins-Job: your-job-name
X-Jenkins-User: sonar (optional)
```

### Example
```bash
curl -X POST \
  -H "X-Jenkins-Url: http://10.0.1.84:8080" \
  -H "X-Jenkins-Token: 11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  -H "X-Jenkins-Job: plb-sonarqube" \
  -H "X-Jenkins-User: sonar" \
  http://localhost:8080/api/v1/jenkins/trigger
```

## Documentation

Complete documentation is available in the `docs/` folder:

- **[API Documentation](docs/API_DOCUMENTATION.md)** - Complete API reference
- **[Integration Guide](docs/INTEGRATION_GUIDE.md)** - Integration examples
- **[Troubleshooting](docs/)** - Various troubleshooting guides

### Common Issues

If you encounter the error:
```json
{
  "error": "Cannot connect to Jenkins server at http://10.0.1.84:8080..."
}
```

Check:
1. Jenkins is running: `systemctl status jenkins`
2. Network connectivity: `ping 10.0.1.84`
3. Port accessibility: `nc -zv 10.0.1.84 8080`

See **[Connection Issues Guide](docs/CONNECTION_ISSUES.md)** for detailed troubleshooting.

## Features

### SonarQube Reports
- Generate PDF and Markdown reports
- Project and branch selection
- Code snippets and fix suggestions
- Report history and management

### Jenkins Integration
- Trigger jobs via API
- CSRF protection support
- Custom username support
- Detailed error handling

### Web Interface
- Dashboard for report management
- Report preview
- Shareable report links
- Admin authentication

## Configuration

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
# Edit .env with your configuration
```

## License

MIT License