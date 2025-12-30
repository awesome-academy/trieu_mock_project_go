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
};
