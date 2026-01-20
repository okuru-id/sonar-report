# Solution Summary

## Your Issue

You're getting this error:
```json
{
  "error": "Cannot connect to Jenkins server at http://10.0.1.84:8080. Please check: 1) Jenkins is running, 2) Network connectivity, 3) Correct URL and port"
}
```

## Root Cause

The application server cannot connect to your Jenkins server at `http://10.0.1.84:8080`.

## Immediate Actions

### 1. Check Jenkins Server

On the Jenkins server (10.0.1.84):
```bash
# Check if Jenkins is running
systemctl status jenkins

# If not running, start it
systemctl start jenkins

# Check listening ports
ss -tuln | grep 8080
```

### 2. Test Network Connectivity

On your application server:
```bash
# Test basic connectivity
ping 10.0.1.84

# Test port connectivity
nc -zv 10.0.1.84 8080

# Test HTTP access
curl -v http://10.0.1.84:8080
```

### 3. Verify Configuration

```bash
# Check Jenkins URL and port
cat /etc/default/jenkins | grep JENKINS_PORT

# Check Jenkins configuration
cat /var/lib/jenkins/config.xml | grep "jenkinsUrl"
```

## Quick Test Script

Run this on your application server:

```bash
#!/bin/bash

echo "Quick Jenkins Connection Test"
echo "============================"

JENKINS_URL="http://10.0.1.84:8080"

echo "Testing connection to: $JENKINS_URL"
echo ""

# Test 1: Network
echo "1. Network connectivity..."
if ping -c 1 10.0.1.84 &> /dev/null; then
    echo "   ✓ Network OK"
else
    echo "   ✗ Network FAILED"
    echo "   → Check network connection to 10.0.1.84"
    exit 1
fi

# Test 2: Port
echo "2. Port 8080..."
if nc -zv 10.0.1.84 8080 &> /dev/null; then
    echo "   ✓ Port OK"
else
    echo "   ✗ Port FAILED"
    echo "   → Check if Jenkins is running on port 8080"
    exit 1
fi

# Test 3: HTTP
echo "3. HTTP access..."
if curl -s -I "$JENKINS_URL" | head -1 | grep -q "200\|302"; then
    echo "   ✓ HTTP OK"
else
    echo "   ✗ HTTP FAILED"
    echo "   → Check Jenkins configuration and logs"
    exit 1
fi

echo ""
echo "✓ All tests passed!"
echo ""
echo "Now test the API:"
echo "curl -X POST -H 'X-Jenkins-Url: $JENKINS_URL' -H 'X-Jenkins-Token: YOUR_TOKEN' -H 'X-Jenkins-Job: plb-sonarqube' -H 'X-Jenkins-User: sonar' http://localhost:8080/api/v1/jenkins/trigger"
```

## Common Solutions

### Solution 1: Start Jenkins

```bash
# On Jenkins server
systemctl start jenkins
systemctl enable jenkins
systemctl status jenkins
```

### Solution 2: Fix Network

```bash
# Check firewall
ufw status
ufw allow 8080

# Check routing
ip route
route -n
```

### Solution 3: Verify URL

```bash
# Try different URL variations
http://10.0.1.84:8080
http://10.0.1.84:8081  # If using custom port
http://jenkins-server:8080  # If using hostname
```

## Expected Success

When everything is working:

```bash
curl -X POST \
  -H "X-Jenkins-Url: http://10.0.1.84:8080" \
  -H "X-Jenkins-Token: 11e408c6ac99dd1bc79ae0dfc4cd7b3f10" \
  -H "X-Jenkins-Job: plb-sonarqube" \
  -H "X-Jenkins-User: sonar" \
  http://localhost:8080/api/v1/jenkins/trigger
```

Response:
```json
{
  "success": true,
  "message": "Jenkins job triggered successfully",
  "job": "plb-sonarqube"
}
```

## Next Steps

1. **Verify Jenkins is running** on 10.0.1.84:8080
2. **Test network connectivity** from application server
3. **Check firewall settings** on both servers
4. **Test with simple curl** commands first
5. **Try the API again** once basic connectivity works

## Need More Help?

Check these files:
- `CONNECTION_ISSUES.md` - Detailed connection troubleshooting
- `TEST_CONNECTION.md` - Step-by-step testing guide
- `QUICK_FIX.md` - Quick solutions for common issues
- `DEBUGGING_GUIDE.md` - Advanced debugging techniques