/**
 * Common notification logic for all user pages
 */
const NotificationApp = {
  /**
   * Load unread notification count and update the navbar
   */
  loadUnreadCount: function () {
    const token = localStorage.getItem("accessToken");
    if (!token) return;

    API.get("/api/notifications/unread-count")
      .done((response) => {
        const count = response.unread_count || 0;
        const $badge = $("#unread-notifications-count");

        if (count > 0) {
          $badge.text(count).show();
        } else {
          $badge.hide();
        }
      })
      .fail((error) => {
        console.error("Failed to load unread notification count:", error);
      });
  },

  /**
   * Initialize notification polling or single load
   */
  init: function () {
    $(document).ready(() => {
      // Don't run on login/admin pages
      const path = window.location.pathname;
      if (path === "/login" || path.startsWith("/admin")) {
        return;
      }

      this.loadUnreadCount();

      // Refresh every 30 seconds
      setInterval(() => {
        this.loadUnreadCount();
      }, 30000);
    });
  },
};

// Start notification monitoring
NotificationApp.init();
