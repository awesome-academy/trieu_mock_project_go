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
   * Get specific user profile
   * @param {number} userId
   * @returns {Promise}
   */
  getUserProfile: function (userId) {
    return API.get(`/api/profile/${userId}`);
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
