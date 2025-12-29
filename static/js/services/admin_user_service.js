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
    return fetch(url).then((res) => {
      if (!res.ok) throw new Error("Failed to fetch users");
      return res.text();
    });
  },

  /**
   * Create a new user
   * @param {Object} data
   * @returns {Promise}
   */
  createUser: function (data) {
    return fetch("/admin/users", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    }).then(async (res) => {
      if (!res.ok) {
        const err = await res.json();
        return Promise.reject(err);
      }
      return res.json();
    });
  },

  /**
   * Update an existing user
   * @param {number|string} userId
   * @param {Object} data
   * @returns {Promise}
   */
  updateUser: function (userId, data) {
    return fetch(`/admin/users/${userId}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    }).then(async (res) => {
      if (!res.ok) {
        const err = await res.json();
        return Promise.reject(err);
      }
      return res.json();
    });
  },

  /**
   * Delete a user
   * @param {number|string} userId
   * @returns {Promise}
   */
  deleteUser: function (userId) {
    return fetch(`/admin/users/${userId}`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
    }).then(async (res) => {
      if (!res.ok) {
        const err = await res.json();
        return Promise.reject(err);
      }
      return res.json();
    });
  },
};
