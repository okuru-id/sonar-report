#!/bin/bash

# Script untuk menguji API secara manual
# Pastikan server sudah berjalan sebelum menjalankan script ini

echo "Jenkins Trigger API Manual Test"
echo "=============================="
echo ""

# Konfigurasi
SERVER_URL="http://localhost:8080"

# Test 1: Missing headers
echo "Test 1: Missing headers"
response=$(curl -s -X POST "$SERVER_URL/api/v1/jenkins/trigger")
echo "Response: $response"
echo ""

# Test 2: Invalid headers (empty values)
echo "Test 2: Empty header values"
response=$(curl -s -X POST \
    -H "X-Jenkins-Url: " \
    -H "X-Jenkins-Token: " \
    -H "X-Jenkins-Job: " \
    "$SERVER_URL/api/v1/jenkins/trigger")
echo "Response: $response"
echo ""

# Test 3: Valid request (will fail because Jenkins URL is not real, but should return proper error)
echo "Test 3: Valid request structure (with fake Jenkins URL)"
response=$(curl -s -X POST \
    -H "X-Jenkins-Url: https://fake-jenkins.example.com" \
    -H "X-Jenkins-Token: fake-token" \
    -H "X-Jenkins-Job: fake-job" \
    -H "X-Jenkins-User: sonar" \
    "$SERVER_URL/api/v1/jenkins/trigger")
echo "Response: $response"
echo ""

echo "Manual testing completed."
echo ""
echo "To test with real Jenkins:"
echo "1. Make sure server is running: ./sonarqube-report"
echo "2. Use the jenkins_trigger_example.sh script with your real Jenkins credentials"