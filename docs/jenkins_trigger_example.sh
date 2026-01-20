#!/bin/bash

# Contoh penggunaan Jenkins Trigger API
# Script ini menunjukkan cara menggunakan API untuk memicu Jenkins job

# Pastikan Anda memiliki curl terinstal
if ! command -v curl &> /dev/null; then
    echo "Error: curl is not installed. Please install curl first."
    exit 1
fi

# Konfigurasi - ganti dengan nilai yang sesuai
SERVER_URL="http://localhost:8080"  # URL server SonarQube Report Generator
JENKINS_URL="https://jenkins.example.com"  # URL Jenkins Anda
JENKINS_TOKEN="your-jenkins-token-here"  # Token Jenkins Anda
JENKINS_JOB="plb-sonarqube"  # Nama job Jenkins yang ingin dipicu
JENKINS_USER="sonar"  # Username untuk autentikasi Jenkins (default: sonar)

# Fungsi untuk menampilkan bantuan
show_help() {
echo "Usage: $0 [options]"
echo ""
echo "Options:"
echo "  -s, --server URL     Server URL (default: $SERVER_URL)"
echo "  -j, --jenkins URL    Jenkins URL (default: $JENKINS_URL)"
echo "  -t, --token TOKEN    Jenkins token (default: $JENKINS_TOKEN)"
echo "  -b, --job JOB        Jenkins job name (default: $JENKINS_JOB)"
echo "  -u, --user USER      Jenkins username (default: $JENKINS_USER)"
echo "  -h, --help           Show this help message"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -s|--server)
            SERVER_URL="$2"
            shift 2
            ;;
        -j|--jenkins)
            JENKINS_URL="$2"
            shift 2
            ;;
        -t|--token)
            JENKINS_TOKEN="$2"
            shift 2
            ;;
        -b|--job)
            JENKINS_JOB="$2"
            shift 2
            ;;
        -u|--user)
            JENKINS_USER="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

echo "Jenkins Trigger API Example"
echo "=========================="
echo ""
echo "Configuration:"
echo "  Server URL:     $SERVER_URL"
echo "  Jenkins URL:    $JENKINS_URL"
echo "  Jenkins Job:    $JENKINS_JOB"
echo "  Jenkins User:   $JENKINS_USER"
echo ""

# Kirim request ke API
echo "Sending request to trigger Jenkins job..."

response=$(curl -s -X POST \
    -H "X-Jenkins-Url: $JENKINS_URL" \
    -H "X-Jenkins-Token: $JENKINS_TOKEN" \
    -H "X-Jenkins-Job: $JENKINS_JOB" \
    -H "X-Jenkins-User: $JENKINS_USER" \
    "$SERVER_URL/api/v1/jenkins/trigger")

echo ""
echo "Response:"
echo "$response | jq ."
echo ""

# Analisis response
if echo "$response" | grep -q "success"; then
    echo "✓ SUCCESS: Jenkins job triggered successfully!"
    exit 0
else
    echo "✗ ERROR: Failed to trigger Jenkins job"
    echo ""
    echo "Possible causes:"
    echo "  1. Invalid Jenkins URL, token, or job name"
    echo "  2. Jenkins server is not accessible"
    echo "  3. Authentication failed"
    echo "  4. Network connectivity issues"
    exit 1
fi