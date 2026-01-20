# Testing Connection from Application Server

## Quick Connection Test

Run this command from your application server to test connectivity to Jenkins:

```bash
# Test basic connectivity
curl -v http://10.0.1.84:8080

# Test with authentication
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" http://10.0.1.84:8080/api/json
```

## Step-by-Step Testing

### 1. Test Basic Connectivity

```bash
# Test if the server can reach Jenkins
telnet 10.0.1.84 8080

# Or use nc (netcat)
nc -zv 10.0.1.84 8080

# Expected output: Connection succeeded
```

### 2. Test HTTP Access

```bash
# Test HTTP response
curl -I http://10.0.1.84:8080

# Expected: HTTP/1.1 200 OK or 302 Found
```

### 3. Test Authentication

```bash
# Test with your credentials
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" http://10.0.1.84:8080/api/json

# Expected: JSON response with Jenkins information
```

### 4. Test Job Access

```bash
# Test accessing the specific job
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" http://10.0.1.84:8080/job/plb-sonarqube/api/json

# Expected: JSON response with job information
```

### 5. Test Crumb Issuer

```bash
# Test CSRF crumb endpoint
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" http://10.0.1.84:8080/crumbIssuer/api/json

# Expected: JSON with crumb and crumbRequestField
```

## Common Issues and Fixes

### Issue: Connection Refused

**Error:** `Connection refused`

**Fix:**
```bash
# Check if Jenkins is running on the target server
systemctl status jenkins

# Start Jenkins if needed
systemctl start jenkins

# Check listening ports
ss -tuln | grep 8080
```

### Issue: Host Not Found

**Error:** `no such host` or `Name or service not known`

**Fix:**
```bash
# Check DNS resolution
nslookup 10.0.1.84

# Or test with ping
ping 10.0.1.84

# Check /etc/hosts
cat /etc/hosts
```

### Issue: Authentication Failed

**Error:** `401 Unauthorized`

**Fix:**
```bash
# Verify credentials
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" http://10.0.1.84:8080/api/json

# Check user permissions in Jenkins
# Ensure user has Overall/Read permission
```

### Issue: Permission Denied

**Error:** `403 Forbidden`

**Fix:**
```bash
# Check user permissions
# In Jenkins: Manage Jenkins > Security > Configure Global Security

# Ensure user has:
# - Overall/Read permission
# - Job/Build permission for the specific job
# - Access to crumb issuer
```

## Complete Test Script

```bash
#!/bin/bash

echo "Jenkins Connection Test"
echo "======================"
echo ""

JENKINS_URL="http://10.0.1.84:8080"
JENKINS_USER="sonar"
JENKINS_TOKEN="11e408c6ac99dd1bc79ae0dfc4cd7b3f10"
JENKINS_JOB="plb-sonarqube"

echo "Testing connection to: $JENKINS_URL"
echo ""

# Test 1: Basic connectivity
echo "Test 1: Basic connectivity..."
if curl -s -I "$JENKINS_URL" | head -1 | grep -q "200\|302"; then
    echo "✓ PASS: Server is accessible"
else
    echo "✗ FAIL: Cannot connect to server"
    exit 1
fi
echo ""

# Test 2: Authentication
echo "Test 2: Authentication..."
auth_test=$(curl -s -u "$JENKINS_USER:$JENKINS_TOKEN" "$JENKINS_URL/api/json")
if echo "$auth_test" | grep -q "jenkins"; then
    echo "✓ PASS: Authentication successful"
else
    echo "✗ FAIL: Authentication failed"
    echo "Response: $auth_test"
    exit 1
fi
echo ""

# Test 3: Crumb issuer
echo "Test 3: Crumb issuer..."
crumb_test=$(curl -s -u "$JENKINS_USER:$JENKINS_TOKEN" "$JENKINS_URL/crumbIssuer/api/json")
if echo "$crumb_test" | grep -q "crumb"; then
    echo "✓ PASS: Crumb issuer accessible"
else
    echo "✗ FAIL: Crumb issuer failed"
    echo "Response: $crumb_test"
    exit 1
fi
echo ""

# Test 4: Job access
echo "Test 4: Job access..."
job_test=$(curl -s -u "$JENKINS_USER:$JENKINS_TOKEN" "$JENKINS_URL/job/$JENKINS_JOB/api/json")
if echo "$job_test" | grep -q "$JENKINS_JOB"; then
    echo "✓ PASS: Job is accessible"
else
    echo "✗ FAIL: Job access failed"
    echo "Response: $job_test"
    exit 1
fi
echo ""

echo "All tests passed! Your configuration is correct."
echo ""
echo "You can now use the API:"
echo "curl -X POST -H 'X-Jenkins-Url: $JENKINS_URL' -H 'X-Jenkins-Token: $JENKINS_TOKEN' -H 'X-Jenkins-Job: $JENKINS_JOB' -H 'X-Jenkins-User: $JENKINS_USER' http://localhost:8080/api/v1/jenkins/trigger"
```

## Network Troubleshooting

### Check Routing

```bash
# Check route to Jenkins
route -n | grep 10.0.1.84

# Or use traceroute
mtr 10.0.1.84
```

### Check Firewall

```bash
# Check firewall rules
ufw status

# Check iptables
iptables -L -n

# Temporarily disable firewall (for testing)
ufw disable
```

### Check DNS

```bash
# Check DNS resolution
cat /etc/resolv.conf

# Test DNS lookup
dig 10.0.1.84

# Check hosts file
cat /etc/hosts
```

## Jenkins Configuration Check

On the Jenkins server, verify:

```bash
# Check Jenkins is running
systemctl status jenkins

# Check Jenkins port
cat /etc/default/jenkins | grep JENKINS_PORT

# Check Jenkins URL
cat /var/lib/jenkins/config.xml | grep "jenkinsUrl"
```

## Final Verification

Once all tests pass, try the API again:

```bash
curl -X POST \
  -H "X-Jenkins-Url: http://10.0.1.84:8080" \
  -H "X-Jenkins-Token: 11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  -H "X-Jenkins-Job: plb-sonarqube" \
  -H "X-Jenkins-User: sonar" \
  http://localhost:8080/api/v1/jenkins/trigger
```

Expected response:
```json
{
  "success": true,
  "message": "Jenkins job triggered successfully",
  "job": "plb-sonarqube"
}
```