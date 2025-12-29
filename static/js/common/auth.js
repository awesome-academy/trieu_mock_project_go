/**
 * Common authentication check utility
 */

/**
 * Check if user is authenticated
 * Redirects to login if no token found
 * @returns {boolean} true if user is authenticated, false otherwise
 */
function checkAuth() {
  const token = localStorage.getItem("accessToken");

  if (!token) {
    // No token found, redirect to login
    window.location.href = "/login";
    return false;
  }

  // Mock logic: verify user is authenticated
  // In a real application, you might validate the token with the server
  return verifyUser();
}

/**
 * Verify user is valid (mock implementation)
 * @returns {boolean} true if user is valid
 */
function verifyUser() {
  // Mock logic - return true for now
  // In a real application, you could:
  // - Validate JWT expiry
  // - Call API to verify token
  // - Check user data in localStorage

  const userEmail = localStorage.getItem("userEmail");
  const userId = localStorage.getItem("userId");

  if (!userEmail || !userId) {
    return false;
  }

  return true;
}

/**
 * Get the current access token
 * @returns {string|null} the access token or null if not found
 */
function getAccessToken() {
  return localStorage.getItem("accessToken");
}

/**
 * Get user info from localStorage
 * @returns {Object} user object with id, email, and name
 */
function getUserInfo() {
  return {
    id: localStorage.getItem("userId"),
    email: localStorage.getItem("userEmail"),
    name: localStorage.getItem("userName"),
  };
}

/**
 * Logout user
 * Clear all user data from localStorage and redirect to login
 */
function logout() {
  localStorage.removeItem("accessToken");
  localStorage.removeItem("userEmail");
  localStorage.removeItem("userName");
  localStorage.removeItem("userId");
  window.location.href = "/login";
}

/**
 * Initialize auth check on page load
 * Call this in jQuery ready handler
 */
function initAuthCheck() {
  $(document).ready(function () {
    checkAuth();
  });
}

initAuthCheck();
