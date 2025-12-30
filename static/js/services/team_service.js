/**
 * Team Service
 */
const TeamService = {
  /**
   * List all teams with pagination
   * @param {number} limit
   * @param {number} offset
   * @returns {Promise}
   */
  listTeams: function (limit = 10, offset = 0) {
    return API.get(`/api/teams?limit=${limit}&offset=${offset}`);
  },

  /**
   * Get team details
   * @param {number} id
   * @returns {Promise}
   */
  getTeamDetails: function (id) {
    return API.get(`/api/teams/${id}`);
  },

  /**
   * Get team members with pagination
   * @param {number} id
   * @param {number} limit
   * @param {number} offset
   * @returns {Promise}
   */
  getTeamMembers: function (id, limit = 10, offset = 0) {
    return API.get(`/api/teams/${id}/members?limit=${limit}&offset=${offset}`);
  },
};
