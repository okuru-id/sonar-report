# Debugging Guide

## Step-by-Step Debugging

### Step 1: Identify the Error

Look at the error message and categorize it:

- **Connection errors**: "Cannot connect to Jenkins server"
- **Authentication errors**: "authentication failed"
- **Permission errors**: "403 Forbidden"
- **Configuration errors**: "Missing headers", "Job not found"

### Step 2: Gather Information

Collect these details:

1. Exact error message
2. Jenkins URL and port
3. User credentials (masked)
4. Job name
5. Network topology
6. Jenkins version

### Step 3: Test Each Component

#### Test 1: Network Connectivity

```bash
# From application server
ping 10.0.1.84
telnet 10.0.1.84 8080
nc -zv 10.0.1.84 8080
```

#### Test 2: HTTP Access

```bash
curl -v http://10.0.1.84:8080
curl -I http://10.0.1.84:8080
```

#### Test 3: Authentication

```bash
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" http://10.0.1.84:8080/api/json
```

#### Test 4: Crumb Issuer

```bash
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" http://10.0.1.84:8080/crumbIssuer/api/json
```

#### Test 5: Job Access

```bash
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" http://10.0.1.84:8080/job/plb-sonarqube/api/json
```

### Step 4: Analyze Results

| Test | Expected | Actual | Issue |
|------|----------|--------|-------|
| Network | Connection succeeded | ? | Network/firewall |
| HTTP | HTTP 200/302 | ? | Jenkins not running |
| Auth | JSON response | ? | Credentials/permissions |
| Crumb | JSON with crumb | ? | CSRF/permissions |
| Job | JSON with job | ? | Job exists/permissions |

### Step 5: Fix the Issue

Based on which test failed:

- **Network test failed**: Check firewall, routing, DNS
- **HTTP test failed**: Start Jenkins, check port
- **Auth test failed**: Fix credentials, check permissions
- **Crumb test failed**: Enable CSRF, check permissions
- **Job test failed**: Verify job name, check permissions

## Debugging Tools

### Network Tools

```bash
# Check routing
route -n
ip route

# Check DNS
dig 10.0.1.84
nslookup 10.0.1.84

# Check connections
ss -tuln
netstat -tuln

# Check firewall
ufw status
iptables -L -n
```

### Jenkins Tools

```bash
# Check Jenkins status
systemctl status jenkins

# Check Jenkins logs
journalctl -u jenkins -f
tail -f /var/log/jenkins/jenkins.log

# Check Jenkins configuration
cat /etc/default/jenkins
cat /var/lib/jenkins/config.xml
```

### API Tools

```bash
# Test API with verbose output
curl -v -X POST \
  -H "X-Jenkins-Url: http://10.0.1.84:8080" \
  -H "X-Jenkins-Token: 11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  -H "X-Jenkins-Job: plb-sonarqube" \
  -H "X-Jenkins-User: sonar" \
  http://localhost:8080/api/v1/jenkins/trigger

# Test with different headers
curl -X POST \
  -H "X-Jenkins-Url: http://10.0.1.84:8080" \
  -H "X-Jenkins-Token: 11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  -H "X-Jenkins-Job: plb-sonarqube" \
  -H "X-Jenkins-User: sonar" \
  -v http://localhost:8080/api/v1/jenkins/trigger
```

## Common Debugging Scenarios

### Scenario 1: Connection Refused

**Symptoms:**
- `Connection refused` error
- No response from Jenkins

**Debugging Steps:**
1. Check if Jenkins is running on target server
2. Verify Jenkins port (default: 8080)
3. Test from Jenkins server itself
4. Check firewall on Jenkins server
5. Test network connectivity

**Commands:**
```bash
# On Jenkins server
systemctl status jenkins
ss -tuln | grep 8080

# On application server
ping 10.0.1.84
nc -zv 10.0.1.84 8080
```

### Scenario 2: Authentication Failed

**Symptoms:**
- `401 Unauthorized` or authentication error
- Credentials work in browser but not API

**Debugging Steps:**
1. Verify token is correct
2. Test with basic auth
3. Check user permissions
4. Test from same machine
5. Regenerate token

**Commands:**
```bash
# Test authentication
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" http://10.0.1.84:8080/api/json

# Check permissions in Jenkins UI
# Manage Jenkins > Security > Configure Global Security
```

### Scenario 3: Permission Denied

**Symptoms:**
- `403 Forbidden` error
- Authentication works but actions fail

**Debugging Steps:**
1. Check CSRF settings
2. Verify user permissions
3. Test crumb endpoint
4. Check job-specific permissions
5. Test with admin user

**Commands:**
```bash
# Test crumb
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" http://10.0.1.84:8080/crumbIssuer/api/json

# Test with admin
curl -u "admin:admin-token" http://10.0.1.84:8080/crumbIssuer/api/json
```

## Advanced Debugging

### Enable Debug Logging

Add to your application startup:

```bash
# Set environment variable
export GIN_MODE=debug

# Or in code
gin.SetMode(gin.DebugMode)
```

### Jenkins Debug Logging

```bash
# Enable debug logging in Jenkins
java -Djava.util.logging.config.file=/path/to/logging.properties -jar jenkins.war

# logging.properties
.handlers=java.util.logging.ConsoleHandler
.level=FINEST
java.util.logging.ConsoleHandler.level=FINEST
```

### Network Packet Capture

```bash
# Capture network traffic
tcpdump -i eth0 -w jenkins.pcap port 8080

# Analyze with Wireshark
wireshark jenkins.pcap
```

## Debugging Checklist

### Connection Issues

- [ ] Jenkins server is running
- [ ] Network connectivity is working
- [ ] Firewall allows the connection
- [ ] Port is correct and open
- [ ] URL is accessible from application server

### Authentication Issues

- [ ] Credentials are correct
- [ ] Token is not expired
- [ ] User exists in Jenkins
- [ ] User has Overall/Read permission
- [ ] Basic auth is working

### Permission Issues

- [ ] CSRF protection is enabled
- [ ] User can access crumb issuer
- [ ] User has Job/Build permission
- [ ] Job exists and is accessible
- [ ] No security plugins blocking access

### Configuration Issues

- [ ] Headers are correctly formatted
- [ ] Job name is correct
- [ ] URL includes proper path
- [ ] No typos in configuration
- [ ] Environment variables are set

## Final Verification

When debugging is complete, verify with:

```bash
# Complete test
curl -v -X POST \
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