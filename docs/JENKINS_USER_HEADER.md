# X-Jenkins-User Header Documentation

## Overview

The `X-Jenkins-User` header allows you to specify a custom username for Jenkins authentication. This provides flexibility when integrating with different Jenkins configurations.

## Default Behavior

If the `X-Jenkins-User` header is not provided, the API will use "sonar" as the default username for Jenkins authentication.

## Usage

### Without X-Jenkins-User (uses default "sonar")

```bash
curl -X POST \
  -H "X-Jenkins-Url: https://jenkins.example.com" \
  -H "X-Jenkins-Token: your-token" \
  -H "X-Jenkins-Job: your-job" \
  http://localhost:8080/api/v1/jenkins/trigger
```

### With X-Jenkins-User (custom username)

```bash
curl -X POST \
  -H "X-Jenkins-Url: https://jenkins.example.com" \
  -H "X-Jenkins-Token: your-token" \
  -H "X-Jenkins-Job: your-job" \
  -H "X-Jenkins-User: custom-username" \
  http://localhost:8080/api/v1/jenkins/trigger
```

## When to Use Custom Username

Use a custom username when:

1. **Different Jenkins Credentials**: Your Jenkins setup requires a specific username different from "sonar"
2. **User-Specific Jobs**: You need to trigger jobs with specific user permissions
3. **Audit Trails**: You want to track which user triggered builds in Jenkins logs
4. **Multiple Integrations**: You have multiple systems integrating with the same Jenkins instance

## Examples

### Example 1: Using Default Username

```bash
# This will use "sonar" as the username
curl -X POST \
  -H "X-Jenkins-Url: https://jenkins.company.com" \
  -H "X-Jenkins-Token: abc123xyz" \
  -H "X-Jenkins-Job: production-deploy" \
  http://api.example.com/api/v1/jenkins/trigger
```

### Example 2: Using Custom Username

```bash
# This will use "deploy-bot" as the username
curl -X POST \
  -H "X-Jenkins-Url: https://jenkins.company.com" \
  -H "X-Jenkins-Token: abc123xyz" \
  -H "X-Jenkins-Job: production-deploy" \
  -H "X-Jenkins-User: deploy-bot" \
  http://api.example.com/api/v1/jenkins/trigger
```

### Example 3: Different Users for Different Jobs

```bash
# Trigger development job with dev user
curl -X POST \
  -H "X-Jenkins-Url: https://jenkins.company.com" \
  -H "X-Jenkins-Token: dev-token" \
  -H "X-Jenkins-Job: dev-build" \
  -H "X-Jenkins-User: dev-bot" \
  http://api.example.com/api/v1/jenkins/trigger

# Trigger production job with prod user
curl -X POST \
  -H "X-Jenkins-Url: https://jenkins.company.com" \
  -H "X-Jenkins-Token: prod-token" \
  -H "X-Jenkins-Job: prod-deploy" \
  -H "X-Jenkins-User: prod-bot" \
  http://api.example.com/api/v1/jenkins/trigger
```

## Implementation Details

The API handles the username as follows:

1. Checks if `X-Jenkins-User` header is present
2. If present and not empty, uses the provided username
3. If not present or empty, defaults to "sonar"
4. Uses the determined username for Basic Authentication with Jenkins

## Security Considerations

1. **Token Security**: The token (not the username) is the primary security mechanism
2. **Username Visibility**: The username is visible in logs and network traffic
3. **Jenkins Permissions**: Ensure the username has appropriate permissions in Jenkins
4. **Avoid Hardcoding**: Don't hardcode usernames in client applications

## Best Practices

1. **Use Environment Variables**: Store usernames in environment variables
2. **Consistent Naming**: Use consistent username patterns (e.g., `app-name-bot`)
3. **Document**: Document which usernames are used for which integrations
4. **Monitor**: Monitor Jenkins logs to track which users are triggering builds