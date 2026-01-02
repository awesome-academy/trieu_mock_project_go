/**
 * Admin Activity Log Service
 */
const AdminActivityLogService = {
  /**
   * Search activity logs with pagination
   * @param {Object} params - { limit, offset }
   * @returns {Promise}
   */
  searchActivityLogs: function (params) {
    let url = `/admin/activity-logs/partial/search?limit=${
      params.limit || 10
    }&offset=${params.offset || 0}`;
    return AdminAPI.get(url, { dataType: "html" });
  },

  /**
   * Delete an activity log
   * @param {number|string} activityLogId
   * @returns {Promise}
   */
  deleteActivityLog: function (activityLogId) {
    return AdminAPI.delete(`/admin/activity-logs/${activityLogId}`);
  },
};
