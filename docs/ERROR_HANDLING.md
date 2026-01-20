# Error Handling Guide

## API Error Responses

The Jenkins Trigger API provides detailed error messages to help troubleshoot issues.

### Common Error Scenarios

#### 1. Missing or Invalid Headers

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/jenkins/trigger
```

**Response:**
```json
{
  "error": "Missing or invalid headers"
}
```

**Solution:** Provide all required headers: `X-Jenkins-Url`, `X-Jenkins-Token`, `X-Jenkins-Job`

#### 2. Jenkins Authentication Failed (403)

**Request:**
```bash
curl -X POST \
  -H "X-Jenkins-Url: https://jenkins.example.com" \
  -H "X-Jenkins-Token: wrong-token" \
  -H "X-Jenkins-Job: my-job" \
  http://localhost:8080/api/v1/jenkins/trigger
```

**Response:**
```json
{
  "error": "Jenkins authentication failed: Invalid credentials or insufficient permissions. Please check your Jenkins URL, token, and user permissions."
}
```

**Solution:**
- Verify Jenkins token is correct
- Check user permissions in Jenkins
- Test credentials manually with `curl -u "user:token" https://jenkins-url/api/json`

#### 3. Jenkins Connection Failed

**Request:**
```bash
curl -X POST \
  -H "X-Jenkins-Url: https://wrong-jenkins-url.com" \
  -H "X-Jenkins-Token: my-token" \
  -H "X-Jenkins-Job: my-job" \
  http://localhost:8080/api/v1/jenkins/trigger
```

**Response:**
```json
{
  "error": "Failed to get Jenkins crumb: Get \"https://wrong-jenkins-url.com/crumbIssuer/api/json\": dial tcp: lookup wrong-jenkins-url.com: no such host"
}
```

**Solution:**
- Verify Jenkins URL is correct
- Check network connectivity
- Test URL in browser first

#### 4. Jenkins Job Not Found

**Request:**
```bash
curl -X POST \
  -H "X-Jenkins-Url: https://jenkins.example.com" \
  -H "X-Jenkins-Token: my-token" \
  -H "X-Jenkins-Job: non-existent-job" \
  http://localhost:8080/api/v1/jenkins/trigger
```

**Response:**
```json
{
  "error": "Jenkins job trigger failed: HTTP 404 - Not Found"
}
```

**Solution:**
- Verify job name is correct
- Check job exists in Jenkins
- Ensure correct case sensitivity

#### 5. CSRF Protection Issues

**Request:**
```bash
curl -X POST \
  -H "X-Jenkins-Url: https://jenkins.example.com" \
  -H "X-Jenkins-Token: my-token" \
  -H "X-Jenkins-Job: my-job" \
  http://localhost:8080/api/v1/jenkins/trigger
```

**Response:**
```json
{
  "error": "Failed to get Jenkins crumb: HTTP 403"
}
```

**Solution:**
- Check CSRF protection is enabled in Jenkins
- Verify user has permission to access crumb issuer
- Test crumb endpoint manually: `curl -u "user:token" https://jenkins-url/crumbIssuer/api/json`

## Error Handling Best Practices

### Client-Side Error Handling

```javascript
// JavaScript example with proper error handling
async function triggerJenkinsJob() {
  try {
    const response = await axios.post(
      '/api/v1/jenkins/trigger',
      {},
      {
        headers: {
          'X-Jenkins-Url': jenkinsUrl,
          'X-Jenkins-Token': jenkinsToken,
          'X-Jenkins-Job': jenkinsJob
        }
      }
    );
    
    if (response.data.success) {
      showSuccess("Jenkins job triggered successfully!");
    }
  } catch (error) {
    if (error.response) {
      // Server responded with error status
      const errorMsg = error.response.data.error;
      
      if (errorMsg.includes("authentication failed")) {
        showError("Authentication failed. Please check your Jenkins credentials.");
      } else if (errorMsg.includes("404")) {
        showError("Jenkins job not found. Please check the job name.");
      } else if (errorMsg.includes("403")) {
        showError("Permission denied. Check your Jenkins user permissions.");
      } else {
        showError(`Jenkins error: ${errorMsg}`);
      }
    } else if (error.request) {
      // No response received
      showError("Cannot connect to server. Please check your network connection.");
    } else {
      // Other errors
      showError(`Error: ${error.message}`);
    }
  }
}
```

### Server-Side Logging

The API logs detailed error information that can help with troubleshooting:

```
// Example log entries

// Missing headers
2024/01/20 10:00:00 [Error] Missing or invalid headers

// Authentication failure
2024/01/20 10:01:00 [Error] Jenkins authentication failed: Invalid credentials or insufficient permissions

// Connection failure
2024/01/20 10:02:00 [Error] Failed to get Jenkins crumb: connection refused

// Job not found
2024/01/20 10:03:00 [Error] Jenkins job trigger failed: HTTP 404 - Not Found
```

## Troubleshooting Flowchart

```
Start
  │
  ▼
Check Error Message
  │
  ├── "Missing or invalid headers"
  │   │
  │   ▼
  │   Provide all required headers
  │
  ├── "authentication failed"
  │   │
  │   ▼
  │   Check credentials & permissions
  │
  ├── "404"
  │   │
  │   ▼
  │   Verify Jenkins URL and job name
  │
  ├── "403"
  │   │
  │   ▼
  │   Check CSRF settings & permissions
  │
  └── Other connection errors
      │
      ▼
      Check network & Jenkins availability
```

## Common Solutions

### 1. Authentication Issues

- **Regenerate token**: Create a new API token in Jenkins
- **Check permissions**: Ensure user has "Build" permission for the job
- **Test manually**: Use curl to test credentials directly with Jenkins

### 2. Connection Issues

- **Verify URL**: Check for typos in Jenkins URL
- **Test connectivity**: Use ping or curl to test connection
- **Check firewall**: Ensure no firewall is blocking the connection

### 3. Permission Issues

- **Check Jenkins security**: Review global security settings
- **Verify user roles**: Ensure user has correct role assignments
- **Test with admin**: Try with admin credentials to isolate permission issues

### 4. CSRF Issues

- **Enable CSRF**: Ensure CSRF protection is enabled in Jenkins
- **Check permissions**: User needs permission to access crumb issuer
- **Test crumb endpoint**: Verify crumb endpoint is accessible

## Monitoring and Alerting

For production environments, consider implementing:

1. **Error rate monitoring**: Track frequency of different error types
2. **Alert thresholds**: Set alerts for high error rates
3. **Automated retries**: For transient errors like network issues
4. **Fallback mechanisms**: Alternative ways to trigger jobs if API fails

## Support Information

When reporting issues, please provide:

1. Exact error message received
2. Jenkins version and configuration
3. User permissions details
4. Network topology (if applicable)
5. Steps to reproduce the issue