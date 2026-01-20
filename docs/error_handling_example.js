/**
 * Example: Jenkins Trigger API Error Handling
 * 
 * This example demonstrates how to properly handle errors when using
 * the Jenkins Trigger API in a Node.js application.
 */

const axios = require('axios');

class JenkinsTrigger {
  constructor(baseUrl) {
    this.baseUrl = baseUrl;
  }

  /**
   * Trigger a Jenkins job with proper error handling
   * 
   * @param {string} jenkinsUrl - Jenkins server URL
   * @param {string} jenkinsToken - Jenkins API token
   * @param {string} jenkinsJob - Jenkins job name
   * @param {string} [jenkinsUser='sonar'] - Jenkins username
   * @returns {Promise<Object>} Result object with success status and message
   */
  async triggerJob(jenkinsUrl, jenkinsToken, jenkinsJob, jenkinsUser = 'sonar') {
    try {
      const response = await axios.post(
        `${this.baseUrl}/api/v1/jenkins/trigger`,
        {},
        {
          headers: {
            'X-Jenkins-Url': jenkinsUrl,
            'X-Jenkins-Token': jenkinsToken,
            'X-Jenkins-Job': jenkinsJob,
            'X-Jenkins-User': jenkinsUser
          },
          timeout: 10000 // 10 second timeout
        }
      );

      return {
        success: true,
        message: 'Jenkins job triggered successfully',
        data: response.data
      };

    } catch (error) {
      return this.handleError(error);
    }
  }

  /**
   * Handle API errors and return user-friendly messages
   * 
   * @param {Error} error - The error object from axios
   * @returns {Object} Error result with details
   */
  handleError(error) {
    let errorType = 'unknown';
    let userMessage = 'Failed to trigger Jenkins job';
    let technicalDetails = error.message;
    let suggestedAction = 'Please try again or contact support';

    if (error.response) {
      // Server responded with error status
      const status = error.response.status;
      const errorMsg = error.response.data.error || 'Unknown error';

      // Classify error based on status code and message
      if (status === 400) {
        if (errorMsg.includes('Missing or invalid headers')) {
          errorType = 'missing_headers';
          userMessage = 'Required headers are missing or invalid';
          suggestedAction = 'Please provide all required headers: X-Jenkins-Url, X-Jenkins-Token, X-Jenkins-Job';
        } else if (errorMsg.includes('authentication failed')) {
          errorType = 'authentication_failed';
          userMessage = 'Jenkins authentication failed';
          suggestedAction = 'Please verify your Jenkins credentials and permissions';
        }
      } else if (status === 403) {
        errorType = 'permission_denied';
        userMessage = 'Permission denied';
        suggestedAction = 'Check your Jenkins user permissions and CSRF settings';
      } else if (status === 404) {
        errorType = 'not_found';
        userMessage = 'Jenkins job not found';
        suggestedAction = 'Please verify the Jenkins job name is correct';
      } else if (status >= 500) {
        errorType = 'server_error';
        userMessage = 'Server error occurred';
        suggestedAction = 'Please try again later or contact support';
      }

      technicalDetails = `HTTP ${status}: ${errorMsg}`;

    } else if (error.request) {
      // No response received
      errorType = 'connection_error';
      userMessage = 'Cannot connect to server';
      suggestedAction = 'Please check your network connection and server availability';
      technicalDetails = 'No response received from server';

    } else {
      // Other errors
      errorType = 'request_error';
      userMessage = 'Request failed';
      suggestedAction = 'Please check your request parameters';
    }

    return {
      success: false,
      errorType,
      userMessage,
      technicalDetails,
      suggestedAction,
      timestamp: new Date().toISOString()
    };
  }

  /**
   * Get user-friendly error message for display
   * 
   * @param {Object} errorResult - Result from handleError
   * @returns {string} User-friendly message
   */
  getUserMessage(errorResult) {
    if (errorResult.success) return '';
    
    return `${errorResult.userMessage}. ${errorResult.suggestedAction}`;
  }

  /**
   * Get technical error details for logging
   * 
   * @param {Object} errorResult - Result from handleError
   * @returns {string} Technical error details
   */
  getTechnicalDetails(errorResult) {
    if (errorResult.success) return '';
    
    return `[${errorResult.errorType}] ${errorResult.technicalDetails}`;
  }
}

// Example Usage
async function exampleUsage() {
  const jenkinsTrigger = new JenkinsTrigger('http://localhost:8080');

  // Example 1: Successful trigger
  try {
    const result = await jenkinsTrigger.triggerJob(
      'https://jenkins.example.com',
      'valid-token',
      'my-job'
    );
    
    if (result.success) {
      console.log('✓ Success:', result.message);
    } else {
      console.log('User message:', jenkinsTrigger.getUserMessage(result));
      console.log('Technical details:', jenkinsTrigger.getTechnicalDetails(result));
    }
  } catch (error) {
    console.error('Unexpected error:', error);
  }

  // Example 2: Handling authentication error
  try {
    const result = await jenkinsTrigger.triggerJob(
      'https://jenkins.example.com',
      'invalid-token',
      'my-job'
    );
    
    if (!result.success) {
      console.log('✗ Error handled gracefully:');
      console.log('  Type:', result.errorType);
      console.log('  User message:', result.userMessage);
      console.log('  Technical:', result.technicalDetails);
      console.log('  Action:', result.suggestedAction);
    }
  } catch (error) {
    console.error('Unexpected error:', error);
  }
}

// Run example
if (require.main === module) {
  exampleUsage();
}

module.exports = JenkinsTrigger;