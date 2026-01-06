document.addEventListener("DOMContentLoaded", async () => {
  const listElement = document.getElementById("notifications-list");
  const markAllBtn = document.getElementById("mark-all-read");
  let currentNotifications = [];
  let selectedNotificationID = null;

  const modal = new bootstrap.Modal(
    document.getElementById("notificationModal")
  );
  const modalTitle = document.getElementById("modal-title");
  const modalContent = document.getElementById("modal-content");
  const modalTime = document.getElementById("modal-time");
  const deleteBtn = document.getElementById("delete-notification");

  async function loadNotifications() {
    try {
      const data = await NotificationService.getNotifications();
      currentNotifications = data.notifications;
      renderNotifications(currentNotifications);
    } catch (error) {
      console.error(error);
      listElement.innerHTML =
        '<div class="alert alert-danger">Error loading notifications.</div>';
    }
  }

  function renderNotifications(notifications) {
    if (notifications.length === 0) {
      listElement.innerHTML =
        '<div class="text-center p-5 text-muted">No notifications found.</div>';
      return;
    }

    listElement.innerHTML = notifications
      .map(
        (n) => `
            <a href="#" class="list-group-item list-group-item-action notification-card ${
              n.is_read ? "" : "notification-unread"
            }" 
               data-id="${n.id}">
                <div class="d-flex w-100 justify-content-between">
                    <h5 class="mb-1">${n.title}</h5>
                    <small class="notification-time">${new Date(
                      n.created_at
                    ).toLocaleString()}</small>
                </div>
                <p class="mb-1 text-truncate">${n.content}</p>
            </a>
        `
      )
      .join("");

    document.querySelectorAll(".notification-card").forEach((card) => {
      card.addEventListener("click", (e) => {
        e.preventDefault();
        const id = card.getAttribute("data-id");
        const notification = currentNotifications.find((n) => n.id == id);
        showNotificationDetail(notification);
      });
    });
  }

  async function showNotificationDetail(n) {
    selectedNotificationID = n.id;
    modalTitle.textContent = n.title;
    modalContent.textContent = n.content;
    modalTime.textContent = new Date(n.created_at).toLocaleString();

    modal.show();

    if (!n.is_read) {
      try {
        await NotificationService.markAsRead(n.id);
        n.is_read = true;
        renderNotifications(currentNotifications);
      } catch (error) {
        console.error("Failed to mark as read", error);
      }
    }
  }

  markAllBtn.addEventListener("click", async () => {
    try {
      await NotificationService.markAllAsRead();
      await loadNotifications();
    } catch (error) {
      alert("Failed to mark all as read");
    }
  });

  deleteBtn.addEventListener("click", async () => {
    if (!selectedNotificationID) return;
    if (confirm("Are you sure you want to delete this notification?")) {
      try {
        await NotificationService.deleteNotification(selectedNotificationID);
        modal.hide();
        await loadNotifications();
      } catch (error) {
        alert("Failed to delete notification");
      }
    }
  });

  await loadNotifications();
});
