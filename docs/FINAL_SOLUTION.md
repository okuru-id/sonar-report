# Final Solution for Your Issue

## Your Current Status

✅ **Jenkins is accessible** - You can connect to `http://10.0.1.84:8080`
✅ **Basic authentication works** - `curl -u "sonar:token" /api/json` returns valid JSON
✅ **Job exists** - `plb-sonarqube` is in the jobs list
❌ **API fails** - Getting "User does not have permission to access crumb issuer"

## Root Cause

The issue is **CSRF crumb permissions**. Your user "sonar" can access basic Jenkins API but **cannot access the crumb issuer endpoint** (`/crumbIssuer/api/json`).

## The Fix

### Step 1: Grant Overall/Read Permission

**Log in to Jenkins as admin** and:

1. Go to: **Manage Jenkins** > **Security** > **Configure Global Security**
2. Find the **"Authorization"** section
3. Locate user **"sonar"** in the matrix
4. Check these permissions:
   - ✅ Overall/Read
   - ✅ Job/Build
   - ✅ Job/Read
5. **Save** the configuration

### Step 2: Test Crumb Access

```bash
# Test if user can now access crumb issuer
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  http://10.0.1.84:8080/crumbIssuer/api/json
```

**Expected response:**
```json
{
  "crumbRequestField": "Jenkins-Crumb",
  "crumb": "abc123xyz456"
}
```

### Step 3: Try the API Again

```bash
curl -X POST \
  -H "X-Jenkins-Url: http://10.0.1.84:8080" \
  -H "X-Jenkins-Token: 11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  -H "X-Jenkins-Job: plb-sonarqube" \
  -H "X-Jenkins-User: sonar" \
  http://localhost:8080/api/v1/jenkins/trigger
```

**Expected success:**
```json
{
  "success": true,
  "message": "Jenkins job triggered successfully",
  "job": "plb-sonarqube"
}
```

## Alternative Solutions

### If You Can't Modify Permissions

**Use admin credentials temporarily:**

```bash
curl -X POST \
  -H "X-Jenkins-Url: http://10.0.1.84:8080" \
  -H "X-Jenkins-Token: admin-token" \
  -H "X-Jenkins-Job: plb-sonarqube" \
  -H "X-Jenkins-User: admin" \
  http://localhost:8080/api/v1/jenkins/trigger
```

### If CSRF is Not Needed

**Disable CSRF temporarily (for testing only):**

1. Go to: Manage Jenkins > Security > Configure Global Security
2. Uncheck "Prevent Cross Site Request Forgery exploits"
3. Save and try the API
4. **Re-enable CSRF** after testing for security

## Verification Steps

### 1. Verify User Permissions

```bash
# Check if user has Overall/Read
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  http://10.0.1.84:8080/crumbIssuer/api/json
```

### 2. Test Manual Job Trigger

```bash
# Get crumb
CRUMB_JSON=$(curl -s -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  "http://10.0.1.84:8080/crumbIssuer/api/json")
CRUMB_FIELD=$(echo $CRUMB_JSON | jq -r .crumbRequestField)
CRUMB=$(echo $CRUMB_JSON | jq -r .crumb)

# Trigger job manually
curl -X POST \
  -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  -H "$CRUMB_FIELD: $CRUMB" \
  http://10.0.1.84:8080/job/plb-sonarqube/build
```

### 3. Check Jenkins Security Configuration

```bash
# Get security config
curl -u "admin:admin-token" \
  http://10.0.1.84:8080/securityRealm/api/json
```

## Documentation References

### For This Issue
- **[CSRF Issues Guide](CSRF_ISSUES.md)** - Detailed CSRF troubleshooting
- **[Authentication Issues Guide](AUTHENTICATION_ISSUES.md)** - Permission troubleshooting

### General Documentation
- **[API Documentation](API_DOCUMENTATION.md)** - Complete API reference
- **[Connection Issues Guide](CONNECTION_ISSUES.md)** - Connection troubleshooting

## Summary

**The Problem:** User "sonar" doesn't have permission to access `/crumbIssuer/api/json`

**The Solution:** Grant "sonar" user **Overall/Read** permission in Jenkins

**Verification:** Test crumb access manually before using the API

**Expected Result:** API will work and trigger Jenkins jobs successfully