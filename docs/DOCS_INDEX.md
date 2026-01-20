# Documentation Index

## Overview

This folder contains comprehensive documentation for the SonarQube Report Generator and Jenkins Trigger API.

## Documentation Structure

```
docs/
├── API Documentation/          # API reference and guides
│   ├── API_DOCUMENTATION.md    # Main API documentation
│   ├── JENKINS_API.md          # Jenkins-specific API details
│   ├── JENKINS_USER_HEADER.md  # X-Jenkins-User header guide
│   └── INTEGRATION_GUIDE.md    # Integration examples
│
├── Troubleshooting/            # Problem-solving guides
│   ├── CONNECTION_ISSUES.md    # Connection problems
│   ├── TEST_CONNECTION.md      # Connection testing
│   ├── QUICK_FIX.md            # Quick solutions
│   ├── TROUBLESHOOTING_403.md   # 403 error guide
│   ├── ERROR_HANDLING.md       # Error handling
│   ├── DEBUGGING_GUIDE.md      # Debugging techniques
│   ├── COMMON_ISSUES.md        # Common problems
│   └── SOLUTION_SUMMARY.md     # Solution summary
│
├── Scripts/                    # Example and test scripts
│   ├── jenkins_trigger_example.sh  # API usage example
│   ├── test_jenkins_connection.sh  # Connection test
│   ├── test_jenkins_direct.sh      # Direct testing
│   ├── test_api_manually.sh        # Manual API test
│   └── error_handling_example.js   # Error handling
│
├── Configuration/              # Configuration files
│   └── .env.example              # Environment example
│
└── DOCS_INDEX.md               # This file
```

## Quick Start

### Your Issue: Cannot Connect to Jenkins

**Error:**
```json
{
  "error": "Cannot connect to Jenkins server at http://10.0.1.84:8080..."
}
```

**Solution:**
1. Read **[CONNECTION_ISSUES.md](CONNECTION_ISSUES.md)**
2. Run `test_jenkins_connection.sh`
3. Check **[QUICK_FIX.md](QUICK_FIX.md)**

### Working Example

```bash
# Test connection
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" http://10.0.1.84:8080/api/json

# Use API
curl -X POST \
  -H "X-Jenkins-Url: http://10.0.1.84:8080" \
  -H "X-Jenkins-Token: 11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  -H "X-Jenkins-Job: plb-sonarqube" \
  -H "X-Jenkins-User: sonar" \
  http://localhost:8080/api/v1/jenkins/trigger
```

## API Documentation

### Main API
- **[API_DOCUMENTATION.md](API_DOCUMENTATION.md)** - Complete API reference
- **[JENKINS_API.md](JENKINS_API.md)** - Jenkins-specific details
- **[JENKINS_USER_HEADER.md](JENKINS_USER_HEADER.md)** - X-Jenkins-User header

### Integration
- **[INTEGRATION_GUIDE.md](INTEGRATION_GUIDE.md)** - Integration examples
- `jenkins_trigger_example.sh` - Bash example
- `error_handling_example.js` - Node.js example

## Troubleshooting

### Connection Issues
- **[CONNECTION_ISSUES.md](CONNECTION_ISSUES.md)** - Cannot connect to Jenkins
- **[TEST_CONNECTION.md](TEST_CONNECTION.md)** - Connection testing guide
- **[QUICK_FIX.md](QUICK_FIX.md)** - Quick solutions

### Authentication Issues
- **[TROUBLESHOOTING_403.md](TROUBLESHOOTING_403.md)** - 403 error guide
- **[ERROR_HANDLING.md](ERROR_HANDLING.md)** - Error handling reference

### Debugging
- **[DEBUGGING_GUIDE.md](DEBUGGING_GUIDE.md)** - Advanced debugging
- **[COMMON_ISSUES.md](COMMON_ISSUES.md)** - Common problems
- **[SOLUTION_SUMMARY.md](SOLUTION_SUMMARY.md)** - Solution summary

## Test Scripts

### Connection Testing
- `test_jenkins_connection.sh` - Comprehensive connection test
- `test_jenkins_direct.sh` - Direct Jenkins testing
- `test_api_manually.sh` - Manual API testing

### Example Usage
- `jenkins_trigger_example.sh` - Complete API example
- `error_handling_example.js` - Error handling example

## Configuration

- `.env.example` - Environment configuration example

## Getting Help

### For Connection Issues
```
1. Read CONNECTION_ISSUES.md
2. Run test_jenkins_connection.sh
3. Check QUICK_FIX.md
```

### For Authentication Issues
```
1. Read TROUBLESHOOTING_403.md
2. Check ERROR_HANDLING.md
3. Review JENKINS_API.md
```

### For General Issues
```
1. Check COMMON_ISSUES.md
2. Follow DEBUGGING_GUIDE.md
3. Read SOLUTION_SUMMARY.md
```

## Support Information

When requesting support, provide:
1. Exact error message
2. Jenkins URL and version
3. Network topology
4. Steps you've tried
5. Relevant log entries

## Quick Reference

### API Endpoint
```
POST /api/v1/jenkins/trigger
```

### Required Headers
```
X-Jenkins-Url: https://jenkins.example.com
X-Jenkins-Token: your-jenkins-token
X-Jenkins-Job: your-job-name
X-Jenkins-User: sonar (optional)
```

### Success Response
```json
{
  "success": true,
  "message": "Jenkins job triggered successfully",
  "job": "your-job-name"
}
```

### Common Errors
```json
// Connection error
{"error": "Cannot connect to Jenkins server..."}

// Authentication error
{"error": "Jenkins authentication failed..."}

// Permission error
{"error": "Jenkins job trigger failed: HTTP 403..."}

// Missing headers
{"error": "Missing or invalid headers"}
```

## Documentation Index

For complete documentation index, see:
- **[DOCUMENTATION.md](DOCUMENTATION.md)** - Complete documentation guide
- **[README_DOCS.md](README_DOCS.md)** - Detailed documentation index