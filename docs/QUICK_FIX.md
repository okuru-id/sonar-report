# Quick Fix Guide

## Error: "Cannot connect to Jenkins server"

### Your Error
```json
{
  "error": "Cannot connect to Jenkins server at http://10.0.1.84:8080. Please check: 1) Jenkins is running, 2) Network connectivity, 3) Correct URL and port"
}
```

### Quick Checks

#### 1. Is Jenkins Running?
On the Jenkins server (10.0.1.84):
```bash
systemctl status jenkins
```

If not running:
```bash
systemctl start jenkins
```

#### 2. Can You Connect?
On your application server:
```bash
ping 10.0.1.84
nc -zv 10.0.1.84 8080
```

#### 3. Is Jenkins Accessible?
```bash
curl http://10.0.1.84:8080
```

### Common Fixes

#### Fix 1: Start Jenkins
```bash
# On Jenkins server
systemctl start jenkins
systemctl enable jenkins
```

#### Fix 2: Check Network
```bash
# On application server
ping 10.0.1.84

# If ping fails, check:
# - Network cables
# - Firewall settings
# - IP address configuration
```

#### Fix 3: Verify URL
```bash
# Make sure URL is correct
# Try different variations:
http://10.0.1.84:8080
http://10.0.1.84:8081  # If using custom port
http://jenkins-server:8080  # If using hostname
```

### Test Script

Run this on your application server:

```bash
#!/bin/bash

JENKINS_URL="http://10.0.1.84:8080"

echo "Testing connection to Jenkins..."

# Test 1: Ping
echo "1. Testing network connectivity..."
if ping -c 1 10.0.1.84 &> /dev/null; then
    echo "   ✓ Network OK"
else
    echo "   ✗ Network FAILED - Check network connection"
    exit 1
fi

# Test 2: Port
echo "2. Testing port 8080..."
if nc -zv 10.0.1.84 8080 &> /dev/null; then
    echo "   ✓ Port OK"
else
    echo "   ✗ Port FAILED - Check if Jenkins is running"
    exit 1
fi

# Test 3: HTTP
echo "3. Testing HTTP access..."
if curl -s -I "$JENKINS_URL" | head -1 | grep -q "200\|302"; then
    echo "   ✓ HTTP OK"
else
    echo "   ✗ HTTP FAILED - Check Jenkins configuration"
    exit 1
fi

echo ""
echo "All tests passed! Connection is working."
```

### If All Else Fails

1. **Check Jenkins logs**: `/var/log/jenkins/jenkins.log`
2. **Verify firewall**: `ufw status` or `iptables -L`
3. **Test from different machine**
4. **Check Jenkins port**: `netstat -tuln | grep 8080`
5. **Restart Jenkins**: `systemctl restart jenkins`

### Working Example

When everything is working:

```bash
# Test connection
curl -u "sonar:11e408c6ac99dd1bc79ae0dfc4cd7b3f10" http://10.0.1.84:8080/api/json

# Should return Jenkins API JSON response

# Then try the API
curl -X POST \
  -H "X-Jenkins-Url: http://10.0.1.84:8080" \
  -H "X-Jenkins-Token: 11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  -H "X-Jenkins-Job: plb-sonarqube" \
  -H "X-Jenkins-User: sonar" \
  http://localhost:8080/api/v1/jenkins/trigger
```

Expected success response:
```json
{
  "success": true,
  "message": "Jenkins job triggered successfully",
  "job": "plb-sonarqube"
}
```