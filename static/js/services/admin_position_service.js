/**
 * Admin Position Service
 */
const AdminPositionService = {
  /**
   * Search positions with pagination
   * @param {Object} params - { limit, offset }
   * @returns {Promise}
   */
  searchPositions: function (params) {
    let url = `/admin/positions/partial/search?limit=${
      params.limit || 10
    }&offset=${params.offset || 0}`;
    return AdminAPI.get(url, { dataType: "html" });
  },

  /**
   * Create a new position
   * @param {Object} data
   * @returns {Promise}
   */
  createPosition: function (data) {
    return AdminAPI.post("/admin/positions", data);
  },

  /**
   * Update an existing position
   * @param {number|string} positionId
   * @param {Object} data
   * @returns {Promise}
   */
  updatePosition: function (positionId, data) {
    return AdminAPI.put(`/admin/positions/${positionId}`, data);
  },

  /**
   * Delete a position
   * @param {number|string} positionId
   * @returns {Promise}
   */
  deletePosition: function (positionId) {
    return AdminAPI.delete(`/admin/positions/${positionId}`);
  },
};
