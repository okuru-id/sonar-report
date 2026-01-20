#!/bin/bash

echo "Testing Jenkins connection directly..."
echo ""

# Konfigurasi dari user
JENKINS_URL="http://10.0.1.84:8080"
JENKINS_TOKEN="11e408c6ac99dd1bc79ae0dfc4cd7b3f10"
JENKINS_USER="sonar"
JENKINS_JOB="plb-sonarqube"

echo "Jenkins URL: $JENKINS_URL"
echo "Jenkins User: $JENKINS_USER"
echo "Jenkins Job: $JENKINS_JOB"
echo ""

# Test 1: Basic connectivity
echo "Test 1: Checking basic connectivity..."
if curl -s -I "$JENKINS_URL" | head -1 | grep -q "200\|302"; then
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
    echo ""
    echo "Possible issues:"
    echo "1. Invalid username or token"
    echo "2. User doesn't have Overall/Read permission"
    echo "3. CSRF protection blocking the request"
    exit 1
fi
echo ""

# Test 3: Crumb issuer
echo "Test 3: Testing crumb issuer..."
crumb_test=$(curl -s -u "$JENKINS_USER:$JENKINS_TOKEN" "$JENKINS_URL/crumbIssuer/api/json")
echo "Crumb response: $crumb_test"
if echo "$crumb_test" | grep -q "crumb"; then
    echo "✓ Crumb issuer accessible"
else
    echo "✗ Crumb issuer failed"
    echo ""
    echo "Possible issues:"
    echo "1. CSRF protection is disabled in Jenkins"
    echo "2. User doesn't have permission to access crumb issuer"
    echo "3. Jenkins version doesn't support crumb issuer"
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
    echo ""
    echo "Possible issues:"
    echo "1. Job doesn't exist"
    echo "2. User doesn't have permission to access the job"
    echo "3. Incorrect job name"
    exit 1
fi
echo ""

echo "All tests passed! The issue might be with:"
echo "1. CSRF configuration in Jenkins"
echo "2. Specific permissions for the job"
echo "3. Network configuration between servers"