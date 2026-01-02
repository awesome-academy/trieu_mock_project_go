/**
 * Admin Skill Service
 */
const AdminSkillService = {
  /**
   * Search skills with pagination
   * @param {Object} params - { limit, offset }
   * @returns {Promise}
   */
  searchSkills: function (params) {
    let url = `/admin/skills/partial/search?limit=${
      params.limit || 10
    }&offset=${params.offset || 0}`;
    return AdminAPI.get(url, { dataType: "html" });
  },

  /**
   * Create a new skill
   * @param {Object} data
   * @returns {Promise}
   */
  createSkill: function (data) {
    return AdminAPI.post("/admin/skills", data);
  },

  /**
   * Update an existing skill
   * @param {number|string} skillId
   * @param {Object} data
   * @returns {Promise}
   */
  updateSkill: function (skillId, data) {
    return AdminAPI.put(`/admin/skills/${skillId}`, data);
  },

  /**
   * Delete a skill
   * @param {number|string} skillId
   * @returns {Promise}
   */
  deleteSkill: function (skillId) {
    return AdminAPI.delete(`/admin/skills/${skillId}`);
  },
};
