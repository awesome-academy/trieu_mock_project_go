/**
 * Notification Service
 */
const NotificationService = {
  /**
   * List all notifications with pagination
   * @param {number} limit
   * @param {number} offset
   * @returns {Promise}
   */
  getNotifications: function (limit = 10, offset = 0) {
    return API.get(`/api/notifications?limit=${limit}&offset=${offset}`);
  },

  /**
   * Mark notification as read
   * @param {number} id
   * @returns {Promise}
   */
  markAsRead: function (id) {
    return API.put(`/api/notifications/${id}/read`);
  },

  /**
   * Mark all notifications as read
   * @returns {Promise}
   */
  markAllAsRead: function () {
    return API.put(`/api/notifications/read-all`);
  },

  /**
   * Delete notification
   * @param {number} id
   * @returns {Promise}
   */
  deleteNotification: function (id) {
    return API.delete(`/api/notifications/${id}`);
  },
};
