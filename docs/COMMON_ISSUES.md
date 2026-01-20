# Common Issues and Solutions

## 1. Cannot Connect to Jenkins Server

**Error:**
```json
{
  "error": "Cannot connect to Jenkins server at http://10.0.1.84:8080..."
}
```

**Solutions:**
- ✅ Check Jenkins is running: `systemctl status jenkins`
- ✅ Test network connectivity: `ping 10.0.1.84`
- ✅ Verify port is open: `nc -zv 10.0.1.84 8080`
- ✅ Check firewall settings: `ufw status`
- ✅ Test URL in browser first

## 2. Authentication Failed

**Error:**
```json
{
  "error": "Jenkins authentication failed: Invalid credentials or insufficient permissions..."
}
```

**Solutions:**
- ✅ Verify Jenkins token is correct
- ✅ Check user has "Overall/Read" permission
- ✅ Ensure user has "Job/Build" permission
- ✅ Test credentials manually: `curl -u "user:token" http://jenkins/api/json`
- ✅ Regenerate API token if needed

## 3. Permission Denied (403)

**Error:**
```json
{
  "error": "Jenkins job trigger failed: HTTP 403 - Forbidden"
}
```

**Solutions:**
- ✅ Check CSRF protection is enabled in Jenkins
- ✅ Verify user can access crumb issuer
- ✅ Ensure user has build permission for the job
- ✅ Test crumb endpoint: `curl -u "user:token" http://jenkins/crumbIssuer/api/json`
- ✅ Check Jenkins security settings

## 4. Job Not Found (404)

**Error:**
```json
{
  "error": "Jenkins job trigger failed: HTTP 404 - Not Found"
}
```

**Solutions:**
- ✅ Verify job name is correct (case-sensitive)
- ✅ Check job exists in Jenkins
- ✅ Ensure correct job path (for folder jobs)
- ✅ Test job URL: `curl -u "user:token" http://jenkins/job/JOB_NAME/api/json`

## 5. Missing Headers

**Error:**
```json
{
  "error": "Missing or invalid headers"
}
```

**Solutions:**
- ✅ Provide all required headers
- ✅ Check header names are correct
- ✅ Ensure headers are not empty
- ✅ Verify header format

## 6. CSRF Crumb Failed

**Error:**
```json
{
  "error": "Failed to get Jenkins crumb: HTTP 403"
}
```

**Solutions:**
- ✅ Enable CSRF protection in Jenkins
- ✅ Grant user permission to access crumb issuer
- ✅ Test crumb endpoint manually
- ✅ Check Jenkins security configuration

## 7. Connection Timeout

**Error:**
```json
{
  "error": "Failed to connect to Jenkins: dial tcp: i/o timeout"
}
```

**Solutions:**
- ✅ Check network connectivity
- ✅ Verify Jenkins is running
- ✅ Test from same network
- ✅ Check firewall timeout settings

## 8. SSL Certificate Issues

**Error:**
```json
{
  "error": "Failed to connect to Jenkins: x509: certificate signed by unknown authority"
}
```

**Solutions:**
- ✅ Add certificate to trusted store
- ✅ Use HTTP instead of HTTPS (for testing)
- ✅ Disable SSL verification (not recommended for production)
- ✅ Update CA certificates

## Quick Troubleshooting Checklist

### Before Contacting Support

1. [ ] Jenkins server is running
2. [ ] Network connectivity is working
3. [ ] Credentials are correct
4. [ ] User has proper permissions
5. [ ] Job name is correct
6. [ ] CSRF protection is enabled
7. [ ] Firewall allows the connection
8. [ ] URL is correct and accessible

### Test Commands

```bash
# Test connectivity
ping 10.0.1.84
nc -zv 10.0.1.84 8080

# Test authentication
curl -u "user:token" http://jenkins/api/json

# Test crumb
curl -u "user:token" http://jenkins/crumbIssuer/api/json

# Test job
curl -u "user:token" http://jenkins/job/JOB_NAME/api/json
```

### Working Configuration Example

```bash
# Working API call example
curl -X POST \
  -H "X-Jenkins-Url: http://10.0.1.84:8080" \
  -H "X-Jenkins-Token: 11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  -H "X-Jenkins-Job: plb-sonarqube" \
  -H "X-Jenkins-User: sonar" \
  http://localhost:8080/api/v1/jenkins/trigger
```

Expected success:
```json
{
  "success": true,
  "message": "Jenkins job triggered successfully",
  "job": "plb-sonarqube"
}
```