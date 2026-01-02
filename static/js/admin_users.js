document.addEventListener("DOMContentLoaded", function () {
  const teamFilter = document.getElementById("teamFilter");
  const searchBtn = document.getElementById("searchBtn");
  const userListContainer = document.getElementById("userListContainer");
  const loadingTemplate = document.getElementById("loadingTemplate");

  async function loadUsers(offset = 0) {
    const teamId = teamFilter.value;
    const limit = 10;

    // Show loading spinner
    userListContainer.innerHTML = loadingTemplate.innerHTML;

    try {
      const html = await AdminUserService.searchUsers({
        limit,
        offset,
        team_id: teamId,
      });
      userListContainer.innerHTML = html;
      attachPaginationEvents();
    } catch (error) {
      console.error("Error loading users:", error);
      Toast.error("Failed to load users list");
      userListContainer.innerHTML =
        '<div class="alert alert-danger">Failed to load users.</div>';
    }
  }

  function attachPaginationEvents() {
    const paginationLinks = userListContainer.querySelectorAll(".page-link");
    paginationLinks.forEach((link) => {
      link.addEventListener("click", function (e) {
        e.preventDefault();

        const offsetAttr = this.getAttribute("data-offset");
        if (offsetAttr !== null) {
          const offset = parseInt(offsetAttr, 10);
          if (!Number.isNaN(offset) && offset >= 0) {
            loadUsers(offset);
          }
        }
      });
    });
  }

  searchBtn.addEventListener("click", () => loadUsers(0));
  teamFilter.addEventListener("change", () => loadUsers(0));

  // Initial load
  loadUsers(0);
});
