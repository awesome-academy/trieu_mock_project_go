/**
 * Common notification logic for all user pages
 */
const NotificationApp = {
  socket: null,
  refreshTimer: null,

  /**
   * Load unread notification count and update the navbar
   */
  loadUnreadCount: function (delay = 0) {
    const token = localStorage.getItem("accessToken");
    if (!token) return;

    if (this.refreshTimer) {
      clearTimeout(this.refreshTimer);
    }

    this.refreshTimer = setTimeout(() => {
      API.get("/api/notifications/unread-count")
        .done((response) => {
          this.updateBadge(response.unread_count);
        })
        .fail((error) => {
          console.error("Failed to load unread notification count:", error);
        });
    }, delay);
  },

  /**
   * Update the unread notification badge
   */
  updateBadge: function (count) {
    const $badge = $("#unread-notifications-count");
    if (count > 0) {
      $badge.text(count).show();
    } else {
      $badge.hide();
    }
  },

  /**
   * Initialize web socket connection
   */
  initWebSocket: function () {
    const token = localStorage.getItem("accessToken");
    if (!token) return;

    API.post("/api/ws-ticket")
      .done((response) => {
        const ticket = response.ticket;
        const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
        const wsUrl = `${protocol}//${window.location.host}/ws?ticket=${ticket}`;

        this.socket = new WebSocket(wsUrl);

        this.socket.onopen = () => {
          console.log("WebSocket connected for notifications");
        };

        this.socket.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data);
            this.showToast(data.title, data.content);

            // Refresh count with debounce and delay to handle DB transaction race condition
            this.loadUnreadCount(500);

            // Notify other components (like the notifications list page)
            // Using a small delay to ensure DB transaction is committed before components refresh
            setTimeout(() => {
              window.dispatchEvent(
                new CustomEvent("notificationReceived", { detail: data })
              );
            }, 500);
          } catch (err) {
            console.error("Error parsing WebSocket message:", err);
          }
        };

        this.socket.onclose = (event) => {
          console.log(
            "WebSocket connection closed. Attempting to reconnect..."
          );
          setTimeout(() => {
            this.initWebSocket();
          }, 5000); // Reconnect after 5 seconds
        };

        this.socket.onerror = (error) => {
          console.error("WebSocket error:", error);
          this.socket.close();
        };
      })
      .fail((error) => {
        console.error("Failed to generate WebSocket ticket:", error);
      });
  },

  /**
   * Show a toast notification
   */
  showToast: function (title, content) {
    const toastId =
      "toast-" + Date.now() + "-" + Math.floor(Math.random() * 1000);
    const toastHtml = `
      <div id="${toastId}" class="toast mb-2" role="alert" aria-live="assertive" aria-atomic="true">
        <div class="toast-header">
          <i class="bi bi-bell-fill me-2 text-primary"></i>
          <strong class="me-auto">${title}</strong>
          <small class="text-muted">just now</small>
          <button type="button" class="btn-close" data-bs-dismiss="toast" aria-label="Close"></button>
        </div>
        <div class="toast-body">
          ${content}
        </div>
      </div>
    `;

    $(".toast-container").prepend(toastHtml);
    const toastElement = document.getElementById(toastId);
    if (toastElement) {
      const toast = new bootstrap.Toast(toastElement, {
        delay: 10000,
        autohide: true,
      });
      toast.show();

      // Remove toast element from DOM after it's hidden
      toastElement.addEventListener("hidden.bs.toast", () => {
        toastElement.remove();
      });
    }
  },

  /**
   * Initialize notification logic
   */
  init: function () {
    $(document).ready(() => {
      const path = window.location.pathname;
      if (path === "/login" || path.startsWith("/admin")) {
        return;
      }

      this.loadUnreadCount();
      this.initWebSocket();
    });
  },
};

// Start notification monitoring
NotificationApp.init();
