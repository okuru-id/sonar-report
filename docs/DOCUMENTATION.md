# Jenkins Trigger API Documentation

## Overview

This documentation provides comprehensive information about the Jenkins Trigger API and troubleshooting guides.

## API Documentation

- **[API_DOCUMENTATION.md](API_DOCUMENTATION.md)** - Complete API reference
- **[JENKINS_API.md](JENKINS_API.md)** - Jenkins-specific API details
- **[JENKINS_USER_HEADER.md](JENKINS_USER_HEADER.md)** - X-Jenkins-User header documentation

## Troubleshooting Guides

### Connection Issues
- **[CONNECTION_ISSUES.md](CONNECTION_ISSUES.md)** - Cannot connect to Jenkins server
- **[TEST_CONNECTION.md](TEST_CONNECTION.md)** - Testing connection from application server
- **[QUICK_FIX.md](QUICK_FIX.md)** - Quick solutions for common problems

### Authentication Issues
- **[TROUBLESHOOTING_403.md](TROUBLESHOOTING_403.md)** - HTTP 403 error troubleshooting
- **[ERROR_HANDLING.md](ERROR_HANDLING.md)** - Complete error handling guide

### Debugging
- **[DEBUGGING_GUIDE.md](DEBUGGING_GUIDE.md)** - Step-by-step debugging guide
- **[COMMON_ISSUES.md](COMMON_ISSUES.md)** - Common issues and solutions

## Quick Start

### Your Current Issue

You're getting:
```json
{
  "error": "Cannot connect to Jenkins server at http://10.0.1.84:8080..."
}
```

### Immediate Solution

1. **Check Jenkins is running**: `systemctl status jenkins`
2. **Test connectivity**: `ping 10.0.1.84`
3. **Test port**: `nc -zv 10.0.1.84 8080`
4. **Test HTTP**: `curl http://10.0.1.84:8080`

### Working Example

```bash
# Test connection first
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" http://10.0.1.84:8080/api/json

# Then use the API
curl -X POST \
  -H "X-Jenkins-Url: http://10.0.1.84:8080" \
  -H "X-Jenkins-Token: 11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  -H "X-Jenkins-Job: plb-sonarqube" \
  -H "X-Jenkins-User: sonar" \
  http://localhost:8080/api/v1/jenkins/trigger
```

## Documentation Structure

```
DOCUMENTATION/
├── API_DOCUMENTATION.md          # Main API reference
├── JENKINS_API.md               # Jenkins-specific details
├── JENKINS_USER_HEADER.md       # X-Jenkins-User header guide
├── 
├── CONNECTION_ISSUES.md         # Connection troubleshooting
├── TEST_CONNECTION.md           # Connection testing guide
├── QUICK_FIX.md                 # Quick solutions
├── 
├── TROUBLESHOOTING_403.md        # 403 error guide
├── ERROR_HANDLING.md            # Error handling reference
├── 
├── DEBUGGING_GUIDE.md           # Advanced debugging
├── COMMON_ISSUES.md             # Common problems
├── SOLUTION_SUMMARY.md          # Solution summary
└── DOCUMENTATION.md             # This file
```

## Getting Help

### For Connection Issues
1. Read **[CONNECTION_ISSUES.md](CONNECTION_ISSUES.md)**
2. Run tests from **[TEST_CONNECTION.md](TEST_CONNECTION.md)**
3. Try solutions from **[QUICK_FIX.md](QUICK_FIX.md)**

### For Authentication Issues
1. Read **[TROUBLESHOOTING_403.md](TROUBLESHOOTING_403.md)**
2. Check **[ERROR_HANDLING.md](ERROR_HANDLING.md)**
3. Review **[JENKINS_API.md](JENKINS_API.md)**

### For General Issues
1. Check **[COMMON_ISSUES.md](COMMON_ISSUES.md)**
2. Follow **[DEBUGGING_GUIDE.md](DEBUGGING_GUIDE.md)**
3. Read **[SOLUTION_SUMMARY.md](SOLUTION_SUMMARY.md)**

## Support Information

When requesting support, please provide:

1. Exact error message
2. Jenkins URL and version
3. Network topology
4. Steps you've tried
5. Relevant log entries

## Example Scripts

- `test_jenkins_connection.sh` - Test Jenkins connection
- `test_jenkins_direct.sh` - Direct Jenkins testing
- `jenkins_trigger_example.sh` - API usage example
- `test_api_manually.sh` - Manual API testing

## Configuration Files

- `.env.example` - Environment configuration example
- `INTEGRATION_GUIDE.md` - Integration examples
