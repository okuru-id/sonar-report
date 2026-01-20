# Connection Issues Troubleshooting

## Cannot Connect to Jenkins Server

If you receive an error like:
```json
{
  "error": "Cannot connect to Jenkins server at http://10.0.1.84:8080. Please check: 1) Jenkins is running, 2) Network connectivity, 3) Correct URL and port"
}
```

This indicates the application cannot establish a connection to your Jenkins server.

## Common Causes

### 1. Jenkins Not Running

**Symptoms:**
- Connection refused errors
- No response from Jenkins URL

**Solutions:**
- Check if Jenkins service is running
- Start Jenkins if it's stopped
- Verify Jenkins port is correct

```bash
# Check if Jenkins is running (on Jenkins server)
systemctl status jenkins

# Start Jenkins if stopped
systemctl start jenkins

# Check Jenkins logs
journalctl -u jenkins -f
```

### 2. Network Connectivity Issues

**Symptoms:**
- Host not found errors
- Connection timeout
- Network unreachable

**Solutions:**
- Verify network connection between servers
- Check firewall settings
- Test connectivity with ping/curl

```bash
# Test basic connectivity
ping 10.0.1.84

# Test Jenkins port
telnet 10.0.1.84 8080

# Test from application server
curl -v http://10.0.1.84:8080
```

### 3. Incorrect URL or Port

**Symptoms:**
- Connection refused
- Invalid URL format

**Solutions:**
- Verify Jenkins URL is correct
- Check if Jenkins uses standard port (8080) or custom port
- Ensure URL includes http:// or https://

```bash
# Check Jenkins configuration
cat /etc/default/jenkins | grep HTTP_PORT

# Or check Jenkins system configuration
grep "jenkins.model.JenkinsLocationConfiguration" /var/lib/jenkins/config.xml
```

### 4. Firewall Blocking Connection

**Symptoms:**
- Connection timeout
- No response
- Silent failure

**Solutions:**
- Check firewall rules on both servers
- Ensure Jenkins port is open
- Test with firewall temporarily disabled (for testing only)

```bash
# Check firewall status
ufw status

# Allow Jenkins port
ufw allow 8080

# Check iptables
iptables -L -n
```

## Step-by-Step Troubleshooting

### Step 1: Verify Jenkins is Running

On the Jenkins server:
```bash
# Check service status
systemctl status jenkins

# Check listening ports
netstat -tuln | grep 8080

# Check Jenkins process
ps aux | grep jenkins
```

### Step 2: Test Local Access

On the Jenkins server:
```bash
# Test local access
curl http://localhost:8080

# Test with authentication
curl -u username:token http://localhost:8080/api/json
```

### Step 3: Test Network Connectivity

On the application server:
```bash
# Test basic connectivity
ping 10.0.1.84

# Test port connectivity
nc -zv 10.0.1.84 8080

# Test HTTP access
curl -v http://10.0.1.84:8080
```

### Step 4: Check DNS Resolution

```bash
# Test DNS resolution
nslookup 10.0.1.84

# Or for hostname
dig jenkins.example.com

# Check /etc/hosts
cat /etc/hosts
```

### Step 5: Verify Jenkins Configuration

```bash
# Check Jenkins URL configuration
cat /var/lib/jenkins/config.xml | grep "jenkinsUrl"

# Check Jenkins port
cat /etc/default/jenkins | grep JENKINS_PORT
```

## Common Solutions

### Solution 1: Start Jenkins Service

```bash
# Start Jenkins
systemctl start jenkins

# Enable auto-start
systemctl enable jenkins

# Check status
systemctl status jenkins
```

### Solution 2: Fix Network Configuration

```bash
# Check network interface
ip a

# Check routing
ip route

# Check DNS
cat /etc/resolv.conf
```

### Solution 3: Configure Firewall

```bash
# Allow Jenkins port
ufw allow 8080/tcp

# Reload firewall
ufw reload

# Check status
ufw status
```

### Solution 4: Use Correct URL

```bash
# If Jenkins uses custom port
JENKINS_URL="http://10.0.1.84:8081"

# If using HTTPS
JENKINS_URL="https://jenkins.example.com"

# If using context path
JENKINS_URL="http://10.0.1.84:8080/jenkins"
```

## Testing from Application Server

Create a simple test script:

```bash
#!/bin/bash

# Test Jenkins connection
JENKINS_URL="http://10.0.1.84:8080"

# Test 1: Basic connectivity
echo "Testing basic connectivity..."
if curl -s -I "$JENKINS_URL" | head -1 | grep -q "200\|302"; then
    echo "✓ Jenkins is accessible"
else
    echo "✗ Cannot connect to Jenkins"
    exit 1
fi

# Test 2: API access
echo "Testing API access..."
if curl -s "$JENKINS_URL/api/json" | grep -q "jenkins"; then
    echo "✓ Jenkins API is accessible"
else
    echo "✗ Jenkins API not accessible"
    exit 1
fi

echo "All tests passed!"
```

## Network Diagnostics

### Check Routing

```bash
# Check route to Jenkins server
route -n | grep 10.0.1.84

# Or use traceroute
mtr 10.0.1.84
```

### Check Port Availability

```bash
# Check if port is listening
netstat -tuln | grep 8080

# Or use ss
ss -tuln | grep 8080
```

### Check Proxy Settings

```bash
# Check environment variables
env | grep -i proxy

# Check system proxy
cat /etc/environment | grep -i proxy
```

## Final Checks

If you've tried all the above and still cannot connect:

1. **Check Jenkins logs** for startup errors
2. **Verify network cables** and physical connections
3. **Test from different machine** to isolate the issue
4. **Check VPN or VLAN** configurations
5. **Contact network administrator** for assistance

## Example Working Configuration

```bash
# Working example
JENKINS_URL="http://10.0.1.84:8080"
JENKINS_USER="sonar"
JENKINS_TOKEN="11e408c6ac99dd1bc79ae0dfc4cd7b3f10"
JENKINS_JOB="plb-sonarqube"

# Test command
curl -X POST \
  -H "X-Jenkins-Url: $JENKINS_URL" \
  -H "X-Jenkins-Token: $JENKINS_TOKEN" \
  -H "X-Jenkins-Job: $JENKINS_JOB" \
  -H "X-Jenkins-User: $JENKINS_USER" \
  http://localhost:8080/api/v1/jenkins/trigger
```