# Documentation Index

## 📚 API Documentation

### Core API
- **[API_DOCUMENTATION.md](API_DOCUMENTATION.md)** - Complete API reference with examples
- **[JENKINS_API.md](JENKINS_API.md)** - Jenkins-specific API details
- **[JENKINS_USER_HEADER.md](JENKINS_USER_HEADER.md)** - X-Jenkins-User header documentation

### Integration
- **[INTEGRATION_GUIDE.md](INTEGRATION_GUIDE.md)** - Integration examples (Node.js, Python, Go, etc.)
- `jenkins_trigger_example.sh` - Bash integration example
- `error_handling_example.js` - Error handling example

## 🔧 Troubleshooting

### Connection Issues
- **[CONNECTION_ISSUES.md](CONNECTION_ISSUES.md)** - Cannot connect to Jenkins server
- **[TEST_CONNECTION.md](TEST_CONNECTION.md)** - Step-by-step connection testing
- **[QUICK_FIX.md](QUICK_FIX.md)** - Quick solutions for common problems

### Authentication & Permissions
- **[TROUBLESHOOTING_403.md](TROUBLESHOOTING_403.md)** - HTTP 403 error troubleshooting
- **[ERROR_HANDLING.md](ERROR_HANDLING.md)** - Complete error handling guide

### Debugging
- **[DEBUGGING_GUIDE.md](DEBUGGING_GUIDE.md)** - Advanced debugging techniques
- **[COMMON_ISSUES.md](COMMON_ISSUES.md)** - Common issues and solutions
- **[SOLUTION_SUMMARY.md](SOLUTION_SUMMARY.md)** - Solution summary for your issue

## 🚀 Quick Start

### Your Issue: Cannot Connect to Jenkins

**Error:**
```json
{
  "error": "Cannot connect to Jenkins server at http://10.0.1.84:8080..."
}
```

**Quick Fix:**
```bash
# 1. Check Jenkins is running (on Jenkins server)
systemctl status jenkins

# 2. Test connectivity (on application server)
ping 10.0.1.84
nc -zv 10.0.1.84 8080

# 3. Test HTTP access
curl http://10.0.1.84:8080

# 4. Try the API
curl -X POST -H "X-Jenkins-Url: http://10.0.1.84:8080" -H "X-Jenkins-Token: YOUR_TOKEN" -H "X-Jenkins-Job: plb-sonarqube" -H "X-Jenkins-User: sonar" http://localhost:8080/api/v1/jenkins/trigger
```

## 📖 Documentation Guide

### For API Users
1. Start with **[API_DOCUMENTATION.md](API_DOCUMENTATION.md)**
2. Read **[JENKINS_API.md](JENKINS_API.md)** for Jenkins specifics
3. Check **[INTEGRATION_GUIDE.md](INTEGRATION_GUIDE.md)** for examples

### For Troubleshooting
1. Identify your error type
2. Read the corresponding troubleshooting guide
3. Follow the step-by-step instructions
4. Use the test scripts provided

### For Developers
1. Review **[ERROR_HANDLING.md](ERROR_HANDLING.md)**
2. Check **[DEBUGGING_GUIDE.md](DEBUGGING_GUIDE.md)**
3. Use the example scripts as reference

## 🔍 Test Scripts

### Connection Testing
- `test_jenkins_connection.sh` - Comprehensive Jenkins connection test
- `test_jenkins_direct.sh` - Direct Jenkins testing
- `test_api_manually.sh` - Manual API testing

### Example Usage
- `jenkins_trigger_example.sh` - Complete API usage example
- `error_handling_example.js` - Error handling in Node.js

## 📝 Configuration

- `.env.example` - Environment configuration example
- `DOCUMENTATION.md` - Complete documentation index

## 🆘 Getting Help

### Connection Issues
```
Cannot connect to Jenkins server
→ Read CONNECTION_ISSUES.md
→ Run test_jenkins_connection.sh
→ Check QUICK_FIX.md
```

### Authentication Issues
```
Authentication failed / 403 Forbidden
→ Read TROUBLESHOOTING_403.md
→ Check ERROR_HANDLING.md
→ Review JENKINS_API.md
```

### General Issues
```
Other errors
→ Check COMMON_ISSUES.md
→ Follow DEBUGGING_GUIDE.md
→ Read SOLUTION_SUMMARY.md
```

## 🎯 Quick Reference

### API Endpoint
```
POST /api/v1/jenkins/trigger
```

### Required Headers
```
X-Jenkins-Url: https://jenkins.example.com
X-Jenkins-Token: your-jenkins-token
X-Jenkins-Job: your-job-name
X-Jenkins-User: sonar (optional, defaults to "sonar")
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

## 📚 Complete Documentation List

### API Reference
- API_DOCUMENTATION.md
- JENKINS_API.md
- JENKINS_USER_HEADER.md
- INTEGRATION_GUIDE.md

### Troubleshooting
- CONNECTION_ISSUES.md
- TEST_CONNECTION.md
- QUICK_FIX.md
- TROUBLESHOOTING_403.md
- ERROR_HANDLING.md
- DEBUGGING_GUIDE.md
- COMMON_ISSUES.md
- SOLUTION_SUMMARY.md

### Documentation
- DOCUMENTATION.md
- README_DOCS.md (this file)

### Example Scripts
- test_jenkins_connection.sh
- test_jenkins_direct.sh
- jenkins_trigger_example.sh
- test_api_manually.sh
- error_handling_example.js

### Configuration
- .env.example