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
    return fetch(url).then((res) => {
      if (!res.ok) throw new Error("Failed to fetch positions");
      return res.text();
    });
  },

  /**
   * Create a new position
   * @param {Object} data
   * @returns {Promise}
   */
  createPosition: function (data) {
    return fetch("/admin/positions", {
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
   * Update an existing position
   * @param {number|string} positionId
   * @param {Object} data
   * @returns {Promise}
   */
  updatePosition: function (positionId, data) {
    return fetch(`/admin/positions/${positionId}`, {
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
   * Delete a position
   * @param {number|string} positionId
   * @returns {Promise}
   */
  deletePosition: function (positionId) {
    return fetch(`/admin/positions/${positionId}`, {
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
