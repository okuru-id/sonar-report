# Authentication Issues

## Your Issue

You're getting this error:
```json
{
  "error": "Jenkins authentication failed: User does not have permission to access crumb issuer..."
}
```

But you've verified that basic authentication works:
```bash
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" http://10.0.1.84:8080/api/json
# ✅ Returns valid JSON
```

## Root Cause

The issue is **not with basic authentication**, but with **specific permissions** required to access the CSRF crumb issuer endpoint.

## Key Differences

### Basic Authentication ✅ (Works)
```bash
curl -u "sonar:token" http://10.0.1.84:8080/api/json
```
- Only requires basic read access
- Works with minimal permissions

### Crumb Issuer ❌ (Fails)
```bash
curl -u "sonar:token" http://10.0.1.84:8080/crumbIssuer/api/json
```
- Requires **Overall/Read** permission
- Requires **CSRF access** permission
- More restrictive than basic API access

## Solutions

### Solution 1: Grant Overall/Read Permission

1. **Log in to Jenkins as admin**
2. **Go to**: Manage Jenkins > Security > Configure Global Security
3. **Find user "sonar"** and grant:
   - Overall/Read ✅
   - Job/Build ✅
   - Job/Read ✅

### Solution 2: Use Admin User for Testing

```bash
# Test with admin first
curl -X POST \
  -H "X-Jenkins-Url: http://10.0.1.84:8080" \
  -H "X-Jenkins-Token: admin-token" \
  -H "X-Jenkins-Job: plb-sonarqube" \
  -H "X-Jenkins-User: admin" \
  http://localhost:8080/api/v1/jenkins/trigger
```

### Solution 3: Check Jenkins Security Realm

1. **Go to**: Manage Jenkins > Security > Configure Global Security
2. **Check Security Realm**: Ensure it's set to "Jenkins' own user database" or appropriate LDAP
3. **Verify user exists**: Check user "sonar" exists in the system

### Solution 4: Test Crumb Issuer Directly

```bash
# Test crumb issuer with your user
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  http://10.0.1.84:8080/crumbIssuer/api/json
```

**If 403:** User doesn't have permission to access crumb issuer.

**If 200:** User has permission, issue is elsewhere.

## Permission Requirements

### Minimum Permissions for "sonar" User

| Endpoint | Required Permission |
|----------|---------------------|
| `/api/json` | Overall/Read |
| `/crumbIssuer/api/json` | Overall/Read + CSRF Access |
| `/job/*/api/json` | Job/Read |
| `/job/*/build` | Job/Build |

### How to Grant Permissions

#### Matrix-Based Security

1. Go to: Manage Jenkins > Security > Configure Global Security
2. Under "Authorization", select "Matrix-based security"
3. Find user "sonar" and check:
   - Overall/Read ✅
   - Job/Build ✅
   - Job/Read ✅
   - Job/Workspace ✅ (optional)

#### Role-Based Security

1. Go to: Manage Jenkins > Security > Manage and Assign Roles
2. Create or edit role with:
   - Overall/Read
   - Job/Build
   - Job/Read
3. Assign role to user "sonar"

## Testing Steps

### Step 1: Verify User Exists

```bash
# List all users
curl -u "admin:admin-token" \
  http://10.0.1.84:8080/asynchPeople/api/json
```

### Step 2: Test Basic API Access

```bash
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  http://10.0.1.84:8080/api/json
```

### Step 3: Test Crumb Issuer

```bash
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  http://10.0.1.84:8080/crumbIssuer/api/json
```

### Step 4: Test Job Access

```bash
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  http://10.0.1.84:8080/job/plb-sonarqube/api/json
```

### Step 5: Test Manual Job Trigger

```bash
# Get crumb
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

## Common Permission Issues

### Issue 1: Missing Overall/Read

**Symptoms:**
- Can access jobs but not crumb issuer
- Gets 403 on `/crumbIssuer/api/json`

**Fix:** Grant Overall/Read permission

### Issue 2: User Doesn't Exist

**Symptoms:**
- Authentication fails completely
- 401 Unauthorized

**Fix:** Create user in Jenkins

### Issue 3: Wrong Security Realm

**Symptoms:**
- Authentication works for some users but not others
- Inconsistent behavior

**Fix:** Check security realm configuration

### Issue 4: Token Expired

**Symptoms:**
- Worked before, now fails
- Intermittent failures

**Fix:** Regenerate API token

## Jenkins Security Configuration

### Check Current Configuration

```bash
# Get security configuration
curl -u "admin:admin-token" \
  http://10.0.1.84:8080/securityRealm/api/json

# Get authorization strategy
curl -u "admin:admin-token" \
  http://10.0.1.84:8080/authorizationStrategy/api/json
```

### Enable Proper Permissions

1. **Log in as admin**
2. **Go to**: Manage Jenkins > Security > Configure Global Security
3. **Under "Authorization"**:
   - Select "Matrix-based security" or "Role-Based Strategy"
   - Grant appropriate permissions to "sonar" user
4. **Save** configuration

## Solution Summary

1. **Grant "sonar" user Overall/Read permission**
2. **Verify user exists in Jenkins**
3. **Test crumb issuer manually**
4. **Check security realm configuration**
5. **Use admin user for testing**

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

## Additional Resources

- **[CSRF Issues Guide](CSRF_ISSUES.md)** - CSRF-specific troubleshooting
- **[Connection Issues Guide](CONNECTION_ISSUES.md)** - Connection troubleshooting
- **[Debugging Guide](DEBUGGING_GUIDE.md)** - Advanced debugging techniques