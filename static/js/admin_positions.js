document.addEventListener("DOMContentLoaded", function () {
  const positionListContainer = document.getElementById(
    "positionListContainer"
  );
  const loadingTemplate = document.getElementById("loadingTemplate");

  async function loadPositions(offset = 0) {
    const limit = 10;

    // Show loading spinner
    positionListContainer.innerHTML = loadingTemplate.innerHTML;

    try {
      const html = await AdminPositionService.searchPositions({
        limit,
        offset,
      });
      positionListContainer.innerHTML = html;
      attachEvents();
    } catch (error) {
      console.error("Error loading positions:", error);
      Toast.error("Failed to load positions list");
      positionListContainer.innerHTML =
        '<div class="alert alert-danger">Failed to load positions.</div>';
    }
  }

  function attachEvents() {
    // Pagination events
    const paginationLinks =
      positionListContainer.querySelectorAll(".page-link");
    paginationLinks.forEach((link) => {
      link.addEventListener("click", function (e) {
        e.preventDefault();
        const offset = this.getAttribute("data-offset");
        if (offset !== null) {
          loadPositions(parseInt(offset));
        }
      });
    });

    // Delete events
    const deleteBtns = positionListContainer.querySelectorAll(
      ".delete-position-btn"
    );
    deleteBtns.forEach((btn) => {
      btn.addEventListener("click", async function () {
        const id = this.getAttribute("data-id");
        const name = this.getAttribute("data-name");

        if (confirm(`Are you sure you want to delete position "${name}"?`)) {
          try {
            const response = await AdminPositionService.deletePosition(id);
            Toast.success(response.message || "Position deleted successfully");
            loadPositions(0);
          } catch (error) {
            console.error("Error deleting position:", error);
            Toast.error(error.message || "Failed to delete position");
          }
        }
      });
    });
  }

  // Initial load
  loadPositions(0);
});
