# CSRF Protection Issues

## Your Issue

You're getting this error:
```json
{
  "error": "Jenkins authentication failed: User does not have permission to access crumb issuer..."
}
```

But you've verified:
- ✅ Jenkins is accessible
- ✅ Authentication works for basic API calls
- ✅ Job exists and is accessible

## Root Cause

The issue is with **CSRF (Cross-Site Request Forgery) protection** in Jenkins. Your user "sonar" doesn't have permission to access the crumb issuer endpoint.

## About CSRF in Jenkins

Jenkins uses CSRF protection to prevent unauthorized requests. When CSRF is enabled:

1. **Crumb Issuer** (`/crumbIssuer/api/json`) provides a token for each request
2. **All POST requests** must include this crumb token
3. **User must have permission** to access crumb issuer

## Solutions

### Solution 1: Grant User Proper Permissions

1. **Log in to Jenkins as admin**
2. **Go to**: Manage Jenkins > Security > Configure Global Security
3. **Find your user** ("sonar") and grant:
   - Overall/Read ✅
   - Job/Build ✅
   - Job/Read ✅

### Solution 2: Test Crumb Access Manually

```bash
# Test if user can access crumb issuer
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

**If you get 403:** User doesn't have permission.

### Solution 3: Use Admin User Temporarily

```bash
# Try with admin credentials
curl -X POST \
  -H "X-Jenkins-Url: http://10.0.1.84:8080" \
  -H "X-Jenkins-Token: admin-token" \
  -H "X-Jenkins-Job: plb-sonarqube" \
  -H "X-Jenkins-User: admin" \
  http://localhost:8080/api/v1/jenkins/trigger
```

### Solution 4: Check Jenkins Security Configuration

1. **Go to**: Manage Jenkins > Security > Configure Global Security
2. **Verify**: "Prevent Cross Site Request Forgery exploits" is checked ✅
3. **Check**: Security Realm configuration

### Solution 5: Grant Specific Permissions

If using **Matrix-based security**:

1. Go to: Manage Jenkins > Security > Configure Global Security
2. Under "Authorization", find your user
3. Grant these permissions:
   - Overall/Read ✅
   - Job/Build ✅
   - Job/Read ✅
   - Job/Workspace ✅ (optional)

## Testing Steps

### Step 1: Test Basic Authentication

```bash
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  http://10.0.1.84:8080/api/json
```

**Should return:** Jenkins API JSON response ✅

### Step 2: Test Crumb Issuer

```bash
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  http://10.0.1.84:8080/crumbIssuer/api/json
```

**Should return:** Crumb JSON with `crumbRequestField` and `crumb` ✅

### Step 3: Test Job Access

```bash
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  http://10.0.1.84:8080/job/plb-sonarqube/api/json
```

**Should return:** Job information JSON ✅

### Step 4: Test Manual Job Trigger

```bash
# Get crumb first
CRUMB_JSON=$(curl -s -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  "http://10.0.1.84:8080/crumbIssuer/api/json")
CRUMB_FIELD=$(echo $CRUMB_JSON | jq -r .crumbRequestField)
CRUMB=$(echo $CRUMB_JSON | jq -r .crumb)

# Trigger job
curl -X POST \
  -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  -H "$CRUMB_FIELD: $CRUMB" \
  http://10.0.1.84:8080/job/plb-sonarqube/build
```

**Should return:** HTTP 201 Created or 200 OK ✅

## Jenkins Permission Matrix

### Required Permissions for "sonar" User

| Permission | Required | Purpose |
|------------|----------|---------|
| Overall/Read | ✅ Yes | Access Jenkins API |
| Job/Build | ✅ Yes | Trigger jobs |
| Job/Read | ✅ Yes | Access job information |
| Job/Workspace | ⚠️ Optional | Access workspace files |
| Credentials/View | ⚠️ Optional | View credentials |
| Run/Scripts | ⚠️ Optional | Run scripts |

### How to Check Permissions

1. **Log in as admin**
2. **Go to**: Manage Jenkins > Security > Configure Global Security
3. **Find** "Authorization" section
4. **Check** your user's permissions

## Common CSRF Issues

### Issue 1: User Doesn't Have Overall/Read

**Symptoms:**
- Can access jobs but not crumb issuer
- Gets 403 on `/crumbIssuer/api/json`

**Fix:** Grant Overall/Read permission

### Issue 2: CSRF Protection Disabled

**Symptoms:**
- Crumb issuer returns 404
- No CSRF protection

**Fix:** Enable CSRF protection in Jenkins security settings

### Issue 3: Wrong User

**Symptoms:**
- Authentication works for some endpoints but not others
- Inconsistent permission errors

**Fix:** Use the correct user with proper permissions

### Issue 4: Token Issues

**Symptoms:**
- Authentication fails intermittently
- Works sometimes, fails others

**Fix:** Regenerate API token

## Advanced Troubleshooting

### Check Jenkins Logs

```bash
# On Jenkins server
cat /var/log/jenkins/jenkins.log | grep -i "crumb\|csrf\|auth"

# Or follow logs
journalctl -u jenkins -f
```

### Check User Permissions via API

```bash
# Get all users
curl -u "admin:admin-token" \
  http://10.0.1.84:8080/asynchPeople/api/json

# Check specific user permissions
curl -u "admin:admin-token" \
  http://10.0.1.84:8080/user/sonar/api/json
```

### Test with Different Users

```bash
# Test with admin
curl -u "admin:admin-token" \
  http://10.0.1.84:8080/crumbIssuer/api/json

# Test with sonar
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  http://10.0.1.84:8080/crumbIssuer/api/json
```

## Solution Summary

1. **Grant "sonar" user Overall/Read permission**
2. **Verify CSRF protection is enabled**
3. **Test crumb issuer manually**
4. **Use admin user temporarily for testing**
5. **Check Jenkins security configuration**

## Working Example

When permissions are correct:

```bash
# Test crumb access
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  http://10.0.1.84:8080/crumbIssuer/api/json

# Should return:
# {
#   "crumbRequestField": "Jenkins-Crumb",
#   "crumb": "abc123xyz456"
# }

# Then use the API
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