#!/bin/bash

# Script to test Jenkins connection and troubleshoot 403 errors

echo "Jenkins Connection Test"
echo "======================"
echo ""

# Read configuration
read -p "Enter Jenkins URL (e.g., https://jenkins.example.com): " JENKINS_URL
read -p "Enter Jenkins Username: " JENKINS_USER
read -s -p "Enter Jenkins Token: " JENKINS_TOKEN
echo ""
read -p "Enter Jenkins Job Name: " JENKINS_JOB

echo ""
echo "Testing Jenkins connection..."
echo ""

# Test 1: Basic connectivity
echo "Test 1: Checking basic connectivity to Jenkins..."
if curl -s -I "$JENKINS_URL" > /dev/null; then
    echo "✓ Jenkins URL is accessible"
else
    echo "✗ Cannot connect to Jenkins URL"
    exit 1
fi
echo ""

# Test 2: Authentication
echo "Test 2: Testing authentication..."
auth_test=$(curl -s -u "$JENKINS_USER:$JENKINS_TOKEN" "$JENKINS_URL/api/json")
if echo "$auth_test" | grep -q "jenkins"; then
    echo "✓ Authentication successful"
else
    echo "✗ Authentication failed"
    echo "Response: $auth_test"
    exit 1
fi
echo ""

# Test 3: Crumb issuer
echo "Test 3: Testing crumb issuer access..."
crumb_test=$(curl -s -u "$JENKINS_USER:$JENKINS_TOKEN" "$JENKINS_URL/crumbIssuer/api/json")
if echo "$crumb_test" | grep -q "crumb"; then
    echo "✓ Crumb issuer accessible"
    echo "Crumb data: $crumb_test"
else
    echo "✗ Crumb issuer failed"
    echo "Response: $crumb_test"
    exit 1
fi
echo ""

# Test 4: Job access
echo "Test 4: Testing job access..."
job_test=$(curl -s -u "$JENKINS_USER:$JENKINS_TOKEN" "$JENKINS_URL/job/$JENKINS_JOB/api/json")
if echo "$job_test" | grep -q "$JENKINS_JOB"; then
    echo "✓ Job is accessible"
else
    echo "✗ Job access failed"
    echo "Response: $job_test"
    exit 1
fi
echo ""

# Test 5: Job trigger (without crumb first)
echo "Test 5: Testing job trigger (without CSRF protection)..."
trigger_test=$(curl -s -X POST -u "$JENKINS_USER:$JENKINS_TOKEN" "$JENKINS_URL/job/$JENKINS_JOB/build")
if [ $? -eq 0 ]; then
    echo "✓ Job trigger successful (or CSRF protection is disabled)"
    echo "Response: $trigger_test"
else
    echo "✗ Job trigger failed"
    echo "Response: $trigger_test"
fi
echo ""

echo "All tests completed!"
echo ""
echo "If all tests passed, your Jenkins configuration is correct."
echo "If any test failed, check the error messages for troubleshooting."