# API Documentation

## Jenkins Trigger API

Endpoint untuk memicu Jenkins job build tanpa autentikasi.

### Endpoint

```
POST /api/v1/jenkins/trigger
```

### Request Headers

| Header | Type | Required | Description |
|--------|------|----------|-------------|
| X-Jenkins-Url | string | Yes | URL dasar Jenkins server (e.g., `https://jenkins.example.com`) |
| X-Jenkins-Token | string | Yes | Token autentikasi Jenkins |
| X-Jenkins-Job | string | Yes | Nama job Jenkins yang akan dipicu (e.g., `plb-sonarqube`) |
| X-Jenkins-User | string | No | Username untuk autentikasi Jenkins (default: `sonar`) |

### Example Request

```bash
curl -X POST \
  -H "X-Jenkins-Url: https://jenkins.example.com" \
  -H "X-Jenkins-Token: your-jenkins-token" \
  -H "X-Jenkins-Job: plb-sonarqube" \
  -H "X-Jenkins-User: sonar" \
  http://localhost:8080/api/v1/jenkins/trigger
```

### Example Response (Success)

```json
{
  "success": true,
  "message": "Jenkins job triggered successfully",
  "job": "plb-sonarqube"
}
```

### Example Response (Error)

```json
{
  "error": "Missing or invalid headers"
}

OR

```json
{
  "error": "Failed to trigger Jenkins job: HTTP 404 - Not Found"
}
```

### Error Codes

| Status Code | Description |
|-------------|-------------|
| 400 | Missing or invalid headers |
| 400 | Cannot connect to Jenkins server (network/connection issues) |
| 400 | Jenkins authentication failed (invalid credentials or insufficient permissions) |
| 500 | Failed to trigger Jenkins job (connection error, authentication error, etc.) |

### Implementation Details

The API performs the following steps:

1. **Header Validation**: Validates that all required headers are present
2. **Crumb Request**: Gets CSRF crumb from Jenkins (`/crumbIssuer/api/json`)
3. **Job Trigger**: Calls Jenkins build endpoint with proper authentication and CSRF headers
4. **Response Handling**: Returns success or detailed error information

### Authentication

- This endpoint does not require application-level authentication
- Jenkins authentication is handled using the provided token and username (default: "sonar")
- CSRF protection is handled automatically by getting and using Jenkins crumb

### Notes

- Ensure the Jenkins URL is accessible from the server running this application
- The Jenkins token must have sufficient permissions to trigger builds
- The API handles both HTTP 200 and 201 responses from Jenkins as success
- Error responses include the HTTP status code and response body from Jenkins for debugging
- For 403 errors, check the troubleshooting guide for common solutions
- For connection errors, verify network connectivity and Jenkins server status