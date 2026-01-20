# Integration Guide

## Jenkins Trigger API Integration

This guide explains how to integrate the Jenkins Trigger API with your existing systems.

### Overview

The Jenkins Trigger API allows you to trigger Jenkins jobs from external systems without requiring authentication to the SonarQube Report Generator application itself. This is useful for:

- CI/CD pipelines that need to trigger Jenkins jobs
- Automated workflows
- External monitoring systems
- Webhook integrations

### Integration Methods

#### 1. Direct HTTP Request

```bash
curl -X POST \
  -H "X-Jenkins-Url: https://your-jenkins-server.com" \
  -H "X-Jenkins-Token: your-jenkins-api-token" \
  -H "X-Jenkins-Job: your-job-name" \
  https://sonar-report-generator/api/v1/jenkins/trigger
```

#### 2. From Node.js

```javascript
const axios = require('axios');

async function triggerJenkinsJob() {
  try {
    const response = await axios.post(
      'https://sonar-report-generator/api/v1/jenkins/trigger',
      {}, // empty body
      {
        headers: {
          'X-Jenkins-Url': 'https://your-jenkins-server.com',
          'X-Jenkins-Token': 'your-jenkins-api-token',
          'X-Jenkins-Job': 'your-job-name',
          'X-Jenkins-User': 'sonar' // optional, defaults to 'sonar'
        }
      }
    );
    console.log('Success:', response.data);
  } catch (error) {
    console.error('Error:', error.response?.data || error.message);
  }
}
```

#### 3. From Python

```python
import requests

url = "https://sonar-report-generator/api/v1/jenkins/trigger"
headers = {
    "X-Jenkins-Url": "https://your-jenkins-server.com",
    "X-Jenkins-Token": "your-jenkins-api-token",
    "X-Jenkins-Job": "your-job-name",
    "X-Jenkins-User": "sonar"  # optional, defaults to 'sonar'
}

response = requests.post(url, headers=headers)

if response.status_code == 200:
    print("Success:", response.json())
else:
    print("Error:", response.json())
```

#### 4. From Go

```go
package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
)

func main() {
    url := "https://sonar-report-generator/api/v1/jenkins/trigger"
    
    client := &http.Client{}
    req, err := http.NewRequest("POST", url, nil)
    if err != nil {
        panic(err)
    }
    
    req.Header.Set("X-Jenkins-Url", "https://your-jenkins-server.com")
    req.Header.Set("X-Jenkins-Token", "your-jenkins-api-token")
    req.Header.Set("X-Jenkins-Job", "your-job-name")
    req.Header.Set("X-Jenkins-User", "sonar") // optional, defaults to 'sonar'
    
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### Webhook Integration

You can set up webhooks in various systems to trigger Jenkins jobs:

#### GitHub/GitLab Webhook Example

```yaml
# .github/workflows/trigger-jenkins.yml
name: Trigger Jenkins Job

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  trigger-jenkins:
    runs-on: ubuntu-latest
    steps:
      - name: Trigger Jenkins Job
        run: |
          curl -X POST \
            -H "X-Jenkins-Url: ${{ secrets.JENKINS_URL }}" \
            -H "X-Jenkins-Token: ${{ secrets.JENKINS_TOKEN }}" \
            -H "X-Jenkins-Job: ${{ secrets.JENKINS_JOB }}" \
            https://sonar-report-generator/api/v1/jenkins/trigger
```

### Security Considerations

1. **Network Security**: Ensure the API endpoint is only accessible from trusted networks
2. **Jenkins Token Security**: Store Jenkins tokens securely and rotate them regularly
3. **HTTPS**: Always use HTTPS for API communications
4. **Rate Limiting**: Consider adding rate limiting to prevent abuse
5. **Logging**: Monitor API usage through server logs

### Error Handling

The API returns detailed error messages that you should handle in your integration:

```json
{
  "error": "Missing or invalid headers"
}
```

```json
{
  "error": "Failed to trigger Jenkins job: HTTP 404 - Not Found"
}
```

### Troubleshooting

**Problem**: "Missing or invalid headers"
- **Solution**: Ensure all three required headers are present and not empty

**Problem**: "Failed to get Jenkins crumb"
- **Solution**: Check Jenkins URL is correct and accessible
- **Solution**: Verify Jenkins CSRF protection is enabled

**Problem**: "Jenkins job trigger failed: HTTP 404"
- **Solution**: Verify the job name is correct
- **Solution**: Check Jenkins token has sufficient permissions

**Problem**: Connection timeout
- **Solution**: Check network connectivity between servers
- **Solution**: Verify Jenkins server is running

### Best Practices

1. **Environment Variables**: Store sensitive information (tokens, URLs) in environment variables
2. **Retry Logic**: Implement retry logic for transient failures
3. **Logging**: Log API calls and responses for auditing
4. **Validation**: Validate inputs before making API calls
5. **Timeouts**: Set appropriate timeout values for HTTP requests