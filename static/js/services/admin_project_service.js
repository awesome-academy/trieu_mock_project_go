/**
 * Admin Project Service
 */
const AdminProjectService = {
  /**
   * Search projects with pagination
   * @param {Object} params - { limit, offset, team_id }
   * @returns {Promise}
   */
  listProjects: function (params) {
    let url = `/admin/projects/partial/search?limit=${
      params.limit || 10
    }&offset=${params.offset || 0}`;
    if (params.team_id) {
      url += `&team_id=${params.team_id}`;
    }
    return AdminAPI.get(url, { dataType: "html" });
  },

  /**
   * Create a new project
   * @param {Object} data
   * @returns {Promise}
   */
  createProject: function (data) {
    return AdminAPI.post("/admin/projects", data);
  },

  /**
   * Update an existing project
   * @param {number|string} id
   * @param {Object} data
   * @returns {Promise}
   */
  updateProject: function (id, data) {
    return AdminAPI.put(`/admin/projects/${id}`, data);
  },

  /**
   * Delete a project
   * @param {number|string} id
   * @returns {Promise}
   */
  deleteProject: function (id) {
    return AdminAPI.delete(`/admin/projects/${id}`);
  },

  /**
   * Get all members of a team
   * @param {number|string} teamId
   * @returns {Promise}
   */
  getTeamMembers: function (teamId) {
    return AdminAPI.get(`/admin/teams/${teamId}/members/all`);
  },
};
