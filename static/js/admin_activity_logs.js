document.addEventListener("DOMContentLoaded", function () {
  const activityLogListContainer = document.getElementById(
    "activityLogListContainer"
  );
  const loadingTemplate = document.getElementById("loadingTemplate");

  async function loadActivityLogs(offset = 0) {
    const limit = 10;

    // Show loading spinner
    activityLogListContainer.innerHTML = loadingTemplate.innerHTML;

    try {
      const html = await AdminActivityLogService.searchActivityLogs({
        limit,
        offset,
      });
      activityLogListContainer.innerHTML = html;
      attachEvents();
    } catch (error) {
      console.error("Error loading activity logs:", error);
      Toast.error("Failed to load activity logs list");
      activityLogListContainer.innerHTML =
        '<div class="alert alert-danger">Failed to load activity logs.</div>';
    }
  }

  function attachEvents() {
    // Pagination events
    const paginationLinks =
      activityLogListContainer.querySelectorAll(".page-link");
    paginationLinks.forEach((link) => {
      link.addEventListener("click", function (e) {
        e.preventDefault();
        const offsetAttr = this.getAttribute("data-offset");
        if (offsetAttr !== null) {
          const offset = parseInt(offsetAttr, 10);
          if (!Number.isNaN(offset) && offset >= 0) {
            loadActivityLogs(offset);
          }
        }
      });
    });

    // Delete events
    const deleteBtns =
      activityLogListContainer.querySelectorAll(".delete-log-btn");
    deleteBtns.forEach((btn) => {
      btn.addEventListener("click", async function () {
        const id = this.getAttribute("data-id");

        if (
          confirm(`Are you sure you want to delete activity log ID "${id}"?`)
        ) {
          try {
            const response = await AdminActivityLogService.deleteActivityLog(
              id
            );
            Toast.success(
              response.message || "Activity log deleted successfully"
            );
            loadActivityLogs(0);
          } catch (error) {
            console.error("Error deleting activity log:", error);
            Toast.error(error.message || "Failed to delete activity log");
          }
        }
      });
    });
  }

  // Initial load
  loadActivityLogs(0);
});
