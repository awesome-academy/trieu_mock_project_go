/**
 * Admin User Service
 */
const AdminUserService = {
  /**
   * Search users with pagination and filters
   * @param {Object} params - { limit, offset, team_id }
   * @returns {Promise}
   */
  searchUsers: function (params) {
    let url = `/admin/users/partial/search?limit=${params.limit || 10}&offset=${
      params.offset || 0
    }`;
    if (params.team_id) {
      url += `&team_id=${params.team_id}`;
    }
    return AdminAPI.get(url, { dataType: "html" });
  },

  /**
   * Create a new user
   * @param {Object} data
   * @returns {Promise}
   */
  createUser: function (data) {
    return AdminAPI.post("/admin/users", data);
  },

  /**
   * Update an existing user
   * @param {number|string} userId
   * @param {Object} data
   * @returns {Promise}
   */
  updateUser: function (userId, data) {
    return AdminAPI.put(`/admin/users/${userId}`, data);
  },

  /**
   * Delete a user
   * @param {number|string} userId
   * @returns {Promise}
   */
  deleteUser: function (userId) {
    return AdminAPI.delete(`/admin/users/${userId}`);
  },
};
