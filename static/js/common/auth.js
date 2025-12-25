/**
 * Common authentication check utility
 * This file acts as a bridge between the legacy global functions and the new AuthService
 */

/**
 * Check if user is authenticated
 * Redirects to login if no token found
 * @returns {boolean} true if user is authenticated, false otherwise
 */
function checkAuth() {
  return AuthService.checkAuthAndRedirect();
}

/**
 * Verify user is valid
 * @returns {boolean} true if user is valid
 */
function verifyUser() {
  return AuthService.isAuthenticated();
}

/**
 * Get the current access token
 * @returns {string|null} the access token or null if not found
 */
function getAccessToken() {
  return AuthService.getToken();
}

/**
 * Get user info from localStorage
 * @returns {Object} user object with id, email, and name
 */
function getUserInfo() {
  return AuthService.getUser();
}

/**
 * Logout user
 * Clear all user data from localStorage and redirect to login
 */
function logout() {
  AuthService.logout();
}

/**
 * Initialize auth check on page load
 * Call this in jQuery ready handler
 */
function initAuthCheck() {
  $(document).ready(function () {
    // Only run checkAuth if we are not on the login page
    const path = window.location.pathname;
    if (path !== "/login" && path !== "/admin/login") {
      checkAuth();
    }
  });
}

// Initialize auth check
initAuthCheck();
