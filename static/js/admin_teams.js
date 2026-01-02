document.addEventListener("DOMContentLoaded", function () {
  const teamListContainer = document.getElementById("teamListContainer");
  const loadingTemplate = document.getElementById("loadingTemplate");

  async function loadTeams(offset = 0) {
    teamListContainer.innerHTML = loadingTemplate.innerHTML;
    try {
      const html = await AdminTeamService.listTeams({
        limit: 10,
        offset: offset,
      });
      teamListContainer.innerHTML = html;
      attachEventListeners();
    } catch (error) {
      teamListContainer.innerHTML = `<div class="alert alert-danger">${error.message}</div>`;
    }
  }

  function attachEventListeners() {
    // Pagination
    document.querySelectorAll(".page-link").forEach((link) => {
      link.addEventListener("click", function (e) {
        e.preventDefault();
        const offset = this.getAttribute("data-offset");
        if (offset !== null) {
          loadTeams(parseInt(offset));
        }
      });
    });

    // Delete buttons
    document.querySelectorAll(".delete-team-btn").forEach((btn) => {
      btn.addEventListener("click", async function () {
        const id = this.getAttribute("data-id");
        const name = this.getAttribute("data-name");

        if (confirm(`Are you sure you want to delete team "${name}"?`)) {
          try {
            await AdminTeamService.deleteTeam(id);
            loadTeams();
          } catch (error) {
            alert(error.message);
          }
        }
      });
    });
  }

  loadTeams();
});
