#!/bin/bash

# Script untuk menguji Jenkins Trigger API

# Konfigurasi
SERVER_URL="http://localhost:8080"
JENKINS_URL="https://jenkins.example.com"
JENKINS_TOKEN="your-jenkins-token"
JENKINS_JOB="plb-sonarqube"

# Fungsi untuk menguji API
test_jenkins_trigger() {
    echo "Testing Jenkins Trigger API..."
    
    # Kirim request ke API
    response=$(curl -s -X POST \
        -H "X-Jenkins-Url: $JENKINS_URL" \
        -H "X-Jenkins-Token: $JENKINS_TOKEN" \
        -H "X-Jenkins-Job: $JENKINS_JOB" \
        "$SERVER_URL/api/v1/jenkins/trigger")
    
    echo "Response: $response"
    
    # Periksa apakah response mengandung success
    if echo "$response" | grep -q "success"; then
        echo "✓ Jenkins Trigger API works successfully!"
        return 0
    else
        echo "✗ Jenkins Trigger API failed!"
        return 1
    fi
}

# Jalankan test
test_jenkins_trigger

exit $?