/**
 * Authentication Service
 */
const AuthService = {
  /**
   * Login user
   * @param {string} email
   * @param {string} password
   * @returns {Promise}
   */
  login: async function (email, password) {
    const response = await API.post("/login", {
      user: { email, password },
    });

    if (response && response.user && response.user.access_token) {
      this.setSession(response.user);
      return response.user;
    }
    throw new Error("Invalid response from server");
  },

  /**
   * Logout user
   */
  logout: async function () {
    const token = this.getToken();
    if (token) {
      try {
        await API.post("/api/logout");
      } catch (error) {
        console.error("Failed to invalidate session on server:", error);
      }
    }

    localStorage.removeItem("accessToken");
    localStorage.removeItem("userEmail");
    localStorage.removeItem("userName");
    localStorage.removeItem("userId");
    window.location.href = "/login";
  },

  /**
   * Set session data in localStorage
   * @param {Object} user
   */
  setSession: function (user) {
    localStorage.setItem("accessToken", user.access_token);
    localStorage.setItem("userEmail", user.email);
    localStorage.setItem("userName", user.name);
    localStorage.setItem("userId", user.id);
  },

  /**
   * Get current access token
   * @returns {string|null}
   */
  getToken: function () {
    return localStorage.getItem("accessToken");
  },

  /**
   * Get current user info
   * @returns {Object|null}
   */
  getUser: function () {
    const id = localStorage.getItem("userId");
    if (!id) return null;

    return {
      id: id,
      email: localStorage.getItem("userEmail"),
      name: localStorage.getItem("userName"),
    };
  },

  /**
   * Check if user is authenticated
   * @returns {boolean}
   */
  isAuthenticated: function () {
    const token = this.getToken();
    const user = this.getUser();
    return !!(token && user && user.email && user.id);
  },

  /**
   * Initialize auth check for protected pages
   */
  checkAuthAndRedirect: function () {
    if (!this.isAuthenticated()) {
      this.logout();
      return false;
    }
    return true;
  },
};
