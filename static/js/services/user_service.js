/**
 * User Service
 */
const UserService = {
  /**
   * Get current user profile
   * @returns {Promise}
   */
  getProfile: function () {
    return API.get("/api/profile");
  },

  /**
   * Update user profile (placeholder for future)
   * @param {Object} data
   * @returns {Promise}
   */
  updateProfile: function (data) {
    return API.put("/api/profile", data);
  },
};
