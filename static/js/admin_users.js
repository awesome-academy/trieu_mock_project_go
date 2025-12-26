document.addEventListener("DOMContentLoaded", function () {
  const teamFilter = document.getElementById("teamFilter");
  const searchBtn = document.getElementById("searchBtn");
  const userListContainer = document.getElementById("userListContainer");

  function loadUsers(offset = 0) {
    const teamId = teamFilter.value;
    const limit = 10;
    let url = `/admin/users/partial/search?limit=${limit}&offset=${offset}`;
    if (teamId) {
      url += `&team_id=${teamId}`;
    }

    fetch(url)
      .then((response) => response.text())
      .then((html) => {
        userListContainer.innerHTML = html;
        attachPaginationEvents();
      })
      .catch((error) => {
        console.error("Error loading users:", error);
        userListContainer.innerHTML =
          '<div class="alert alert-danger">Failed to load users.</div>';
      });
  }

  function attachPaginationEvents() {
    const paginationLinks = userListContainer.querySelectorAll(".page-link");
    paginationLinks.forEach((link) => {
      link.addEventListener("click", function (e) {
        e.preventDefault();
        const offset = this.getAttribute("data-offset");
        if (offset !== null) {
          loadUsers(parseInt(offset));
        }
      });
    });
  }

  searchBtn.addEventListener("click", function () {
    loadUsers(0);
  });

  teamFilter.addEventListener("change", function () {
    loadUsers(0);
  });

  // Initial load
  loadUsers(0);
});
