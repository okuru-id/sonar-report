# Troubleshooting HTTP 403 Errors

## Common Causes of 403 Errors

When you receive a `HTTP 403` error from the Jenkins Trigger API, it typically indicates authentication or permission issues. Here are the most common causes:

### 1. Invalid Jenkins Credentials

**Symptoms:**
- Error: `Jenkins authentication failed: Invalid credentials or insufficient permissions`
- Occurs when trying to get crumb or trigger job

**Solutions:**
- Verify your Jenkins token is correct
- Check if the token has expired
- Ensure the token belongs to the correct user
- Regenerate the token if needed

### 2. Insufficient User Permissions

**Symptoms:**
- Error: `Jenkins job trigger failed: Invalid credentials or insufficient permissions`
- Authentication succeeds but job trigger fails

**Solutions:**
- Check Jenkins user permissions in Jenkins security settings
- Ensure the user has "Build" permission for the specific job
- Verify the user has "Overall/Read" permission
- Check if the user has "Job/Build" permission

### 3. CSRF Protection Issues

**Symptoms:**
- Error occurs when getting crumb (`/crumbIssuer/api/json`)
- May see "CSRF" related messages in Jenkins logs

**Solutions:**
- Ensure CSRF protection is enabled in Jenkins (Manage Jenkins > Security > CSRF Protection)
- Verify the user has permission to access crumb issuer
- Check Jenkins system logs for CSRF-related errors

### 4. Incorrect Jenkins URL

**Symptoms:**
- 403 error with additional HTML content
- URL might be pointing to wrong Jenkins instance

**Solutions:**
- Verify the Jenkins URL is correct and accessible
- Check if URL includes correct path/port
- Test URL in browser first
- Ensure no typos in the URL

### 5. Jenkins Security Configuration

**Symptoms:**
- Consistent 403 errors across different endpoints
- May work with some users but not others

**Solutions:**
- Check Jenkins global security settings
- Verify "Prevent Cross Site Request Forgery exploits" is enabled
- Check "Authorization" strategy (e.g., Matrix-based, Role-based)
- Review security realm configuration

## Step-by-Step Troubleshooting

### Step 1: Verify Basic Connectivity

```bash
# Test if Jenkins URL is accessible
curl -I "https://your-jenkins-url.com"

# Should return HTTP 200 or 302
```

### Step 2: Test Authentication

```bash
# Test authentication with your credentials
curl -u "username:token" "https://your-jenkins-url.com/api/json"

# Should return Jenkins API response, not 403
```

### Step 3: Test Crumb Access

```bash
# Test crumb issuer endpoint
curl -u "username:token" "https://your-jenkins-url.com/crumbIssuer/api/json"

# Should return crumb JSON, not 403
```

### Step 4: Test Job Access

```bash
# Test if you can access the job API
curl -u "username:token" "https://your-jenkins-url.com/job/your-job/api/json"

# Should return job information, not 403
```

### Step 5: Test Job Trigger

```bash
# Test if you can manually trigger the job
curl -X POST -u "username:token" "https://your-jenkins-url.com/job/your-job/build"

# Should return 201 Created or 200 OK
```

## Jenkins Configuration Checklist

### 1. User Permissions

- [ ] User has "Overall/Read" permission
- [ ] User has "Job/Build" permission for the specific job
- [ ] User has "Job/Read" permission for the specific job
- [ ] User has "Job/Workspace" permission if needed

### 2. Token Configuration

- [ ] Token is generated for the correct user
- [ ] Token has not expired
- [ ] Token is not revoked
- [ ] Token has sufficient scope

### 3. CSRF Settings

- [ ] CSRF protection is enabled in Jenkins
- [ ] User has permission to access crumb issuer
- [ ] Crumb issuer endpoint is accessible

### 4. Network Configuration

- [ ] No firewall blocking access
- [ ] No IP restrictions in place
- [ ] Jenkins URL is correct and complete
- [ ] No proxy issues

## Common Solutions

### Solution 1: Regenerate Jenkins Token

1. Log in to Jenkins
2. Go to your user profile
3. Click "Configure"
4. Find "API Token" section
5. Click "Add new Token"
6. Give it a name and click "Generate"
7. Copy the new token and use it in your API calls

### Solution 2: Check User Permissions

1. Log in to Jenkins as admin
2. Go to "Manage Jenkins" > "Security" > "Configure Global Security"
3. Check authorization strategy
4. Ensure your user has proper permissions
5. For matrix-based security, add appropriate permissions

### Solution 3: Disable CSRF Temporarily (for testing only)

⚠️ **Warning:** Only do this for testing in development environments

1. Go to "Manage Jenkins" > "Security" > "Configure Global Security"
2. Find "Prevent Cross Site Request Forgery exploits"
3. Uncheck the option temporarily
4. Test your API call
5. Re-enable CSRF protection after testing

### Solution 4: Use Different Authentication Method

If basic auth with token doesn't work, try:

```bash
# Using API token directly
curl -u "username:api-token" "https://jenkins-url/..."

# Using password instead of token
curl -u "username:password" "https://jenkins-url/..."
```

## Debugging Tips

### 1. Check Jenkins Logs

- Go to "Manage Jenkins" > "System Log"
- Look for authentication-related errors
- Check for permission denied messages

### 2. Enable Debug Logging

Add these Java options to Jenkins startup:

```
-Djava.util.logging.config.file=/path/to/logging.properties
```

With logging.properties containing:

```
handlers=java.util.logging.ConsoleHandler
.level=FINEST
java.util.logging.ConsoleHandler.level=FINEST
```

### 3. Test with Different Users

Try with an admin user first, then gradually reduce permissions:

```bash
# Test with admin
curl -u "admin:admin-token" "https://jenkins-url/crumbIssuer/api/json"

# Test with regular user
curl -u "user:user-token" "https://jenkins-url/crumbIssuer/api/json"
```

### 4. Check Jenkins Version

Some Jenkins versions have different security behaviors:

```bash
curl -u "username:token" "https://jenkins-url/api/json?tree=version"
```

## Final Checks

If you've tried all the above and still get 403 errors:

1. **Check Jenkins plugins**: Some security plugins might interfere
2. **Check reverse proxy**: If using nginx/Apache, check authentication headers
3. **Check Jenkins URL**: Ensure it's the correct URL (no typos)
4. **Test from different network**: Rule out network/firewall issues
5. **Check Jenkins system logs**: For detailed error information